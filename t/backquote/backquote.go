package main


func f1() {
	var s gostring = S(`hello
`)
	fmtPrintf(s)
}

func f2() {
	var s gostring = S(`h"e"llo
`)
	fmtPrintf(s)
}

func main() {
	f1()
	f2()
}
