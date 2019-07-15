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

	var i int
	i = 1
	if i == 1 {
		fmt.Printf("5\n")
	} else if i == 2 {
		fmt.Printf("Error\n")
	} else {
		fmt.Printf("Error\n")
	}

	i = 2
	if i == 1 {
		fmt.Printf("Error\n")
	} else if i == 2 {
		fmt.Printf("6\n")
	} else {
		fmt.Printf("Error\n")
	}

	if i = 3; i == 1 {
		fmt.Printf("Error\n")
	} else if i == 2 {
		fmt.Printf("Error\n")
	} else {
		fmt.Printf("7\n")
	}
}
