package inv

import (
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/json-iterator/go"
)

// JSONInvoices is used for JSON file type
type JSONInvoices struct {
	FileType   string
	FileVesion int
	Invoices   []*Invoice
}

// JSONMarshaller collects the mathods marshalling or unmarshalling the csv-type data
type JSONMarshaller struct{}

// MarshalInvoices marshalls the csv-type data of invoices
func (JSONMarshaller) MarshalInvoices(fn string, invoices []*Invoice) error {
	prun("  > Writing data to .jsn or .json file %q ...\n", fn)
	j := JSONInvoices{
		FileType:   fileType,
		FileVesion: fileVesion,
		Invoices:   invoices,
	}
	// b, err := jsoniter.MarshalIndent(&j, "", "    ")
	b, err := jsoniter.Marshal(&j)
	if err != nil {
		return err
	}
	io.WriteBytesToFile(fn, b)
	return nil
}

// UnmarshalInvoices unmarshalls the csv-type data of invoices
func (JSONMarshaller) UnmarshalInvoices(fn string) ([]*Invoice, error) {
	prun("  > Reading data from .jsn or .json file %q ...\n", fn)
	b, err := io.ReadFile(fn)
	if err != nil {
		return nil, err
	}
	var j JSONInvoices
	err = jsoniter.Unmarshal(b, &j)
	if err != nil {
		return nil, err
	}
	if j.FileType != fileType {
		return nil, chk.Err("cannot read non-invoices json file")
	}
	if j.FileVesion > fileVesion {
		return nil, chk.Err("version %d is too new to read", j.FileVesion)
	}
	plog(GetInvoicesTable(j.Invoices))
	prun("    updating database ...\n")
	DBInsertFrom(j.Invoices)
	return j.Invoices, nil
}
