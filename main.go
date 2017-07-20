package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	"./inv"
	"github.com/cpmech/gosl/io"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

// Version :
var Version = "v0.0.3"

func init() {
	log.SetPrefix("LOG: ")
	// log.SetFlags(log.Ldate | log.Lmicroseconds | log.Llongfile)
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)
	// log.Println("init started")
	io.Verbose = true
}

func main() {
	start := time.Now()
	//
	inv.ConfigCmds(Version)
	//
	inv.Opts.GetOptions()
	//
	if inv.Cfg.IsInitializing {
		inv.InitDB()
	}
	initDb()
	//
	pvs, err := inv.ReadInvoices()
	checkErr(err)
	if inv.Opts.IsOutput {
		err = inv.WriteInvoices(pvs)
		checkErr(err)
	}
	// pfields()
	duration := time.Since(start) //.Seconds()
	fmt.Fprintf(os.Stdout, "run-time elapsed : %v\n", duration)
}

func initDb() {
	//初始化并保持连接
	var err error
	inv.DB, err = gorm.Open("sqlite3", inv.Cfg.DBPath)
	//    DB.LogMode(true)//打印sql语句
	if err != nil {
		log.Fatalf("database connect is err: %s", err.Error())
	} else {
		// log.Print("connect database is success")
		io.Pfyel("* connect database is success\n")
	}
	err = inv.DB.DB().Ping()
	if err != nil {
		inv.DB.DB().Close()
		log.Fatalf("Error on opening database connection: %s", err.Error())
	}
	inv.DB.Model(&inv.Invoice{}).Related(&inv.Detail{}, "uin")
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func pfields() {
	print(inv.Invoice{}, "cht")
	io.Pf("%v", io.StrThinLine(60))
	print(inv.Detail{}, "cht")
}

func print(obj interface{}, tag string) {
	fields, types, kinds, tags := inv.GetFieldsInfo(obj, tag)
	n := len(fields)
	nc, ne := inv.CountChEngChar(tags[0])
	sizes := []int{len(fields[0]), len(types[0]), len(kinds[0]), nc*2 + ne, nc + ne}
	for i := 1; i < n; i++ {
		sizes[0] = imax(sizes[0], len(fields[i]))
		sizes[1] = imax(sizes[1], len(types[i]))
		sizes[2] = imax(sizes[2], len(kinds[i]))
		nc, ne = inv.CountChEngChar(tags[i])
		sizes[3] = imax(sizes[3], nc*2+ne) //
		sizes[4] = imax(sizes[4], nc+ne)
	}
	io.Pf("%v\n", sizes)
	m := sizes[0] + sizes[1] + sizes[2] + sizes[3] + 6
	tab := io.StrThickLine(m)
	tab += io.Sf("%*s  %*s  %*s  %*s\n",
		sizes[0], "Field", sizes[1], "Type", sizes[2], "Kind", sizes[3], "Tag")
	tab += io.StrThinLine(m)
	for i := 0; i < n; i++ {
		stag := io.Sf("%*s", sizes[4], tags[i])
		nc, ne = inv.CountChEngChar(stag)
		// if nc < sizes[3] {
		// 	stag = io.StrSpaces((sizes[3]-nc)*2) + tags[i]
		// }
		tab += io.Sf("%*s  %*s  %*s  %*s\n",
			sizes[0], fields[i], sizes[1], types[i], sizes[2], kinds[i], nc+ne*2, tags[i])
	}
	tab += io.StrThickLine(m)
	io.Pf("%s", tab)
}

func imax(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func testch() {
	str := "中文文文文"
	io.Pf("str = %q -- %d\n", str, len(str))
	for i := 0; i < len(str); i++ {
		io.Pf("%q ", str[i])
	}
	io.Pl()
	var hzRegexp = regexp.MustCompile("^[\u4e00-\u9fa5]{3,8}$")
	fmt.Println(hzRegexp.MatchString(str))
	fmt.Println(inv.IsChineseChar(str))
	//
	str = "abc 123,_-中文文文文"
	io.Pf("str = %q -- %d\n", str, len(str))
	c, e := inv.CountChEngChar(str)
	io.Pl()
	io.Pf("nc = %d | ne = %d\n", c, e)
	l := c*2 + e
	tab := io.StrThinLine(l)
	tab += io.Sf("%*s\n", c+e, str)
	tab += io.StrThinLine(l)
	io.Pf("%s\n", tab)
}
