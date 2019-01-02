package main

import "fmt"

func multi() (int, int, int,int ) {
	return 1,2,3,4
}

func main() {
	var i int = 0
	var j int = 0
	var k int = 0
	var l int = 0
	i, j,k,l = multi()
	fmt.Printf("%d\n", i)
	fmt.Printf("%d\n", j)
	fmt.Printf("%d\n", k)
	fmt.Printf("%d\n", l)
}
