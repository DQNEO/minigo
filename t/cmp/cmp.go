package main


func f1() {
	var l = 1
	var g = 2

	if 1 != 1 {
		fmtPrintf(S("Error\n"))
	}
	if l < g {
		fmtPrintf(S("%d\n"), 1)
	}
	if l > g {
		fmtPrintf(S("Error\n"))
	}
	fmtPrintf(S("%d\n"), 2)
	if 1 == l {
		fmtPrintf(S("%d\n"), 3)
	}

	if g == 2 {
		fmtPrintf(S("%d\n"), 4)
	}

	if 1 <= l {
		fmtPrintf(S("%d\n"), 5)
	}

	if g >= 2 {
		fmtPrintf(S("%d\n"), 6)
	}
}

func f2() {
	if 1 == 0 || 1 == 1 {
		fmtPrintf(S("7\n"))
	} else {
		fmtPrintf(S("ERROR\n"))
	}

	if 1 == 1 && 1 == 0 {
		fmtPrintf(S("ERROR\n"))
	} else {
		fmtPrintf(S("8\n"))
	}
}

func f3() {
	var flg bool
	flg = true
	if flg {
		fmtPrintf(S("9\n"))
	}
	if !flg {
		fmtPrintf(S("ERROR\n"))
	}
	flg = false
	if !flg {
		fmtPrintf(S("10\n"))
	}
}

func f4() {
	if 0 > 20-1 {
		fmtPrintf(S("ERROR\n"))
	} else {
		fmtPrintf(S("11\n"))
	}
}
func main() {
	f1()
	f2()
	f3()
	f4()
}
