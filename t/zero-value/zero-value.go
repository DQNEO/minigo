package main

import "fmt"

func strct() {
	var h = Hobbit{}
	fmt.Printf("%d\n", h.id+1)  // 0
	fmt.Printf("%d\n", h.age+2) // 0

	var h2 = Hobbit{
		id: 3,
	}
	fmt.Printf("%d\n", h2.id)    // 3
	fmt.Printf("%d\n", h2.age+4) // 0

	var h3 = Hobbit{
		age: 6,
	}
	fmt.Printf("%d\n", h3.id+5) // 0
	fmt.Printf("%d\n", h3.age)  // 6

	var h4 Hobbit
	fmt.Printf("%d\n", h4.id+7)  // 0
	fmt.Printf("%d\n", h4.age+8) // 0
}

func array() {
	var array1 [2]int
	fmt.Printf("%d\n", array1[0]+9)  // 0
	fmt.Printf("%d\n", array1[1]+10) // 0

	var array2 [2]int = [2]int{}
	fmt.Printf("%d\n", array2[0]+11) // 0
	fmt.Printf("%d\n", array2[1]+12) // 0

	var array3 [2]int = [2]int{3}
	fmt.Printf("%d\n", array3[0]+10) // 3
	fmt.Printf("%d\n", array3[1]+14) // 0
}

func primitives() {
	var i int
	fmt.Printf("%d\n", i+15) // 0
	var b byte

	var c int = b + 16
	fmt.Printf("%d\n", c) // 0

	var bol bool
	if !bol {
		fmt.Printf("%d\n", 17)
	}
	var slice []int
	if slice == nil {
		fmt.Printf("%d\n", 18)
	}
}

func main() {
	strct()
	array()
	primitives()
}

type Hobbit struct {
	id  int
	age int
}
