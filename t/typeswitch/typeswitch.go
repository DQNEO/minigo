package main

import "fmt"

func typeswtch(flg bool) int {
	var r int
	var p *Point = &Point{
		a: 1,
		b: 2,
	}

	var p2 *Point2 = &Point2{
		c: 2,
		d: 3,
	}

	var i myInterface
	if flg {
		i = p
	} else {
		i = p2
	}
	switch i.(type) {
	case *Point:
		r = 1
	case *Point2:
		r = 2
	default:
		r = 3
	}

	return r
}

func f3() {
	var x int
	x = typeswtch(true)
	fmt.Printf("%d\n", x) // 1

	x = typeswtch(false)
	fmt.Printf("%d\n", x) // 2
}

type myInterface interface {
	sum() int
}

type Point struct {
	a int
	b int
}

type Point2 struct {
	c int
	d int
}

func (p *Point) sum() int {
	return p.a + p.b
}

func (p *Point2) sum() int {
	return p.c + p.d
}

func main() {
	f3()
}
