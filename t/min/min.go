package main

import "fmt"

func main() {
	var t bool = true
	if t {
		fmt.Printf("1\n")
	}
	t = false
	if t {
		fmt.Printf("-1\n")
	}
	fmt.Printf("2\n")
}
