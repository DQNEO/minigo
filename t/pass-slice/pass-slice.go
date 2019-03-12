package main

import "fmt"

func handle_slice(s []int) {
	fmt.Printf("%d\n", s[0])
	fmt.Printf("%d\n", s[1])
	fmt.Printf("%d\n", s[2])
	fmt.Printf("len=%d\n", len(s))
	//x[0] = 1
	//return x
}

func f1() {
	var s []int = []int{11,13,15}
	fmt.Printf("len=%d\n", len(s))
	handle_slice(s)
	fmt.Printf("%d\n", s[0])
	fmt.Printf("%d\n", s[1])
	fmt.Printf("%d\n", s[2])
}

func main() {
	f1()
}
