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
		111: map[int]int{
			11:2,
		},
		112: map[int]int{
			12:3,
		},
	}

	y := x[111]
	z := y[11]
	fmt.Printf("%d\n", z)

	y = x[112]
	z = y[12]
	fmt.Printf("%d\n", z)
}

type MapIntInt map[int]int

func main() {
	f1()
	f2()
}
