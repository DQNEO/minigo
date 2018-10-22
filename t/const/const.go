package main

import "fmt"

const a = b
const b = 1

const x = iota
const iota = 2

const sum = 1 + 2

func main() {
	fmt.Printf("%d\n", a)
	fmt.Printf("%d\n", x)
	fmt.Printf("%d\n", sum)
}
