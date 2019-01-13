package main

import "fmt"

func main() {
	var i int = 1
	var u user = user{
		id:  3,
		age: 2,
	}
	fmt.Printf("%d\n", i)
	fmt.Printf("%d\n", u.age)
	fmt.Printf("%d\n", u.id)

	u.id = 4
	fmt.Printf("%d\n", u.id)

	u = user{id: 3, age: 5}
	fmt.Printf("%d\n", u.age)

	u.age = u.age + 1
	fmt.Printf("%d\n", u.age)
}

type user struct {
	id  int
	age int
}
