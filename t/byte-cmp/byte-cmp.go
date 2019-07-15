package main

import "fmt"

func f1() {
	var x byte = 'a'
	var e byte = 'e'

	if e <= 'z' {
		fmt.Printf("1\n")
	} else {
		fmt.Printf("%s\n", x)
	}
}

func f2() {
	var c1 byte = 'p'
	var c2 byte = 'a'

	if 'a' <= c1 && c1 <= 'z' {
		fmt.Printf("2\n")
	}

	if 'a' <= c2 && c2 <= 'z' {
		fmt.Printf("3\n")
	}
}

func main() {
	f1()
	f2()
}
