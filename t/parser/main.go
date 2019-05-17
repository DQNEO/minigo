package main

import "os"

var GENERATION int = 2

var debugMode = true
var debugToken = false
var debugParser = false

func f1() {
	os.Stderr = os.Stdout
	path := "t2/min.go"
	bs := NewByteStreamFromFile(path)
	// parser to parse imported
	p := &parser{}
	p.scopes = map[identifier]*scope{}

	universe := newUniverse()
	astFile := p.parseSourceFile(bs, universe, false)
	debugNest = 0
	astFile.dump()
}

func main() {
	f1()
}
