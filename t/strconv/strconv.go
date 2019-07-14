package main

import "strconv"

func f1() {
	var a string = "10485760"
	var i int
	var e error
	i, e = strconv.Atoi(a)
	fmtPrintf("%d\n", i-10485760) // 0

	a = "1"
	i ,e  = strconv.Atoi(a)
	fmtPrintf("%d\n", i) // 1

	a = "-2"
	i ,e  = strconv.Atoi(a)
	fmtPrintf("%d\n", i+4) // 1
}

func f2() {
	var s string
	s = strconv.Itoa(0)
	fmtPrintf("%s\n", s)

	s = strconv.Itoa(7)
	fmtPrintf("%s\n", s)

	s = strconv.Itoa(10)
	fmtPrintf("%s\n", s)

	s = strconv.Itoa(100)
	fmtPrintf("%s\n", s)

	s = strconv.Itoa(1234567890)
	fmtPrintf("%s\n", s)

	s = strconv.Itoa(-1)
	fmtPrintf("%s\n", s)

	s = strconv.Itoa(-7)
	fmtPrintf("%s\n", s)

	s = strconv.Itoa(-10)
	fmtPrintf("%s\n", s)

	s = strconv.Itoa(-100)
	fmtPrintf("%s\n", s)

	s = strconv.Itoa(-1234567890)
	fmtPrintf("%s\n", s)

}

func main() {
	f1()
	f2()
}
