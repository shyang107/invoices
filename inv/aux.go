package inv

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strings"
	"unicode"
	"unsafe"

	"github.com/cpmech/gosl/io"
)

var (
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
func GetFieldsInfo(obj interface{}, tag string) (fields, types, kinds, tags []string) {
	d := reflect.ValueOf(obj) // vals[0] = d.Field[0]
	t := d.Type()
	n := t.NumField()
	fields = make([]string, n)
	types = make([]string, n)
	kinds = make([]string, n)
	tags = make([]string, n)
	for i := 0; i < n; i++ {
		fields[i], types[i], kinds[i], tags[i] =
			t.Field(i).Name,
			t.Field(i).Type.String(),
			d.Field(i).Kind().String(),
			t.Field(i).Tag.Get(tag)
	}
	return fields, types, kinds, tags
}

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
	pfsep("%s", io.StrThinLine(60))
}

// BtoStr convert []byte to string
func BtoStr(bs []byte) string {
	return *(*string)(unsafe.Pointer(&bs))
}

// Imin reports the minimum value of a and b
func Imin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Imax reports the maximum value of a and b
func Imax(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Isum reports the summation of args...
func Isum(args ...int) int {
	n := 0
	for _, a := range args {
		n += a
	}
	return n
}

func isOpened(path string) bool {
	fpath := os.ExpandEnv(path)
	f, err := os.Open(fpath)
	defer f.Close()
	if err != nil {
		perr("!!! File %q does not open !!!\n", fpath)
		return false
	}
	return true
}

func isNotExist(path string) bool {
	fpath := os.ExpandEnv(path)
	_, err := os.Stat(fpath)
	if os.IsNotExist(err) {
		perr("!!! File %q does not exist !!!\n", fpath)
		return true
	}
	return false
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

// GetColStr return string use in field
func GetColStr(s string, size int, left bool) string {
	_, _, n := CountChars(s)
	spaces := strings.Repeat(" ", size-n)
	// size := nc*2 + ne // s 實際佔位數
	var tab string
	if left {
		tab = fmt.Sprintf("%v%s", s, spaces)
	} else {
		tab = fmt.Sprintf("%s%v", spaces, s)
	}
	return " " + tab
}

// CountChars returns the number of each other of chinses and english characters
func CountChars(str string) (nc, ne, n int) {
	for _, r := range str {
		lchar := len(string(r))
		// n += lchar
		if lchar > 1 {
			nc++
		} else {
			ne++
		}
	}
	n = 2*nc + ne
	return nc, ne, n
}

// IsChineseChar judges whether the chinese character exists ?
func IsChineseChar(str string) bool {
	// n := 0
	for _, r := range str {
		// io.Pf("%q ", r)
		if unicode.Is(unicode.Scripts["Han"], r) {
			// n++
			return true
		}
	}
	return false
}

// ArgsTable prints a nice table with input arguments
//  Input:
//   title -- title of table; e.g. INPUT ARGUMENTS
//   data  -- sets of THREE items in the following order:
//                 description, key, value, ...
//                 description, key, value, ...
//                      ...
//                 description, key, value, ...
func ArgsTable(title string, data ...interface{}) string {
	heads := []string{"description", "key", "value"}
	return ArgsTableN(title, 0, heads, data...)
}

func setupMaxSize(psizes *[]int, compared []int) {
	for i := 0; i < len(*psizes); i++ {
		(*psizes)[i] = Imax((*psizes)[i], compared[i])
	}
}

// ArgsTableN prints a nice table with input arguments
//  Input:
//   title -- title of table; e.g. INPUT ARGUMENTS
//	 heads -- heads of table; e.g. []string{ col1,  col2, ... }
//	 lensp -- length of leading spaces in every row
//   data  -- sets of THREE items in the following order:
//                 column1, column2, column3, ...
//                 column1, column2, column3, ...
//                      ...
//                 column1, column2, column3, ...
func ArgsTableN(title string, lensp int, heads []string, data ...interface{}) string {
	Sf, StrSpaces, StrThickLine, StrThinLine := io.Sf, io.StrSpaces, io.StrThickLine, io.StrThinLine
	nf := len(heads)
	ndat := len(data)
	if ndat < nf {
		return ""
	}
	if lensp < 0 {
		lensp = 0
	}
	lspaces := io.StrSpaces(lensp)
	nlines := ndat / nf
	sizes := make([]int, nf)
	for i := 0; i < nf; i++ {
		_, _, sizes[i] = CountChars(heads[i])
	}
	for i := 0; i < nlines; i++ {
		if i*nf+(nf-1) >= ndat {
			return Sf("ArgsTable: input arguments are not a multiple of %d\n", nf)
		}
		for j := 0; j < nf; j++ {
			str := Sf("%v", data[i*nf+j])
			_, _, nmix := CountChars(str)
			sizes[j] = Imax(sizes[j], nmix)
		}
	}
	// strfmt := Sf("%%v  %%v  %%v\n")
	n := Isum(sizes...) + nf + (nf-1)*2 + 1 // sizes[0] + sizes[1] + sizes[2] + 3 + 4
	_, _, l := CountChars(title)
	m := (n - l) / 2
	//
	var b bytes.Buffer
	bw := b.WriteString
	//
	bw(StrSpaces(m+lensp) + title + "\n")
	bw(lspaces + StrThickLine(n))
	isleft := true
	sfields := make([]string, nf)
	for i := 0; i < nf; i++ {
		sfields[i] = GetColStr(heads[i], sizes[i], isleft)
		switch i {
		case 0:
			bw(Sf("%v", lspaces+sfields[i]))
		default:
			bw(Sf("  %v", sfields[i]))
		}
	}
	bw("\n")
	bw(lspaces + StrThinLine(n))
	for i := 0; i < nlines; i++ {
		for j := 0; j < nf; j++ {
			sfields[j] = GetColStr(Sf("%v", data[i*nf+j]), sizes[j], isleft)
			switch j {
			case 0:
				bw(Sf("%v", lspaces+sfields[j]))
			default:
				bw(Sf("  %v", sfields[j]))
			}
		}
		bw("\n")
	}
	bw(lspaces + StrThickLine(n))
	return b.String()
}
