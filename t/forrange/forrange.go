package main

import "fmt"

// for range test
func main() {
	var array1 [3]int = [3]int{9,9,9}
	var array2 [3]int = [3]int{4,6,8}

	var v int
	var i int
	for i = range array1 {
		fmt.Printf("%d\n", i)
	}

	for i,v = range array2 {
		fmt.Printf("%d\n", i * 2 + 3)
		fmt.Printf("%d\n", v)
	}
}
