package main

import "fmt"

var a = b
var b = 0

var x = iota
var iota = 0

func main() {
	fmt.Printf("%d\n", a)
	fmt.Printf("%d\n", x)
}
