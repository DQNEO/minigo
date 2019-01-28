package main

import "fmt"


func main() {
	var heapHead *int
	var address *int
	heapHead = malloc(0)

	address = malloc(8)
	*address = 1
	fmt.Printf("%d\n", *address)
	address = malloc(8)
	*address = 2
	fmt.Printf("%d\n", *address)
	address = malloc(8)
	*address = 3
	fmt.Printf("%d\n", *address)

	fmt.Printf("%d\n", (address - heapHead)  - 20) // 4
}
