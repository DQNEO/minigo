package main

import "fmt"

// for range test
func f1() {
	var array1 [3]int = [3]int{9, 9, 9}
	var array2 [3]int = [3]int{4, 6, 8}

	var v int
	var i int
	for i = range array1 {
		fmt.Printf("%d\n", i)
	}

	for i, v = range array2 {
		fmt.Printf("%d\n", i*2+3)
		fmt.Printf("%d\n", v)
	}
}

func f2() {
	bilbo := Hobbit{
		id:    1,
		age:   111,
		items: [3]int{9, 10, 11},
	}
	for _, v := range bilbo.items {
		fmt.Printf("%d\n", v)
	}
}

func main() {
	f1()
	f2()
}

type Hobbit struct {
	id    int
	age   int
	items [3]int
}
