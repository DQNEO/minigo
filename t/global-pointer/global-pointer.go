package main

import "fmt"

var GlobalInt int = 1
var GlobalPtr *int = &GlobalInt

func f1() {
	fmt.Printf("%d\n", *GlobalPtr)
}

func main() {
	f1()
}

