package main

import "fmt"

func main() {
	var h = Hobbit{}
	fmt.Printf("%d\n", h.id + 1) // 0
	fmt.Printf("%d\n", h.age + 2) // 0

	var h2 = Hobbit{
		id:3,
	}
	fmt.Printf("%d\n", h2.id) // 3
	fmt.Printf("%d\n", h2.age + 4) // 0

	var h3 = Hobbit{
		age:6,
	}
	fmt.Printf("%d\n", h3.id + 5) // 0
	fmt.Printf("%d\n", h3.age) // 6

	var h4 Hobbit
	fmt.Printf("%d\n", h4.id + 7) // 0
	fmt.Printf("%d\n", h4.age + 8) // 0
}

type Hobbit struct {
	id int
	age int
}
