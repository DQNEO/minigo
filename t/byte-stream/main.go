package main

import "fmt"

func f1() {
	path := "t/min/min.go"
	s := readFile(path)
	_bs := ByteStream{
		filename:  path,
		source:    s,
		nextIndex: 0,
		line:      1,
		column:    0,
	}
	bs := &_bs
	len1 := len(bs.source)

	fmt.Printf("%d\n", len1-64) // 1
	var c byte
	c, _ = bs.get()
	fmt.Printf("%d\n", c-'p'+2)        // 2
	fmt.Printf("%d\n", bs.nextIndex+2) // 3
	c, _ = bs.get()
	fmt.Printf("%d\n", c-'a'+4)        // 4
	fmt.Printf("%d\n", bs.nextIndex+3) // 5
}

func main() {
	f1()
}
