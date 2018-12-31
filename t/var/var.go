package main

import "fmt"

var a = b
var b = 1
var c = '2'

var x = iota
var iota = 0


func main() {
	fmt.Printf("%d\n", x)
	fmt.Printf("%d\n", a)
	localvar := 1
	fmt.Printf("%d\n", localvar)
	fmt.Printf("%c\n", c)
	a = 3
	fmt.Printf("%d\n", a)
}
