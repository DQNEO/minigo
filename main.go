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
	println("minigo 0.1.0")
	println("Copyright (C) 2019 @DQNEO")
}

func parseOpts(args []gostring) []gostring {
	var r []gostring

	for _, opt := range args {
		if eq(opt,"--version") {
			printVersion()
			return nil
		}
		if eq(opt, "-t") {
			debugToken = true
		}
		if eq(opt, "-a") {
			debugAst = true
		}
		if eq(opt, "-p") {
			debugParser = true
		}
		if eq(opt, "--position") {
			emitPosition = true
		}
		if eq(opt, "-d") {
			debugMode = true
		}
		if eq(opt, "--tokenize-only") {
			tokenizeOnly = true
		}
		if eq(opt, "--parse-only") {
			parseOnly = true
		}
		if eq(opt, "--resolve-only") {
			resolveOnly = true
		}
		if strings.HasSuffix(string(opt), ".go") {
			r = append(r, gostring(opt))
		} else if eq(opt, "-") {
			return []gostring{gostring("/dev/stdin")}
		}
	}

	return r
}

func main() {
	gobinops = convertCstringsToGostrings(binops)
	gokeywords = convertCstringsToGostrings(keywords)

	// parsing arguments
	var sourceFiles []gostring
	osArgs := convertCstringsToGostrings(os.Args)
	assert(len(osArgs) > 0, nil, "os.Args should not be empty")
	if len(os.Args) > 1 {
		sourceFiles = parseOpts(osArgs[1:len(osArgs)])
	}

	if len(sourceFiles) == 0 {
		println("No input files.")
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
