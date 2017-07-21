package inv

import "os"

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
