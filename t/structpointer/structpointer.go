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

type S struct {
	dummy *int
	id int
}

func f1() {
	var p *S
	p = &S{
		id:123,
	}

	fmt.Printf("%d\n", p.id - 115) // 8

	p.dummy = nil
	fmt.Printf("%d\n", p.id - 114) // 9
}

func main() {
	literal()
	assign()
	f1()
}

type User struct {
	id  int
	age int
}
