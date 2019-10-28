package main

import (
	"./stdlib/fmt"
	"os"
	"./stdlib/strings"
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

var _printVersion = false
func printVersion() {
	fmt.Println("minigo 0.3.0")
	fmt.Println("Copyright (C) 2019 @DQNEO")
	_printVersion = true
}

func parseOpts(args []string) []string {
	var r []string

	for _, opt := range args {
		if opt == "--version" {
			printVersion()
			return nil
		}
		if opt == "-t" {
			debugToken = true
		}
		if opt == "-a" {
			debugAst = true
		}
		if opt == "-p" {
			debugParser = true
		}
		if opt == "--position" {
			emitPosition = true
		}
		if opt == "-d" {
			debugMode = true
		}
		if opt == "--tokenize-only" {
			tokenizeOnly = true
		}
		if opt == "--parse-only" {
			parseOnly = true
		}
		if opt == "--resolve-only" {
			resolveOnly = true
		}
		if strings.HasSuffix(opt, ".go") {
			r = append(r, opt)
		} else if opt == "-" {
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
		if _printVersion {
			return
		}
		fmt.Println("No input files.")
		return
	}

	if tokenizeOnly {
		dumpTokenForFiles(sourceFiles)
		return
	}

	// setup the universe scope
	universe := newUniverse()

	pkgUniverse := compileUniverse(universe)
	pkgIRuntime := compileRuntime(universe)


	symbolTable = &SymbolTable{}
	symbolTable.allScopes = map[identifier]*Scope{}

	directDependencies := parseImports(sourceFiles)
	libs := compileStdLibs(universe, directDependencies)

	pkgMain := compileFiles(universe, sourceFiles)
	if pkgMain == nil {
		return
	}

	program := build(pkgUniverse, pkgIRuntime, libs, pkgMain)
	program.emit()
}
