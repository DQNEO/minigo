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
	fmt.Printf("%d\n", *GlobalPtr)
}

func f2() {
	fmt.Printf("%d\n", pp1.y)
	fmt.Printf("%d\n", pp2.x)
}

func main() {
	f1()
	f2()
}

