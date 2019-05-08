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
	keyFoo := "keyfoo"
	var lmap map[string]string = map[string]string{
		keyFoo:   "valuefoo",
		"keybar": "valuebar",
	}

	fmt.Printf("9%s\n", lmap["noexists"])
	fmt.Printf("%s\n", lmap[keyFoo]) // "valuefoo"
	fmt.Printf("%s\n", lmap["keybar"]) // "valuebar"


	fmt.Printf("len=%d\n", len(lmap))

	lmap["keyadded"] = "valueadded"
	fmt.Printf("%s\n", lmap["keyadded"])

	fmt.Printf("len=%d\n", len(lmap))
	for k, v := range lmap {
		fmt.Printf("%s\n", k)
		fmt.Printf("%s\n", v)
	}
}

func main() {
	f1()
	f2()
}
