package inv

import (
	"bytes"
	"fmt"
	"reflect"
	"time"

	"github.com/cpmech/gosl/io"
	"github.com/jinzhu/gorm"
)

const (
	appversion    = "0.0.2"
	fileType      = "INVOICES"   // using in text file
	magicNumber   = 0x125D       // using in binary file
	fileVesion    = 100          // using in all filetype
	dateFormat    = "20060102"   // allways using the date
	strDateFormat = "2006-01-02" // allways using the date
	//
	fcstart = 101
	fcstop  = 102
	fostart = 111
	fostop  = 112
	ffstart = 21
	ffstop  = 22
	csvSep  = "|"
)

var (
	// Opts : configuration
	// Opts   = DefaultOptions
	format = map[int]string{
		// config
		fcstart: "# Start to configure. -- %q\n",
		fcstop:  "# Configuration has been concluded. -- %q\n",
		// option
		fostart: "# Start to get case-options. -- %q\n",
		fostop:  "# Case-options has been concluded. -- %q\n",
		// start/end function
		ffstart: "* Function %q start.\n",
		ffstop:  "* Function %q stop.\n",
	}
)

// Invoice : 消費發票
// 表頭=M	發票狀態 發票號碼 發票日期 商店統編 商店店 載具名稱 載具號碼 總金額
// 範例：
// M 開立、作廢 ZZ00000050 20130111 97162640 新北市第1000號門市 手機條碼	/WYY+.,HG 97
type Invoice struct {
	// auto-populate columns: id, created_at, updated_at, deleted_at
	// gorm.Model
	// Or alternatively write:
	Model gorm.Model `json:"-" gorm:"embedded"`
	// ID    int    `json:"-" sql:"AUTO_INCREMENT" gorm:"primary_key"`
	Head  string `cht:"表頭" json:"head" sql:"DEFAULT:'M'"`
	State string `cht:"發票狀態" json:"state"`
	// Uniform-Invoice Number or  tax ID number
	UINumber string    `cht:"發票號碼" json:"uniform_invoice_number" sql:"size:10;unique;index" gorm:"column:uin"`
	Date     time.Time `cht:"發票日期" json:"purchase_date" sql:"index"`
	// Date    string     `cht:"發票日期" json:"date" sql:"index"`
	SUN     string  `cht:"商店統編" json:"store_uniform_number"`
	SName   string  `cht:"商店店名" json:"store_name"`
	CName   string  `cht:"載具名稱" json:"carrier_name"`
	CNumber string  `cht:"載具號碼" json:"carrier_number"`
	Total   float64 `cht:"總金額" json:"total_amount"`
	// one-to-many relationship
	Details []*Detail `cht:"明細清單" json:"Details" gorm:"ForeignKey:UINumber;AssociationForeignKey:UINumber"`
}

func (pv Invoice) String() string {
	var b bytes.Buffer
	ff := fmt.Fprintf
	ds := reflect.ValueOf(pv)
	t := ds.Type()
	nh := t.NumField()
	for i := 0; i < nh; i++ {
		field := t.Field(i)
		val := ds.Field(i).Interface()
		var v string
		switch val.(type) {
		case gorm.Model, []*Detail:
			continue
		case time.Time:
			v = io.Sf("%v", val.(time.Time).Format(strDateFormat))
		case float64:
			v = io.Sf("%.1f", val.(float64))
		default:
			if field.Name == "UINumber" {
				v = val.(string)[0:2] + "-" + val.(string)[2:]
				break
			}
			v = val.(string)
		}
		h := field.Tag.Get("cht")
		ff(&b, " %s : %s |", h, v)
	}
	ff(&b, "\n")
	lspaces := io.StrSpaces(4)
	for i, d := range pv.Details {
		ff(&b, "%s> %2d. %s", lspaces, i+1, d)
	}
	return b.String()
	// re, _ := regexp.Compile("^[\u4e00-\u9fa5]")
}

// TableName : set Invoice's table name to be `invoices`
func (Invoice) TableName() string {
	// custom table name, this is default
	return "invoices"
}

// GetArgsTable :
func (pv *Invoice) GetArgsTable(args ...string) string {
	var title string
	if len(args) > 0 {
		title = args[0]
	} else {
		title = "Invoice"
	}
	tab := ArgsTable(
		title,
		"表頭", "Head", pv.Head,
		"發票狀態", "State", pv.State,
		"發票號碼", "UINumber", pv.UINumber,
		"發票日期", "Date", pv.Date.Format(strDateFormat),
		"商店統編", "SUN", pv.SUN,
		"商店店名", "SName", pv.SName,
		"載具名稱", "CName", pv.CName,
		"載具號碼", "CNumber", pv.CNumber,
		"總金額", "Total", io.Sf("%.1f", pv.Total),
		"明細清單", "Details", "[]*Detail",
	)
	return tab
}

type invoiceSlcie struct {
	data    []string
	details []detailSlcie
}
type detailSlcie struct {
	data []string
}

