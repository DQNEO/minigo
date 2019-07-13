package main


func f1() {
	var a bytes = S("10485760")
	var i int
	var e error
	i, e = Atoi(a)
	fmtPrintf(S("%d\n"), i-10485760) // 0

	a = S("1")
	i ,e  = Atoi(a)
	fmtPrintf(S("%d\n"), i) // 1

	a = S("-2")
	i ,e  = Atoi(a)
	fmtPrintf(S("%d\n"), i+4) // 1
}

func f2() {
	var s []byte
	s = Itoa(0)
	fmtPrintf(S("%s\n"), s)

	s = Itoa(7)
	fmtPrintf(S("%s\n"), s)

	s = Itoa(10)
	fmtPrintf(S("%s\n"), s)

	s = Itoa(100)
	fmtPrintf(S("%s\n"), s)

	s = Itoa(1234567890)
	fmtPrintf(S("%s\n"), s)

	s = Itoa(-1)
	fmtPrintf(S("%s\n"), s)

	s = Itoa(-7)
	fmtPrintf(S("%s\n"), s)

	s = Itoa(-10)
	fmtPrintf(S("%s\n"), s)

	s = Itoa(-100)
	fmtPrintf(S("%s\n"), s)

	s = Itoa(-1234567890)
	fmtPrintf(S("%s\n"), s)

}

func main() {
	f1()
	f2()
}
