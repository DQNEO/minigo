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
		if strings.HasSuffix(opt, ".go") {
			r = append(r, opt)
		} else if opt == "-" {
			return []string{"/dev/stdin"}
		}
	}

	return r
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
		files: []*AstSourceFile{asf},
	}
}

func main() {
	var sourceFiles []string

	if len(os.Args) > 1 {
		sourceFiles = parseOpts(os.Args[1:len(os.Args)])
	}

	packageblockscope := newScope(nil)

	// parse
	p := &parser{}
	p.namedTypes = make(map[identifier]methods)
	p.scopes = make(map[identifier]*scope)

	var bs *ByteStream
	var astFiles []*AstSourceFile

	universe := newScope(nil)

	bs = NewByteStreamFromString("builtin.memory", builtinCode)
	astFileBuiltin := p.parseSourceFile(bs, universe)
	astFiles = append(astFiles, astFileBuiltin)

	setPredeclaredIdentifiers(universe)

	// add std packages
	var stdpkgs []*stdpkg

	for _, pkgsrc := range pkgsources {
		pkg := parseStdPkg(p, universe, pkgsrc.name, pkgsrc.code)
		stdpkgs = append(stdpkgs, pkg)
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

	ir := ast2ir(stdpkgs,  astFiles, p.stringLiterals)
	ir.emit()
}

type stdpkg struct {
	name identifier
	files []*AstSourceFile
}

func ast2ir(stdpkgs []*stdpkg, files []*AstSourceFile, stringLiterals []*ExprStringLiteral) *IrRoot {

	root := &IrRoot{
		stringLiterals:stringLiterals,
	}

	for _, pkg := range stdpkgs {
		for _, f := range pkg.files {
			for _, decl := range f.decls {
				if decl.vardecl != nil {
					root.vars = append(root.vars, decl.vardecl)
				} else if decl.funcdecl != nil {
					root.funcs = append(root.funcs, decl.funcdecl)
				}
			}
		}
	}

	for _, f := range files {
		for _, decl := range f.decls {
			if decl.vardecl != nil {
				root.vars = append(root.vars, decl.vardecl)
			} else if decl.funcdecl != nil {
				root.funcs = append(root.funcs, decl.funcdecl)
			}
		}
	}
	return root
}
