package main


func sum(a int, b int) int {
	fmtPrintf(S("%d\n"), sum)
	return a + b
}

func main() {
	fmtPrintf(S("%d\n"), sum)
	s := sum(1, 2)
	fmtPrintf(S("%d\n"), s)
}
