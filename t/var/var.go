package main


var a = 1 // @TODO try a = b when possible
var b = 1
var c = 'A'

var x = 0

func main() {
	fmtPrintf(S("%d\n"), x)
	fmtPrintf(S("%d\n"), a)
	localvar := 1
	fmtPrintf(S("%d\n"), localvar)
	fmtPrintf(S("%d\n"), c)
	a = 3
	fmtPrintf(S("%d\n"), a)
}
