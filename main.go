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
var sourceFile string

func parseOpts(args []string) {
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
			sourceFile = opt
		} else if opt == "-" {
			sourceFile = "/dev/stdin"
		}
	}
}

var internalCode  = `
const MiniGo int = 1
`

func main() {

	if len(os.Args) > 1 {
		parseOpts(os.Args[1:len(os.Args)])
	}

	packageblockscope := newScope(nil)

	// parse
	p := &parser{}
	p.namedTypes = make(map[identifier]methods)
	astFile := p.parseSourceFile(sourceFile, packageblockscope)

	if debugAst {
		astFile.dump()
	}
	debugf("methods=%v", p.namedTypes)

	p.parseInternalCode(internalCode, astFile)

	p.resolve()

	if debugAst {
		astFile.dump()
	}

	if parseOnly {
		return
	}
	// generate
	generate(astFile)
}
