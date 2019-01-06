package main

import "fmt"

func main() {
	var x int
	x = -1
	switch x {
	case 0:
		x = 0
	case -1:
		x = 1
	default:
		x = 2
	}

	fmt.Printf("%d\n", x)
}
