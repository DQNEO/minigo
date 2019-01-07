package main

import "fmt"

func main() {
	var i int = 1
	var u user = user{
		id:  3,
		age: 2,
		p: point{
			x:6,
			y:7,
		},
	}
	fmt.Printf("%d\n", i)
	fmt.Printf("%d\n", u.age)
	fmt.Printf("%d\n", u.id)

	u.id = 4
	fmt.Printf("%d\n", u.id)

	u = user{id:  3, age: 5}
	fmt.Printf("%d\n", u.age)

	fmt.Printf("%d\n", u.p.x)
	u.p.x = 6
	fmt.Printf("%d\n", u.p.x)
}

type user struct {
	id  int
	age int
	p point
}

type point struct {
	x int
	y int
}
