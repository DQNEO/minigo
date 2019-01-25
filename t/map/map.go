package main

import "fmt"

var gmap map[string]bool

func f1() {
	var lmap map[string]bool
	fmt.Printf("%d\n", len(gmap) + 1) // 1
	fmt.Printf("%d\n", len(lmap) + 2) // 2
}

func f2() {
	var lmap map[int]int = map[int]int{
		4:3,
		5:9,
		6:25,
	}
	fmt.Printf("%d\n", len(lmap)) // 3
	for i := range lmap {
		fmt.Printf("%d\n", i) // 4,5,6
	}
}

func f3() {
	var lmap map[string]int
	value := lmap["hello"]
	fmt.Printf("%d\n", value + 7) // 7
	//lmap["x"] = 1
	value = lmap["hello"]
	fmt.Printf("%d\n", value + 8) // 8
}

func main() {
	f1()
	f2()
	f3()
}
