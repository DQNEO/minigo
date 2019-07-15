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

func (stmt *StmtSwitch) isTypeSwitch() bool {
	_, isTypeSwitch := stmt.cond.(*ExprTypeSwitchGuard)
	return isTypeSwitch
}

func emitConvertNilToEmptyString() {
	emit("# emitConvertNilToEmptyString")
	emit("POP_8")
	emit("PUSH_8")
	emit("# convert nil to an empty string")
	emit("TEST_IT")
	emit("pop %%rax")
	labelEnd := makeLabel()
	emit("jne %s # jump if not nil", labelEnd)
	emit("# if nil then")
	emitEmptyString()
	emit("%s:", labelEnd)
}

func emitCompareDynamicTypeFromStack(gtype *Gtype) {
	emitConvertNilToEmptyString()
	emit("PUSH_8")

	if gtype.isNil() {
		emitEmptyString()
	} else {
		typeLabel := symbolTable.getTypeLabel(gtype)
		emit("LOAD_STRING_LITERAL .%s # type: %s", typeLabel, gtype.String())
	}

	emit("PUSH_8")
	emit("CMP_FROM_STACK sete") // compare addresses
}

func (stmt *StmtSwitch) emit() {

	emit("# switch statement")
	labelEnd := makeLabel()
	var labels []string
	// switch (expr) {
	var cond Expr
	if stmt.cond != nil {
		cond = stmt.cond
		emit("# the cond expression")
		cond.emit()
		if cond.getGtype().is24WidthType() {
			emit("PUSH_24 # the cond value")
		} else {
			emit("PUSH_8 # the cond value")
		}
	} else {
		// switch {
		emit("# no condition")
	}

	// case exp1,exp2,..:
	//     stmt1;
	//     stmt2;
	//     ...
	for i, caseClause := range stmt.cases {
		var j int = i
		emit("# case %d", j)
		myCaseLabel := makeLabel()
		labels = append(labels, myCaseLabel)
		if stmt.cond == nil {
			for _, e := range caseClause.exprs {
				e.emit()
				emit("TEST_IT")
				emit("jne %s # jump if matches", myCaseLabel)
			}
		} else if stmt.isTypeSwitch() {
			// compare type
			for _, gtype := range caseClause.gtypes {
				emit("# Duplicate the cond value in stack")
				emit("POP_24")
				emit("PUSH_24")

				emit("push %%rcx # push dynamic type addr")
				emitCompareDynamicTypeFromStack(gtype)

				emit("TEST_IT")
				emit("jne %s # jump if matches", myCaseLabel)
			}
		} else {
			for _, e := range caseClause.exprs {
				emit("# Duplicate the cond value in stack")
				if e.getGtype().isString() {
					emit("POP_SLICE # the cond value")
					emit("PUSH_SLICE # the cond value")

					emit("PUSH_SLICE # the cond valiue")

					e.emit()
					emit("PUSH_SLICE")

					emitGoStringsEqualFromStack()
				} else {
					emit("POP_8 # the cond value")
					emit("PUSH_8 # the cond value")

					emit("PUSH_8 # arg1: the cond value")
					e.emit()
					emit("PUSH_8 # arg2: case value")
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

	if cond != nil && cond.getGtype().is24WidthType() {
		emit("POP_24 # destroy the cond value")
	} else {
		emit("POP_8 # destroy the cond value")

	}
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

func (f *IrStmtForRangeList) emit() {
	// i = 0
	emit("# init index")
	f.init.emit()

	emit("%s: # begin loop ", f.labels.labelBegin)

	f.cond.emit()
	emit("TEST_IT")
	emit("je %s  # if false, go to loop end", f.labels.labelEndLoop)

	if f.assignVar != nil {
		f.assignVar.emit()
	}

	f.block.emit()
	emit("%s: # end block", f.labels.labelEndBlock)

	f.cond2.emit()
	emit("TEST_IT")
	emit("jne %s  # if this iteration is final, go to loop end", f.labels.labelEndLoop)

	f.incr.emit()

	emit("jmp %s", f.labels.labelBegin)
	emit("%s: # end loop", f.labels.labelEndLoop)
}

func (f *IrStmtClikeFor) emit() {
	emit("# emit IrStmtClikeFor")
	if f.cls.init != nil {
		f.cls.init.emit()
	}
	emit("%s: # begin loop ", f.labels.labelBegin)
	if f.cls.cond != nil {
		f.cls.cond.emit()
		emit("TEST_IT")
		emit("je %s  # jump if false", f.labels.labelEndLoop)
	}
	f.block.emit()
	emit("%s: # end block", f.labels.labelEndBlock)
	if f.cls.post != nil {
		f.cls.post.emit()
	}
	emit("jmp %s", f.labels.labelBegin)
	emit("%s: # end loop", f.labels.labelEndLoop)
}

func (f *StmtFor) emit() {
	assertNotReached(f.token())
}

func (f *StmtFor) convert() Stmt {
	// Determine kind
	if f.rng != nil {
		if f.rng.rangeexpr.getGtype().getKind() == G_MAP {
			f.kind = FOR_KIND_RANGE_MAP
		} else {
			f.kind = FOR_KIND_RANGE_LIST
		}
	} else {
		f.kind = FOR_KIND_CLIKE
	}

	l1 := makeLabel()
	l2 := makeLabel()
	l3 := makeLabel()

	f.labels.labelBegin = l1

	lbls := f.labels
	lbls.labelEndBlock = l2
	lbls.labelEndLoop = l3

	var em Stmt

	switch f.kind {
	case FOR_KIND_RANGE_MAP:
		assertNotNil(f.rng.indexvar != nil, f.rng.tok)
		em = &IrStmtRangeMap{
			tok:        f.token(),
			block:      f.block,
			labels:     f.labels,
			rangeexpr:  f.rng.rangeexpr,
			indexvar:   f.rng.indexvar,
			valuevar:   f.rng.valuevar,
			mapCounter: f.rng.invisibleMapCounter,
		}
	case FOR_KIND_RANGE_LIST:
		emit("# for range list")
		assertNotNil(f.rng.indexvar != nil, f.rng.tok)
		assert(f.rng.rangeexpr.getGtype().isArrayLike(), f.rng.tok, "rangeexpr should be G_ARRAY or G_SLICE, but got ", f.rng.rangeexpr.getGtype().String())

		var init = &StmtAssignment{
			lefts: []Expr{
				f.rng.indexvar,
			},
			rights: []Expr{
				&ExprNumberLiteral{
					val: 0,
				},
			},
		}
		// i < len(list)
		var cond = &ExprBinop{
			op:   "<",
			left: f.rng.indexvar, // i
			// @TODO
			// The range expression x is evaluated once before beginning the loop
			right: &ExprLen{
				arg: f.rng.rangeexpr, // len(expr)
			},
		}

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
		}

		// break if i == len(list) - 1
		var cond2 = &ExprBinop{
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

		// i++
		var incr = &StmtInc{
			operand: f.rng.indexvar,
		}

		em = &IrStmtForRangeList{
			init:      init,
			cond:      cond,
			assignVar: assignVar,
			cond2:     cond2,
			incr:      incr,
			block:     f.block,
			labels:    f.labels,
		}
	case FOR_KIND_CLIKE:
		em = &IrStmtClikeFor{
			tok:    f.token(),
			labels: f.labels,
			cls:    f.cls,
			block:  f.block,
		}
	default:
		assertNotReached(f.token())
	}

	return em
}

func (stmt *StmtReturn) emitDeferAndReturn() {
	if len(stmt.labelDeferHandler) != 0 {
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
	labelStart := string(makeLabel()) + "_defer"
	labelEnd := string(makeLabel()) + "_defer"
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
	assert(len(ast.labels.labelEndBlock) > 0, ast.token(), "labelEndLoop should not be empty")
	emit("jmp %s # continue", ast.labels.labelEndBlock)
}

func (ast *StmtBreak) emit() {
	assert(len(ast.labels.labelEndLoop) > 0, ast.token(), "labelEndLoop should not be empty")
	emit("jmp %s # break", ast.labels.labelEndLoop)
}
