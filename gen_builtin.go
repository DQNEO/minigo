package main

import "fmt"

func (e *ExprLen) emit() {
	emit2("# emit len()")
	arg := unwrapRel(e.arg)
	gtype := arg.getGtype()
	assert(gtype != nil, e.token(), "gtype should not be  nil:\n"+fmt.Sprintf("%#v", arg))

	switch gtype.getKind() {
	case G_ARRAY:
		emit2("LOAD_NUMBER %d", gtype.length)
	case G_SLICE:
		emit2("# len(slice)")
		switch arg.(type) {
		case *ExprVariable, *ExprStructField, *ExprIndex:
			emitOffsetLoad(arg, 8, ptrSize)
		case *ExprSliceLiteral:
			emit2("# ExprSliceLiteral")
			_arg := arg.(*ExprSliceLiteral)
			var length int = len(_arg.values)
			emit2("LOAD_NUMBER %d", length)
		case *ExprSlice:
			sliceExpr := arg.(*ExprSlice)
			uop := &ExprBinop{
				op:    gostring("-"),
				left:  sliceExpr.high,
				right: sliceExpr.low,
			}
			uop.emit()
		default:
			TBI(arg.token(), "unable to handle %T", arg)
		}
	case G_MAP:
		emit2("# emit len(map)")
		arg.emit()

		// if not nil
		// then 0
		// else len
		labelNil := makeLabel()
		labelEnd := makeLabel()
		emit2("TEST_IT # map && map (check if map is nil)")
		emit2("je %s # jump if map is nil", labelNil)
		// not nil case
		emit2("mov 8(%%rax), %%rax # load map len")
		emit2("jmp %s", labelEnd)
		// nil case
		emit2("%s:", labelNil)
		emit2("LOAD_NUMBER 0")
		emit2("%s:", labelEnd)
	case G_STRING:
		arg.emit()
		emit2("PUSH_8")
		eStrLen := &IrLowLevelCall{
			symbol:        "strlen",
			argsFromStack: 1,
		}
		eStrLen.emit()
	default:
		TBI(arg.token(), "unable to handle %s", gtype)
	}
}

type IrLowLevelCall struct {
	token         *Token
	symbol        cstring
	argsFromStack int // args are taken from the stack
}

func (e *IrLowLevelCall) emit() {
	var i int
	for i=e.argsFromStack - 1;i>=0;i-- {
		emit2("POP_TO_ARG_%d", i)
	}
	emit2("FUNCALL %s", gostring(e.symbol))
}


func (e *ExprCap) emit() {
	emit2("# emit cap()")
	arg := unwrapRel(e.arg)
	gtype := arg.getGtype()
	switch gtype.getKind() {
	case G_ARRAY:
		emit2("LOAD_NUMBER %d", gtype.length)
	case G_SLICE:
		switch arg.(type) {
		case *ExprVariable, *ExprStructField, *ExprIndex:
			emitOffsetLoad(arg, 8, ptrSize*2)
		case *ExprSliceLiteral:
			emit2("# ExprSliceLiteral")
			_arg := arg.(*ExprSliceLiteral)
			var length int = len(_arg.values)
			emit2("LOAD_NUMBER %d", length)
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
	emitWithoutIndent("%s:", gostring("iruntime.makeSlice"))
	emit2("FUNC_PROLOGUE")
	emitNewline()

	emit2("PUSH_ARG_2") // -8
	emit2("PUSH_ARG_1") // -16
	emit2("PUSH_ARG_0") // -24

	emit2("LOAD_8_FROM_LOCAL -16 # newcap")
	emit2("PUSH_8")
	emit2("LOAD_8_FROM_LOCAL -8 # unit")
	emit2("PUSH_8")
	emit2("IMUL_FROM_STACK")
	emit2("ADD_NUMBER 1 # 1 byte buffer")

	emit2("PUSH_8")
	emit2("POP_TO_ARG_0")
	emit2("FUNCALL iruntime.malloc")

	emit2("mov -24(%%rbp), %%rbx # newlen")
	emit2("mov -16(%%rbp), %%rcx # newcap")

	emit2("LEAVE_AND_RET")
	emitNewline()
}
