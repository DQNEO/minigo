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
var exit = false
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
			exit = true
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

	if exit {
		return
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
	// analyze imports of the given go files
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

	var importOS bool
	importOS = in_array("os", imported)
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
	bs = NewByteStreamFromString("internalcode.memory", internalRuntimeCode)
	astFiles = append(astFiles, p.parseSourceFile(bs, universe, false))
	p.resolve(nil)
	if debugAst {
		astFiles[0].dump()
	}

	// add std packages
	var compiledPackages map[identifier]*stdpkg = map[identifier]*stdpkg{}
	var uniqPackageNames []string
	// parse std packages

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
		compiledPackages[pkgName] = pkg
		if !in_array(string(pkgName), uniqPackageNames) {
			uniqPackageNames = append(uniqPackageNames, string(pkgName))
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

	debugf("resolve done")
	var importedPackages []*stdpkg
	for _, pkgName := range uniqPackageNames {
		compiledPkg := compiledPackages[identifier(pkgName)]
		importedPackages = append(importedPackages, compiledPkg)
	}
	ir := ast2ir(importedPackages, astFiles, p.stringLiterals)
	debugf("ir is created")
	ir.setDynamicTypes(p.allDynamicTypes)
	debugf("set uniquedDtypes")

	var typeId = 1 // start with 1 because we want to zero as error
	for _, concreteNamedType := range p.allNamedTypes {
		concreteNamedType.gtype.receiverTypeId = typeId
		//debugf("concreteNamedType: id=%d, name=%s", receiverTypeId, concreteNamedType.name)
		typeId++
	}
	debugf("set concreteNamedType")

	ir.composeMethodTable()
	ir.importOS = importOS

	ir.emit()
}

type stdpkg struct {
	name  identifier
	files []*SourceFile
}

func (ir *IrRoot) composeMethodTable() {
	var methodTable map[int][]string = map[int][]string{} // receiverTypeId : []methodTable
	for _, funcdecl := range ir.funcs {
		if funcdecl.receiver != nil {
			//debugf("funcdecl:%v", funcdecl)
			gtype := funcdecl.receiver.getGtype()
			if gtype.kind == G_POINTER {
				gtype = gtype.origType
			}
			if gtype.relation == nil {
				errorf("no relation for %#v", funcdecl.receiver.getGtype())
			}
			typeId := gtype.relation.gtype.receiverTypeId
			symbol := funcdecl.getSymbol()
			methods := methodTable[typeId]
			methods = append(methods, symbol)
			methodTable[typeId] = methods
		}
	}
	debugf("set methodTable")

	ir.methodTable = methodTable
}

func ast2ir(stdpkgs []*stdpkg, files []*SourceFile, stringLiterals []*ExprStringLiteral) *IrRoot {

	root := &IrRoot{}

	var declvars []*DeclVar
	for _, pkg := range stdpkgs {
		for _, f := range pkg.files {
			for _, decl := range f.topLevelDecls {
				if decl.vardecl != nil {
					declvars = append(declvars, decl.vardecl)
				} else if decl.funcdecl != nil {
					root.funcs = append(root.funcs, decl.funcdecl)
				}
			}
		}
	}

	for _, f := range files {
		for _, decl := range f.topLevelDecls {
			if decl.vardecl != nil {
				declvars = append(declvars, decl.vardecl)
			} else if decl.funcdecl != nil {
				root.funcs = append(root.funcs, decl.funcdecl)
			}
		}
	}

	root.stringLiterals = stringLiterals // a dirtyworkaround
	root.vars = declvars
	return root
}

func (ir *IrRoot) setDynamicTypes(allDynamicTypes []*Gtype) {
	var uniquedDTypes []string = builtinTypesAsString
	for _, gtype := range allDynamicTypes {
		gs := gtype.String()
		if !in_array(gs, uniquedDTypes) {
			uniquedDTypes = append(uniquedDTypes, gs)
		}
	}

	ir.uniquedDTypes = uniquedDTypes
}
