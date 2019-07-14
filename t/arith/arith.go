package main


func divmode() {
	var a = 5
	var b = 3
	fmtPrintf("%d\n", a/b)
	fmtPrintf("%d\n", a%b)

	fmtPrintf("%d\n", 3/1)
	fmtPrintf("%d\n", 4%5)
}

func uop_minus() {
	i := -3
	j := -i
	j += 2
	fmtPrintf("%d\n", j)
}

func paren() {
	x := 3 * (1 + 1)
	y := (1+1)*3 - (1 - 2)
	fmtPrintf("%d\n", x)
	fmtPrintf("%d\n", y)
}

func main() {
	divmode()
	uop_minus()
	paren()
}
