package main

import "fmt"

func structfield() {
	bilbo := Hobbit{
		id: 1,
		age: 2,
		items: [3]int{3,4,5},
	}

	fmt.Printf("%d\n", bilbo.id)
	fmt.Printf("%d\n", bilbo.age)
	fmt.Printf("%d\n", bilbo.items)
}

type Hobbit struct {
	id int
	age int
	items [3]int
}

func main() {
	structfield()
}
