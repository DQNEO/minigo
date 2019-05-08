package main

import "fmt"

func f1() {
	var lmap map[int]int = map[int]int{
		1: 2,
		3: 4,
	}

	for i, v := range lmap {
		fmt.Printf("%d\n", i)
		fmt.Printf("%d\n", v)
	}

	fmt.Printf("%d\n", lmap[1]+3) // 5
	fmt.Printf("%d\n", lmap[3]+2) // 6

	lmap[7] = 8
	fmt.Printf("%d\n", lmap[4]+7) // 7
	fmt.Printf("%d\n", lmap[7])   // 8
}

func f2() {
	keyFoo := "15"
	var lmap map[string]string = map[string]string{
		keyFoo:   "10",
		"17": "11",
	}

	fmt.Printf("9%s\n", lmap["noexists"])
	fmt.Printf("%s\n", lmap["15"]) // 10
	fmt.Printf("%s\n", lmap["17"]) // 11


	fmt.Printf("%d\n", len(lmap) + 10) // 12

	lmap["19"] = "13"
	fmt.Printf("%s\n", lmap["19"]) // 13

	fmt.Printf("%d\n", len(lmap) + 11 ) // 14
	lmap["15"] = "16"
	lmap["17"] = "18"
	lmap["19"] = "20"
	for k, v := range lmap {
		fmt.Printf("%s\n%s\n", k,v) // 15,16,17,18,19,20
	}
}

func main() {
	f1()
	f2()
}
