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

func printVersion() {
	println("minigo 0.1.0")
	println("Copyright (C) 2019 @DQNEO")
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
		sourceFiles = parseOpts(os.Args[1:len(os.Args)])
	}

	if len(sourceFiles) == 0 {
		println("No input files.")
		return
	}

	if tokenizeOnly {
		dumpTokenForFiles(sourceFiles)
		return
	}

	// parser starts
	p := &parser{}
	p.initPackage("")

	// setup universe scopes
	universe := newUniverse()
	var globalStringLiterals []*ExprStringLiteral
	var allDynamicTypes []*Gtype

	// inject builtin functions into the universe
	p.packageStringLiterals = nil
	internalUniverse := p.parseString("internal_universe.go", internalUniverseCode, universe, false)
	p.resolve(nil)
	inferTypes(p.packageUninferredGlobals, p.packageUninferredLocals)
	pUniverse := &AstPackage{
		name:"",
		files:[]*AstFile{internalUniverse},
		stringLiterals:p.packageStringLiterals,
		dynamicTypes:p.packageDynamicTypes,
	}

	for _, sl := range pUniverse.stringLiterals {
		globalStringLiterals = append(globalStringLiterals, sl)
	}

	for _, dt := range pUniverse.dynamicTypes {
		allDynamicTypes = append(allDynamicTypes, dt)
	}

	// inject runtime things into the universe
	p = &parser{}
	p.initPackage("")
	internalRuntime := p.parseString("internal_runtime.go", internalRuntimeCode, universe, false)
	p.resolve(nil)
	inferTypes(p.packageUninferredGlobals, p.packageUninferredLocals)
	for _, sl := range p.packageStringLiterals {
		globalStringLiterals = append(globalStringLiterals, sl)
	}
	for _, dt := range p.packageDynamicTypes {
		allDynamicTypes = append(allDynamicTypes, dt)
	}

	imported := parseImports(sourceFiles)
	allScopes = map[identifier]*Scope{}
	stdlibs := compileStdLibs(p, universe, imported)
	mainPkg := compileMainPkg(universe,sourceFiles)
	if mainPkg == nil {
		return
	}
	ir := makeIR(pUniverse, internalRuntime, stdlibs, mainPkg, globalStringLiterals, allDynamicTypes)
	ir.emit()
}
