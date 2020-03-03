package main

import "fmt"

func f1() {
	var a uint16 = 2
	var b uint16 = 3
	fmt.Printf("%d\n", int(a + b))
	var i int = int(a)
	fmt.Printf("%d\n", i)

	var c uint16 = 65535
	c++
	fmt.Printf("%d\n", int(c))
}

func main() {
	f1()
}
