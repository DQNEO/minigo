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
var resolveOnly = false

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

type pkgsource struct {
	name identifier
	code string
}

func parseStdPkg(p *parser, universe *scope, pkgname identifier, code string) *stdpkg {
	filename := string(pkgname) + ".memory"
	bs = NewByteStreamFromString(filename, code)
	scp := newScope(nil)
	p.scopes[pkgname] = scp
	p.currentPackageName = pkgname
	asf := p.parseSourceFile(bs, scp)
	p.resolve(universe)
	return &stdpkg{
		name:  pkgname,
		files: []*SourceFile{asf},
	}
}

var gp *parser // for debug
var importOS bool

func main() {
	var sourceFiles []string

	if len(os.Args) > 1 {
		sourceFiles = parseOpts(os.Args[1:len(os.Args)])
	}

	packageblockscope := newScope(nil)

	// parse
	p := &parser{}
	gp = p
	p.methods = make(map[identifier]methods)
	p.scopes = make(map[identifier]*scope)

	var bs *ByteStream
	var astFiles []*SourceFile

	universe := newScope(nil)

	bs = NewByteStreamFromString("internalcode.memory", internalRuntimeCode)
	astFiles = append(astFiles, p.parseSourceFile(bs, universe))

	bs = NewByteStreamFromString("builtin.memory", builtinCode)
	astFiles = append(astFiles, p.parseSourceFile(bs, universe))

	setPredeclaredIdentifiers(universe)

	// add std packages
	var compiledPackages map[identifier]*stdpkg = map[identifier]*stdpkg{}

	for pkgName, pkgCode := range pkgMap {
		pkg := parseStdPkg(p, universe, pkgName, pkgCode)
		compiledPackages[pkgName] = pkg
	}

	p.currentPackageName = "main"
	for _, sourceFile := range sourceFiles {
		bs := NewByteStreamFromFile(sourceFile)
		astFiles = append(astFiles, p.parseSourceFile(bs, packageblockscope))
	}

	if parseOnly {
		return
	}
	p.resolve(universe)
	if resolveOnly {
		return
	}

	if debugAst {
		astFiles[len(astFiles)-1].dump()
	}

	var importedPackages []*stdpkg
	for importedName := range p.importedNames {
		compiledPkg, ok := compiledPackages[importedName]
		if !ok {
			errorf("package not found")
		}
		importedPackages = append(importedPackages, compiledPkg)
	}
	_, importOS = p.importedNames["os"]
	ir := ast2ir(importedPackages, astFiles, p.stringLiterals)

	var uniquedDynamicTypes map[string]int = map[string]int{}
	for _, gtype := range p.alldynamictypes {
		uniquedDynamicTypes[gtype.String()] = 0
	}
	ir.hashedTypes = uniquedDynamicTypes
	var methods map[int][]string = map[int][]string{} // typeId : []methods

	var typeId = 1 // start with 1 because we want to zero as error
	for _, concreteNamedType := range p.concreteNamedTypes {
		concreteNamedType.gtype.typeId = typeId
		typeId++
	}

	for _, funcdecl := range ir.funcs {
		if funcdecl.receiver != nil {
			gtype := funcdecl.receiver.getGtype()
			if gtype.typ == G_POINTER {
				gtype = gtype.origType
			}
			if gtype.relation == nil {
				errorf("no relation for %#v", funcdecl.receiver.getGtype())
			}
			typeId := gtype.relation.gtype.typeId
			methods[typeId] = append(methods[typeId], string(funcdecl.getUniqueName()))
		}
	}
	ir.methodTable = methods
	debugf("methods=%v", methods)
	ir.emit()
}

type stdpkg struct {
	name  identifier
	files []*SourceFile
}

func ast2ir(stdpkgs []*stdpkg, files []*SourceFile, stringLiterals []*ExprStringLiteral) *IrRoot {

	root := &IrRoot{
		stringLiterals: stringLiterals,
	}

	for _, pkg := range stdpkgs {
		for _, f := range pkg.files {
			for _, decl := range f.topLevelDecls {
				if decl.vardecl != nil {
					root.vars = append(root.vars, decl.vardecl)
				} else if decl.funcdecl != nil {
					root.funcs = append(root.funcs, decl.funcdecl)
				}
			}
		}
	}

	for _, f := range files {
		for _, decl := range f.topLevelDecls {
			if decl.vardecl != nil {
				root.vars = append(root.vars, decl.vardecl)
			} else if decl.funcdecl != nil {
				root.funcs = append(root.funcs, decl.funcdecl)
			}
		}
	}
	return root
}
