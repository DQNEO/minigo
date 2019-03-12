package main

import "fmt"

func handle_slice(s []int) {
	fmt.Printf("%d\n", len(s)-2) // 1
	fmt.Printf("%d\n", cap(s)-3) // 2
	fmt.Printf("%d\n", s[0]-12)  // 3
	fmt.Printf("%d\n", s[1]-10)  // 4
	s[2] = 3
	fmt.Printf("%d\n", s[2]+2) // 5
}

var array [5]int = [...]int{15, 14, 13, 12, 11}

func f1() {
	var s []int = array[0:3]
	handle_slice(s)
	fmt.Printf("%d\n", len(s)+3) // 6
	fmt.Printf("%d\n", cap(s)+2) // 7
	fmt.Printf("%d\n", s[0]-7)   // 8
	fmt.Printf("%d\n", s[1]-5)   // 9
	fmt.Printf("%d\n", s[2]+7)   // 10
}

func main() {
	f1()
}
