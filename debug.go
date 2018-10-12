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

var nest = 0

func dumpAst(ast *Ast) {
	switch ast.typ {
	case "package":
		debugf("(package %s)", ast.pkgname)
	case "funcdef":
		debugf("funcdef %s", ast.fname)
		nest++
		for _, stmt := range ast.body.stmts {
			dumpAst(stmt)
		}
		nest--
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
	case "binop":
		debugf("binop %s", ast.op)
		nest++
		dumpAst(ast.left)
		dumpAst(ast.right)
		nest--
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
