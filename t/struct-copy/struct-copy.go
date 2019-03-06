package main

import "fmt"

func f1() {
	var h Hobbit = Hobbit{
		age:1,
		height:2,
	}

	fmt.Printf("%d\n", h.age)
	fmt.Printf("%d\n", h.height)

	var h2 Hobbit = h
	fmt.Printf("%d\n", h2.age + 2) // 3

	h.height = 100
	fmt.Printf("%d\n", h2.height + 2) // 4
}

type Hobbit struct {
	age int
	height int
}

func main() {
	f1()
}
