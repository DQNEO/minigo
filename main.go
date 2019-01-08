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

var fmtCode  = `
package fmt

func Printf(format string, param interface{}) {
	printf(format, param)
}
`
var builtinCode  = `
package main

const MiniGo int = 1

// println should be a "Predeclared identifiers"
// https://golang.org/ref/spec#Predeclared_identifiers
func println(s string) {
	puts(s)
}
`

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

	bs = NewByteStreamFromString("builtin.memory", builtinCode)
	astFileBuiltin := p.parseSourceFile(bs, packageblockscope)
	astFiles = append(astFiles, astFileBuiltin)

	universe := newScope(nil)
	setPredeclaredIdentifiers(universe)

	bs = NewByteStreamFromString("fmt.memory", fmtCode)
	fmtScope := newScope(nil)
	p.scopes["fmt"] = fmtScope
	astFileFmt := p.parseSourceFile(bs, fmtScope)
	p.resolve(universe)
	astFiles = append(astFiles, astFileFmt)

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

	ir := ast2ir(astFiles, p.stringLiterals)
	ir.emit()
}

func ast2ir(files []*AstSourceFile, stringLiterals []*ExprStringLiteral) *IrRoot {

	root := &IrRoot{
		stringLiterals:stringLiterals,
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
