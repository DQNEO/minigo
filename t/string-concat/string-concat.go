package main

import "fmt"

func f1() {
	var a = "abc"
	var b = "defg"
	var x string
	x = a + b
	fmt.Printf("%s\n", x)
}

func f2() {
	spaces := "> "
	for i := 0; i < 3; i++ {
		spaces += "xx"
	}
	fmt.Printf("%s\n", spaces)
}

func main() {
	f1()
	f2()
}
