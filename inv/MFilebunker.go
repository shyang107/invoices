package inv

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/cpmech/gosl/io"
	"github.com/jinzhu/gorm"
)

// FileBunker use to backup original file of invoices
type FileBunker struct {
	Model    gorm.Model `json:"-" gorm:"embedded"`
	Name     string     `cht:"檔案名稱" json:"name"`
	Size     int        `cht:"檔案大小" json:"size"`
	ModAt    time.Time  `cht:"修改時間" json:"modtime_at" sql:"index"` // modification time
	Encoding string     `cht:"編碼" json:"encoding"`
	Checksum string     `cht:"檢查碼" json:"checksum"` //sha256
	Contents []byte     `cht:"內容" json:"-"`
}

func (f FileBunker) String() string {
	Sf := fmt.Sprintf
	location, _ := time.LoadLocation("Local")
	val := reflect.ValueOf(f) //.Elem()
	fld := val.Type()
	var str string
	var cols = make([]string, 0)
	for i := 0; i < val.NumField(); i++ {
		switch fld.Field(i).Name {
		case "Model":
			continue
		case "ModAt":
			str = Sf("%v", val.Field(i).Interface().(time.Time).In(location))
		case "Size":
			str = BytesSizeToString(val.Field(i).Interface().(int))
		case "Contents":
			str = "[略...]"
		default:
			// str = val.Field(i).Interface().(string)
			str = Sf("%v", val.Field(i).Interface().(string))
		}
		cols = append(cols, Sf("%s:%s", fld.Field(i).Tag.Get("cht"), str))
	}
	return strings.Join(cols, csvSep)
}

// TableName : set Detail's table name to be `details`
func (FileBunker) TableName() string {
	// custom table name, this is default
	return "filebunker"
}

// GetArgsTable :
func (f *FileBunker) GetArgsTable(title string, lensp int) string {
	// Sf := fmt.Sprintf
	location, _ := time.LoadLocation("Local")
	if len(title) == 0 {
		title = "原始發票檔案清單"
	}
	// var heads = []string{"項次"}
	_, _, _, heads := GetFieldsInfo(FileBunker{}, "cht", "Model")
	if lensp < 0 {
		lensp = 0
	}
	// heads = append(heads, tmp...)
	strSize := BytesSizeToString(f.Size)
	table := ArgsTableN(title, lensp, heads,
		f.Name, strSize, f.ModAt.In(location), f.Encoding, f.Checksum, "[略...]")
	return table
}

// GetFileBunkerTable returns the table string of the list of []*Detail
func GetFileBunkerTable(pfbs []*FileBunker, lensp int) string {
	Sf := fmt.Sprintf
	location, _ := time.LoadLocation("Local")
	title := "原始發票檔案清單"
	heads := []string{"項次"}
	_, _, _, tmp := GetFieldsInfo(FileBunker{}, "cht", "Model")
	heads = append(heads, tmp...)
	if lensp < 0 {
		lensp = 0
	}
	var data []interface{}
	for i, f := range pfbs {
		strSize := BytesSizeToString(f.Size)
		data = append(data, i+1,
			f.Name, strSize, Sf("%v", f.ModAt.In(location)), f.Encoding, f.Checksum, "[略...]")
	}
	table := ArgsTableN(title, lensp, heads, data...)
	return table
}

// UpdateFileBunker updates DB
func (o *Option) UpdateFileBunker() error {
	fi, err := os.Stat(o.InpFn)
	if err != nil {
		return err
	}
	if strings.ToLower(o.IfnSuffix) == ".csv" && strings.ToLower(o.IfnEncoding) == "big5" {
		b, err := io.ReadFile(o.InpFn)
		if err != nil {
			return err
		}
		sum := fmt.Sprintf("%x", sha256.Sum256(b))
		fn := filepath.Base(o.InpFn)
		fb := FileBunker{
			Name:     fn,
			Size:     int(fi.Size()),
			ModAt:    fi.ModTime(),
			Encoding: o.IfnEncoding,
			Checksum: sum,
			Contents: b,
		}
		DB.Where(&fb).FirstOrCreate(&fb)
		// fbs = append(fbs, &fb)
	}
	// plog((&fb).GetArgsTable("", 0))
	return nil
}
