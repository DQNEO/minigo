package main

import "fmt"

func main() {
	var t bool = true
	if t {
		fmt.Printf("1\n")
	}
	t = false
	if t {
		fmt.Printf("Error\n")
	}
	fmt.Printf("2\n")

	t = true
	if t {
		fmt.Printf("3\n")
	} else {
		fmt.Printf("Error\n")
	}

	t = false
	if t {
		fmt.Printf("Error\n")
	} else {
		fmt.Printf("4\n")
	}
}
