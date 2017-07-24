package inv

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"time"

	"github.com/cpmech/gosl/io"
)

type (
	xmlInvoices struct {
		XMLName  xml.Name      `xml:"INVOICES"`
		Version  int           `xml:"version,attr"`
		Invoices []*xmlInvoice `xml:"Invoice"`
	}
	xmlInvoice struct {
		XMLName  xml.Name    `xml:"Invoice"`
		Head     string      `xml:"head,attr"`
		State    string      `xml:"state,attr"`
		UINumber string      `xml:"uniform_invoice_number,attr"`
		Date     string      `xml:"purchase_date,attr"`
		SUN      string      `xml:"store_uniform_number,attr"`
		SName    string      `xml:"store_name,attr"`
		CName    string      `xml:"carrier_name,attr"`
		CNumber  string      `xml:"carrier_number,attr"`
		Total    float64     `xml:"total_amount,attr"`
		Details  []xmlDetail `xml:"Details"`
	}
	xmlDetail struct {
		XMLName  xml.Name `xml:"Detail"`
		Head     string   `xml:"head,attr"`
		UINumber string   `xml:"uniform_invoice_number,attr"`
		Subtotal float64  `xml:"subtotal_amount,attr"`
		Name     string   `xml:"name,attr"`
	}
)

func (x *xmlInvoice) toInvoice() *Invoice {
	date, err := time.Parse(dateFormat, x.Date)
	location, _ := time.LoadLocation("Local")
	if err != nil {
		return nil
	}
	inv := Invoice{
		Head:     x.Head,
		State:    x.State,
		UINumber: x.UINumber,
		Date:     date.In(location),
		SUN:      x.SUN,
		SName:    x.SName,
		CName:    x.CName,
		CNumber:  x.CNumber,
		Total:    x.Total,
	}
	inv.Details = make([]*Detail, 0, len(x.Details))
	for _, d := range x.Details {
		inv.Details = append(inv.Details, d.toDetail())
	}
	return &inv
}

func (v *Invoice) toXMLInvoice() *xmlInvoice {
	xv := xmlInvoice{
		Head:     v.Head,
		State:    v.State,
		UINumber: v.UINumber,
		Date:     v.Date.Format(dateFormat),
		SUN:      v.SUN,
		SName:    v.SName,
		CName:    v.CName,
		CNumber:  v.CNumber,
		Total:    v.Total,
	}
	xv.Details = make([]xmlDetail, 0, len(v.Details))
	for _, d := range v.Details {
		xv.Details = append(xv.Details, d.toXMLDetail())
	}
	return &xv
}

func (xd *xmlDetail) toDetail() *Detail {
	return &Detail{
		Head:     xd.Head,
		UINumber: xd.UINumber,
		Subtotal: xd.Subtotal,
		Name:     xd.Name,
	}
}

func (d *Detail) toXMLDetail() xmlDetail {
	return xmlDetail{
		Head:     d.Head,
		UINumber: d.UINumber,
		Subtotal: d.Subtotal,
		Name:     d.Name,
	}
}

// XMLMarshaller :
type XMLMarshaller struct{}

// MarshalInvoices :
func (XMLMarshaller) MarshalInvoices(fn string, vs []*Invoice) error {
	prun("  > Writing data to .xml file %q ...\n", fn)
	xvs := xmlInvoices{Version: fileVesion}
	xvs.Invoices = make([]*xmlInvoice, 0, len(vs))
	for _, v := range vs {
		xvs.Invoices = append(xvs.Invoices, v.toXMLInvoice())
	}
	io.WriteBytesToFile(fn, []byte(xml.Header))
	// b, err := xml.MarshalIndent(&xvs, "", "    ")
	b, err := xml.Marshal(&xvs)
	if err != nil {
		return err
	}
	io.AppendToFile(fn, bytes.NewBuffer(b))
	return nil
}

// UnmarshalInvoices is used to unmarshal invoices
func (XMLMarshaller) UnmarshalInvoices(fn string) ([]*Invoice, error) {
	prun("  > Reading data from .xml file %q ...\n", fn)
	b, err := io.ReadFile(fn)
	if err != nil {
		return nil, err
	}
	pxvs := &xmlInvoices{}
	err = xml.Unmarshal(b, pxvs)
	if err != nil {
		return nil, err
	}
	if pxvs.Version > fileVesion {
		return nil, fmt.Errorf("version %d is too new to read", pxvs.Version)
	}
	pvs := make([]*Invoice, 0, len(pxvs.Invoices))
	for _, xinv := range pxvs.Invoices {
		pvs = append(pvs, xinv.toInvoice())
	}
	// for _, p := range pvs {
	// 	pchk("%#q\n", *p)
	// }
	plog(GetInvoicesTable(pvs))
	// pchk("%v\n", vsToTable(pvs))
	prun("    updating database ...\n")
	DBInsertFrom(pvs)
	return pvs, nil
}
