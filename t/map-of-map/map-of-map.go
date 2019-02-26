package main

import "fmt"

func f1() {
	var x map[int]map[int]int =   map[int]map[int]int{}
	fmt.Printf("%d\n", x[0][0])
}

func f2() {
	var mi MapIntInt = map[int]int{
		2:3,
	}

	fmt.Printf("%d\n", mi[2])

	var x map[int]map[int]int =  map[int]map[int]int{
		1: map[int]int{
			2:3,
		},
	}
	fmt.Printf("%d\n", x[1][2])
}

type MapIntInt map[int]int

func main() {
	f1()
	f2()
}
