package main

import "fmt"

var myarray [2]myint = [2]myint{3, 2}

func main() {
	var a myint = 1
	fmt.Printf("%d\n", a)
	fmt.Printf("%d\n", myarray[1])
}

type myint int

type int byte
