package main

import (
	"os"
	"strings"
)

var debugAst = false
var debugToken = false
var debugParser = false
var debugMode = false
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

		if strings.HasSuffix(opt, ".go") {
			sourceFile = opt
		} else if opt == "-" {
			sourceFile = "/dev/stdin"
		}
	}
}

func main() {

	if len(os.Args) > 1 {
		parseOpts(os.Args[1:len(os.Args)])
	}

	packageblockscope := newScope(nil)

	// parse
	p := &parser{}
	astFile := p.parseSourceFile(sourceFile, packageblockscope)

	if debugAst {
		astFile.dump()
	}

	p.resolve()

	if debugAst {
		astFile.dump()
	}

	// generate
	generate(astFile)
}
