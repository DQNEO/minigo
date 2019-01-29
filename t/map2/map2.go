package main

import "fmt"

func f1() {
	var lmap map[int]int = map[int]int {
		1: 2,
		3: 4,
	}

	for i,v := range lmap {
		fmt.Printf("%d\n", i)
		fmt.Printf("%d\n", v)
	}

	fmt.Printf("%d\n", lmap[1] + 3) // 5
	fmt.Printf("%d\n", lmap[3] + 2) // 6

	lmap[7] = 8
	fmt.Printf("%d\n", lmap[4] + 7 ) // 7
	fmt.Printf("%d\n", lmap[7]) // 8
}

func main() {
	f1()
}
