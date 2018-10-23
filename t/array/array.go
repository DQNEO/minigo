package main

import "fmt"

var ary [3]int = [3]int{3,2,1}

func main() {
	fmt.Printf("%d\n", ary[2])
	fmt.Printf("%d\n", ary[1])
	fmt.Printf("%d\n", ary[0])
	//ary[2] = 7
	//fmt.Printf("%d\n", ary[2])
}
