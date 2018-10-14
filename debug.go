package main

import "fmt"
import "os"

func debugf(format string, v... interface{}) {
	debugPrint(fmt.Sprintf(format, v...))
}

func debugPrint(s string) {
	spaces := ""
	for i:=0;i<nest;i++ {
		spaces += "  "
	}

	fmt.Fprintf(os.Stderr, "|%s %s\n", spaces, s)
}

func debugPrintV(v interface{}) {
	debugPrint(fmt.Sprintf("%v", v))
}

func debugPrintVar(name string, v interface{}) {
	debugPrint(fmt.Sprintf("%s=%v", name, v))
}

func dumpToken(tok *Token) {
	debugPrint(fmt.Sprintf("tok: type=%-8s, sval=\"%s\"", tok.typ, tok.sval))
}

var nest int

func (a *AstPkgDecl) dump() {
	debugf("package %s", a.name)
}

func (a *AstFuncDef) dump() {
	debugf("funcdef %s", a.fname)
	nest++
	for _, stmt := range a.body.stmts {
		dumpAst(stmt.expr)
	}
	nest--
}

func (a *AstFile) dump() {
	debugPrint("==== Dump AstExpr Start ===")
	a.pkg.dump()
	for _, imprt := range a.imports {
		debugf("import \"%v\"", imprt.paths)
	}
	for _, f := range a.funcdefs {
		f.dump()
	}
	debugPrint("==== Dump AstExpr End ===")
}


func dumpAst(ast *AstExpr) {
	switch ast.typ {
	case "funcall":
		debugf(ast.fname)
		nest++
		for _, arg := range ast.args {
			dumpAst(arg)
		}
		nest--
	case "int":
		debugf("int %d", ast.ival)
	case "string":
		debugf("\"%s\"", ast.sval)
	case "assign":
		debugf("assign")
		nest++
		dumpAst(ast.left)
		dumpAst(ast.right)
		nest--
	case "binop":
		debugf("binop %s", ast.op)
		nest++
		dumpAst(ast.left)
		dumpAst(ast.right)
		nest--
	case "decl":
		debugf("decl")
	case "lvar":
		debugf("lvar")
	default:
		errorf("Unknown ast type:%v", ast.typ)
	}
}

func errorf(format string, v ...interface{}) {
	/*
		currentTokenIndex := ts.index - 1
		fmt.Printf("%v %v %v\n",
			ts.getToken(currentTokenIndex-2), ts.getToken(currentTokenIndex-1), ts.getToken(currentTokenIndex))
	*/
	s := fmt.Sprintf(format, v...)
	panic(s)
}

func assert(cond bool, msg string) {
	if !cond {
		panic(fmt.Sprintf("assertion failed: %s", msg))
	}
}
