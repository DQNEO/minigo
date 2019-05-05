package main

func (a *PackageClause) dump() {
	debugf("package %s", a.name)
}

func (a *DeclFunc) dump() {
	debugf("funcdef %s", a.fname)
	debugNest++
	for _, stmt := range a.body.stmts {
		stmt.dump()
	}
	debugNest--
}

func (a *StmtShortVarDecl) dump() {
	debugf("shot var decl")
	debugf("left")
	debugNest++
	for _, left := range a.lefts {
		left.dump()
	}
	debugNest--
	debugf("=")
	debugf("right")
	debugNest++
	for _, right := range a.rights {
		right.dump()
	}
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
			a.variable.varname, a.variable.gtype.String())
	} else {
		debugf("decl var")
		debugNest++
		a.variable.dump()
		a.initval.dump()
		debugNest--
	}
}

func (a *DeclConst) dump() {
	debugf("decl consts")
	debugNest++
	for _, cnst := range a.consts {
		debugf("const %s", cnst.name)
		cnst.val.dump()
	}
	debugNest--
}

func (a *DeclType) dump() {
	debugf("decl type def %s %s",
		a.name, a.gtype.String())
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
	debugf("==== AST DUMP START ===")
	a.packageClause.dump()
	for _, imprt := range a.importDecls {
		for _, spec := range imprt.specs {
			debugf("import \"%s\"", spec.path)
		}
	}
	for _, decl := range a.topLevelDecls {
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
	debugf("==== AST DUMP END ===")
}

func (ast *ExprFuncallOrConversion) dump() {
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
	debugf("var %s T %s", ast.varname, ast.gtype.String())
}

func (ast *ExprConstVariable) dump() {
	debugf("var %s", ast.name)
}

func (e *ExprArrayLiteral) dump() {
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
	assert(a != nil, nil, "ident shoud not be nil ")
	assert(a.expr != nil, nil, "ident.expr shoud not be nil for " + string(a.name))
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

func (e *ExprSlice) dump() {
	debugf("ExprIndex:")
	debugNest++
	e.collection.dump()
	e.low.dump()
	e.high.dump()
	e.max.dump()
	debugNest--
}

func (e *ExprIndex) dump() {
	debugf("ExprIndex;")
	debugNest++
	e.collection.dump()
	e.index.dump()
	debugNest--
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

func (f *StmtFor) dump() {
	if f.rng != nil {
		debugf("for range")
		debugNest++
		f.rng.indexvar.dump()
		if f.rng.valuevar != nil {
			f.rng.valuevar.dump()
		}
		debugf("range")
		f.rng.rangeexpr.dump()
		debugNest--
	} else if f.cls != nil {
		debugf("for clause")
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
	panic("implement me")
}

func (stmt *StmtReturn) dump() {
	debugf("return")
	debugNest++
	for _, e := range stmt.exprs {
		e.dump()
	}
	debugNest--
}

func (ast *StmtInc) dump() {
	debugf("++")
	ast.operand.dump()
}

func (ast *StmtDec) dump() {
	debugf("--")
	ast.operand.dump()
}

func (ast *StmtSatementList) dump() {
	for _, stmt := range ast.stmts {
		stmt.dump()
	}
}

func (ast *StmtContinue) dump() {
	panic("implement me")
}

func (ast *StmtBreak) dump() {
	panic("implement me")
}

func (ast *StmtExpr) dump() {
	ast.expr.dump()
}

func (ast *StmtDefer) dump() {
	panic("implement me")
}

func (e *ExprMapLiteral) dump() {
	panic("implement me")
}

func (e *ExprConversionToInterface) dump() {
	panic("implement me")
}
