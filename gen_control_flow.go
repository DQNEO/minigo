package main

func (stmt *StmtIf) emit() {
	emit(S("# if"))
	if stmt.simplestmt != nil {
		stmt.simplestmt.emit()
	}
	stmt.cond.emit()
	emit(S("TEST_IT"))
	if stmt.els != nil {
		labelElse := makeLabel()
		labelEndif := makeLabel()
		emit(S("je %s  # jump if 0"), labelElse)
		emit(S("# then block"))
		stmt.then.emit()
		emit(S("jmp %s # jump to endif"), labelEndif)
		emit(S("# else block"))
		emit(S("%s:"), labelElse)
		stmt.els.emit()
		emit(S("# endif"))
		emit(S("%s:"), labelEndif)
	} else {
		// no else block
		labelEndif := makeLabel()
		emit(S("je %s  # jump if 0"), labelEndif)
		emit(S("# then block"))
		stmt.then.emit()
		emit(S("# endif"))
		emit(S("%s:"), labelEndif)
	}
}

func (stmt *StmtSwitch) isTypeSwitch() bool {
	_, isTypeSwitch := stmt.cond.(*ExprTypeSwitchGuard)
	return isTypeSwitch
}

func emitConvertNilToEmptyString() {
	emit(S("# emitConvertNilToEmptyString"))
	emit(S("POP_8"))
	emit(S("PUSH_8"))
	emit(S("# convert nil to an empty string"))
	emit(S("TEST_IT"))
	emit(S("pop %%rax"))
	labelEnd := makeLabel()
	emit(S("jne %s # jump if not nil"), labelEnd)
	emit(S("# if nil then"))
	emitEmptyString()
	emit(S("%s:"), labelEnd)
}

func emitCompareDynamicTypeFromStack(gtype *Gtype) {
	emitConvertNilToEmptyString()
	emit(S("PUSH_8"))

	if gtype.isNil() {
		emitEmptyString()
	} else {
		typeLabel := symbolTable.getTypeLabel(gtype)
		emit(S("LOAD_STRING_LITERAL .%s # type: %s"), typeLabel, gtype.String())
	}

	emit(S("PUSH_8"))
	emit(S("CMP_FROM_STACK sete")) // compare addresses
}

func (stmt *StmtSwitch) needStringToSliceConversion() bool {
	return ! stmt.isTypeSwitch() && stmt.cond.getGtype().isClikeString() && !gString.is24WidthType()
}

func (stmt *StmtSwitch) emit() {

	emit(S("# switch statement"))
	labelEnd := makeLabel()
	var labels []gostring
	// switch (expr) {
	var cond Expr
	if stmt.cond != nil {
		cond = stmt.cond
		emit(S("# the cond expression"))
		if stmt.needStringToSliceConversion() {
			irConversion, ok := stmt.cond.(*IrExprConversion)
			assert(ok, nil, S("should be IrExprConversion"))
			origType := irConversion.arg.getGtype()
			assert(origType.getKind() == G_SLICE, nil, S("must be slice"))
			cond = irConversion.arg // set original slice
		}

		cond.emit()
		if cond.getGtype().is24WidthType() {
			emit(S("PUSH_24 # the cond value"))
		} else {
			emit(S("PUSH_8 # the cond value"))
		}
	} else {
		// switch {
		emit(S("# no condition"))
	}

	// case exp1,exp2,..:
	//     stmt1;
	//     stmt2;
	//     ...
	for i, caseClause := range stmt.cases {
		var j int = i
		emit(S("# case %d"), j)
		myCaseLabel := makeLabel()
		labels = append(labels, myCaseLabel)
		if stmt.cond == nil {
			for _, e := range caseClause.exprs {
				e.emit()
				emit(S("TEST_IT"))
				emit(S("jne %s # jump if matches"), myCaseLabel)
			}
		} else if stmt.isTypeSwitch() {
			// compare type
			for _, gtype := range caseClause.gtypes {
				emit(S("# Duplicate the cond value in stack"))
				emit(S("POP_24"))
				emit(S("PUSH_24"))

				emit(S("push %%rcx # push dynamic type addr"))
				emitCompareDynamicTypeFromStack(gtype)

				emit(S("TEST_IT"))
				emit(S("jne %s # jump if matches"), myCaseLabel)
			}
		} else {
			for _, e := range caseClause.exprs {
				emit(S("# Duplicate the cond value in stack"))

				if e.getGtype().isClikeString() {
					assert(e.getGtype().isClikeString(), e.token(), S("caseClause should be string"))
					emit(S("POP_SLICE # the cond value"))
					emit(S("PUSH_SLICE # the cond value"))

					emit(S("PUSH_SLICE # the cond valiue"))

					emitConvertCstringToSlice(e)
					emit(S("PUSH_SLICE"))

					emitGoStringsEqualFromStack()
				} else {
					emit(S("POP_8 # the cond value"))
					emit(S("PUSH_8 # the cond value"))

					emit(S("PUSH_8 # arg1: the cond value"))
					e.emit()
					emit(S("PUSH_8 # arg2: case value"))
					emit(S("CMP_FROM_STACK sete"))
				}

				emit(S("TEST_IT"))
				emit(S("jne %s # jump if matches"), myCaseLabel)
			}
		}
	}

	var defaultLabel gostring
	if stmt.dflt == nil {
		emit(S("jmp %s"), labelEnd)
	} else {
		emit(S("# default"))
		defaultLabel = makeLabel()
		emit(S("jmp %s"), defaultLabel)
	}

	if cond != nil && cond.getGtype().is24WidthType() {
		emit(S("POP_24 # destroy the cond value"))
	} else {
		emit(S("POP_8 # destroy the cond value"))

	}
	emit(S("#"))
	for i, caseClause := range stmt.cases {
		emit(S("# case stmts"))
		emit(S("%s:"), labels[i])
		caseClause.compound.emit()
		emit(S("jmp %s"), labelEnd)
	}

	if stmt.dflt != nil {
		emit(S("%s:"), defaultLabel)
		stmt.dflt.emit()
	}

	emit(S("%s: # end of switch"), labelEnd)
}

