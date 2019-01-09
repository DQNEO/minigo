package main

import "fmt"
import "os"

func debugf(format string, v ...interface{}) {
	if !debugMode {
		return
	}
	spaces := "> "
	for i := 0; i < debugNest; i++ {
		spaces += "  "
	}

	fmt.Fprintf(os.Stderr, spaces+format+"\n", v...)
}

func (tok *Token) errorf(format string, v ...interface{}) {
	errorf(tok.String()+": "+format, v...)
}

func dumpToken(tok *Token) {
	debugf(fmt.Sprintf("tok: type=%-8s, sval=\"%s\"", tok.typ, tok.sval))
}

var debugNest int

func (a *PackageClause) dump() {
	debugf("package %s", a.name)
}

func (a *DeclFunc) dump() {
	debugf("funcdef %s", a.fname)
	debugNest++
	//for _, stmt := range a.body.stmts {
	//stmt.dump()
	//}
	debugNest--
}

func (ast *StmtAssignment) dump() {
	debugf("assign")
	debugNest++
	for _, e := range ast.lefts {
		e.dump()
	}
	for _, e := range ast.rights {
		e.dump()
	}
	debugNest--
}

func (a *DeclVar) dump() {
	if a.initval == nil {
		debugf("decl var %s %s",
			a.variable.varname, a.variable.gtype)
	} else {
		debugf("decl var")
		debugNest++
		a.variable.dump()
		a.initval.dump()
		debugNest--
	}
}

func (a *DeclConst) dump() {
	debugf("decl consts %v", a.consts)
}

func (a *DeclType) dump() {
	debugf("decl type def %v gtype(%v)",
		a.name, a.gtype)
}

func (stmt *StmtIf) dump() {
	debugf("if")
	debugNest++
	stmt.cond.dump()
	//stmt.then.dump()
	//stmt.els.dump()
	debugNest--
}

/*
func (s *StmtSatementList) dump() {
	for _, stmt := range s.stmts {
		stmt.dump()
	}
}
*/

func (a *SourceFile) dump() {
	debugf("==== Dump AstExpr Start ===")
	a.pkg.dump()
	for _, imprt := range a.imports {
		debugf("import \"%v\"", imprt.specs)
	}
	for _, decl := range a.decls {
		if decl.funcdecl != nil {
			decl.funcdecl.dump()
		} else if decl.typedecl != nil {
			decl.typedecl.dump()
		} else if decl.vardecl != nil {
			decl.vardecl.dump()
		} else if decl.constdecl != nil {
			decl.constdecl.dump()
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

func (ast *ExprMethodcall) dump() {
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

func (e ExprArrayLiteral) dump() {
	debugNest++
	for _, v := range e.values {
		v.dump()
	}
	debugNest--
}

func (ast *ExprNumberLiteral) dump() {
	debugf("int %d", ast.val)
}

func (ast *ExprStringLiteral) dump() {
	debugf("\"%s\"", ast.val)
}

func (a *Relation) dump() {
	assert(a.expr != nil, nil, "ident.expr is set for " + string(a.name))
	a.expr.dump()
}

func (ast *ExprBinop) dump() {
	debugf("binop %s", ast.op)
	debugNest++
	ast.left.dump()
	ast.right.dump()
	debugNest--
}

func (ast *ExprUop) dump() {
	debugf("unop %s", ast.op)
	debugNest++
	ast.operand.dump()
	debugNest--
}

func (a *ExprStructField) dump() {
	debugf("%s.%s", a.strct, a.fieldname)
}

func (stmt *StmtSwitch) emit() {
	panic("implement me")
}

func (stmt *ExprCaseClause) dump() {
	//stmt.exprs.dump()
	//stmt.compound.dump()
}

func (stmt *StmtSwitch) dump() {
	stmt.cond.dump()
	for _, c := range stmt.cases {
		c.dump()
	}
	//stmt.dflt.dump()
}

func (e *ExprNilLiteral) dump() {
	debugf("nil")
}

func (f *ExprFuncRef) dump() {
	f.funcdef.dump()
}

func (e *ExprSliced) dump() {
	errorf("TBD")
}

func (e *ExprArrayIndex) dump() {
	errorf("TBD")

}

func (e *ExprTypeAssertion) dump() {
	panic("implement me")
}

func (e *ExprVaArg) dump() {
	panic("implement me")
}

func (e *ExprConversion) dump() {
	panic("implement me")
}

func (e *ExprStructLiteral) dump() {
	debugf("%s{", e.strctname.name)
	for _, field := range e.fields {
		debugf("  %v:%v", field.key, field.value)
	}
	debugf("}")
}

func (e *ExprTypeSwitchGuard) dump() {
	panic("implement me")
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

func assert(cond bool, tok *Token, msg string) {
	if !cond {
		panic(fmt.Sprintf("assertion failed: %s %s", msg, tok))
	}
}

func assertNotNil(cond bool, tok *Token) {
	assert(cond, tok, "should not be nil")
}
