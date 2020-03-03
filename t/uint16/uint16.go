package main

import "fmt"

func f1() {
	var a uint16 = 1024
	var b uint16 = 1024
	fmt.Printf("%d\n", int(a + b))
	var i int = int(a)
	fmt.Printf("%d\n", i)

	var c uint16 = 0
	c--
	i = int(c)
	i++
	fmt.Printf("%d\n", i)
}

func main() {
	f1()
}
