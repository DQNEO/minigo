package main

import "fmt"

var globalarray [3]int = [3]int{3,2,1}

func main() {
	fmt.Printf("%d\n", globalarray[2])
	fmt.Printf("%d\n", globalarray[1])
	fmt.Printf("%d\n", globalarray[0])
	//globalarray[2] = 7
	//fmt.Printf("%d\n", globalarray[2])
}