func (f *IrStmtForRangeList) emit() {
	// i = 0
	emit(S("# init index"))
	f.init.emit()

	emit(S("%s: # begin loop "), f.labels.labelBegin)

	f.cond.emit()
	emit(S("TEST_IT"))
	emit(S("je %s  # if false, go to loop end"), f.labels.labelEndLoop)

	if f.assignVar != nil {
		f.assignVar.emit()
	}

	f.block.emit()
	emit(S("%s: # end block"), f.labels.labelEndBlock)

	f.cond2.emit()
	emit(S("TEST_IT"))
	emit(S("jne %s  # if this iteration is final, go to loop end"), f.labels.labelEndLoop)

	f.incr.emit()

	emit(S("jmp %s"), f.labels.labelBegin)
	emit(S("%s: # end loop"), f.labels.labelEndLoop)
}

func (f *IrStmtClikeFor) emit() {
	emit(S("# emit IrStmtClikeFor"))
	if f.cls.init != nil {
		f.cls.init.emit()
	}
	emit(S("%s: # begin loop "), f.labels.labelBegin)
	if f.cls.cond != nil {
		f.cls.cond.emit()
		emit(S("TEST_IT"))
		emit(S("je %s  # jump if false"), f.labels.labelEndLoop)
	}
	f.block.emit()
	emit(S("%s: # end block"), f.labels.labelEndBlock)
	if f.cls.post != nil {
		f.cls.post.emit()
	}
	emit(S("jmp %s"), f.labels.labelBegin)
	emit(S("%s: # end loop"), f.labels.labelEndLoop)
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
	l2  := makeLabel()
	l3 := makeLabel()

	lbls := f.labels
	// @FIXME:  f.labels.labelBegin = l1  does not work!!!
	lbls.labelBegin = l1
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
		emit(S("# for range list"))
		assertNotNil(f.rng.indexvar != nil, f.rng.tok)
		assert(f.rng.rangeexpr.getGtype().isArrayLike(), f.rng.tok, S("rangeexpr should be G_ARRAY or G_SLICE, but got "), f.rng.rangeexpr.getGtype().String())

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
			op:   gostring("<"),
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
			op:   gostring("=="),
			left: f.rng.indexvar, // i
			// @TODO2
			// The range expression x is evaluated once before beginning the loop
			right: &ExprBinop{
				op: gostring("-"),
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
		emit(S("# defer and return"))
		emit(S("jmp %s"), stmt.labelDeferHandler)
	}
}

func (ast *StmtDefer) emit() {
	emit(S("# defer"))
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
			errorft(ast.token(), S("defer should be a funcall"))
		}
	*/
	labelStart := concat(makeLabel() , S("_defer"))
	labelEnd := concat(makeLabel() , S("_defer"))
	ast.label = labelStart

	emit(S("jmp %s"), labelEnd)
	emit(S("%s: # defer start"), labelStart)

	for i := 0; i < len(retRegi); i++ {
		emit(S("push %%%s"), retRegi[i])
	}

	ast.expr.emit()

	for i := len(retRegi) - 1; i >= 0; i-- {
		emit(S("pop %%%s"), retRegi[i])
	}

	emit(S("leave"))
	emit(S("ret"))
	emit(S("%s: # defer end"), labelEnd)

}

func (ast *StmtContinue) emit() {
	assert(len(ast.labels.labelEndBlock) > 0, ast.token(), S("labelEndLoop should not be empty"))
	emit(S("jmp %s # continue"), ast.labels.labelEndBlock)
}

func (ast *StmtBreak) emit() {
	assert(len(ast.labels.labelEndLoop) > 0, ast.token(), S("labelEndLoop should not be empty"))
	emit(S("jmp %s # break"), ast.labels.labelEndLoop)
}
