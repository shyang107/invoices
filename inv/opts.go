package inv

import (
	"os"
	"strings"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	// "github.com/shyang/invoices/inv/goini"
	"github.com/widuu/goini"
)

var (
	// Opts is options for this application
	Opts Options
)

// Options setup the application
type Options struct {
	// [input]
	InpFn       string
	IfnSuffix   string
	IsNative    bool
	IfnEncoding string
	// [output]
	OutFn     string
	OfnSuffix string
	IsOutput  bool
	// [punch]
	PunchFn string
}

func (o Options) String() string {
	strdashk := strings.Repeat("-", 12)
	strdashv := strings.Repeat("-", 15)
	tab := ArgsTable(
		"Options",
		"Input:", strdashk, strdashv,
		"input file name", "InpFn", o.InpFn,
		"input file type", "IfnSuffix", o.IfnSuffix,
		"is official invoices file from government?", "IsNative", o.IsNative,
		"encoding of input file", "INFencoding", o.IfnEncoding,
		"Output:", strdashk, strdashv,
		"output file name (if you want)", "OutFn", o.OutFn,
		"output file type", "OfnSuffix", o.OfnSuffix,
		"do output?", "IsOutput", o.IsOutput,
		"Punch:", strdashk, strdashv,
		"punch file name (not use for now)", "PunchFn", o.PunchFn,
	)
	return tab
}

// DefaultOptions sets a list of safe recommended options. Feel free to modify these to suit your needs.
var DefaultOptions = Options{
	InpFn:       "./data/inp/09751085061.csv",
	IfnSuffix:   ".csv",
	IsNative:    true,
	IfnEncoding: "Big5",
	OutFn:       "./data/out/09751085061.json",
	OfnSuffix:   ".json",
	PunchFn:     "./data/out/punch.out",
}

// GetOptions gets the configuration from cfgFN
// // [input]
// // inputFile   = ./09751085061.csv
// // is_native   = false
// // encoding    = big5
// // [output]
// // outputFile  = ./09751085061.json
// // is_output   = true
// // [punch]
// // punchFileName = ./punch.out
func (o *Options) GetOptions() {
	startfunc(fostart)
	cfn := Cfg.CaseFn
	if !isOpened(cfn) {
		panic(chk.Err("config-file %q can not open", cfn))
	}
	c := goini.SetConfig(cfn)

	// [input]
	o.InpFn = os.ExpandEnv(c.GetValue("input", "input_file"))
	o.IfnSuffix = io.FnExt(o.InpFn)
	o.IsNative = io.Atob(c.GetValue("input", "is_native"))
	o.IfnEncoding = c.GetValue("input", "encoding")
	// [output]
	o.OutFn = os.ExpandEnv(c.GetValue("output", "output_file"))
	o.OfnSuffix = io.FnExt(o.OutFn)
	o.IsOutput = io.Atob(c.GetValue("output", "is_output"))
	// [punch]
	o.PunchFn = c.GetValue("punch", "punchFile")
	//
	plog("%s", *o)
	stopfunc(fostop) //, "GetFromConfig")
	// os.Exit(1)
}
