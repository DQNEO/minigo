package main

func (stmt *StmtIf) emit() {
	emit("# if")
	if stmt.simplestmt != nil {
		stmt.simplestmt.emit()
	}
	stmt.cond.emit()
	emit("TEST_IT")
	if stmt.els != nil {
		labelElse := makeLabel()
		labelEndif := makeLabel()
		emit("je %s  # jump if 0", labelElse)
		emit("# then block")
		stmt.then.emit()
		emit("jmp %s # jump to endif", labelEndif)
		emit("# else block")
		emit("%s:", labelElse)
		stmt.els.emit()
		emit("# endif")
		emit("%s:", labelEndif)
	} else {
		// no else block
		labelEndif := makeLabel()
		emit("je %s  # jump if 0", labelEndif)
		emit("# then block")
		stmt.then.emit()
		emit("# endif")
		emit("%s:", labelEndif)
	}
}

func (stmt *StmtSwitch) emit() {

	emit("#")
	emit("# switch statement")
	labelEnd := makeLabel()
	var labels []string

	// switch (expr) {
	if stmt.cond != nil {
		emit("# the subject expression")
		stmt.cond.emit()
		emit("PUSH_8 # the subject value")
		emit("#")
	} else {
		// switch {
		emit("# no condition")
	}

	// case exp1,exp2,..:
	//     stmt1;
	//     stmt2;
	//     ...
	for i, caseClause := range stmt.cases {
		emit("# case %d", i)
		myCaseLabel := makeLabel()
		labels = append(labels, myCaseLabel)
		if stmt.cond == nil {
			for _, e := range caseClause.exprs {
				e.emit()
				emit("TEST_IT")
				emit("jne %s # jump if matches", myCaseLabel)
			}
		} else if stmt.isTypeSwitch {
			// compare type
			for _, gtype := range caseClause.gtypes {
				emit("# Duplicate the subject value in stack")
				emit("POP_8")
				emit("PUSH_8")
				emit("PUSH_8")

				if gtype.isNil() {
					emit("mov $0, %%rax # nil")
				} else {
					typeLabel := symbolTable.getTypeLabel(gtype)
					emit("LOAD_STRING_LITERAL .%s # type: %s", typeLabel, gtype.String())
				}
				emit("PUSH_8")
				emitStringsEqualFromStack(true)

				emit("TEST_IT")
				emit("jne %s # jump if matches", myCaseLabel)
			}
		} else {
			for _, e := range caseClause.exprs {
				emit("# Duplicate the subject value in stack")
				emit("POP_8")
				emit("PUSH_8")
				emit("PUSH_8")

				e.emit()
				emit("PUSH_8")
				if e.getGtype().isString() {
					emitStringsEqualFromStack(true)
				} else {
					emit("CMP_FROM_STACK sete")
				}

				emit("TEST_IT")
				emit("jne %s # jump if matches", myCaseLabel)
			}
		}
	}

	var defaultLabel string
	if stmt.dflt == nil {
		emit("jmp %s", labelEnd)
	} else {
		emit("# default")
		defaultLabel = makeLabel()
		emit("jmp %s", defaultLabel)
	}

	emit("POP_8 # destroy the subject value")
	emit("#")
	for i, caseClause := range stmt.cases {
		emit("# case stmts")
		emit("%s:", labels[i])
		caseClause.compound.emit()
		emit("jmp %s", labelEnd)
	}

	if stmt.dflt != nil {
		emit("%s:", defaultLabel)
		stmt.dflt.emit()
	}

	emit("%s: # end of switch", labelEnd)
}

