package main

import "fmt"

func (e *ExprLen) emit() {
	emit("# emit len()")
	arg := e.arg
	gtype := arg.getGtype()
	assert(gtype != nil, e.token(), "gtype should not be  nil:\n"+fmt.Sprintf("%#v", arg))

	switch gtype.getKind() {
	case G_ARRAY:
		emit("LOAD_NUMBER %d", gtype.length)
	case G_SLICE:
		emit("# len(slice)")
		switch arg.(type) {
		case *Relation:
			emit("# Relation")
			emitOffsetLoad(arg, 8, ptrSize)
		case *ExprStructField:
			emit("# ExprStructField")
			emitOffsetLoad(arg, 8, ptrSize)
		case *ExprIndex:
			emitOffsetLoad(arg, 8, ptrSize)
		case *ExprSliceLiteral:
			emit("# ExprSliceLiteral")
			_arg := arg.(*ExprSliceLiteral)
			length := len(_arg.values)
			emit("LOAD_NUMBER %d", length)
		case *ExprSlice:
			sliceExpr := arg.(*ExprSlice)
			uop := &ExprBinop{
				op:    "-",
				left:  sliceExpr.high,
				right: sliceExpr.low,
			}
			uop.emit()
		default:
			TBI(arg.token(), "unable to handle %T", arg)
		}
	case G_MAP:
		emit("# emit len(map)")
		arg.emit()

		// if not nil
		// then 0
		// else len
		labelNil := makeLabel()
		labelEnd := makeLabel()
		emit("TEST_IT # map && map (check if map is nil)")
		emit("je %s # jump if map is nil", labelNil)
		// not nil case
		emit("mov 8(%%rax), %%rax # load map len")
		emit("jmp %s", labelEnd)
		// nil case
		emit("%s:", labelNil)
		emit("LOAD_NUMBER 0")
		emit("%s:", labelEnd)
	case G_STRING:
		arg.emit()
		emit("PUSH_8")
		emit("POP_TO_ARG_0")
		emit("FUNCALL strlen")
	default:
		TBI(arg.token(), "unable to handle %s", gtype)
	}
}

func (e *ExprCap) emit() {
	emit("# emit cap()")
	arg := e.arg
	gtype := arg.getGtype()
	switch gtype.getKind() {
	case G_ARRAY:
		emit("LOAD_NUMBER %d", gtype.length)
	case G_SLICE:
		switch arg.(type) {
		case *Relation:
			emit("# Relation")
			emitOffsetLoad(arg, 8, ptrSize*2)
		case *ExprStructField:
			emit("# ExprStructField")
			emitOffsetLoad(arg, 8, ptrSize*2)
		case *ExprIndex:
			emitOffsetLoad(arg, 8, ptrSize*2)
		case *ExprSliceLiteral:
			emit("# ExprSliceLiteral")
			_arg := arg.(*ExprSliceLiteral)
			length := len(_arg.values)
			emit("LOAD_NUMBER %d", length)
		case *ExprSlice:
			sliceExpr := arg.(*ExprSlice)
			if sliceExpr.collection.getGtype().getKind() == G_ARRAY {
				cp := &ExprBinop{
					tok: e.tok,
					op:  "-",
					left: &ExprLen{
						tok: e.tok,
						arg: sliceExpr.collection,
					},
					right: sliceExpr.low,
				}
				cp.emit()
			} else {
				TBI(arg.token(), "unable to handle %T", arg)
			}
		default:
			TBI(arg.token(), "unable to handle %T", arg)
		}
	case G_MAP:
		errorft(arg.token(), "invalid argument for cap")
	case G_STRING:
		TBI(arg.token(), "unable to handle %T", arg)
	default:
		TBI(arg.token(), "unable to handle %s", gtype)
	}
}

func emitMakeSliceFunc() {
	// makeSlice
	emitWithoutIndent("%s:", "iruntime.makeSlice")
	emit("FUNC_PROLOGUE")
	emitNewline()

	emit("PUSH_ARG_2") // -8
	emit("PUSH_ARG_1") // -16
	emit("PUSH_ARG_0") // -24

	emit("LOAD_8_FROM_LOCAL -16 # newcap")
	emit("PUSH_8")
	emit("LOAD_8_FROM_LOCAL -8 # unit")
	emit("PUSH_8")
	emit("IMUL_FROM_STACK")
	emit("ADD_NUMBER 1 # 1 byte buffer")

	emit("PUSH_8")
	emit("POP_TO_ARG_0")
	emit("FUNCALL iruntime.malloc")

	emit("mov -24(%%rbp), %%rbx # newlen")
	emit("mov -16(%%rbp), %%rcx # newcap")

	emit("LEAVE_AND_RET")
	emitNewline()
}


