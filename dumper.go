package main

func dumpTokenForFiles(sourceFiles []gostring) {
	for _, sourceFile := range sourceFiles {
		debugf(S("--- file:%s"), sourceFile)
		bs := NewByteStreamFromFile(sourceFile)
		NewTokenStream(bs)
	}
}

func (pkg *AstPackage) dump() {
	for _, f := range pkg.files {
		f.dump()
	}
}

func (a *PackageClause) dump() {
	debugf(S("package %s"), a.name)
}

func (a *DeclFunc) dump() {
	debugf(S("funcdef %s"), a.fname)
	debugNest++
	for _, stmt := range a.body.stmts {
		stmt.dump()
	}
	debugNest--
}

func (a *StmtShortVarDecl) dump() {
	debugf(S("shot var decl"))
	debugf(S("left"))
	debugNest++
	for _, left := range a.lefts {
		left.dump()
	}
	debugNest--
	debugf(S("="))
	debugf(S("right"))
	debugNest++
	for _, right := range a.rights {
		right.dump()
	}
	debugNest--
}
func (ast *StmtAssignment) dump() {
	debugf(S("assign"))
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
		debugf(S("decl var %s %s"),
			a.variable.varname, a.variable.gtype.String())
	} else {
		debugf(S("decl var"))
		debugNest++
		a.variable.dump()
		a.initval.dump()
		debugNest--
	}
}

func (a *DeclConst) dump() {
	debugf(S("decl consts"))
	debugNest++
	for _, cnst := range a.consts {
		debugf(S("const %s"), cnst.name)
		cnst.val.dump()
	}
	debugNest--
}

func (a *DeclType) dump() {
	debugf(S("decl type def %s %s"),
		a.name, a.gtype.String())
}

