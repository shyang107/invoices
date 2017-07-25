package inv

import "os"

// func isOpened(path string) bool {
// 	fpath := os.ExpandEnv(path)
// 	f, err := os.Open(fpath)
// 	defer f.Close()
// 	if err != nil {
// 		perr("!!! File %q does not open !!!\n", fpath)
// 		return false
// 	}
// 	return true
// }

func isNotExist(path string) bool {
	return !CheckFileIsExist(path)
}

func isExist(path string) bool {
	return CheckFileIsExist(path)
}

// CheckFileIsExist 判断文件是否存在  存在返回 true 不存在返回false
func CheckFileIsExist(filename string) bool {
	path := os.ExpandEnv(filename)
	var isExist = true
	if _, err := os.Stat(path); os.IsNotExist(err) {
		isExist = false
	}
	return isExist
}
