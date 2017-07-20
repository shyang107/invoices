package inv

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"runtime"
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

// []byte to string
func btostr(bs []byte) string {
	return *(*string)(unsafe.Pointer(&bs))
}

func imin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func imax(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func isum(args ...int) int {
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
func getIsMixCh(ismix bool, nc, ne int) bool {
	if ismix {
		return true
	}
	return ((nc > 0) && (ne > 0))
}

func checkMixCh(ismix bool, nc, ne int) bool {
	if ismix {
		return true
	}
	return ((nc > 0) && (ne > 0))
}

func checkMixChs(ismix *[]bool, nc, ne []int) {
	for i := 0; i < len(*ismix); i++ {
		(*ismix)[i] = checkMixCh((*ismix)[i], nc[i], ne[i])
	}
}

// func getColStr(s string, ncmax, nemax int, ismix bool, left bool) string {
// 	nc, ne := CountChEngChar(s)
// 	size := nc*2 + ne // s 實際佔位數
// 	var sizemax int
// 	if ismix {
// 		sizemax = ncmax*2 + nemax // 欄位最大實際佔位數
// 	} else {
// 		sizemax = imax(ncmax*2, nemax)
// 	}
// 	// io.Pfred("!!!    nc = %d,    ne = %d,     size = %d, mix = %t <- %q\n", nc, ne, size, ismix, s)
// 	var tab string
// 	if left {
// 		tab = s + io.StrSpaces(sizemax-size)
// 	} else {
// 		tab = io.StrSpaces(sizemax-size) + s
// 	}
// 	// io.Pforan("!!! ncmax = %d, nemax = %d, sizemax = %d, mix = %t <- %q\n", ncmax, nemax, sizemax, ismix, tab)
// 	return " " + tab
// }

func getColStr(s string, ncmax, nemax, sizemax int, ismix bool, left bool) string {
	nc, ne := CountChEngChar(s)
	size := nc*2 + ne // s 實際佔位數
	var tab string
	if left {
		tab = s + io.StrSpaces(sizemax-size)
	} else {
		tab = io.StrSpaces(sizemax-size) + s
	}
	return " " + tab
}

// CountChEngChar returns the number of each other of chinses and english characters
func CountChEngChar(str string) (nc int, ne int) {
	for _, r := range str {
		// io.Pf("%q ", r)
		if unicode.Is(unicode.Scripts["Han"], r) {
			nc++
		} else {
			ne++
		}
	}
	return nc, ne
}

// CountChEngChar0 returns the number of each other of chinses and english characters
func CountChEngChar0(str string) (nc int, ne int) {
	for _, r := range str {
		plog("%#U ", r)
		if len(string(r)) > 1 {
			plog("<- Han")
			nc++
		} else {
			plog("<- eng.")
			ne++
		}
		// if unicode.Is(unicode.Scripts["Han"], r) {
		// 	plog("<- Han")
		// 	nc++
		// } else {
		// 	plog("<- eng.")
		// 	ne++
		// }
		io.Pl()
	}
	return nc, ne
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

func countNPos(nc, ne int) (np int) {
	ismix := ((nc > 0) && (ne > 0))
	if ismix {
		np = 2*nc + ne
	} else {
		if nc > 0 {
			np = 2 * nc
		} else {
			np = ne
		}
	}
	return np
}

// CountNPos :
func CountNPos(nc, ne int) (np int) {
	ismix := ((nc > 0) && (ne > 0))
	if ismix {
		np = 2*nc + ne
	} else {
		if nc > 0 {
			np = 2 * nc
		} else {
			np = ne
		}
	}
	return np
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
	return ArgsTableN(title, heads, data...)
}

func setupMaxSize(psizes *[]int, compared []int) {
	for i := 0; i < len(*psizes); i++ {
		(*psizes)[i] = imax((*psizes)[i], compared[i])
	}
}

// ArgsTableN prints a nice table with input arguments
//  Input:
//   title -- title of table; e.g. INPUT ARGUMENTS
//	 heads -- heads of table; e.g. []string{ col1,  col2, ... }
//   data  -- sets of THREE items in the following order:
//                 column1, column2, column3, ...
//                 column1, column2, column3, ...
//                      ...
//                 column1, column2, column3, ...
func ArgsTableN(title string, heads []string, data ...interface{}) string {
	Sf, StrSpaces, StrThickLine, StrThinLine := io.Sf, io.StrSpaces, io.StrThickLine, io.StrThinLine
	nf := len(heads)
	ndat := len(data)
	if ndat < nf {
		return ""
	}
	nlines := ndat / nf
	sizes := make([]int, nf)
	csizes := make([]int, nf)
	esizes := make([]int, nf)
	ismixch := make([]bool, nf)
	for i := 0; i < nf; i++ {
		csizes[i], esizes[i], sizes[i] = CountChars(heads[i])
		ismixch[i] = getIsMixCh(false, csizes[i], esizes[i])
	}
	for i := 0; i < nlines; i++ {
		if i*nf+(nf-1) >= ndat {
			return Sf("ArgsTable: input arguments are not a multiple of %d\n", nf)
		}
		for j := 0; j < nf; j++ {
			str := Sf("%v", data[i*nf+j])
			nc, ne, nmix := CountChars(str)
			csizes[j] = imax(csizes[j], nc)
			esizes[j] = imax(esizes[j], ne)
			sizes[j] = imax(sizes[j], nmix)
			ismixch[j] = getIsMixCh(ismixch[j], nc, ne)
		}
	}
	// strfmt := Sf("%%v  %%v  %%v\n")
	n := isum(sizes...) + nf + (nf-1)*2 + 1 // sizes[0] + sizes[1] + sizes[2] + 3 + 4
	tnc, tne := CountChEngChar(title)
	l := countNPos(tnc, tne)
	m := (n - l) / 2
	//
	var b bytes.Buffer
	bw := b.WriteString
	//
	bw(StrSpaces(m) + title + "\n")
	bw(StrThickLine(n))
	isleft := true
	sfields := make([]string, nf)
	for i := 0; i < nf; i++ {
		sfields[i] = getColStr(heads[i], csizes[i], esizes[i], sizes[i], ismixch[i], isleft)
		switch i {
		case 0:
			bw(Sf("%v", sfields[i]))
		default:
			bw(Sf("  %v", sfields[i]))
		}
	}
	bw("\n")
	bw(StrThinLine(n))
	for i := 0; i < nlines; i++ {
		for j := 0; j < nf; j++ {
			sfields[j] = getColStr(Sf("%v", data[i*nf+j]),
				csizes[j], esizes[j], sizes[j], ismixch[j], isleft)
			switch j {
			case 0:
				bw(Sf("%v", sfields[j]))
			default:
				bw(Sf("  %v", sfields[j]))
			}
		}
		bw("\n")
	}
	bw(StrThickLine(n))
	return b.String()
}
