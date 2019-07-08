package main


func divmode() {
	var a = 5
	var b = 3
	fmtPrintf(S("%d\n"), a/b)
	fmtPrintf(S("%d\n"), a%b)

	fmtPrintf(S("%d\n"), 3/1)
	fmtPrintf(S("%d\n"), 4%5)
}

func uop_minus() {
	i := -3
	j := -i
	j += 2
	fmtPrintf(S("%d\n"), j)
}

func paren() {
	x := 3 * (1 + 1)
	y := (1+1)*3 - (1 - 2)
	fmtPrintf(S("%d\n"), x)
	fmtPrintf(S("%d\n"), y)
}

func main() {
	divmode()
	uop_minus()
	paren()
}
