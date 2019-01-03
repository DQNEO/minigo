package main

import "fmt"

func main() {
	f,g := getIntInt()
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
	fmt.Printf("%d\n", f)
	fmt.Printf("%d\n", g)
}

type User struct {
	id int
	age int
}

func getIntInt() (int, int) {
	return 6,7
}
