package main

import "fmt"

func main() {
	var t bool = true
	if t {
		fmt.Println("1")
	}
	t = false
	if t {
		fmt.Println("Error")
	}
	fmt.Println("2")

	t = true
	if t {
		fmt.Println("3")
	} else {
		fmt.Println("Error")
	}

	t = false
	if t {
		fmt.Println("Error")
	} else {
		fmt.Println("4")
	}

	var i int
	i = 1
	if i == 1 {
		fmt.Println("5")
	} else if i == 2 {
		fmt.Println("Error")
	} else {
		fmt.Println("Error")
	}

	i = 2
	if i == 1 {
		fmt.Println("Error")
	} else if i == 2 {
		fmt.Println("6")
	} else {
		fmt.Println("Error")
	}

	if i = 3; i == 1 {
		fmt.Println("Error")
	} else if i == 2 {
		fmt.Println("Error")
	} else {
		fmt.Println("7")
	}
}
