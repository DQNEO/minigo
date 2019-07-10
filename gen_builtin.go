package main

func (e *ExprLen) emit() {
	emit(S("# emit len()"))
	arg := unwrapRel(e.arg)
	gtype := arg.getGtype()
	assert(gtype != nil, e.token(), S("gtype should not be  nil:\n"))

	switch gtype.getKind() {
	case G_ARRAY:
		emit(S("LOAD_NUMBER %d"), gtype.length)
	case G_SLICE:
		emit(S("# len(slice)"))
		switch arg.(type) {
		case *ExprVariable, *ExprStructField, *ExprIndex:
			emitOffsetLoad(arg, 8, ptrSize)
		case *ExprSliceLiteral:
			emit(S("# ExprSliceLiteral"))
			_arg := arg.(*ExprSliceLiteral)
			var length int = len(_arg.values)
			emit(S("LOAD_NUMBER %d"), length)
		case *ExprSlice:
			sliceExpr := arg.(*ExprSlice)
			uop := &ExprBinop{
				op:    gostring("-"),
				left:  sliceExpr.high,
				right: sliceExpr.low,
			}
			uop.emit()
		default:
			TBI(arg.token(), S("unable to handle %T"), arg)
		}
	case G_MAP:
		emit(S("# emit len(map)"))
		arg.emit()

		// if not nil
		// then 0
		// else len
		labelNil := makeLabel()
		labelEnd := makeLabel()
		emit(S("TEST_IT # map && map (check if map is nil)"))
		emit(S("je %s # jump if map is nil"), labelNil)
		// not nil case
		emit(S("mov 8(%%rax), %%rax # load map len"))
		emit(S("jmp %s"), labelEnd)
		// nil case
		emit(S("%s:"), labelNil)
		emit(S("LOAD_NUMBER 0"))
		emit(S("%s:"), labelEnd)
	case G_CLIKE_STRING:
		arg.emit()
		emit(S("PUSH_8"))
		eStrLen := &IrLowLevelCall{
			symbol:        S("strlen"),
			argsFromStack: 1,
		}
		eStrLen.emit()
	default:
		TBI(arg.token(), S("unable to handle %s"), gtype)
	}
}

type IrLowLevelCall struct {
	token         *Token
	symbol        gostring
	argsFromStack int // args are taken from the stack
}

func (e *IrLowLevelCall) emit() {
	var i int
	for i=e.argsFromStack - 1;i>=0;i-- {
		emit(S("POP_TO_ARG_%d"), i)
	}
	emit(S("FUNCALL %s"), gostring(e.symbol))
}


func (e *ExprCap) emit() {
	emit(S("# emit cap()"))
	arg := unwrapRel(e.arg)
	gtype := arg.getGtype()
	switch gtype.getKind() {
	case G_ARRAY:
		emit(S("LOAD_NUMBER %d"), gtype.length)
	case G_SLICE:
		switch arg.(type) {
		case *ExprVariable, *ExprStructField, *ExprIndex:
			emitOffsetLoad(arg, 8, ptrSize*2)
		case *ExprSliceLiteral:
			emit(S("# ExprSliceLiteral"))
			_arg := arg.(*ExprSliceLiteral)
			var length int = len(_arg.values)
			emit(S("LOAD_NUMBER %d"), length)
		case *ExprSlice:
			sliceExpr := arg.(*ExprSlice)
			if sliceExpr.collection.getGtype().getKind() == G_ARRAY {
				cp := &ExprBinop{
					tok: e.tok,
					op:  gostring("-"),
					left: &ExprLen{
						tok: e.tok,
						arg: sliceExpr.collection,
					},
					right: sliceExpr.low,
				}
				cp.emit()
			} else {
				TBI(arg.token(), S("unable to handle %T"), arg)
			}
		default:
			TBI(arg.token(), S("unable to handle %T"), arg)
		}
	case G_MAP:
		errorft(arg.token(), S("invalid argument for cap"))
	case G_CLIKE_STRING:
		TBI(arg.token(), S("unable to handle %T"), arg)
	default:
		TBI(arg.token(), S("unable to handle %s"), gtype)
	}
}

func emitMakeSliceFunc() {
	// makeSlice
	emitWithoutIndent(S("%s:"), gostring("iruntime.makeSlice"))
	emit(S("FUNC_PROLOGUE"))
	emitNewline()

	emit(S("PUSH_ARG_2")) // -8
	emit(S("PUSH_ARG_1")) // -16
	emit(S("PUSH_ARG_0")) // -24

	emit(S("LOAD_8_FROM_LOCAL -16 # newcap"))
	emit(S("PUSH_8"))
	emit(S("LOAD_8_FROM_LOCAL -8 # unit"))
	emit(S("PUSH_8"))
	emit(S("IMUL_FROM_STACK"))
	emit(S("ADD_NUMBER 1 # 1 byte buffer"))

	emit(S("PUSH_8"))
	emit(S("POP_TO_ARG_0"))
	emit(S("FUNCALL iruntime.malloc"))

	emit(S("mov -24(%%rbp), %%rbx # newlen"))
	emit(S("mov -16(%%rbp), %%rcx # newcap"))

	emit(S("LEAVE_AND_RET"))
	emitNewline()
}
