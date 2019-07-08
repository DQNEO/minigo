package main


var var2 int = 2

func main() {
	fmtPrintf(S("%s\n"), const1)
	fmtPrintf(S("%d\n"), var2)
	fmtPrintf(S("%d\n"), const3)
	func4()
	func5() // mutuall dependent
}

func func5sub() {
	fmtPrintf(S("5\n"))
}
