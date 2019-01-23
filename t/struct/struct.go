package main

import "fmt"

func incr() {
	var u User = User{
		id:  0,
		age: 5,
	}

	//u.age = u.age + 1
	u.age++
	fmt.Printf("%d\n", u.age)
}

func decr() {
	var u User = User{
		id:  0,
		age: 8,
	}

	u.age--
	fmt.Printf("%d\n", u.age) // 7
}

func f1() {
	var i int = 1
	var u User = User{
		id:  3,
		age: 2,
		p: point{
			x:8,
			y:7,
		},
	}
	fmt.Printf("%d\n", i)
	fmt.Printf("%d\n", u.age)
	fmt.Printf("%d\n", u.id)

	u.id = 4
	fmt.Printf("%d\n", u.id)

	u = User{id: 3, age: 5}
	fmt.Printf("%d\n", u.age)

	incr()
	decr()

}

func f2() {
	var u User = User{
		id:  3,
		age: 2,
		p: point{
			x:8,
			y:9,
		},
	}

	fmt.Printf("%d\n", u.p.x) // 8
	fmt.Printf("%d\n", u.p.y) // 9
	u.p.y = 10
	fmt.Printf("%d\n", u.p.y) // 10
}

func main() {
	f1()
	f2()
}
type User struct {
	id  int
	age int
	p point
}

type point struct {
	x int
	y int
}
