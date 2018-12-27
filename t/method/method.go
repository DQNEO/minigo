package main

import "fmt"

func f1() {
	var a myint = 0
	var b int
	b = a.add3(1)
	fmt.Printf("%d\n", b)
}

func f2() {
	var u *User
	u = &User{
		id:1,
		age:5,
	}

	fmt.Printf("%d\n", u.getAge())
}

func main() {
	f1()
	f2()
}

func (x myint) add3(y int) int {
	return 3 + y
}

type myint int

func (u *User) getAge() int {
	return u.age
}

type User struct {
	id int
	age int
}
