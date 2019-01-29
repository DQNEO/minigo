package main

import "fmt"

var gmap map[string]bool

//var debug [6]int // r10, r11,...


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

	for _,v := range lmap {
		fmt.Printf("%d\n", v) // 7,8,9
	}
}

func f3() {
	var lmap map[int]int = map[int]int{
		10:  11,
		12: 13,
	}

	lmap[14] = 15
	lmap[16] = 17
	for i, v := range lmap {
		fmt.Printf("%d\n", i)
		fmt.Printf("%d\n", v)
	}
}

func f4() {
	var lmap map[int]int = map[int]int{
		7:  17,
		9:  10,
		11: 12,
		0:  18,
	}

	fmt.Printf("%d\n", lmap[0]) // 18

	fmt.Printf("%d\n", lmap[999]+19) // 19
	lmap[9] = 21
	fmt.Printf("%d\n", len(lmap)+16) // 20
	fmt.Printf("%d\n", lmap[9]) // 21

	lmap[2] = 23
	fmt.Printf("%d\n", len(lmap)+17) // 22
	fmt.Printf("%d\n", lmap[2]) // 23

	var lmap2 map[int]int = map[int]int{
		0:  1,
		1:  1,
		2: 1,
		3:  1,
	}

	fmt.Printf("%d\n", lmap[7] + 7) // 24
	fmt.Printf("%d\n", lmap2[0] + 24) // 25
}



func main() {
	f1()
	f2()
	f3()
	f4()
}
