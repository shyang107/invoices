package inv

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	iconv "github.com/djimenez/iconv-go"
)

// CsvMarshaller collects the mathods marshalling or unmarshalling the .csv data
type CsvMarshaller struct{}

// MarshalInvoices marshalls the .csv data of invoices
func (CsvMarshaller) MarshalInvoices(fn string, pvs []*Invoice) error {
	prun("  > Writing data to .csv file %q ...\n", fn)
	var b bytes.Buffer
	fmt.Fprintln(&b, fileType)
	fmt.Fprintln(&b, io.Sf("%v", fileVesion))
	for _, pv := range pvs {
		fmt.Fprintln(&b, pv.toCSVString())
		for _, d := range pv.Details {
			fmt.Fprintln(&b, d.toCSVString())
		}
	}
	io.WriteFile(fn, &b)
	return nil
}

func (v *Invoice) toCSVString() string {
	csv := []string{
		v.Head,
		v.State,
		v.UINumber,
		v.Date.Format(dateFormat),
		v.SUN,
		v.SName,
		v.CName,
		v.CNumber,
		io.Sf("%v", v.Total),
	}
	return strings.Join(csv, csvSep)
}

func (d *Detail) toCSVString() string {
	csv := []string{
		d.Head,
		d.UINumber,
		io.Sf("%v", d.Subtotal),
		d.Name,
	}
	return strings.Join(csv, csvSep)
}

// UnmarshalInvoices unmarshalls the .csv data of invoices
func (CsvMarshaller) UnmarshalInvoices(fn string) ([]*Invoice, error) {
	pstat("  > Reading data from .csv file %q ...\n", fn)
	f, err := io.OpenFileR(fn)
	if err != nil {
		return nil, err
	}
	var pinvs []*Invoice
	var pdets []*Detail
	err = io.ReadLinesFile(f, func(idx int, line string) (stop bool) {
		// plog("line = %v\n", line)
		if Opt.IsNative {
			line = big5ToUtf8(line)
		} else {
			switch idx {
			case 0:
				ft := strings.Trim(line, " ")
				if ft != fileType {
					panic(chk.Err("type of .csv file is not matched (%q)", fileType))
				}
			case 1:
				fv := io.Atoi(strings.Trim(line, " "))
				if fv != fileVesion {
					panic(chk.Err("version (%v) of .csv file is not matched (%v)", fv, fileVesion))
				}
			}
		}
		recs := strings.Split(line, csvSep)
		head := recs[0]
		switch head {
		case "M": // invoice
			pinv := unmarshalCSVInvoice(recs)
			// io.Pf("%s%v\n", io.StrSpaces(4), *pinv)
			pinvs = append(pinvs, pinv)
		case "D": // deltail of invoice
			pdet := unmarshalCSVDetail(recs)
			// io.Pf("%s%v\n", io.StrSpaces(4), det)
			pdets = append(pdets, pdet)
		}
		return
	})
	if err != nil {
		return nil, err
	}
	prun("    combining invoices ...\n")
	combineInvoice(pinvs, pdets)
	plog(GetInvoicesTable(pinvs))
	// printInvList(pinvs)
	prun("    updating database ...\n")
	DBInsertFrom(pinvs)
	return pinvs, nil
}

func combineInvoice(pvs []*Invoice, pds []*Detail) {
	for _, d := range pds {
		no := d.UINumber
		for _, p := range pvs {
			if p.UINumber == no {
				// d.Invoice = p
				p.Details = append(p.Details, d)
				break
			}
		}
	}
}

func unmarshalCSVDetail(recs []string) *Detail {
	// pchk("%sDetail : %#v\n", io.StrSpaces(4), recs)
	det := Detail{
		Head:     recs[0],
		UINumber: recs[1],
		Subtotal: io.Atof(recs[2]),
		Name:     recs[3],
	}
	return &det
}

func unmarshalCSVInvoice(recs []string) *Invoice {
	// pchk("%sInvoice : %#v\n", io.StrSpaces(4), recs)
	date, err := time.Parse(dateFormat, recs[3])
	location, _ := time.LoadLocation("Local")
	if err != nil {
		panic(chk.Err("%v : %s", err, recs[3]))
	}
	inv := Invoice{
		Head:     recs[0],
		State:    recs[1],
		UINumber: recs[2],
		Date:     date.In(location),
		SUN:      recs[4],
		SName:    recs[5],
		CName:    recs[6],
		CNumber:  recs[7],
		Total:    io.Atof(recs[8]),
	}
	return &inv
}

func big5ToUtf8(str string) string {
	res, err := iconv.ConvertString(str, Opt.IfnEncoding, "utf-8")
	checkErr(err)
	return res
}
