package main

import "fmt"

// for range test
func main() {
	var array1 [3]int = [3]int{9,9,9}
	var array2 [3]int = [3]int{4,6,8}

	for i := range array1 {
		fmt.Printf("%d\n", i)
	}

	for k,v := range array2 {
		fmt.Printf("%d\n", k * 2 + 3)
		fmt.Printf("%d\n", v)
	}
}
