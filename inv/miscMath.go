package inv

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
