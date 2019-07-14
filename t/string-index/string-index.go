package main


func f1() {
	var s []byte = []byte("543210")
	var c byte = s[5]
	fmtPrintf("%c\n", c)
	fmtPrintf("%c\n", s[4])
	fmtPrintf("%c\n", s[3])
	fmtPrintf("%c\n", s[2])
	fmtPrintf("%c\n", s[1])
	fmtPrintf("%c\n", s[0])
}

func main() {
	f1()
}
