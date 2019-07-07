package main

import "fmt"

func f1() {
	var x byte = 'a'
	var e byte = 'e'

	if e <= 'z' {
		fmt.Println("1")
	} else {
		fmt.Println(x)
	}
}

func f2() {
	var c1 byte = 'p'
	var c2 byte = 'a'

	if 'a' <= c1 && c1 <= 'z' {
		fmt.Println("2")
	}

	if 'a' <= c2 && c2 <= 'z' {
		fmt.Println("3")
	}
}

func main() {
	f1()
	f2()
}
