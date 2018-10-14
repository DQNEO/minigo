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
		stmt.dump()
	}
	nest--
}

func (a *AstAssignment) dump() {
	debugf("assign")
	nest++
	a.left.dump()
	a.right.dump()
	nest--
}

func (a *AstDeclLocalVar) dump() {
	if a.initval == nil {
		debugf("var %s", a.localvar.varname)
	} else {
		debugf("var %s =", a.localvar.varname)
		a.initval.dump()
	}
}

func (a *AstStmt) dump() {
	if a.decllocalvar != nil {
		a.decllocalvar.dump()
	} else if a.assignment != nil {
		a.assignment.dump()
	} else if a.expr != nil {
		a.expr.dump()
	}
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


func (ast *AstExpr) dump() {
	switch ast.typ {
	case "funcall":
		debugf(ast.fname)
		nest++
		for _, arg := range ast.args {
			arg.dump()
		}
		nest--
	case "int":
		debugf("int %d", ast.ival)
	case "string":
		debugf("\"%s\"", ast.sval)
	case "assign":
		debugf("assign")
		nest++
		ast.left.dump()
		ast.right.dump()
		nest--
	case "binop":
		debugf("binop %s", ast.op)
		nest++
		ast.left.dump()
		ast.right.dump()
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
