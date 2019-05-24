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
var slientForStdlib = true

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

// analyze imports of given go files
func parseImports(sourceFiles []string) []string {

	pForImport := &parser{}
	// "fmt" depends on "os. So inject it in advance.
	// Actually, dependency graph should be analyzed.
	var imported []string = []string{"os"}
	for _, sourceFile := range sourceFiles {
		bs := NewByteStreamFromFile(sourceFile)
		astFile := pForImport.parseSourceFile(bs, nil, true)
		for _, importDecl := range astFile.importDecls {
			for _, spec := range importDecl.specs {
				baseName := getBaseNameFromImport(spec.path)
				if !in_array(baseName, imported) {
					imported = append(imported, baseName)
				}
			}
		}
	}

	return imported
}

func parseStdPkg(p *parser, universe *Scope, pkgname identifier, code string) *stdpkg {
	filename := string(pkgname) + ".memory"
	bs := NewByteStreamFromString(filename, code)

	// initialize a package
	p.initPackage(pkgname)
	p.scopes[pkgname] = newScope(nil, string(pkgname))

	asf := p.parseSourceFile(bs, p.scopes[pkgname], false)

	p.resolve(universe)
	if debugAst {
		asf.dump()
	}
	return &stdpkg{
		name:  pkgname,
		files: []*SourceFile{asf},
	}
}

func main() {
	// parsing arguments
	var sourceFiles []string

	if len(os.Args) == 0 {
		println("ERROR: os.Args should not be empty")
		return
	}

	if len(os.Args) > 1 {
		sourceFiles = parseOpts(os.Args[1:len(os.Args)])
	}

	if len(sourceFiles) == 0 {
		println("No input files.")
		return
	}

	if tokenizeOnly {
		for _, sourceFile := range sourceFiles {
			debugf("--- file:%s", sourceFile)
			bs := NewByteStreamFromFile(sourceFile)
			NewTokenStream(bs)
		}
		return
	}

	imported := parseImports(sourceFiles)
	// parser starts
	p := &parser{}
	p.scopes = map[identifier]*Scope{}
	p.initPackage("")

	allScopes = p.scopes

	var bs *ByteStream
	var astFiles []*SourceFile

	var _debugAst bool
	var _debugParer bool
	if slientForStdlib {
		_debugAst = debugAst
		_debugParer = debugParser
		debugAst = false
		debugParser = false
	}

	// setup universe scopes
	universe := newUniverse()
	// inject runtime things into the universes
	bs = NewByteStreamFromString("internal_universe.go", internalUniverseCode)
	astFiles = append(astFiles, p.parseSourceFile(bs, universe, false))
	bs = NewByteStreamFromString("internal_runtime.go", internalRuntimeCode)
	astFiles = append(astFiles, p.parseSourceFile(bs, universe, false))
	p.resolve(nil)
	if debugAst {
		astFiles[0].dump()
	}

	// add std packages
	// parse std packages
	var csl *compiledStdlib = &compiledStdlib{
		compiledPackages: map[identifier]*stdpkg{},
		uniqImportedPackageNames:nil,
	}
	stdPkgs := makeStdLib()

	for _, spkgName := range imported {
		pkgName := identifier(spkgName)
		var pkgCode string
		var ok bool
		pkgCode, ok = stdPkgs[pkgName]
		if !ok {
			errorf("package '" + string(pkgName) + "' is not a standard library.")
		}
		pkg := parseStdPkg(p, universe, pkgName, pkgCode)
		csl.compiledPackages[pkgName] = pkg
		if !in_array(string(pkgName), csl.uniqImportedPackageNames) {
			csl.uniqImportedPackageNames = append(csl.uniqImportedPackageNames, string(pkgName))
		}
	}

	if slientForStdlib {
		debugAst = _debugAst
		debugParser = _debugParer
	}
	// initialize main package
	var pkgname identifier = "main"
	p.initPackage(pkgname)
	p.scopes[pkgname] = newScope(nil, string(pkgname))

	for _, sourceFile := range sourceFiles {
		bs := NewByteStreamFromFile(sourceFile)
		asf := p.parseSourceFile(bs, p.scopes[pkgname], false)
		astFiles = append(astFiles, asf)
	}

	if parseOnly {
		if debugAst {
			for _, af := range astFiles {
				af.dump()
			}
		}
		return
	}
	p.resolve(universe)
	if debugAst {
		for _, af := range astFiles {
			af.dump()
		}
	}

	if resolveOnly {
		return
	}

	setTypeIds(p.allNamedTypes)

	debugf("resolve done")
	ir := makeIR(csl , astFiles, p.stringLiterals, p.allDynamicTypes)
	ir.emit()
}

type compiledStdlib struct {
	compiledPackages map[identifier]*stdpkg
	uniqImportedPackageNames []string
}

type stdpkg struct {
	name  identifier
	files []*SourceFile
}
