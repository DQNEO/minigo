package main

import (
	"os"
	"strings"
)

var debugAst = false
var debugToken = false
var debugParser = false
var debugMode = false
var parseOnly = false

func parseOpts(args []string) []string {
	var r []string

	for _, opt := range args {
		if opt == "-t" {
			debugToken = true
		}
		if opt == "-a" {
			debugAst = true
		}
		if opt == "-p" {
			debugParser = true
		}
		if opt == "-d" {
			debugMode = true
		}
		if opt == "--parse-only" {
			parseOnly = true
		}
		if strings.HasSuffix(opt, ".go") {
			r = append(r, opt)
		} else if opt == "-" {
			return []string{"/dev/stdin"}
		}
	}

	return r
}

var internalCode  = `
package main

const MiniGo int = 1

// println should be a "Predeclared identifiers"
// https://golang.org/ref/spec#Predeclared_identifiers
func println(s string) {
	puts(s)
}

func Printf(format string, param any) {
	printf(format, param)
}
`

func main() {
	var sourceFiles []string

	if len(os.Args) > 1 {
		sourceFiles = parseOpts(os.Args[1:len(os.Args)])
	}

	packageblockscope := newScope(nil)

	var astFiles []*AstSourceFile
	// parse

	p := &parser{}
	p.namedTypes = make(map[identifier]methods)

	bs := &ByteStream{
		filename:  "memory",
		source:    []byte(internalCode),
		nextIndex: 0,
		line:      1,
		column:    0,
	}
	astFile0 := p.parseSourceFile(bs, packageblockscope)
	for _, sourceFile := range sourceFiles {
		bs := NewByteStream(sourceFile)
		astFile := p.parseSourceFile(bs, packageblockscope)

		if debugAst {
			astFile.dump()
		}
		debugf("methods=%v", p.namedTypes)

		astFiles = append(astFiles, astFile)
	}

	p.resolve()
	if parseOnly {
		return
	}
	astFiles = append(astFiles, astFile0)

	// generate code
	emitDataSection()
	emitRetGlobals()
	generateFiles(astFiles)
}
