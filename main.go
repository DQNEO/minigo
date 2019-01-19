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

	bs = NewByteStreamFromString("builtin.memory", builtinCode)
	astFileBuiltin := p.parseSourceFile(bs, universe)
	astFiles = append(astFiles, astFileBuiltin)

	setPredeclaredIdentifiers(universe)

	// add std packages
	var compiledPackages map[identifier]*stdpkg = make(map[identifier]*stdpkg)

	for pkgName, pkgCode := range pkgMap {
		pkg := parseStdPkg(p, universe, pkgName, pkgCode)
		compiledPackages[pkgName] = pkg
	}

	p.currentPackageName = "main"
	for _, sourceFile := range sourceFiles {
		bs := NewByteStreamFromFile(sourceFile)
		astFile := p.parseSourceFile(bs, packageblockscope)

		if debugAst {
			astFile.dump()
		}

		astFiles = append(astFiles, astFile)
	}

	if parseOnly {
		return
	}
	p.resolve(universe)
	if resolveOnly {
		return
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
