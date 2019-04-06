package main

import "fmt"

var GENERATION int = 2

var debugMode = true
var debugToken = true
var debugParser = true
var allScopes map[identifier]*scope

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
	bs = &_bs

	// parser to parse imported
	p := &parser{}
	p.methods = map[identifier]methods{}
	p.scopes = map[identifier]*scope{}

	universe := newScope(nil)

	astFile := p.parseSourceFile(bs, universe, false)
	fmt.Printf("%s\n", astFile.packageClause.name)

}

func main() {
	f1()
}