func getInvoicesTable(pinvs []*Invoice) string {
	Sf, StrSpaces, StrThickLine, StrThinLine := io.Sf, io.StrSpaces, io.StrThickLine, io.StrThinLine
	vheads := []string{"項次", "表頭", "發票狀態", "發票號碼", "發票日期",
		"商店統編", "商店店名", "載具名稱", "載具號碼", "總金額"}
	dheads := []string{"項次", "表頭", "發票號碼", "小計", "品項名稱"}
	vnf := len(vheads)
	dnf := len(dheads)
	vsizes, vcsizes, vesizes, vismix :=
		make([]int, vnf), make([]int, vnf), make([]int, vnf), make([]bool, vnf)
	dsizes, dcsizes, desizes, dismix :=
		make([]int, dnf), make([]int, dnf), make([]int, dnf), make([]bool, dnf)
	for i := 0; i < vnf; i++ {
		vcsizes[i], vesizes[i], vsizes[i] = CountChars(vheads[i])
		vismix[i] = checkMixCh(false, vcsizes[i], vesizes[i])
	}
	for i := 0; i < dnf; i++ {
		dcsizes[i], desizes[i], dsizes[i] = CountChars(dheads[i])
		dismix[i] = checkMixCh(false, dcsizes[i], desizes[i])
	}
	//
	invs := make([]invoiceSlcie, len(pinvs))
	for i := 0; i < len(pinvs); i++ {
		p := pinvs[i]
		invs[i].data = []string{
			Sf("%d", i+1), p.Head, p.State, p.UINumber[0:2] + "-" + p.UINumber[2:],
			p.Date.Format(strDateFormat),
			p.SUN, p.SName, p.CName, p.CNumber, Sf("%.1f", p.Total),
		}
		for j := 0; j < vnf; j++ {
			str := Sf("%v", invs[i].data[j])
			nc, ne, nmix := CountChars(str)
			vcsizes[j] = imax(vcsizes[j], nc)
			vesizes[j] = imax(vesizes[j], ne)
			vsizes[j] = imax(vsizes[j], nmix)
			vismix[j] = getIsMixCh(vismix[j], nc, ne)
		}
		for j := 0; j < len(p.Details); j++ {
			d := p.Details[j]
			detail := detailSlcie{
				data: []string{
					Sf("%d", j+1), d.Head, d.UINumber[0:2] + "-" + d.UINumber[2:],
					Sf("%.1f", d.Subtotal), d.Name,
				},
			}
			invs[i].details = append(invs[i].details, detail)
			for k := 0; k < dnf; k++ {
				str := Sf("%v", detail.data[k])
				nc, ne, nmix := CountChars(str)
				dcsizes[k] = imax(dcsizes[k], nc)
				desizes[k] = imax(desizes[k], ne)
				dsizes[k] = imax(dsizes[k], nmix)
				dismix[k] = getIsMixCh(dismix[k], nc, ne)
			}
		}
	}
	vn := isum(vsizes...) + vnf + (vnf-1)*2 + 1
	title := "發票清單"
	_, _, vl := CountChars(title)
	vm := (vn - vl) / 2
	isleft := true
	//
	var b bytes.Buffer
	bws := b.WriteString
	//
	bws(StrSpaces(vm) + title + "\n")
	//
	vhtab := StrThickLine(vn)
	svfields := make([]string, vnf)
	for i := 0; i < vnf; i++ {
		svfields[i] = getColStr(vheads[i], vcsizes[i], vesizes[i], vsizes[i], vismix[i], isleft)
		switch i {
		case 0:
			vhtab += Sf("%v", svfields[i])
		default:
			vhtab += Sf("  %v", svfields[i])
		}
	}
	vhtab += "\n" + StrThinLine(vn)
	lspaces := io.StrSpaces(7)
	dn := isum(dsizes...) + dnf + (dnf-1)*2 + 1
	dhtab := lspaces + StrThickLine(dn)
	sdfields := make([]string, dnf)
	for i := 0; i < dnf; i++ {
		sdfields[i] = getColStr(dheads[i], dcsizes[i], desizes[i], dsizes[i], dismix[i], isleft)
		switch i {
		case 0:
			dhtab += lspaces + Sf("%v", sdfields[i])
		default:
			dhtab += Sf("  %v", sdfields[i])
		}
	}
	dhtab += "\n" + lspaces + StrThinLine(dn)
	//
	for i := 0; i < len(invs); i++ {
		v := &invs[i].data
		bws(vhtab)
		// pchk("%v : %v \n", vnf, v)
		for j := 0; j < vnf; j++ {
			svfields[j] = getColStr((*v)[j], vcsizes[j], vesizes[j], vsizes[j], vismix[j], isleft)
			switch j {
			case 0:
				bws(Sf("%v", svfields[j]))
			default:
				bws(Sf("  %v", svfields[j]))
			}
		}
		bws("\n")
		//
		details := invs[i].details
		ndetails := len(details)
		if ndetails > 0 {
			bws(dhtab)
			for k := 0; k < len(details); k++ {
				d := &details[k].data
				for j := 0; j < dnf; j++ {
					sdfields[j] = getColStr((*d)[j], dcsizes[j], desizes[j], dsizes[j], dismix[j], isleft)
					switch j {
					case 0:
						bws(lspaces + Sf("%v", sdfields[j]))
					default:
						bws(Sf("  %v", sdfields[j]))
					}
				}
				bws("\n")
			}
			bws(lspaces + StrThickLine(dn))
		}
	}
	return b.String()
}

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
	var b bytes.Buffer
	bws := b.WriteString
	ds := reflect.ValueOf(d)
	t := ds.Type()
	nh := t.NumField()
	for i := 0; i < nh; i++ {
		fld := t.Field(i)
		val := ds.Field(i).Interface()
		var v string
		switch val.(type) {
		case gorm.Model:
			continue
		case float64:
			v = io.Sf("%.1f", val.(float64))
		default:
			if fld.Name == "UINumber" {
				v = val.(string)[0:2] + "-" + val.(string)[2:]
				break
			}
			v = val.(string)
		}
		h := fld.Tag.Get("cht")
		bws(io.Sf(" %s : %s |", h, v))
	}
	bws("\n")
	return b.String()
}

