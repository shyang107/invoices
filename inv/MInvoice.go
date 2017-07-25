package inv

import (
	"bytes"
	"fmt"
	"reflect"
	"time"

	"github.com/cpmech/gosl/io"
	"github.com/jinzhu/gorm"
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
	Sf, Ff := fmt.Sprintf, fmt.Fprintf
	var b bytes.Buffer
	val := reflect.ValueOf(pv) //.Elem()
	fld := val.Type()
	var str string
	for i := 0; i < val.NumField(); i++ {
		switch fld.Field(i).Name {
		case "Model", "Details":
			continue
		case "Date":
			str = val.Field(i).Interface().(time.Time).Format(strDateFormat)
		case "Total":
			str = Sf("%.1f", val.Field(i).Interface().(float64))
		case "UINumber":
			str = val.Field(i).Interface().(string)[0:2] + "-" + val.Field(i).Interface().(string)[2:]
		default:
			str = val.Field(i).Interface().(string)
		}
		Ff(&b, " %s : %s |", fld.Field(i).Tag.Get("cht"), str)
	}
	Ff(&b, "\n")
	lspaces := StrSpaces(4)
	for i, d := range pv.Details {
		Ff(&b, "%s> %2d. %s", lspaces, i+1, d)
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
func (pv *Invoice) GetArgsTable(title string) string {
	Sf := fmt.Sprintf
	if len(title) == 0 {
		title = "發票清單"
	}
	// heads := []string{"表頭", "發票狀態", "發票號碼", "發票日期",
	// "商店統編", "商店店名", "載具名稱", "載具號碼", "總金額", "明細清單"}
	_, _, _, heads := GetFieldsInfo(Invoice{}, "cht", "Model")
	lensp := 0
	table := ArgsTableN(title, lensp, heads, pv.Head, pv.State,
		pv.UINumber[0:2]+"-"+pv.UINumber[2:], pv.Date.Format(strDateFormat),
		pv.SUN, pv.SName, pv.CName, pv.CNumber,
		Sf("%.1f", pv.Total), "[如下...]")
	lensp = 7
	table += GetDetailsTable(pv.Details, lensp)
	return table
	// var title string
	// if len(args) > 0 {
	// 	title = args[0]
	// } else {
	// 	title = "Invoice"
	// }
	// tab := ArgsTable(
	// 	title,
	// 	"表頭", "Head", pv.Head,
	// 	"發票狀態", "State", pv.State,
	// 	"發票號碼", "UINumber", pv.UINumber,
	// 	"發票日期", "Date", pv.Date.Format(strDateFormat),
	// 	"商店統編", "SUN", pv.SUN,
	// 	"商店店名", "SName", pv.SName,
	// 	"載具名稱", "CName", pv.CName,
	// 	"載具號碼", "CNumber", pv.CNumber,
	// 	"總金額", "Total", io.Sf("%.1f", pv.Total),
	// 	"明細清單", "Details", "[]*Detail",
	// )
	// return tab
}

type invoiceSlcie struct {
	data    []string
	details []detailSlcie
}
type detailSlcie struct {
	data []string
}

// GetInvoicesTable returns the table string of the list of []*Invoice
func GetInvoicesTable(pinvs []*Invoice) string {
	Sf, StrSpaces, StrThickLine, StrThinLine := io.Sf, io.StrSpaces, io.StrThickLine, io.StrThinLine
	vheads := []string{"項次", "表頭", "發票狀態", "發票號碼", "發票日期",
		"商店統編", "商店店名", "載具名稱", "載具號碼", "總金額"}
	dheads := []string{"項次", "表頭", "發票號碼", "小計", "品項名稱"}
	vnf := len(vheads)
	dnf := len(dheads)
	vsizes := make([]int, vnf)
	dsizes := make([]int, vnf)
	for i := 0; i < vnf; i++ {
		_, _, vsizes[i] = CountChars(vheads[i])
	}
	for i := 0; i < dnf; i++ {
		_, _, dsizes[i] = CountChars(vheads[i])
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
			_, _, nmix := CountChars(str)
			vsizes[j] = Imax(vsizes[j], nmix)
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
				_, _, nmix := CountChars(str)
				dsizes[k] = Imax(dsizes[k], nmix)
			}
		}
	}
	vn := Isum(vsizes...) + vnf + (vnf-1)*2 + 1
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
		svfields[i] = GetColStr(vheads[i], vsizes[i], isleft)
		switch i {
		case 0:
			vhtab += Sf("%v", svfields[i])
		default:
			vhtab += Sf("  %v", svfields[i])
		}
	}
	vhtab += "\n" + StrThinLine(vn)
	lspaces := io.StrSpaces(7)
	dn := Isum(dsizes...) + dnf + (dnf-1)*2 + 1
	dhtab := lspaces + StrThickLine(dn)
	sdfields := make([]string, dnf)
	for i := 0; i < dnf; i++ {
		sdfields[i] = GetColStr(dheads[i], dsizes[i], isleft)
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
			svfields[j] = GetColStr((*v)[j], vsizes[j], isleft)
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
					sdfields[j] = GetColStr((*d)[j], dsizes[j], isleft)
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

func printInvList(pvs []*Invoice) {
	var b bytes.Buffer
	fp := fmt.Fprintf
	for ip, pv := range pvs {
		// fp(&b, "%d : %s", ip+1, pv)
		fp(&b, "%s", pv.GetArgsTable(io.Sf("發票 %d", ip+1)))
		// for id, pd := range pv.Details {
		// 	fp(&b, "%s", pd.GetArgsTable(io.Sf("Invoices[%d] -- Details[%d]", ip, id), 7))
		// }
		// fp(&b, "\n")
	}
	pchk("%s", b.String())
}
