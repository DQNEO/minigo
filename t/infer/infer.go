package main

import "fmt"

func main() {
	a := 1
	b := '2'
	c := "3"
	d := User{
		id:4,
		age:0,
	}
	e := &User{
		id:5,
		age:0,
	}
	fmt.Printf("%d\n", a)
	fmt.Printf("%c\n", b)
	fmt.Printf("%s\n", c)
	fmt.Printf("%d\n", d.id)
	fmt.Printf("%d\n", e.id)
}

type User struct {
	id int
	age int
}
