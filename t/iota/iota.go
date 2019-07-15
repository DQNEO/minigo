package main

import "fmt"

const (
	a0 int = iota
	a1
	a2
)

const (
	b0 int = iota
	b1     = iota
	b2     = iota
)

const (
	c0 int = iota
	c1
	c2
)

const (
	d0 int = 7
	d1
	d2 = iota
	d3
)

func f0() {
	fmt.Printf("%d\n", a0)
}

func f1() {
	fmt.Printf("%d\n", a1)
	fmt.Printf("%d\n", a2)
}

func f3() {
	fmt.Printf("%d\n", b0)
	fmt.Printf("%d\n", b1)
	fmt.Printf("%d\n", b2)
}

func main() {
	f0()
	f1()
	f3()
	fmt.Printf("%d\n", c0)
	fmt.Printf("%d\n", c1)
	fmt.Printf("%d\n", c2)
	fmt.Printf("%d\n", d0)
	fmt.Printf("%d\n", d1)
	fmt.Printf("%d\n", d2)
	fmt.Printf("%d\n", d3)
}
