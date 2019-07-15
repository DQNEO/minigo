package main

import "fmt"

func f1() {
	var s string = `hello
`
	fmt.Printf(s)
}

func f2() {
	var s string = `h"e"llo
`
	fmt.Printf(s)
}

func main() {
	f1()
	f2()
}
