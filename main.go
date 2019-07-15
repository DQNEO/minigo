package main

import (
	"fmt"
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
	fmt.Println("minigo 0.2.1")
	fmt.Println("Copyright (C) 2019 @DQNEO")
}

func parseOpts(args []string) []string {
	var r []string

	for _, _opt := range args {
		opt := _opt
		if opt == "--version" {
			printVersion()
			return nil
		}
		if opt ==  "-t" {
			debugToken = true
		}
		if opt ==  "-a" {
			debugAst = true
		}
		if opt ==  "-p" {
			debugParser = true
		}
		if opt ==  "--position" {
			emitPosition = true
		}
		if opt ==  "-d" {
			debugMode = true
		}
		if opt ==  "--tokenize-only" {
			tokenizeOnly = true
		}
		if opt ==  "--parse-only" {
			parseOnly = true
		}
		if opt ==  "--resolve-only" {
			resolveOnly = true
		}
		if strings.HasSuffix(string(opt), ".go") {
			r = append(r, string(opt))
		} else if opt ==  "-" {
			return []string{"/dev/stdin"}
		}
	}

	return r
}

func main() {
	// parsing arguments
	var sourceFiles []string
	assert(len(os.Args) > 0, nil, "os.Args should not be empty")
	if len(os.Args) > 1 {
		sourceFiles = parseOpts(os.Args[1:])
	}

	if len(sourceFiles) == 0 {
		fmt.Println("No input files.")
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
