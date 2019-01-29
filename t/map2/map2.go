package main

import "fmt"

func f5() {
	var lmap map[int]string = map[int]string {
		27: "twenty seven",
		26: "twenty six",
	}

	fmt.Printf("%s\n", lmap[26])
	fmt.Printf("%s\n", lmap[27])

	lmap[1] = "one"
	fmt.Printf("%s\n", lmap[1])
}


func main() {
	f5()
}
