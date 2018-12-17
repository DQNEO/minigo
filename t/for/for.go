package main

import "fmt"

func main() {
	// C style
	for i:= 0; i < 10; i = i + 1 {
		fmt.Printf("%d\n", i)
	}
}
