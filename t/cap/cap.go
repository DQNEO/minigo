package main

import "fmt"

func f0() {
	var x []int
	l := cap(x)
	fmt.Printf("%d\n", l)
}

func f1() {
	var a [1]int
	fmt.Printf("%d\n", cap(a))

	var b [2]int
	fmt.Printf("%d\n", cap(b))

	var c []int = b[:]
	fmt.Printf("%d\n", cap(c)+1) // 3

	c = b[0:1]
	fmt.Printf("%d\n", cap(c)+2) // 4

	c = b[1:2]
	fmt.Printf("%d\n", cap(c)+3) // 5

	var d []int = []int{1, 2, 3, 4, 5, 6}
	fmt.Printf("%d\n", cap(d)) // 6
}

func f2() {
	type Hobbit struct {
		id    int
		items []int
	}
	var h Hobbit
	h.items = []int{1}
	fmt.Printf("%d\n", cap(h.items)+6)          // 7
	fmt.Printf("%d\n", cap([]byte{'a', 'b'})+6) //8
	var array [10]int
	fmt.Printf("%d\n", cap(array[2:7])+1) // 9
}

func f3() {
	var array = [...]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}
	fmt.Printf("%d\n", cap(array))
}

func main() {
	f0()
	f1()
	f2()
	f3()
}
