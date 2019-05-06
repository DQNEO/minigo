package main

import "fmt"

var GlobalInt int = 1
var GlobalPtr *int = &GlobalInt

var pp1 *Point = &Point{x: 1, y: 2,}
var pp2 *Point = &Point{x: 3, y: 4,}

type Point struct {
	x int
	y int
}


func f1() {
	fmt.Printf("%d\n", *GlobalPtr) // 1
}

func f2() {
	fmt.Printf("%d\n", pp1.y) // 2
	fmt.Printf("%d\n", pp2.x) // 3
}

type Gtype struct {
	typ int
	size int
}

var gInt  = &Gtype{typ: 1, size: 8}

type DeclFunc struct {
	tok               string
	pkg               string
	receiver          *int
	fname             string
	rettypes          []*Gtype
}

var builinLen = &DeclFunc{
	rettypes: []*Gtype{gInt},
}

func f3() {
	fmt.Printf("%d\n", len(builinLen.rettypes) + 3) // 4
}

func main() {
	f1()
	f2()
	f3()
}

