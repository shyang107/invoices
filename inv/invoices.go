package inv

import (
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
func ReadInvoices() ([]*Invoice, error) {
	var unmarshaller InvoiceUnmarshaller
	startfunc(ffstart) //, "ReadInvoices")
	// pstat("file-type : %q\n", Opts.IfnSuffix)
	switch Opts.IfnSuffix {
	case ".csv":
		pstat("%q\n", "CsvMarshaller")
		unmarshaller = CsvMarshaller{}
	case ".jsn", ".json":
		pstat("%q\n", "JSONMarshaller")
		unmarshaller = JSONMarshaller{}
	case ".xml":
		pstat("%q\n", "XMLMarshaller")
		unmarshaller = XMLMarshaller{}
	case ".xls":
		pstat("%q\n", "XlsMarshaller")
		unmarshaller = XlsMarshaller{}
	}
	if unmarshaller != nil {
		invs, err := unmarshaller.UnmarshalInvoices(Opts.InpFn)
		stopfunc(ffstop) //, "ReadInvoices")
		return invs, err
	}
	return nil, chk.Err("not supprted file-type : %s (%s)", Opts.IfnSuffix, Opts.InpFn)
}

// WriteInvoices reads invoice-record from fn
func WriteInvoices(invs []*Invoice) error {
	var marshaller InvoiceMarshaller
	startfunc(ffstart) //, "ReadInvoices")
	switch Opts.OfnSuffix {
	case ".csv":
		pstat("%q\n", "CsvMarshaller")
		marshaller = CsvMarshaller{}
	case ".jsn", ".json":
		pstat("%q\n", "JSONMarshaller")
		marshaller = JSONMarshaller{}
	case ".xml":
		pstat("%q\n", "XMLMarshaller")
		marshaller = XMLMarshaller{}
	case ".xls", ".xlsx":
		pstat("%q\n", "XlsMarshaller")
		marshaller = XlsMarshaller{}
	}
	if marshaller != nil {
		err := marshaller.MarshalInvoices(Opts.OutFn, invs)
		stopfunc(ffstop) //, "ReadInvoices")
		return err
	}
	return chk.Err("not supprted file-type : %s (%s)", Opts.IfnSuffix, Opts.InpFn)
}
