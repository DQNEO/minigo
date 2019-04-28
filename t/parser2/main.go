package main

var GENERATION int = 2

var debugMode = true
var debugToken = false
var debugParser = true
var allScopes map[identifier]*scope

func f1() {
	bs = NewByteStreamFromFile("t2/bootstrap.go")

	p := &parser{}
	p.methods = map[identifier]methods{}
	p.scopes = map[identifier]*scope{}
	universe := newUniverse()
	astFile := p.parseSourceFile(bs, universe, false)
	debugNest = 0
	astFile.dump()

	filename := "stdlib/fmt/fmt.go"
	bs = NewByteStreamFromFile(filename)

	// initialize a package
	p.methods = map[identifier]methods{}
	p.scopes["fmt"] = newScope(nil, "fmt")
	p.currentPackageName = "fmt"
	debugParser = true
	asf := p.parseSourceFile(bs, p.scopes["fmt"], false)
	p.resolve(universe)
	asf.dump()
}

func main() {
	f1()
}
