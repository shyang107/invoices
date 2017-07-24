package main

import (
	"log"
	"time"

	"./inv"

	gio "github.com/cpmech/gosl/io"
	_ "github.com/mattn/go-sqlite3"
)

// Version :
var Version = "v0.0.3"

func init() {
	log.SetPrefix("LOG: ")
	// log.SetFlags(log.Ldate | log.Lmicroseconds | log.Llongfile)
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)
	// log.Println("init started")
	// io.Verbose = true
	gio.Verbose = false
}

func main() {
	start := time.Now()
	c := inv.NewConfig()
	c.ReadCommandLine()
	//
	// inv.ConfigCmds(Version)
	// //
	// if inv.CFG.IsInitializing {
	// 	inv.InitDB()
	// }
	// initDb()
	// //
	// inv.Opt.GetOptions()
	// //
	// // if inv.CFG.Dump {
	// // 	if err := inv.DumpData(); err != nil {
	// // 		panic(err)
	// // 	}
	// // } else {
	// //
	// pvs, err := inv.ReadInvoices()
	// checkErr(err)
	// if inv.Opt.IsOutput {
	// 	err = inv.WriteInvoices(pvs)
	// 	checkErr(err)
	// }
	// }
	// pfields()
	duration := time.Since(start) //.Seconds()
	gio.Pf("run-time elapsed : %v\n", duration)
}

// func initDb() {
// 	//初始化并保持连接
// 	var err error
// 	inv.DB, err = gorm.Open("sqlite3", inv.CFG.DBPath)
// 	//    DB.LogMode(true)//打印sql语句
// 	if err != nil {
// 		log.Fatalf("database connect is err: %s", err.Error())
// 	} else {
// 		// log.Print("connect database is success")
// 		io.Pfyel("* connect database is success\n")
// 	}
// 	err = inv.DB.DB().Ping()
// 	if err != nil {
// 		inv.DB.DB().Close()
// 		log.Fatalf("Error on opening database connection: %s", err.Error())
// 	}
// 	inv.DB.Model(&inv.Invoice{}).Related(&inv.Detail{}, "uin")
// }

// func checkErr(err error) {
// 	if err != nil {
// 		panic(err)
// 	}
// }
