package main

import "fmt"

func assign() {
	var u *User
	u = &User{
		id:  0,
		age: 4,
	}
	//u.age = 5
	fmt.Printf("%d\n", u.age)
}

func main() {
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
	assign()
}

type User struct {
	id  int
	age int
}
