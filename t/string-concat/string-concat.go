package main

import "fmt"

func f1() {
	var a = "abc"
	var b = "defg"
	var x string
	x = a + b
	fmt.Printf("%s\n", x)
}

func main() {
	f1()
}
