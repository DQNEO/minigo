package main


var var2 int = 2

func main() {
	fmtPrintf("%s\n", const1)
	fmtPrintf("%d\n", var2)
	fmtPrintf("%d\n", const3)
	func4()
	func5() // mutuall dependent
}

func func5sub() {
	fmtPrintf("5\n")
}
