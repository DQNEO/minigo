package main

import "fmt"

type Ifc interface {
	getId() int
}

type Point struct {
	x int
	y int
}

func (p *Point) getId() int {
	return p.x
}

func f0() {
	var slice []int = []int{1, 2, 3}
	fmt.Printf("%d\n", slice[2]-2)   // 1
	fmt.Printf("%d\n", len(slice)-1) // 2
}

func f1() {
	var e Ifc = &Point{
		x: 1,
		y: 2,
	}
	fmt.Printf("%d\n", e.getId()+2) // 3
	var slice []Ifc = []Ifc{e, e, e}
	fmt.Printf("%d\n", len(slice)+1)       // 4
	fmt.Printf("%d\n", slice[2].getId()+4) // 5
}

func main() {
	f0()
	f1()
}
