package main

import "fmt"

func f1() {
	var x map[int]map[int]int =   map[int]map[int]int{}
	fmt.Printf("%d\n", x[0][0])
}

func f2() {
	var mi MapIntInt = map[int]int{
		5:1,
	}

	fmt.Printf("%d\n", mi[5])

	var x map[int]map[int]int =  map[int]map[int]int{
		11: map[int]int{
			13:2,
		},
	}

	y := x[11]
	z := y[13]
	fmt.Printf("%d\n", z)
}

type MapIntInt map[int]int

func main() {
	f1()
	f2()
}
