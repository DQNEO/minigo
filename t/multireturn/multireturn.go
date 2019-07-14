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
	fmtPrintf("%d\n", i)
	fmtPrintf("%d\n", j)
	fmtPrintf("%d\n", k)
	fmtPrintf("%d\n", l)

	a, b, c, d := multireverse(8, 7, 6, 5)
	fmtPrintf("%d\n", a)
	fmtPrintf("%d\n", b)
	fmtPrintf("%d\n", c)
	fmtPrintf("%d\n", d)
}

var bytes = [3]byte{'a', 'b', 'c'}

func ReadFile() ([]byte, int) {
	return bytes[:], 9
}

func getReadFile() {
	b, i := ReadFile()
	fmtPrintf("%d\n", i)
	fmtPrintf("%c\n", b[0])
	fmtPrintf("%c\n", b[1])
	fmtPrintf("%c\n", b[2])
}

func main() {
	getMulti()
	getReadFile()
}
