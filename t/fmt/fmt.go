package main

import "fmt"

func f1() {
	var i int = 1
	fmt.Printf("%d\n", i)

	var c byte = 'a'
	fmt.Printf("%d\n", c)

	i = int(c)
	fmt.Printf("%d\n", i)
}

func main() {
	f1()
}
