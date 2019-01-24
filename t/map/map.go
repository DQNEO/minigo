package main

import "fmt"

var gmap map[string]bool

func f1() {
	var lmap map[string]bool
	fmt.Printf("%d\n", len(gmap) + 1) // 1
	fmt.Printf("%d\n", len(lmap) + 2) // 2
}

func f2() {
	var lmap map[string]int = map[string]int{
		"foo":1,
		"bar":2,
		"piyo":2,
	}
	fmt.Printf("%d\n", len(lmap)) // 3
}

func f3() {
	var lmap map[string]int
	value := lmap["hello"]
	fmt.Printf("%d\n", value + 4) // 4
	lmap["x"] = 1
	value = lmap["hello"]
	fmt.Printf("%d\n", value + 5) // 5
}

func main() {
	f1()
	f2()
	f3()
}
