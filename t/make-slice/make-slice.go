package main

import "fmt"

func main() {
	var slc[]int = make([]int, 3, 5)
	fmt.Printf("%d\n", len(slc))
	fmt.Printf("%d\n", cap(slc))
}