func (f *StmtFor) emitRangeForList() {
	emitNewline()
	emit("# for range %s", f.rng.rangeexpr.getGtype().String())
	assertNotNil(f.rng.indexvar != nil, f.rng.tok)
	assert(f.rng.rangeexpr.getGtype().isArrayLike(), f.rng.tok, "rangeexpr should be G_ARRAY or G_SLICE, but got "+f.rng.rangeexpr.getGtype().String())

	labelBegin := makeLabel()
	f.labelEndBlock = makeLabel()
	f.labelEndLoop = makeLabel()

	// i = 0
	emit("# init index")
	initstmt := &StmtAssignment{
		lefts: []Expr{
			f.rng.indexvar,
		},
		rights: []Expr{
			&ExprNumberLiteral{
				val: 0,
			},
		},
	}
	initstmt.emit()

	emit("%s: # begin loop ", labelBegin)

	// i < len(list)
	condition := &ExprBinop{
		op:   "<",
		left: f.rng.indexvar, // i
		// @TODO
		// The range expression x is evaluated once before beginning the loop
		right: &ExprLen{
			arg: f.rng.rangeexpr, // len(expr)
		},
	}
	condition.emit()
	emit("TEST_IT")
	emit("je %s  # if false, go to loop end", f.labelEndLoop)

	// v = s[i]
	var assignVar *StmtAssignment
	if f.rng.valuevar != nil {
		assignVar = &StmtAssignment{
			lefts: []Expr{
				f.rng.valuevar,
			},
			rights: []Expr{
				&ExprIndex{
					collection: f.rng.rangeexpr,
					index:      f.rng.indexvar,
				},
			},
		}
		assignVar.emit()
	}

	f.block.emit()
	emit("%s: # end block", f.labelEndBlock)

	// break if i == len(list) - 1
	condition2 := &ExprBinop{
		op:   "==",
		left: f.rng.indexvar, // i
		// @TODO2
		// The range expression x is evaluated once before beginning the loop
		right: &ExprBinop{
			op: "-",
			left: &ExprLen{
				arg: f.rng.rangeexpr, // len(expr)
			},
			right: &ExprNumberLiteral{
				val: 1,
			},
		},
	}
	condition2.emit()
	emit("TEST_IT")
	emit("jne %s  # if this iteration is final, go to loop end", f.labelEndLoop)

	// i++
	indexIncr := &StmtInc{
		operand: f.rng.indexvar,
	}
	indexIncr.emit()

	emit("jmp %s", labelBegin)
	emit("%s: # end loop", f.labelEndLoop)
}

func (f *StmtFor) emitForClause() {
	assertNotNil(f.cls != nil, nil)
	labelBegin := makeLabel()
	f.labelEndBlock = makeLabel()
	f.labelEndLoop = makeLabel()

	if f.cls.init != nil {
		init := f.cls.init
		switch init.(type) {
		case *StmtAssignment:
			s := init.(*StmtAssignment)
			init = walkStmt(s)
		case *StmtShortVarDecl:
			s := init.(*StmtShortVarDecl)
			init = walkStmt(s)
		default:
			errorft(f.token(), "should not reach here")
		}
		//f.cls.init = walkStmt(f.cls.init)

		f.cls.init.emit()
	}
	emit("%s: # begin loop ", labelBegin)
	if f.cls.cond != nil {
		f.cls.cond.emit()
		emit("TEST_IT")
		emit("je %s  # jump if false", f.labelEndLoop)
	}
	f.block.emit()
	emit("%s: # end block", f.labelEndBlock)
	if f.cls.post != nil {
		f.cls.post.emit()
	}
	emit("jmp %s", labelBegin)
	emit("%s: # end loop", f.labelEndLoop)
}

func (f *StmtFor) emit() {
	switch f.kind {
	case FOR_KIND_RANGE_MAP:
		f.emitRangeForMap()
	case FOR_KIND_RANGE_LIST:
		f.emitRangeForList()
	case FOR_KIND_PLAIN:
		f.emitForClause()
	default:
		errorft(f.token(), "unexpected case")
	}
}

func (stmt *StmtReturn) emitDeferAndReturn() {
	if stmt.labelDeferHandler != "" {
		emit("# defer and return")
		emit("jmp %s", stmt.labelDeferHandler)
	}
}

func (ast *StmtDefer) emit() {
	emit("# defer")
	/*
		// arguments should be evaluated immediately
		var args []Expr
		switch ast.expr.(type) {
		case *ExprMethodcall:
			call := ast.expr.(*ExprMethodcall)
			args = call.args
		case *ExprFuncallOrConversion:
			call := ast.expr.(*ExprFuncallOrConversion)
			args = call.args
		default:
			errorft(ast.token(), "defer should be a funcall")
		}
	*/
	labelStart := makeLabel() + "_defer"
	labelEnd := makeLabel() + "_defer"
	ast.label = labelStart

	emit("jmp %s", labelEnd)
	emit("%s: # defer start", labelStart)

	for i := 0; i < len(retRegi); i++ {
		emit("push %%%s", retRegi[i])
	}

	ast.expr.emit()

	for i := len(retRegi) - 1; i >= 0; i-- {
		emit("pop %%%s", retRegi[i])
	}

	emit("leave")
	emit("ret")
	emit("%s: # defer end", labelEnd)

}

func (ast *StmtContinue) emit() {
	assert(ast.stmtFor.labelEndBlock != "", ast.token(), "labelEndLoop should not be empty")
	emit("jmp %s # continue", ast.stmtFor.labelEndBlock)
}

func (ast *StmtBreak) emit() {
	assert(ast.stmtFor.labelEndLoop != "", ast.token(), "labelEndLoop should not be empty")
	emit("jmp %s # break", ast.stmtFor.labelEndLoop)
}


