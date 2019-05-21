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

/*
func f0() {
	var slice []int = []int{1, 2, 3}
	fmt.Printf("%d\n", slice[2]-2)   // 1
	fmt.Printf("%d\n", len(slice)-1) // 2
}
 */

func f1() {
	pp := &Point{
		x: 1,
		y: 2,
	}
	var e Ifc = pp
	printf("pp=%p\n", pp)
	dumpInterface(e)
	printf("*e=%p\n", *e)
	fmt.Printf("3=%d\n", e.getId()+2) // 3
	asComment("AAA")
	var slice []Ifc = []Ifc{e, e, e}
	asComment("BBB")
	fmt.Printf("4=%d\n", len(slice)+1)       // 4
	var e2 Ifc
	e2 = slice[1]
	dumpInterface(e2)
	z := e2.getId()
	printf("z=%d\n", z)
	fmt.Printf("%d\n", z) // 5
	return
}

func main() {
	//f0()
	f1()
}
