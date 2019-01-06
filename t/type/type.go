package main

import "fmt"

var myarray [2]myint = [2]myint{3, 2}

func anytype(x interface{}, y interface{}) {
	fmt.Printf("%d\n", x)
	fmt.Printf("%d\n", y)
}

func main() {
	var a myint = 1
	fmt.Printf("%d\n", a)
	fmt.Printf("%d\n", myarray[1])
	anytype(3, 4)
}

type myint int

type int byte
