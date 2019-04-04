package main

import "fmt"

var GENERATION int = 2

var debugMode = true
var debugToken = true
var debugParser = true
var allScopes map[identifier]*scope

func f1() {
	path := "t/hello/hello.go"
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
	p.importedNames = map[identifier]bool{}
	packageClause := p.parsePackageClause()
	fmt.Printf("%s\n", packageClause.name)
	importDecls := p.parseImportDecls()
	// regsiter imported names
	for _, importdecl := range importDecls {
		for _, spec := range importdecl.specs {
			var pkgName identifier
			pkgName = getBaseNameFromImport(spec.path)
			p.importedNames[pkgName] = true
			fmt.Printf("%s\n", pkgName)
		}
	}

}

func main() {
	f1()
}
