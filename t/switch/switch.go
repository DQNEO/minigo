package main


func swtch(x int) int {
	var y int
	switch x {
	case 1:
		y = 1
	case 2, 3:
		y = 2
	case 5:
		y = 5
		y = 2
	default:
		y = 7
	}

	return y
}

func f1() {
	var i int
	i = swtch(1)
	fmtPrintf("%d\n", i) // 1
	i = swtch(2)
	fmtPrintf("%d\n", i) // 2
	i = swtch(3)
	fmtPrintf("%d\n", i+1) // 3
	i = swtch(999)
	fmtPrintf("%d\n", i-3) // 4
}

func swtch2(x int) int {
	var y int
	switch x {
	case 1:
		y = 1
	}

	return y
}

func f2() {
	i := swtch2(3)
	fmtPrintf("%d\n", i+5) // 5
}

func f3() {
	switch {
	case 1+1 == 3:
		fmtPrintf("Error\n")
	case 1+1 == 2:
		fmtPrintf("%d\n", 6)
	default:
		fmtPrintf("Error\n")
	}
}

func main() {
	f1()
	f2()
	f3()
}
