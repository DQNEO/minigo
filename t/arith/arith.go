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

func paren() {
	x := 3 * (1 + 1)
	y := (1+1)*3 - (1 - 2)
	fmt.Printf("%d\n", x)
	fmt.Printf("%d\n", y)
}

func main() {
	divmode()
	uop_minus()
	paren()
}
