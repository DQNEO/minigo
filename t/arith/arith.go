package main

import "fmt"

func divmode() {
	var a = 5
	var b = 3
	fmt.Printf("%d\n", a/b)
	fmt.Printf("%d\n", a%b)

	fmt.Printf("%d\n", 3/1)
	fmt.Printf("%d\n", 4%5)
}

func uop_minus() {
	i := -3
	j := -i
	j += 2
	fmt.Printf("%d\n", j)
}

func main() {
	divmode()
	uop_minus()
}
