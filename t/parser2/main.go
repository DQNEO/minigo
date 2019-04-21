package main

var GENERATION int = 2

var debugMode = true
var debugToken = false
var debugParser = true
var allScopes map[identifier]*scope

func f1() {
	bs = NewByteStreamFromString("internalcode.memory", internalRuntimeCode)

	p := &parser{}
	p.methods = map[identifier]methods{}
	p.scopes = map[identifier]*scope{}
	universe := newScope(nil)

	astFile := p.parseSourceFile(bs, universe, false)
	debugNest = 0
	astFile.dump()
}

func main() {
	f1()
}
