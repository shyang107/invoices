package main

import (
	"fmt"
	"log"
	"time"

	"./inv"

	gio "github.com/cpmech/gosl/io"
	_ "github.com/mattn/go-sqlite3"
)

// Version :
var Version = "v0.0.4"

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
	c.RunCommands()
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
	fmt.Printf("run-time elapsed : %v\n", time.Since(start))
}
