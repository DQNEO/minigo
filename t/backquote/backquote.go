package main


func f1() {
	var s string = `hello
`
	fmtPrintf(s)
}

func f2() {
	var s string = `h"e"llo
`
	fmtPrintf(s)
}

func main() {
	f1()
	f2()
}
