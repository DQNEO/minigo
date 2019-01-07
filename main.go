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

var fmtCode  = `
package main

func Printf(format string, param any) {
	printf(format, param)
}
`
var builtinCode  = `
package main

const MiniGo int = 1

// println should be a "Predeclared identifiers"
// https://golang.org/ref/spec#Predeclared_identifiers
func println(s string) {
	puts(s)
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

	var bs *ByteStream
	bs = &ByteStream{
		filename:  "memory",
		source:    []byte(builtinCode),
		nextIndex: 0,
		line:      1,
		column:    0,
	}
	astFile0 := p.parseSourceFile(bs, packageblockscope)

	bs = &ByteStream{
		filename:  "memory",
		source:    []byte(fmtCode),
		nextIndex: 0,
		line:      1,
		column:    0,
	}
	astFile1 := p.parseSourceFile(bs, packageblockscope)

	for _, sourceFile := range sourceFiles {
		bs := NewByteStream(sourceFile)
		astFile := p.parseSourceFile(bs, packageblockscope)

		if debugAst {
			astFile.dump()
		}

		astFiles = append(astFiles, astFile)
	}

	if parseOnly {
		return
	}
	p.resolve()
	astFiles = append(astFiles, astFile0)
	astFiles = append(astFiles, astFile1)

	// generate code
	emitDataSection()
	emitRetGlobals()
	generateFiles(astFiles)
}
