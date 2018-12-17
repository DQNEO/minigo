package main

import "fmt"

// For statements with single condition
func main() {
	var i int = 0
	for i < 5 {
		fmt.Printf("%d\n", i)
		i =  i + 1
	}
	var j int = i
	for 10 > j {
		fmt.Printf("%d\n", j)
		j =  j + 1
	}
}
