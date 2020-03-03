package main

import "fmt"

func f1() {
	var i int = 96
	fmt.Printf("%d\n", i)

	var c byte
	c = 'a'
	i = int(c)
	fmt.Printf("%d\n", i)

	c = 'b'
	fmt.Printf("%c\n", c)

	c = 'c'
	fmt.Printf("%d\n", c)

}

func main() {
	f1()
}
