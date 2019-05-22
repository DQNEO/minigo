package main

import "fmt"

func main() {
	var heapHead *int
	var heapTail *int

	var address *int

	address = malloc(8)
	*address = 1
	fmt.Printf("%d\n", *address)
	address = malloc(8)
	*address = 2
	fmt.Printf("%d\n", *address)
	address = malloc(8)
	*address = 3
	fmt.Printf("%d\n", *address)

	heapA := malloc(8)
	heapB := malloc(0)

	fmt.Printf("%d\n", (heapB-heapA) - 4) // 4
}
