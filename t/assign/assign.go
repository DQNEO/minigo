package main

import "fmt"

func main() {
	var i int = 1
	fmt.Printf("%d\n", i)
	i = 0
	i += 2
	fmt.Printf("%d\n", i)
	i = 5
	i -= 2
	fmt.Printf("%d\n", i)
	i = 2
	i *= 2
	fmt.Printf("%d\n", i)
}
