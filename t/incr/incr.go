package main

import "fmt"

func local() {
	var i int = 1
	fmt.Printf("%d\n", i)
	i++
	fmt.Printf("%d\n", i)
	i = 4
	i--
	fmt.Printf("%d\n", i)
}

func main() {
	local()
}
