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
	for i, sourceFile := range sourceFiles {
		p := &parser{}
		p.namedTypes = make(map[identifier]methods)
		astFile := p.parseSourceFile(sourceFile, packageblockscope)
		if i == 0 {
			p.parseInternalCode(internalCode, astFile)
		}

		if debugAst {
			astFile.dump()
		}
		debugf("methods=%v", p.namedTypes)

		p.resolve()

		if debugAst {
			astFile.dump()
		}

		if parseOnly {
			return
		}
		astFiles = append(astFiles, astFile)
	}

	generateFiles(astFiles)
}