func (stmt *StmtIf) dump() {
	debugf(S("if"))
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

func (a *AstFile) dump() {
	debugf(S("=== AST File %s ==="), a.name)
	a.packageClause.dump()
	for _, imprt := range a.importDecls {
		for _, spec := range imprt.specs {
			debugf(S("import \"%s\""), spec.path)
		}
	}
	for _, decl := range a.DeclList {
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
	debugf(S("==="))
}

func (ast *ExprFuncallOrConversion) dump() {
	debugf(gostring(ast.fname))
	debugNest++
	for _, arg := range ast.args {
		arg.dump()
	}
	debugNest--
}

func (ast *ExprMethodcall) dump() {
	debugf(gostring(ast.fname))
	debugNest++
	for _, arg := range ast.args {
		arg.dump()
	}
	debugNest--
}

func (ast *ExprVariable) dump() {
	debugf(S("var %s T %s"), gostring(ast.varname), ast.getGtype().String())
}

func (ast *ExprConstVariable) dump() {
	debugf(S("var %s"), ast.name)
}

func (e *ExprArrayLiteral) dump() {
	debugNest++
	for _, v := range e.values {
		v.dump()
	}
	debugNest--
}

func (ast *ExprNumberLiteral) dump() {
	debugf(S("int %d"), ast.val)
}

func (ast *ExprStringLiteral) dump() {
	debugf(S("\"%s\""), ast.val)
}

func (a *Relation) dump() {
	assert(a != nil, nil, S("ident shoud not be nil "))
	//assert(a.expr != nil, nil, "ident.expr shoud not be nil for " + string(a.name))
	if a.expr == nil && a.gtype == nil {
		debugf(S("rel %s (UNRESOLVED)"), a.name)
		return
	}
	a.expr.dump()
}

func (ast *ExprBinop) dump() {
	debugf(S("binop %s"), ast.op)
	debugNest++
	ast.left.dump()
	ast.right.dump()
	debugNest--
}

func (ast *ExprUop) dump() {
	debugf(S("unop %s"), ast.op)
	debugNest++
	ast.operand.dump()
	debugNest--
}

func (a *ExprStructField) dump() {
	a.strct.dump()
	debugf(S("  .%s"), a.fieldname)
}

func (stmt *ExprCaseClause) dump() {
	debugf(S("case"))
	debugNest++
	for _, expr := range stmt.exprs {
		expr.dump()
	}
	for _, gtype := range stmt.gtypes {
		debugf(S("%s"), gtype.String())
	}
	stmt.compound.dump()
	debugNest--
}

func (stmt *StmtSwitch) dump() {
	debugf(S("switch"))
	if stmt.cond != nil {
		stmt.cond.dump()
	}
	for _, _case := range stmt.cases {
		_case.dump()
	}
	if stmt.dflt != nil {
		debugf(S("default"))
		stmt.dflt.dump()
	}
}

func (e *ExprNilLiteral) dump() {
	debugf(S("nil"))
}

func (f *ExprFuncRef) dump() {
	f.funcdef.dump()
}

func (e *ExprSlice) dump() {
	debugf(S("ExprSlice:"))
	debugNest++
	e.collection.dump()
	if e.low != nil {
		e.low.dump()
	}
	if e.high != nil {
		e.high.dump()
	}
	if e.max != nil {
		e.max.dump()
	}
	debugNest--
}

func (e *ExprIndex) dump() {
	debugf(S("ExprIndex;"))
	debugNest++
	e.collection.dump()
	e.index.dump()
	debugNest--
}

func (e *ExprTypeAssertion) dump() {
	debugf(S("type assertion"))
	e.expr.dump()
	debugf(S(".(%s)"), e.gtype.String())
}

func (e *ExprVaArg) dump() {
	debugf(S("..."))
	e.expr.dump()
}

func (e *IrExprConversion) dump() {
	debugf(S("conversion"))
	debugNest++
	debugf(S("toType:%s"), e.toGtype.String())
	e.arg.dump()
	debugNest--
}

func (e *ExprStructLiteral) dump() {
	debugf(S("%s{"), e.strctname.name)
	for _, field := range e.fields {
		debugf(S("  field %s:"), field.key)
		debugNest++
		field.value.dump()
		debugNest--
	}
	debugf(S("}"))
}

func (e *ExprTypeSwitchGuard) dump() {
	debugf(S("switch"))
	e.expr.dump()
}

func (f *IrStmtForRangeList) dump() {
	debugf(S("for range list"))
	debugNest++
	f.block.dump()
	debugNest--
}

func (f *IrStmtRangeMap) dump() {
	debugf(S("for range map"))
	debugNest++
	f.block.dump()
	debugNest--
}

func (f *IrStmtClikeFor) dump() {
	debugf(S("for clause"))
	if f.cls.init != nil {
		f.cls.init.dump()
	}
	if f.cls.cond != nil {
		f.cls.cond.dump()
	}
	if f.cls.post != nil {
		f.cls.post.dump()
	}
	debugNest++
	f.block.dump()
	debugNest--
}

func (f *StmtFor) dump() {
	if f.rng != nil {
		debugf(S("for range"))
		debugNest++
		f.rng.indexvar.dump()
		if f.rng.valuevar != nil {
			f.rng.valuevar.dump()
		}
		debugf(S("range"))
		f.rng.rangeexpr.dump()
		debugNest--
	} else if f.cls != nil {
		debugf(S("for clause"))
		if f.cls.init != nil {
			f.cls.init.dump()
		}
		if f.cls.cond != nil {
			f.cls.cond.dump()
		}
		if f.cls.post != nil {
			f.cls.post.dump()
		}
	}
	debugNest++
	f.block.dump()
	debugNest--
}

func (e *ExprLen) dump() {
	TBI(e.token(), "")
}

func (e *ExprCap) dump() {
	TBI(e.token(), "")
}

func (e *ExprSliceLiteral) dump() {
	debugf(S("slice %s"), e.gtype.String())
	debugNest++
	for _, v := range e.values {
		v.dump()
	}
	debugNest--
}

func (stmt *StmtReturn) dump() {
	debugf(S("return"))
	debugNest++
	for _, e := range stmt.exprs {
		e.dump()
	}
	debugNest--
}

func (ast *StmtInc) dump() {
	debugf(S("++"))
	ast.operand.dump()
}

func (ast *StmtDec) dump() {
	debugf(S("--"))
	ast.operand.dump()
}

func (ast *StmtSatementList) dump() {
	for _, stmt := range ast.stmts {
		stmt.dump()
	}
}

func (ast *StmtContinue) dump() {
	debugf(S("continue"))
}

func (ast *StmtBreak) dump() {
	debugf(S("break"))
}

func (ast *StmtExpr) dump() {
	ast.expr.dump()
}

func (ast *StmtDefer) dump() {
	debugf(S("defer"))
	debugNest++
	ast.expr.dump()
	debugNest--
}

func (e *ExprMapLiteral) dump() {
	debugf(S("map literal T %s"), e.gtype.String())
	debugNest++
	for _, element := range e.elements {
		debugf(S("element key:"))
		debugNest++
		element.key.dump()
		debugNest--
		debugf(S("element value:"))
		debugNest++
		element.value.dump()
		debugNest--
	}
	debugNest--
}

func (e *IrExprConversionToInterface) dump() {
	panic("implement me: IrExprConversionToInterface")
}
