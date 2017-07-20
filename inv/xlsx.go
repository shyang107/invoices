package inv

import (
	"reflect"
	"time"

	"github.com/cpmech/gosl/chk"
	"github.com/stanim/xlsxtra"
	// "github.com/stanim/xlsxtra"
	"github.com/tealeg/xlsx"
)

const (
	black          = "FF000000"
	white          = "FFFFFFFF"
	red            = "FFFF0000"
	blue           = "FF0000FF"
	yellow         = "FFFFFF00"
	green          = "FF008000"
	pink           = "FFFF00FF"
	turquoise      = "FF00FFFF" // cyan
	darkRed        = "FF800000"
	darkBlue       = "FF000080"
	darkYellow     = "FF808000"
	darkPurple     = "FF660066"
	oceanBlue      = "FF0066CC"
	violet         = "FF800080"
	teal           = "FF008080"
	gray25         = "FFC0C0C0"
	gray40         = "FF969696"
	gray50         = "FF808080"
	gray80         = "FF333333"
	periwinkle     = "FF993366"
	ivory          = "FFFFFFCC"
	coral          = "FFFF8080"
	brightGreen    = "FF00FF00"
	lightGreen     = "FFCCFFCC"
	iceBlue        = "FFCCCCFF"
	lightBlue      = "FF3366FF"
	lightTurquoise = "FFCCFFFF" // light cyan
	lightYellow    = "FFFFFF99"
	//
	numfmtAccountant = `_($* #,##0.0_);_($* (#,##0.0);_($* "-"??_);_(@_)`
	numfmtDollar     = `"NT$"#,##0.0_);[red]"NT$"-#,##0.0`
	numfmt           = `#,##0 ;[red]-#,##0`
)

var (
	letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	numbers = []rune("0123456789")
)

// XlsMarshaller :
type XlsMarshaller struct{}

// MarshalInvoices marshal the records of invoice using in .xls file
func (XlsMarshaller) MarshalInvoices(fn string, pvs []*Invoice) error {
	prun("  > Writing data to .xlsx file %q ...\n", fn)
	if pvs == nil || len(pvs) == 0 {
		return chk.Err("pvs []*Invoice = nil or it's len = 0 ")
	}
	var vh, dh headType
	_, vh.head = getFieldNameAndChtag(Invoice{})
	_, dh.head = getFieldNameAndChtag(Detail{})
	//
	fx := xlsx.NewFile()
	s, _ := fx.AddSheet("消費發票")
	sht := &xlsxtra.Sheet{Sheet: s}
	for i := 0; i < len(pvs); i++ {
		vh.addTo(sht.AddRow(), false)
		rowi := sht.AddRow()
		pvs[i].addTo(rowi, i+1)
		if len(pvs[i].Details) > 0 {
			dh.addTo(sht.AddRow(), true)
			for j := 0; j < len(pvs[i].Details); j++ {
				rowd := sht.AddRow()
				pvs[i].Details[j].addTo(rowd, j+1)
			}
		}
	}
	fx.Save(fn)
	return nil
}

type headType struct {
	head []string
}

func (ht headType) addTo(r *xlsxtra.Row, isDetail bool) {
	border := xlsx.NewBorder("", "", "thin", "thin")
	style := xlsxtra.NewStyle(
		"",
		nil,
		border,
		nil,
	)
	if isDetail {
		r.AddString("")
	}
	cell := r.AddCell()
	cell.SetString("項次")
	cell.SetStyle(style)
	for i := 0; i < len(ht.head); i++ {
		// r.AddString(ht.head[i])
		cell := r.AddCell()
		cell.SetString(ht.head[i])
		cell.SetStyle(style)

	}
}

func getDefaultDetailCellStyle() *xlsx.Style {
	s := xlsx.NewStyle()
	fill := *xlsx.NewFill("solid", lightTurquoise, "")
	s.Fill = fill
	s.ApplyFill = true
	// s.Alignment.ShrinkToFit = true
	// s.ApplyAlignment = true
	border := *xlsx.NewBorder("", "", "thin", "thin")
	s.Border = border
	s.ApplyBorder = true
	return s
}

func (d *Detail) addTo(r *xlsxtra.Row, id int) {
	style := getDefaultDetailCellStyle()
	//
	r.AddString("")
	cell := r.AddCell()
	cell.SetStyle(style)
	cell.SetInt(id)
	//
	val := reflect.ValueOf(*d)
	typ := val.Type()
	n := typ.NumField()
	for i := 0; i < n; i++ {
		typi := typ.Field(i)
		typename := typi.Type.String()
		if typename == "gorm.Model" {
			continue
		}
		v := val.Field(i)
		cell := r.AddCell()
		cell.SetStyle(style)
		switch typename {
		case "float64":
			cell.SetFloatWithFormat(v.Interface().(float64), numfmtAccountant)
		default:
			cell.SetString(v.Interface().(string)) //(v.String())
		}
	}
}

func getDefaultInvoiceCellStyle() *xlsx.Style {
	s := xlsx.NewStyle()
	fill := *xlsx.NewFill("solid", lightGreen, "")
	s.Fill = fill
	s.ApplyFill = true
	// s.Alignment.ShrinkToFit = true
	// s.ApplyAlignment = true
	border := *xlsx.NewBorder("", "", "thin", "thin")
	s.Border = border
	s.ApplyBorder = true
	return s
}

func (v *Invoice) addTo(r *xlsxtra.Row, id int) {
	style := getDefaultInvoiceCellStyle()
	//
	r.AddInt(id)
	//
	val := reflect.ValueOf(*v)
	typ := val.Type()
	n := typ.NumField()
	for i := 0; i < n; i++ {
		typi := typ.Field(i).Type
		typename := typi.String()
		if typename == "gorm.Model" || typename == "[]*inv.Detail" {
			continue
		}
		vv := val.Field(i)
		// cell := r.AddCell()
		// cell.SetStyle(style)
		switch typename {
		case "time.Time":
			// cell.SetDate(v.Date)
			r.AddCell().SetDate(vv.Interface().(time.Time))
		case "float64":
			r.AddFloat(numfmtAccountant, vv.Interface().(float64))
		default:
			// cell.SetString(vv.String())
			r.AddString(vv.Interface().(string))
		}
	}
	r.SetStyle(style)
	return
}

func getFieldNameAndChtag(obj interface{}) (fldn, cfldn []string) {
	vv := reflect.ValueOf(obj)
	tv := vv.Type()
	for i := 0; i < tv.NumField(); i++ {
		field := tv.Field(i)
		typename := field.Type.String()
		switch typename {
		case "gorm.Model", "[]*inv.Detail":
			continue
		default:
			fldn = append(fldn, field.Name)
			cname := tv.Field(i).Tag.Get("cht")
			cfldn = append(cfldn, cname)
		}
	}
	return
}

// UnmarshalInvoices unmarshal the records of invoice using in .xls file
func (XlsMarshaller) UnmarshalInvoices(fn string) ([]*Invoice, error) {
	prun("  > Writing data to .xlsx file %q ...\n", fn)
	perr("!!! Warning !!! wating %q TODO ...\n", callerName(1))
	return nil, nil
}
