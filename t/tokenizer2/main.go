package main

import "fmt"

var GENERATION int = 2

var debugMode = true
var debugToken = false

func f1() {
	filename := "t/data/gen.go.txt"
	bs := NewByteStreamFromFile(filename)
	tokens := Tokenize(bs)
	expectedLen := 17977
	if len(tokens) == expectedLen {
		fmt.Println("1")
	} else {
		panic("ERROR")
	}
	/*
		for _, tok := range tokens {
			fmt.Printf("%s:%s\n", string(tok.typ), tok.sval)
		}
	*/
}

func main() {
	f1()
}
