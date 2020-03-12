package main

import "fmt"

func f1() {
	fmt.Printf("%d\n", 11)
	fmt.Printf("%s\n", "22")
}

func f2() {
	var s string = "hello"
	var i int = 123
	var b byte = 'b'
	var bl bool
	var _uintptr uintptr
	fmt.Printf("%T\n", s)
	fmt.Printf("%T\n", i)
	fmt.Printf("%T\n", b) // should be uint8
	fmt.Printf("%T\n", bl)
	fmt.Printf("%T\n", _uintptr)
}

func main() {
	f1()
	f2()
}
