package main

import "fmt"

func plus(a int, b int) int {
	return a + b
}

func main() {
	fmt.Printf("%d\n", plus(0, 1))
	fmt.Printf("%d\n", plus(1, 1))
	fmt.Printf("%d\n", plus(2, 1))
	fmt.Printf("%d\n", plus(2, 2))
}

