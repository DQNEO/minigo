package main


func main() {
	var t bool = true
	if t {
		fmtPrintf("1\n")
	}
	t = false
	if t {
		fmtPrintf("Error\n")
	}
	fmtPrintf("2\n")

	t = true
	if t {
		fmtPrintf("3\n")
	} else {
		fmtPrintf("Error\n")
	}

	t = false
	if t {
		fmtPrintf("Error\n")
	} else {
		fmtPrintf("4\n")
	}

	var i int
	i = 1
	if i == 1 {
		fmtPrintf("5\n")
	} else if i == 2 {
		fmtPrintf("Error\n")
	} else {
		fmtPrintf("Error\n")
	}

	i = 2
	if i == 1 {
		fmtPrintf("Error\n")
	} else if i == 2 {
		fmtPrintf("6\n")
	} else {
		fmtPrintf("Error\n")
	}

	if i = 3; i == 1 {
		fmtPrintf("Error\n")
	} else if i == 2 {
		fmtPrintf("Error\n")
	} else {
		fmtPrintf("7\n")
	}
}
