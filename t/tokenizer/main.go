package main

import (
	"fmt"
)

var GENERATION int = 2

var debugMode = true
var debugToken = false

func f3() {
	path := "t/min/min.go"
	bs := NewByteStreamFromFile(path)

	var c byte
	c, _ = bs.get()
	tn := &Tokenizer{
		bs: bs,
	}
	ident := tn.readIdentifier(c)
	fmt.Printf("%s\n", ident)
}

func f4() {
	path := "t/min/min.go"
	bs := NewByteStreamFromFile(path)
	tokens := Tokenize(bs)
	fmt.Printf("%d\n", len(tokens)) // 26
	fmt.Printf("----------\n")
	for _, tok := range tokens {
		fmt.Printf("%s:%s\n", string(tok.typ), tok.sval)
	}
}

func f5() {
	debugToken = false
	path := "t/data/string.txt"
	bs := NewByteStreamFromFile(path)

	tokens := Tokenize(bs)
	tok := tokens[0]
	fmt.Printf("----------\n")
	fmt.Printf("[%s]\n", tok.sval)
}

func main() {
	f3()
	f4()
	f5()
}
