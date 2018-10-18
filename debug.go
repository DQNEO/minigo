package main

import "fmt"
import "os"

func debugf(format string, v... interface{}) {
	if !debugMode {
		return
	}
	spaces := ""
	for i:=0;i< debugNest;i++ {
		spaces += "  "
	}

	fmt.Fprintf(os.Stderr, spaces + format + "\n", v...)
}

func debugPrintV(v interface{}) {
	debugf("%v", v)
}

func debugPrintVar(name string, v interface{}) {
	debugf("%s = %v", name, v)
}

func (tok *Token) errorf(format string, v... interface{}) {
	errorf(tok.String() + ": " + format, v...)
}

func dumpToken(tok *Token) {
	debugf(fmt.Sprintf("tok: type=%-8s, sval=\"%s\"", tok.typ, tok.sval))
}

var debugNest int

func (a *AstPackageClause) dump() {
	debugf("package %s", a.name)
}

func (a *AstFuncDecl) dump() {
	debugf("funcdef %s", a.fname)
	debugNest++
	for _, stmt := range a.body.stmts {
		stmt.dump()
	}
	debugNest--
}

func (ast *AstAssignment) dump() {
	debugf("assign")
	debugNest++
	ast.left.dump()
	ast.right.dump()
	debugNest--
}

func (a *AstVarDecl) dump() {
	if a.initval == nil {
		debugf("var %s", a.variable.varname)
	} else {
		debugf("var %s =", a.variable.varname)
		a.initval.dump()
	}
}

func (a *AstStmt) dump() {
	if a.declvar != nil {
		a.declvar.dump()
	} else if a.assignment != nil {
		a.assignment.dump()
	} else if a.expr != nil {
		a.expr.dump()
	}
}
func (a *AstSourceFile) dump() {
	debugf("==== Dump AstExpr Start ===")
	a.pkg.dump()
	for _, imprt := range a.imports {
		debugf("import \"%v\"", imprt.paths)
	}
	for _, decl := range a.decls {
		if decl.funcdecl != nil {
			decl.funcdecl.dump()
		} else if decl.vardecl != nil {
			decl.vardecl.dump()
		}
	}
	debugf("==== Dump AstExpr End ===")
}

func (ast *ExprFuncall) dump() {
	debugf(string(ast.fname))
	debugNest++
	for _, arg := range ast.args {
		arg.dump()
	}
	debugNest--
}

func (ast *ExprVariable) dump() {
	debugf("var %s", ast.varname)
}

func (ast *ExprConstVariable) dump() {
	debugf("var %s", ast.name)
}

func (ast *ExprNumberLiteral) dump() {
	debugf("int %d", ast.val)
}

func (ast *ExprStringLiteral) dump() {
	debugf("\"%s\"", ast.val)
}

func (ast *ExprBinop) dump() {
	debugf("binop %s", ast.op)
	debugNest++
	ast.left.dump()
	ast.right.dump()
	debugNest--
}

func errorf(format string, v ...interface{}) {
	/*
		currentTokenIndex := ts.index - 1
		fmt.Printf("%v %v %v\n",
			ts.getToken(currentTokenIndex-2), ts.getToken(currentTokenIndex-1), ts.getToken(currentTokenIndex))
	*/
	var s string
	//s += bs.location() + ": "
	s += fmt.Sprintf(format, v...)
	panic(s)
}

func assert(cond bool, msg string) {
	if !cond {
		panic(fmt.Sprintf("assertion failed: %s", msg))
	}
}
