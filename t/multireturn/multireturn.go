package main

import "fmt"

func multi(a int, b int, c int, d int) (int, int, int,int ) {
	return b,a,d,c
}

func main() {
	var i int = 0
	var j int = 0
	var k int = 0
	var l int = 0
	i, j,k,l = multi(2,1,4,3)
	fmt.Printf("%d\n", i)
	fmt.Printf("%d\n", j)
	fmt.Printf("%d\n", k)
	fmt.Printf("%d\n", l)
}
