package main

import (
	"fmt"
	"strconv"
)

func f1() {
	var a string = "10485760"
	var i int
	i, _ = strconv.Atoi(a)
	fmt.Printf("%d\n", i-10485760) // 0

	a = "1"
	i, _ = strconv.Atoi(a)
	fmt.Printf("%d\n", i) // 1
}

func main() {
	f1()
}
