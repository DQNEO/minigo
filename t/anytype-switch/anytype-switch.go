package main

import "fmt"

func f1() {
	var x interface{}
	var i int = 1
	x = i

	var z int
	switch x.(type) {
	case nil:
		z = -1
	case int:
		z = 1
	case string:
		z = 2
	default:
		z = 5
	}

	fmt.Printf("%d\n", z)

	var s string = "hello"
	x = s
	switch x.(type) {
	case nil:
		z = -1
	case int:
		z = 1
	case string:
		z = 2
	default:
		z = 5
	}
	fmt.Printf("%d\n", z)
}

func main() {
	f1()
}
