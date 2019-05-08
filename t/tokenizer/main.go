package main

import (
	"fmt"
)

var GENERATION int = 2

var debugMode = true
var debugToken = false

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

func f2() {
	var chars []byte
	chars = append(chars, '7')
	chars = append(chars, '8')
	fmt.Printf("%d\n", len(chars)+4) // 6
	fmt.Printf("%c\n", chars[0])     // 7
	fmt.Printf("%c\n", chars[1])     // 8
	fmt.Printf("9\n")                // 9

	chars[0] = '1'
	chars[1] = '0'
	fmt.Printf("%s\n", chars) // 10
}

func f3() {
	path := "t/min/min.go"
	s := readFile(path)
	_bs := ByteStream{
		filename:  path,
		source:    s,
		nextIndex: 0,
		line:      1,
		column:    0,
	}
	bs = &_bs

	var c byte
	c, _ = bs.get()
	ident := readIdentifier(c)
	if ident == "package" {
		fmt.Printf("11\n")
	} else {
		println("error")
	}
}

func f4() {
	path := "t/min/min.go"
	s := readFile(path)
	_bs := ByteStream{
		filename:  path,
		source:    s,
		nextIndex: 0,
		line:      1,
		column:    0,
	}
	bs = &_bs

	tokens := tokenize(bs)
	fmt.Printf("%d\n", len(tokens)-14) // 26 - 14 = 12
	fmt.Printf("----------\n")
	for _, tok := range tokens {
		fmt.Printf("%s:%s\n", string(tok.typ), tok.sval)
	}
}

func f5() {
	debugToken = false
	path := "t/data/string.txt"
	s := readFile(path)
	_bs := ByteStream{
		filename:  path,
		source:    s,
		nextIndex: 0,
		line:      1,
		column:    0,
	}
	bs = &_bs

	tokens := tokenize(bs)
	tok := tokens[0]
	fmt.Printf("----------\n")
	fmt.Printf("[%s]\n", tok.sval)
}

func main() {
	f1()
	f2()
	f3()
	f4()
	f5()
}
