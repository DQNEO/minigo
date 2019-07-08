package main


func main() {
	var t bool = true
	if t {
		fmtPrintf(S("1\n"))
	}
	t = false
	if t {
		fmtPrintf(S("Error\n"))
	}
	fmtPrintf(S("2\n"))

	t = true
	if t {
		fmtPrintf(S("3\n"))
	} else {
		fmtPrintf(S("Error\n"))
	}

	t = false
	if t {
		fmtPrintf(S("Error\n"))
	} else {
		fmtPrintf(S("4\n"))
	}

	var i int
	i = 1
	if i == 1 {
		fmtPrintf(S("5\n"))
	} else if i == 2 {
		fmtPrintf(S("Error\n"))
	} else {
		fmtPrintf(S("Error\n"))
	}

	i = 2
	if i == 1 {
		fmtPrintf(S("Error\n"))
	} else if i == 2 {
		fmtPrintf(S("6\n"))
	} else {
		fmtPrintf(S("Error\n"))
	}

	if i = 3; i == 1 {
		fmtPrintf(S("Error\n"))
	} else if i == 2 {
		fmtPrintf(S("Error\n"))
	} else {
		fmtPrintf(S("7\n"))
	}
}
