package main

import "fmt"

var GENERATION int = 2

var debugMode = true
var debugToken = true
var debugParser = true
var importOS = true
var gp *parser

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

	p.tokenStream = NewTokenStream(bs)
	p.packageBlockScope = nil
	p.currentScope = nil
	p.importedNames = nil
	packageClause := p.parsePackageClause()
	fmt.Printf("%s\n", packageClause.name)
}

func main() {
	f1()
}
