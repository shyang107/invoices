package inv

import (
	"path/filepath"

	"github.com/cpmech/gosl/chk"
)

// InvoiceMarshaller is marshal-operator of invoices
type InvoiceMarshaller interface {
	MarshalInvoices(fn string, pvs []*Invoice) error
}

// InvoiceUnmarshaller is unmarshal-operator of invoices
type InvoiceUnmarshaller interface {
	UnmarshalInvoices(fn string) ([]*Invoice, error)
}

// ReadInvoices reads invoice-record from fn
func (o *Option) ReadInvoices() ([]*Invoice, error) {
	var unmarshaller InvoiceUnmarshaller
	startfunc(ffstart) //, "ReadInvoices")
	// pstat("file-type : %q\n", Opts.IfnSuffix)
	//
	var fb = FileBunker{Name: filepath.Base(o.InpFn)}
	DB.Where(&fb).First(&fb)
	plog((&fb).GetArgsTable("", 0))
	//
	switch o.IfnSuffix {
	case ".csv":
		pstat("%q\n", "CsvMarshaller")
		unmarshaller = CsvMarshaller{}
	case ".jsn", ".json":
		pstat("%q\n", "JSONMarshaller")
		unmarshaller = JSONMarshaller{}
	case ".xml":
		pstat("%q\n", "XMLMarshaller")
		unmarshaller = XMLMarshaller{}
	case ".xlsx":
		pstat("%q\n", "XlsMarshaller")
		unmarshaller = XlsMarshaller{}
	}
	if unmarshaller != nil {
		invs, err := unmarshaller.UnmarshalInvoices(o.InpFn)
		stopfunc(ffstop) //, "ReadInvoices")
		return invs, err
	}
	return nil, chk.Err("not supprted file-type : %s (%s)", o.IfnSuffix, o.InpFn)
}

// WriteInvoices reads invoice-record from fn
func (o *Option) WriteInvoices(invs []*Invoice) error {
	var marshaller InvoiceMarshaller
	startfunc(ffstart) //, "ReadInvoices")
	switch o.OfnSuffix {
	case ".csv":
		pstat("%q\n", "CsvMarshaller")
		marshaller = CsvMarshaller{}
	case ".jsn", ".json":
		pstat("%q\n", "JSONMarshaller")
		marshaller = JSONMarshaller{}
	case ".xml":
		pstat("%q\n", "XMLMarshaller")
		marshaller = XMLMarshaller{}
	case ".xlsx":
		pstat("%q\n", "XlsMarshaller")
		marshaller = XlsMarshaller{}
	}
	if marshaller != nil {
		err := marshaller.MarshalInvoices(o.OutFn, invs)
		stopfunc(ffstop) //, "ReadInvoices")
		return err
	}
	return chk.Err("not supprted file-type : %s (%s)", o.IfnSuffix, o.InpFn)
}
