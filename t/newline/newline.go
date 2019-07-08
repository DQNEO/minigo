package main


func f1() {
	var x byte = 10

	if x == '\n' {
		fmtPrintf(S("%d\n"), 1)
	} else {
		fmtPrintf(S("error\n"))
	}
}

func main() {
	f1()
}
