package main

import "fmt"

var GENERATION int = 2

var debugMode = true
var debugToken = true
var debugParser = true
var allScopes map[identifier]*scope

func f1() {
	path := "t/data/gen.go.txt"
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
	astFile := p.parseSourceFile(bs, nil, true)
	fmt.Printf("%s\n", astFile.packageClause.name)
	// regsiter imported names
	for _, importdecl := range astFile.importDecls {
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
