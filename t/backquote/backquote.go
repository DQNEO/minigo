package main


func f1() {
	var s string = `hello
`
	fmtPrintf(S(s))
}

func f2() {
	var s string = `h"e"llo
`
	fmtPrintf(S(s))
}

func main() {
	f1()
	f2()
}
