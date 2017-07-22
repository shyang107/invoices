package inv

import (
	"os"

	"github.com/cpmech/gosl/io"
	"github.com/jinzhu/gorm"
	jsoniter "github.com/json-iterator/go"
	// use for sqlite
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// DB is database
var DB *gorm.DB

// InitDB initialize database
func InitDB() {
	os.Remove(Cfg.DBPath)
	db, err := gorm.Open("sqlite3", os.ExpandEnv(Cfg.DBPath))
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	//
	// Migrate the schema
	db.AutoMigrate(&Invoice{}, &Detail{})
	db.Model(&Invoice{}).Related(&Detail{}, "uin")
	// db.Model(&Invoice{}).AddUniqueIndex("idx_invoices_number", "uin")
	// db.Model(&Invoice{}).AddForeignKey("uin", "details(id)", "RESTRICT", "RESTRICT")
}

// GetInvoiceList get the list from database
func GetInvoiceList() ([]Invoice, error) {
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

// InsertFrom creats records from []*Invoice into database
func InsertFrom(pvs []*Invoice) {
	for _, v := range pvs {
		// io.Pforan("# %v", *v)
		// DB.FirstOrCreate(v)
		DB.Where(Invoice{UINumber: v.UINumber}).FirstOrCreate(v)
	}
}

// DumpData dumps all data from db
func DumpData() error {
	pstat("  > Dumping data from database %q ...\n", Cfg.DBPath)
	pvs, err := GetInvoiceList()
	if err != nil {
		return err
	}
	// for i, p := range pvs {
	// 	io.Pfgreen2("Rec. %d : %v\n", i+1, p)
	// }
	fn := io.PathKey(Opt.PunchFn) + ".json"
	pstat("  >> Marshall data in JSON-type, and then write to %q ...\n", fn)
	b, err := jsoniter.Marshal(&pvs)
	if err != nil {
		return err
	}
	io.WriteBytesToFile(fn, b)
	printSepline(60)
	return nil
}
