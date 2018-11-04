package main

import "fmt"

func f1() {
	// assign multi
	var a int
	var b int
	var c int
	a, b, c = 1 ,2, 3
	fmt.Printf("%d\n", a)
	fmt.Printf("%d\n", b)
	fmt.Printf("%d\n", c)
}

func f2() {
	// swap
	var a int
	var b int
	a, b = 5, 4
	a, b = b, a
	fmt.Printf("%d\n", a)
	fmt.Printf("%d\n", b)
}

func main() {
	f1()
	f2()
}
