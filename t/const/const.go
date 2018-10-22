package main

import "fmt"

const a = b
const b = 1

const x = iota
const iota = 7


func main() {
	fmt.Printf("1 == %d\n", a)
	fmt.Printf("7 == %d\n", x)
}
