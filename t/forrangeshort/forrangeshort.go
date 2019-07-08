package main


// for range test
func main() {
	var array1 [3]int = [3]int{9, 9, 9}
	var array2 [3]byte = [3]byte{'4', '6', '8'}

	for i := range array1 {
		fmtPrintf(S("%d\n"), i)
	}

	for k, v := range array2 {
		fmtPrintf(S("%d\n"), k*2+3)
		fmtPrintf(S("%c\n"), v)
	}
}
