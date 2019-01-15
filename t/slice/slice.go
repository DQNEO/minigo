package main

import "fmt"

func f1() {
	var array [3]int = [3]int{1,2,3}
	fmt.Printf("%d\n", array[0] - 1) // 1

	var slice []int

	slice = array[:] // {1,2,3}
	fmt.Printf("%d\n", slice[0])
	fmt.Printf("%d\n", slice[1])
	fmt.Printf("%d\n", slice[2])

	slice = array[:3] // {1,2,3}
	fmt.Printf("%d\n", slice[0] + 3)
	fmt.Printf("%d\n", slice[1] + 3)
	fmt.Printf("%d\n", slice[2] + 3)

	slice = array[1:3] // {2,3}
	fmt.Printf("%d\n", slice[0] + 5)
	fmt.Printf("%d\n", slice[1] + 5)

	slice = array[2:3] // {3}
	fmt.Printf("%d\n", slice[0] + 6)

	slice = array[2:] // {3}
	fmt.Printf("%d\n", slice[0] + 7)
}

func f2() {

	var slice []int = []int{1,2,3}
	fmt.Printf("%d\n", slice[2] + 8)

	/*
	var bilbo = Hobbit{
		id:0,
		items:nil,
	}
	bilbo.items = []int{1,2,3}
	fmt.Printf("%d\n", bilbo.items[2])
	*/
}

func main() {
	f1()
	f2()
}

type Hobbit struct {
	id int
	items []int
}
