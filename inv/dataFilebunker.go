package inv

import (
	"bytes"
	"fmt"
	"reflect"
	"time"

	"github.com/jinzhu/gorm"
)

// FileBunker use to backup original file of invoices
type FileBunker struct {
	Model    gorm.Model
	Name     string    `cht:"檔案名稱"`
	ModTime  time.Time `cht:"修改時間" sql:"index"` // modification time
	Encoding string    `cht:"編碼"`
	Contents []byte    `cht:"內容"`
}

func (f FileBunker) String() string {
	Sf, Ff := fmt.Sprintf, fmt.Fprintf
	var b bytes.Buffer
	val := reflect.ValueOf(f) //.Elem()
	fld := val.Type()
	var str string
	for i := 0; i < val.NumField(); i++ {
		switch fld.Field(i).Name {
		case "Model", "Contents":
			continue
		default:
			// str = val.Field(i).Interface().(string)
			str = Sf("%v", val.Field(i).Interface().(string))
		}
		Ff(&b, " %s : %s |", fld.Field(i).Tag.Get("cht"), str)
	}
	Ff(&b, "\n")
	return b.String()
}

// TableName : set Detail's table name to be `details`
func (FileBunker) TableName() string {
	// custom table name, this is default
	return "file_bunker"
}

// GetArgsTable :
func (f *FileBunker) GetArgsTable(title string, lensp int) string {
	// Sf := fmt.Sprintf
	if len(title) == 0 {
		title = "原始發票檔案清單"
	}
	// var heads = []string{"項次"}
	_, _, _, heads := GetFieldsInfo(Detail{}, "cht", "Model", "Contents")
	if lensp < 0 {
		lensp = 0
	}
	// heads = append(heads, tmp...)
	table := ArgsTableN(title, lensp, heads, f.Name, f.ModTime, f.Encoding)
	return table
}

// GetFileBunkerTable returns the table string of the list of []*Detail
func GetFileBunkerTable(pfbs []*FileBunker, lensp int) string {
	Sf := fmt.Sprintf
	title := "原始發票檔案清單"
	heads := []string{"項次"} //, "表頭", "發票號碼", "小計", "品項名稱"}
	_, _, _, tmp := GetFieldsInfo(Detail{}, "cht", "Model", "Contents")
	heads = append(heads, tmp...)
	if lensp < 0 {
		lensp = 0
	}
	var data []interface{}
	for i, f := range pfbs {
		data = append(data, i+1, f.Name, Sf("%v", f.ModTime), f.Encoding)
	}
	table := ArgsTableN(title, lensp, heads, data...)
	return table
}
