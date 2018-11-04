package main

import "fmt"

func f1() {
	var l = 1
	var g = 2

	if 1 != 1 {
		fmt.Printf("Error\n")
	}
	if l < g {
		fmt.Printf("%d\n", 1)
	}
	if l > g {
		fmt.Printf("Error\n")
	}
	fmt.Printf("%d\n", 2)
	if 1 == l {
		fmt.Printf("%d\n", 3)
	}

	if g == 2 {
		fmt.Printf("%d\n", 4)
	}

	if 1 <= l {
		fmt.Printf("%d\n", 5)
	}

	if g >= 2 {
		fmt.Printf("%d\n", 6)
	}
}

func f2() {
	if 1 == 0 || 1 == 1 {
		println("7")
	} else {
		println("ERROR")
	}

	if 1 == 1 && 1 == 0 {
		println("ERROR")
	} else {
		println("8")
	}
}

func main() {
	f1()
	f2()
}
