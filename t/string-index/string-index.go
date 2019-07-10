package main


func f1() {
	var s gostring = S("543210")
	var c byte = s[5]
	fmtPrintf(S("%c\n"), c)
	fmtPrintf(S("%c\n"), s[4])
	fmtPrintf(S("%c\n"), s[3])
	fmtPrintf(S("%c\n"), s[2])
	fmtPrintf(S("%c\n"), s[1])
	fmtPrintf(S("%c\n"), s[0])
}

func main() {
	f1()
}
