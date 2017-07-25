package inv

import (
	"bytes"
	"fmt"
	"math"
	"strings"
	"unicode"
	"unsafe"
)

// BytesSizeToString convert bytes to a human-readable size
func BytesSizeToString(byteCount int) string {
	suf := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"} //Longs run out around EB
	if byteCount == 0 {
		return "0" + suf[0]
	}
	bytes := math.Abs(float64(byteCount))
	place := int32(math.Floor(math.Log2(bytes) / 10))
	num := bytes / math.Pow(1024.0, float64(place))
	var strnum string
	if place == 0 {
		strnum = fmt.Sprintf("%.0f", num) + suf[place]
	} else {
		strnum = fmt.Sprintf("%.1f", num) + suf[place]
	}
	return strnum
}

// BytesToString convert []byte to string
func BytesToString(bs []byte) string {
	return *(*string)(unsafe.Pointer(&bs))
}

// GetColStr return string use in field
func GetColStr(s string, size int, left bool) string {
	_, _, n := CountChars(s)
	spaces := strings.Repeat(" ", size-n)
	// size := nc*2 + ne // s 實際佔位數
	var tab string
	if left {
		tab = fmt.Sprintf("%[1]s%[2]s", s, spaces)
	} else {
		tab = fmt.Sprintf("%[2]s%[1]s", s, spaces)
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
	Sf := fmt.Sprintf
	nf := len(heads)
	ndat := len(data)
	if ndat < nf {
		return ""
	}
	if lensp < 0 {
		lensp = 0
	}
	lspaces := StrSpaces(lensp)
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

// StrThickLine returns a thick line (using '=')
func StrThickLine(n int) (l string) {
	l = strings.Repeat("=", n)
	return l + "\n"
}

// StrThinLine returns a thin line (using '-')
func StrThinLine(n int) (l string) {
	l = strings.Repeat("-", n)
	return l + "\n"
}

// StrSpaces returns a line with spaces
func StrSpaces(n int) (l string) {
	l = strings.Repeat(" ", n)
	return
}
