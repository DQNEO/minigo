package main

import "fmt"

func f1() {
	var p *Point
	p = &Point{
		x: 1,
		y: 2,
	}
	fmt.Printf("%d\n", p.sum())
}

func f2() {
	//var s summer
	var s *Point
	s = &Point{
		x:2,
		y:3,
	}
	fmt.Printf("%d\n", s.sum())
}

func main() {
	f1()
	f2()
}

type summer interface {
	sum() int
}

type Point struct {
	x int
	y int
}

func (p *Point) sum() int {
	return p.x + p.y
}
