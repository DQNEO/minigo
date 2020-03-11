package main

import "fmt"

func f1() {
	var r []int
	r = append(r, 2)

	fmt.Printf("%d\n", len(r))
	fmt.Printf("%d\n", r[0])
}

func main() {
	f1()
}
