package main

import "fmt"

func literal() {
	var u *User
	u = &User{
		id:  1,
		age: 2,
	}
	fmt.Printf("%d\n", u.id)
	fmt.Printf("%d\n", u.age)

	u = &User{
		id:  3,
		age: 4,
	}
	fmt.Printf("%d\n", u.id)
	fmt.Printf("%d\n", u.age)
}

func assign() {
	var u *User
	u = &User{
		id:  0,
		age: 4,
	}
	u.age = 5
	fmt.Printf("%d\n", u.age)
	u.age++
	fmt.Printf("%d\n", u.age)
	u.age = 8
	u.age--
	fmt.Printf("%d\n", u.age)
}

func main() {
	literal()
	assign()
}

type User struct {
	id  int
	age int
}
