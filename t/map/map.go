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
		4:7,
		5:8,
		6:9,
	}
	fmt.Printf("%d\n", len(lmap)) // 3
	for i := range lmap {
		fmt.Printf("%d\n", i) // 4,5,6
	}
}

func f3() {
	var lmap map[int]int = map[int]int{
		7:8,
		9:10,
		11:12,
	}

	lmap[13] = 14
	lmap[15] = 16
	for i,v := range lmap {
		fmt.Printf("%d\n", i)
		fmt.Printf("%d\n", v)
	}

}

/*
func f3() {
	var lmap map[string]int
	value := lmap["hello"]
	fmt.Printf("%d\n", value + 13)
	//lmap["x"] = 1
	value = lmap["hello"]
	fmt.Printf("%d\n", value + 14)
}
*/

func main() {
	f1()
	f2()
	f3()
	//f3()
}
