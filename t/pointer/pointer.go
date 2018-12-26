package main

import "fmt"

func testint() {
	var i int = 1
	var j *int = &i

	fmt.Printf("%d\n", *j)
	i = 2
	fmt.Printf("%d\n", *j)
}

func testbyte() {
	var a byte = '3'
	var pa *byte = &a

	fmt.Printf("%c\n", a)
	a = '4'
	fmt.Printf("%c\n", *pa)
}

func testmixed() {
	var a byte = '5'
	var pa *byte = &a

	fmt.Printf("%c\n", a)
	a = '6'
	fmt.Printf("%c\n", *pa)

	var i int = 7
	var j *int = &i

	fmt.Printf("%d\n", *j)
	i = 8
	fmt.Printf("%d\n", *j)
}

func main() {
	testint()
	testbyte()
	testmixed()
}
