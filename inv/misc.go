package inv

import (
	"fmt"
	"reflect"
	"runtime"

	"github.com/cpmech/gosl/io"
)

const (
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

	pfstart = io.PfCyan
	pfstop  = io.PfBlue
	pfsep   = io.Pfdyel2
	plog    = io.Pf
	pwarn   = io.Pforan
	perr    = io.Pflmag
	prun    = io.PfYel
	pchk    = io.Pfgreen2
	pstat   = io.Pfyel
)

func callerName(idx int) string {
	pc, _, _, _ := runtime.Caller(idx) //idx = 0 self, 1 for caller, 2 for upper caller
	return runtime.FuncForPC(pc).Name()
}

func startfunc(fid int) {
	// io.Pfdyel2("%s", io.StrThinLine(60))
	// switch fid {
	// case ostart:
	// 	io.PfCyan(format[fid])
	// default:
	// 	io.PfCyan(format[fid], callerName(2))
	// }
	// io.PfCyan(format[fid], callerName(2))
	pfstart(format[fid], callerName(2))
}

func stopfunc(fid int) {
	// switch fid {
	// case oend:
	// 	io.PfBlue(format[fid])
	// default:
	// 	io.PfBlue(format[fid], callerName(2))
	// }
	// io.PfBlue(format[fid], callerName(2))
	// io.Pfdyel2("%s", io.StrThinLine(60))
	pfstop(format[fid], callerName(2))
	printSepline(60)
}

func printSepline(n int) {
	if n <= 0 {
		n = 60
	}
	pfsep("%s", StrThinLine(n))
}

// getErrMessage get error message if error
func getErrMessage(err error) string {
	if err != nil {
		return fmt.Sprintf("Error is :'%s'", err.Error())
	}
	return "Notfound this error"
}

func checkErr(err error) {
	if err != nil {
		// perr(getErrMessage(err))
		panic(err)
	}
}

// GetTypes returns map["Field"]="Type"
func GetTypes(obj interface{}) map[string]string {
	d := reflect.ValueOf(obj)
	t := d.Type()
	n := t.NumField()
	types := make(map[string]string)
	for i := 0; i < n; i++ {
		types[t.Field(i).Name] = t.Field(i).Type.String() //io.Sf("%v", t.Field(i).Type)
	}
	return types
}

// GetTags returns map["Field"]="tag_value"
func GetTags(obj interface{}, tag string) map[string]string {
	d := reflect.ValueOf(obj)
	t := d.Type()
	n := t.NumField()
	tags := make(map[string]string)
	for i := 0; i < n; i++ {
		tags[t.Field(i).Name] = t.Field(i).Tag.Get(tag) //io.Sf("%v", t.Field(i).Tag.Get(tag))
	}
	return tags
}

// GetFieldsInfo return information of fields
func GetFieldsInfo(obj interface{}, tagname string, ignoreFields ...string) (fields, types, kinds, tags []string) {
	Sf := fmt.Sprintf
	t := reflect.TypeOf(obj)
	for i := 0; i < t.NumField(); i++ {
		if isIgnore(t.Field(i).Name, ignoreFields...) {
			continue
		}
		fields = append(fields, t.Field(i).Name)
		types = append(types, Sf("%v", t.Field(i).Type))
		kinds = append(kinds, Sf("%v", t.Field(i).Type.Kind()))
		tags = append(tags, t.Field(i).Tag.Get(tagname))
	}
	return fields, types, kinds, tags
}

func isIgnore(fieldName string, ignoreFields ...string) bool {
	for i := 0; i < len(ignoreFields); i++ {
		if fieldName == ignoreFields[i] {
			return true
		}
	}
	return false
}
