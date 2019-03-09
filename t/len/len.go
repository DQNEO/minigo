package main

import "fmt"

func f0() {
	var x []int
	l := len(x)
	fmt.Printf("%d\n", l)
}

func f1() {
	var a [1]int
	fmt.Printf("%d\n", len(a))

	var b [2]int
	fmt.Printf("%d\n", len(b))

	var c []int = b[:]
	fmt.Printf("%d\n", len(c)+1) // 3

	c = b[0:1]
	fmt.Printf("%d\n", len(c)+3) // 4

	c = b[1:2]
	fmt.Printf("%d\n", len(c)+4) // 5

	var d []int = []int{1, 2, 3, 4, 5, 6}
	fmt.Printf("%d\n", len(d)) // 6
}

func f2() {
	type Hobbit struct {
		id    int
		items []int
	}
	var h Hobbit
	h.items = []int{1}
	fmt.Printf("%d\n", len(h.items)+6)          // 7
	fmt.Printf("%d\n", len([]byte{'a', 'b'})+6) //8
	var array [10]int
	fmt.Printf("%d\n", len(array[2:7])+4) // 9
}

func f3() {
	var array = [...]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}
	fmt.Printf("%d\n", len(array))
}

func main() {
	f0()
	f1()
	f2()
	f3()
}
