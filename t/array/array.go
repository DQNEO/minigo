package main

import "fmt"

var garray [3]int = [3]int{1,2,3}

func main() {
	fmt.Printf("%d\n", garray[0])
	fmt.Printf("%d\n", garray[1])
	fmt.Printf("%d\n", garray[2])

	var larray [4]int = [4]int{4,5,6,7}
	fmt.Printf("%d\n", larray[0])
	fmt.Printf("%d\n", larray[1])
	fmt.Printf("%d\n", larray[2])
	fmt.Printf("%d\n", larray[3])
}
