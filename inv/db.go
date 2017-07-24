package inv

import (
	"os"

	"github.com/jinzhu/gorm"
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
func GetInvoiceList() (invList []Invoice, err error) {
	err = DB.Find(&invList).Error
	if err == nil {
		for i := 0; i < len(invList); i++ {
			DB.Model(&invList[i]).Association("details").Find(&invList[i].Details)
		}
	}
	return invList, err
}

// InsertFrom creats records from []*Invoice into database
func InsertFrom(pvs []*Invoice) {
	for _, v := range pvs {
		// io.Pforan("# %v", *v)
		// DB.FirstOrCreate(v)
		DB.Where(Invoice{UINumber: v.UINumber}).FirstOrCreate(v)
	}
}
