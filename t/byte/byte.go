package main

import "fmt"

var hello [6]byte = [6]byte{'h','e','l','l','o',0}
func main() {
	fmt.Printf("%c", hello[0])
	fmt.Printf("%c", hello[1])
	fmt.Printf("%c", hello[2])
	fmt.Printf("%c", hello[3])
	fmt.Printf("%c\n", hello[4])

	fmt.Printf("%s\n", hello)

	/*

		var world [5]byte = [5]byte{'w','o','r','l','d'}
		fmt.Printf("%s\n", world)
	*/
}
