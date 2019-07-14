package main


func main() {
	var i int = 1
	fmtPrintf("%d\n", i)
	i = 0
	i += 2
	fmtPrintf("%d\n", i)
	i = 5
	i -= 2
	fmtPrintf("%d\n", i)
	i = 2
	i *= 2
	fmtPrintf("%d\n", i)
}
