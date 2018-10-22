package main

import "fmt"

const a = b
const b = 1

const x = iota
const iota = 2

func main() {
	fmt.Printf("1 == %d\n", a)
	fmt.Printf("2 == %d\n", x)
}
