package main


func f1() {
	var s bytes = S(`hello
`)
	fmtPrintf(s)
}

func f2() {
	var s bytes = S(`h"e"llo
`)
	fmtPrintf(s)
}

func main() {
	f1()
	f2()
}
