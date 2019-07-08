package main


func multi(a int, b int, c int, d int) (int, int, int, int) {
	return a, b, c, d
}

func multireverse(a int, b int, c int, d int) (int, int, int, int) {
	return d, c, b, a
}

func getMulti() {
	var i int = 0
	var j int = 0
	var k int = 0
	var l int = 0
	i, j, k, l = multi(1, 2, 3, 4)
	fmtPrintf(S("%d\n"), i)
	fmtPrintf(S("%d\n"), j)
	fmtPrintf(S("%d\n"), k)
	fmtPrintf(S("%d\n"), l)

	a, b, c, d := multireverse(8, 7, 6, 5)
	fmtPrintf(S("%d\n"), a)
	fmtPrintf(S("%d\n"), b)
	fmtPrintf(S("%d\n"), c)
	fmtPrintf(S("%d\n"), d)
}

var bytes = [3]byte{'a', 'b', 'c'}

func ReadFile() ([]byte, int) {
	return bytes[:], 9
}

func getReadFile() {
	b, i := ReadFile()
	fmtPrintf(S("%d\n"), i)
	fmtPrintf(S("%c\n"), b[0])
	fmtPrintf(S("%c\n"), b[1])
	fmtPrintf(S("%c\n"), b[2])
}

func main() {
	getMulti()
	getReadFile()
}
