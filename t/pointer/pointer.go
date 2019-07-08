package main


func testint() {
	var i int = 1
	var j *int = &i

	fmtPrintf(S("%d\n"), *j)
	i = 2
	fmtPrintf(S("%d\n"), *j)
}

func testbyte() {
	var a byte = '3'
	var pa *byte = &a

	fmtPrintf(S("%c\n"), a)
	a = '4'
	fmtPrintf(S("%c\n"), *pa)
}

func testmixed() {
	var a byte = '5'
	var pa *byte = &a

	fmtPrintf(S("%c\n"), a)
	a = '6'
	fmtPrintf(S("%c\n"), *pa)

	var i int = 7
	var j *int = &i

	fmtPrintf(S("%d\n"), *j)
	i = 8
	fmtPrintf(S("%d\n"), *j)
	*j = 9
	fmtPrintf(S("%d\n"), i)
}

func main() {
	testint()
	testbyte()
	testmixed()
}
