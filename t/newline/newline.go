package main

import "fmt"

func f1() {
	var x byte = 10

	if x == '\n' {
		fmt.Printf("%d\n", 1)
	} else {
		println("error")
	}
}

func main() {
	f1()
}
