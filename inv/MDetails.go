package inv

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/jinzhu/gorm"
)

// Detail : 消費發票明細
// 明細=D	發票號碼 小計 品項名稱
// 範例：
// D ZZ00000050 42.00 拿鐵熱咖啡(中)
// D ZZ00000050 55.00 拿鐵冰咖啡(大)
type Detail struct {
	// auto-populate columns: id, created_at, updated_at, deleted_at
	// gorm.Model
	// Or alternatively write:
	Model gorm.Model `json:"-" gorm:"embedded"`
	// ID       int     `json:"-" sql:"AUTO_INCREMENT" gorm:"primary_key"`
	Head     string  `cht:"表頭" json:"head" sql:"DEFAULT:'D'"`
	UINumber string  `cht:"發票號碼" json:"uniform_invoice_number" sql:"size:10;index" gorm:"column:uin"`
	Subtotal float64 `cht:"小計" json:"subtotal_amount"`
	Name     string  `cht:"品項名稱" json:"name"`
	// Invoice  *Invoice   `json:"-"`
}

func (d Detail) String() string {
	Sf, Ff := fmt.Sprintf, fmt.Fprintf
	var b bytes.Buffer
	val := reflect.ValueOf(d) //.Elem()
	fld := val.Type()
	var str string
	for i := 0; i < val.NumField(); i++ {
		switch fld.Field(i).Name {
		case "Model":
			continue
		case "Subtotal":
			str = Sf("%.1f", val.Field(i).Interface().(float64))
		case "UINumber":
			str = val.Field(i).Interface().(string)[0:2] + "-" + val.Field(i).Interface().(string)[2:]
		default:
			str = val.Field(i).Interface().(string)
		}
		Ff(&b, " %s : %s |", fld.Field(i).Tag.Get("cht"), str)
	}
	Ff(&b, "\n")
	return b.String()
}

// GetArgsTable :
func (d *Detail) GetArgsTable(title string, lensp int) string {
	Sf := fmt.Sprintf
	if len(title) == 0 {
		title = "明細清單"
	}
	// dheads := []string{"表頭", "發票號碼", "小計", "品項名稱"}
	_, _, _, dheads := GetFieldsInfo(Detail{}, "cht", "Model")
	if lensp < 0 {
		lensp = 0
	}
	table := ArgsTableN(title, lensp, dheads, d.Head,
		d.UINumber[0:2]+"-", d.UINumber[2:], Sf("%.1f", d.Subtotal), d.Name)
	return table
	// var title string
	// if len(args) > 0 {
	// 	title = args[0]
	// } else {
	// 	title = "Detail"
	// }
	// tab := ArgsTable(
	// 	title,
	// 	"表頭", "Head", d.Head,
	// 	"發票號碼", "UINumber", d.UINumber,
	// 	"小計", "Subtotal", io.Sf("%.1f", d.Subtotal),
	// 	"品項名稱", "Name", d.Name,
	// )
	// return tab
}

// TableName : set Detail's table name to be `details`
func (Detail) TableName() string {
	// custom table name, this is default
	return "details"
}

// GetDetailsTable returns the table string of the list of []*Detail
func GetDetailsTable(pds []*Detail, lensp int) string {
	Sf := fmt.Sprintf
	title := "明細清單"
	dheads := []string{"項次"} //, "表頭", "發票號碼", "小計", "品項名稱"}
	_, _, _, tmp := GetFieldsInfo(Detail{}, "cht", "Model")
	dheads = append(dheads, tmp...)
	if lensp < 0 {
		lensp = 0
	}
	var data []interface{}
	for i, d := range pds {
		data = append(data, i+1, d.Head,
			d.UINumber[0:2]+"-"+d.UINumber[2:], Sf("%.1f", d.Subtotal), d.Name)
	}
	table := ArgsTableN(title, lensp, dheads, data...)
	return table
}
