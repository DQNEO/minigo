package main


func f1() {
	var x byte = 10

	if x == '\n' {
		fmtPrintf("%d\n", 1)
	} else {
		fmtPrintf("error\n")
	}
}

func main() {
	f1()
}
