package main

import "fmt"

func f0() {
	fmt.Printf("%d\n", add(1,1))
}

func f1() {
	var a myint = 1
	var b myint
	b = a.add(2)
	fmt.Printf("%d\n", b)
}

func f2() {
	var u *User
	u = &User{
		id:1,
		age:4,
	}

	fmt.Printf("%d\n", u.getAge())
}

func f3() {
	p := &Point{
		x: 2,
		y: 3,
	}
	fmt.Printf("%d\n", p.sum())
}

func main() {
	f0()
	f1()
	f2()
	f3()
}

func add(a myint, b myint) myint {
	return a + b
}

func (x myint) add(y myint) myint {
	return x + y
}

type myint int

func (u *User) getAge() int {
	return u.age
}

type User struct {
	id int
	age int
}


type Point struct {
	x int
	y int
}

func (p *Point) sum() int {
	return p.x + p.y
}
