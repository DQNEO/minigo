package main


func f1() {
	var l = 1
	var g = 2

	if 1 != 1 {
		fmtPrintf("Error\n")
	}
	if l < g {
		fmtPrintf("%d\n", 1)
	}
	if l > g {
		fmtPrintf("Error\n")
	}
	fmtPrintf("%d\n", 2)
	if 1 == l {
		fmtPrintf("%d\n", 3)
	}

	if g == 2 {
		fmtPrintf("%d\n", 4)
	}

	if 1 <= l {
		fmtPrintf("%d\n", 5)
	}

	if g >= 2 {
		fmtPrintf("%d\n", 6)
	}
}

func f2() {
	if 1 == 0 || 1 == 1 {
		fmtPrintf("7\n")
	} else {
		fmtPrintf("ERROR\n")
	}

	if 1 == 1 && 1 == 0 {
		fmtPrintf("ERROR\n")
	} else {
		fmtPrintf("8\n")
	}
}

func f3() {
	var flg bool
	flg = true
	if flg {
		fmtPrintf("9\n")
	}
	if !flg {
		fmtPrintf("ERROR\n")
	}
	flg = false
	if !flg {
		fmtPrintf("10\n")
	}
}

func f4() {
	if 0 > 20-1 {
		fmtPrintf("ERROR\n")
	} else {
		fmtPrintf("11\n")
	}
}
func main() {
	f1()
	f2()
	f3()
	f4()
}
