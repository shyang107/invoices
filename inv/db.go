package inv

import (
	"log"
	"os"

	"github.com/cpmech/gosl/io"
	"github.com/jinzhu/gorm"
	jsoniter "github.com/json-iterator/go"
	// use for sqlite
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// DB is database
var DB *gorm.DB

// initialdb initialize database
func initialdb() {
	io.Verbose = true
	if isExist(cfg.DBPath) {
		pstat("  > Removing file %q ...\n", cfg.DBPath)
		err := os.Remove(cfg.DBPath)
		if err != nil {
			panic(err)
		}
	}
	db, err := gorm.Open("sqlite3", os.ExpandEnv(cfg.DBPath))
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	//
	// Migrate the schema
	db.AutoMigrate(&Invoice{}, &Detail{}, &FileBunker{})
	db.Model(&Invoice{}).Related(&Detail{}, "uin")
	// db.Model(&Invoice{}).AddUniqueIndex("idx_invoices_number", "uin")
	// db.Model(&Invoice{}).AddForeignKey("uin", "details(id)", "RESTRICT", "RESTRICT")
}

func connectdb() {
	//初始化并保持连接
	var err error
	DB, err = gorm.Open("sqlite3", cfg.DBPath)
	//    DB.LogMode(true)//打印sql语句
	if err != nil {
		log.Fatalf("database connect is err: %s", err.Error())
	} else {
		// log.Print("connect database is success")
		io.Pfyel("* connect database is success\n")
	}
	err = DB.DB().Ping()
	if err != nil {
		DB.DB().Close()
		log.Fatalf("Error on opening database connection: %s", err.Error())
	}
	DB.Model(&Invoice{}).Related(&Detail{}, "uin")
}

// DBDumpInvoices get the list from database
func DBDumpInvoices() ([]Invoice, error) {
	// err = DB.Find(pinvs).Error
	// if err == nil {
	// 	n := DB.Count(&Invoice{}).
	// 	for i := 0; i < len(pinvs); i++ {
	// 		DB.Model(pinvs[i]).Association("details").Find(pinvs[i].Details)
	// 	}
	// }
	invs := []Invoice{}
	DB.Find(&invs)
	for i := range invs {
		// DB.Model(invs[i]).Related(&invs[i].UINumber)
		DB.Model(&invs[i]).Association("details").Find(&invs[i].Details)
	}
	return invs, nil
}

// DBInsertFrom creats records from []*Invoice into database
func DBInsertFrom(pvs []*Invoice) {
	for _, v := range pvs {
		// io.Pforan("# %v", *v)
		// DB.FirstOrCreate(v)
		DB.Where(Invoice{UINumber: v.UINumber}).FirstOrCreate(v)
	}
}

// DBDumpData dumps all data from db
func DBDumpData(dumpFilename string) error {
	pstat("  > Dumping data from database %q ...\n", cfg.DBPath)
	pvs, err := DBDumpInvoices()
	if err != nil {
		return err
	}
	// for i, p := range pvs {
	// 	io.Pfgreen2("Rec. %d : %v\n", i+1, p)
	// }
	fn := io.PathKey(dumpFilename) + ".json"
	pstat("  >> Marshall data in JSON-type, and then write to %q ...\n", fn)
	b, err := jsoniter.Marshal(&pvs)
	if err != nil {
		return err
	}
	io.WriteBytesToFile(fn, b)
	printSepline(60)
	return nil
}
