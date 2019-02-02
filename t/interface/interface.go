package main

import "fmt"

func f1() {
	var p *Point
	p = &Point{
		x: 1,
		y: 2,
	}
	sum := p.sum()
	fmt.Printf("%d\n", sum - 2) // 1
}

func f2() {
	var myInterface MyInterface
	ptr := &Point{
		x: 2,
		y: 3,
	}
	myInterface = ptr
	sum := myInterface.sum()
	fmt.Printf("%d\n", sum - 3) // 2
}

func f3(bol bool) {
	var myInterface MyInterface
	ptr := &Point{
		x: 2,
		y: 3,
	}

	asset := &Asset{
		money: 4,
		stock: 2,
	}

	if bol {
		myInterface = ptr
	} else {
		myInterface = asset
	}

	sum := myInterface.sum()

	fmt.Printf("%d\n", sum - 3) // 2
}

func main() {
	f1()
	f2()
	f3(true)
	f3(false)
}

type MyInterface interface {
	sum() int
}

type Point struct {
	x int
	y int
}

func (p *Point) sum() int {
	return p.x + p.y
}

type Asset struct {
	money int
	stock int
}

func (a *Asset) sum() int {
	return a.money + a.stock
}
