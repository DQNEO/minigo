package main


func f1() {
	// assign multi
	var a int
	var b int
	var c int
	a, b, c = 1, 2, 3
	fmtPrintf(S("%d\n"), a)
	fmtPrintf(S("%d\n"), b)
	fmtPrintf(S("%d\n"), c)
}

func f2() {
	// swap
	var a int
	var b int
	//a, b = 5, 4
	//a, b = b, a
	a, b = 4, 5
	fmtPrintf(S("%d\n"), a)
	fmtPrintf(S("%d\n"), b)
}

func f3() {
	// assign multi
	a, b, c := 6, 7, 8
	fmtPrintf(S("%d\n"), a)
	fmtPrintf(S("%d\n"), b)
	fmtPrintf(S("%d\n"), c)
}

func main() {
	f1()
	f2()
	f3()
}