// GetArgsTable :
func (d *Detail) GetArgsTable(args ...string) string {
	var title string
	if len(args) > 0 {
		title = args[0]
	} else {
		title = "Detail"
	}
	tab := ArgsTable(
		title,
		"表頭", "Head", d.Head,
		"發票號碼", "UINumber", d.UINumber,
		"小計", "Subtotal", io.Sf("%.1f", d.Subtotal),
		"品項名稱", "Name", d.Name,
	)
	return tab
}

// TableName : set Detail's table name to be `details`
func (Detail) TableName() string {
	// custom table name, this is default
	return "details"
}

// getDetailsTable
func getDetailsTable(pds []*Detail) string {
	Sf, StrSpaces, StrThickLine, StrThinLine := io.Sf, io.StrSpaces, io.StrThickLine, io.StrThinLine
	dheads := []string{"項次", "表頭", "發票號碼", "小計", "品項名稱"}
	dnf := len(dheads)
	dsizes, dcsizes, desizes, dismix :=
		make([]int, dnf), make([]int, dnf), make([]int, dnf), make([]bool, dnf)
	for i := 0; i < dnf; i++ {
		dcsizes[i], desizes[i], dsizes[i] = CountChars(dheads[i])
		dismix[i] = checkMixCh(false, dcsizes[i], desizes[i])
	}
	//
	details := make([]detailSlcie, dnf)
	for i := 0; i < len(pds); i++ {
		d := pds[i]
		details[i].data = []string{
			Sf("%d", i+1), d.Head, d.UINumber[0:2] + "-" + d.UINumber[2:],
			Sf("%.1f", d.Subtotal), d.Name,
		}
		for k := 0; k < dnf; k++ {
			str := Sf("%v", details[i].data[k])
			nc, ne, nmix := CountChars(str)
			dcsizes[k] = imax(dcsizes[k], nc)
			desizes[k] = imax(desizes[k], ne)
			dsizes[k] = imax(dsizes[k], nmix)
			dismix[k] = getIsMixCh(dismix[k], nc, ne)
		}
	}
	dn := isum(dsizes...) + dnf + (dnf-1)*2 + 1
	title := "明細清單"
	_, _, dl := CountChars(title)
	dm := (dn - dl) / 2
	isleft := true
	//
	var b bytes.Buffer
	bws := b.WriteString
	//
	bws(StrSpaces(dm) + title + "\n")
	bws(StrThickLine(dn))
	sdfields := make([]string, dnf)
	for i := 0; i < dnf; i++ {
		sdfields[i] = getColStr(dheads[i], dcsizes[i], desizes[i], dsizes[i], dismix[i], isleft)
		switch i {
		case 0:
			bws(Sf("%v", sdfields[i]))
		default:
			bws(Sf("  %v", sdfields[i]))
		}
	}
	bws("\n")
	bws(StrThinLine(dn))
	for i := 0; i < len(details); i++ {
		d := &details[i].data
		for j := 0; j < dnf; j++ {
			sdfields[j] = getColStr((*d)[j], dcsizes[j], desizes[j], dsizes[j], dismix[j], isleft)
			switch j {
			case 0:
				bws(Sf("%v", sdfields[j]))
			default:
				bws(Sf("  %v", sdfields[j]))
			}
		}
		bws("\n")
	}
	bws(StrThickLine(dn))
	return b.String()
}

func printInvList(pvs []*Invoice) {
	var b bytes.Buffer
	fp := fmt.Fprintf
	for ip, pv := range pvs {
		// fp(&b, "Invoice[%d] :\n", ip)
		fp(&b, "%s", pv.GetArgsTable(io.Sf("Invoices[%d]", ip)))
		for id, pd := range pv.Details {
			// fp(&b, "Invoice[%d] -- Details[%d] :\n", ip, id)
			fp(&b, "%s", pd.GetArgsTable(io.Sf("Invoices[%d] -- Details[%d]", ip, id)))
		}
		fp(&b, "\n")
	}
	pchk("%s", b.String())
}
