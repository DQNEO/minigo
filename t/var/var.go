package main


var a = 1 // @TODO try a = b when possible
var b = 1
var c = 'A'

var x = 0

func main() {
	fmtPrintf("%d\n", x)
	fmtPrintf("%d\n", a)
	localvar := 1
	fmtPrintf("%d\n", localvar)
	fmtPrintf("%d\n", c)
	a = 3
	fmtPrintf("%d\n", a)
}
