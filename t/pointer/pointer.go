package main


func testint() {
	var i int = 1
	var j *int = &i

	fmtPrintf("%d\n", *j)
	i = 2
	fmtPrintf("%d\n", *j)
}

func testbyte() {
	var a byte = '3'
	var pa *byte = &a

	fmtPrintf("%c\n", a)
	a = '4'
	fmtPrintf("%c\n", *pa)
}

func testmixed() {
	var a byte = '5'
	var pa *byte = &a

	fmtPrintf("%c\n", a)
	a = '6'
	fmtPrintf("%c\n", *pa)

	var i int = 7
	var j *int = &i

	fmtPrintf("%d\n", *j)
	i = 8
	fmtPrintf("%d\n", *j)
	*j = 9
	fmtPrintf("%d\n", i)
}

func main() {
	testint()
	testbyte()
	testmixed()
}
