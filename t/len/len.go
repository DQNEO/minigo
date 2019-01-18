package main

import "fmt"

func f1() {
	var a [1]int
	fmt.Printf("%d\n", len(a))

	var b [2]int
	fmt.Printf("%d\n", len(b))

	var c []int = b[:]
	fmt.Printf("%d\n", len(c) + 1) // 3

	c = b[0:1]
	fmt.Printf("%d\n", len(c) + 3) // 4

	c = b[1:2]
	fmt.Printf("%d\n", len(c) + 4) // 5

	var d []int = []int{1,2,3,4,5,6}
	fmt.Printf("%d\n", len(d)) // 6
}

func main() {
	f1()
}
