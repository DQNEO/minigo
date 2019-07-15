package main

import "fmt"

const a int = b
const b int = 1

const x int = iota
const iota int = 2

const sum int = 1 + 2

const (
	c = 4
	d = 5
)

func f1() {
	fmt.Printf("%d\n", 0)
}

func f2() {
	fmt.Printf("%d\n", a)
}

func main() {
	f1()
	f2()
	fmt.Printf("%d\n", x)
	fmt.Printf("%d\n", sum)
}
