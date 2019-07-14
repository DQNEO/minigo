package main

import (
	"os"
	"strings"
)

var GENERATION int = 1

var debugMode = false // execute debugf() or not
var debugToken = false

var debugAst = false
var debugParser = false
var tokenizeOnly = false
var parseOnly = false
var resolveOnly = false
var emitPosition = false

func printVersion() {
	fmtPrintln(S("minigo 0.1.0"))
	fmtPrintln(S("Copyright (C) 2019 @DQNEO"))
}

func parseOpts(args []string) []bytes {
	var r []bytes

	for _, _opt := range args {
		opt := bytes(_opt)
		if eq(opt,S("--version")) {
			printVersion()
			return nil
		}
		if eq(opt, S("-t")) {
			debugToken = true
		}
		if eq(opt, S("-a")) {
			debugAst = true
		}
		if eq(opt, S("-p")) {
			debugParser = true
		}
		if eq(opt, S("--position")) {
			emitPosition = true
		}
		if eq(opt, S("-d")) {
			debugMode = true
		}
		if eq(opt, S("--tokenize-only")) {
			tokenizeOnly = true
		}
		if eq(opt, S("--parse-only")) {
			parseOnly = true
		}
		if eq(opt, S("--resolve-only")) {
			resolveOnly = true
		}
		if strings.HasSuffix(string(opt), ".go") {
			r = append(r, opt)
		} else if eq(opt, S("-")) {
			return []bytes{S("/dev/stdin")}
		}
	}

	return r
}

func main() {
	// parsing arguments
	var sourceFiles []bytes
	osArgs := os.Args
	assert(len(osArgs) > 0, nil, S("os.Args should not be empty"))
	if len(os.Args) > 1 {
		sourceFiles = parseOpts(osArgs[1:len(osArgs)])
	}

	if len(sourceFiles) == 0 {
		fmtPrintln(S("No input files."))
		return
	}

	if tokenizeOnly {
		dumpTokenForFiles(sourceFiles)
		return
	}

	// setup the universe scope
	universe := newUniverse()

	u := compileUniverse(universe)
	r := compileRuntime(universe)

	imported := parseImports(sourceFiles)

	symbolTable = &SymbolTable{}

	var allScopes map[identifier]*Scope
	allScopes = map[identifier]*Scope{}
	symbolTable.allScopes = allScopes
	libs := compileStdLibs(universe, imported)

	m := compileFiles(universe, sourceFiles)
	if m == nil {
		return
	}

	program := build(u, r, libs, m)
	program.emit()
}
