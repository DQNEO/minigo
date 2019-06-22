package main

var eEmptyString = ExprStringLiteral{
	val: "",
}

func emitEmptyString() {
	eEmpty := &eEmptyString
	eEmpty.emit()
}

func (e *IrExprStringComparison) token() *Token {
	return e.tok
}

func (binop *IrExprStringComparison) emit() {
	emit("# emitCompareStrings")
	var equal bool
	switch cstring(binop.op) {
	case "<":
		TBI(binop.token(), "")
	case ">":
		TBI(binop.token(), "")
	case "<=":
		TBI(binop.token(), "")
	case ">=":
		TBI(binop.token(), "")
	case "!=":
		equal = false
	case "==":
		equal = true
	}

	// 3rd arg
	var eEqual Expr
	if equal {
		eEqual = &ExprNumberLiteral{
			val:1,
		}
	} else {
		eEqual = &ExprNumberLiteral{
			val:0,
		}
	}

	left := &IrExprConversion{
		tok: binop.cstringLeft.token(),
		toGtype: &Gtype{
			kind: G_SLICE,
			elementType:gByte,
		},
		arg: binop.cstringLeft,
	}

	right := &IrExprConversion{
		tok: binop.cstringRight.token(),
		toGtype: &Gtype{
			kind: G_SLICE,
			elementType:gByte,
		},
		arg: binop.cstringRight,
	}

	var args []Expr
	args = append(args, left)
	args = append(args, right)
	args = append(args, eEqual)

	var params []*ExprVariable
	// func eqGostring(a []byte, b []byte, eq bool) bool
	// @TODO get params by DeclFunc dynamically
	params = append(params, &ExprVariable{}) // a []byte
	params = append(params, &ExprVariable{}) // b []byte
	params = append(params, &ExprVariable{}) // eq bool
	// eqGostring(left, right, eEqual)
	call := &IrStaticCall{
		tok: binop.token(),
		symbol: "iruntime.eqGostring",
		isMethodCall:false,
		args: args,
		callee: &DeclFunc{
			params: params,
		},
	}
	call.emit()
}

func emitConvertNilToEmptyString() {
	emit("# emitConvertNilToEmptyString")

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

func emitCStringsEqualFromStack(equal bool) {
	emit("pop %%rax") // left

	emitConvertNilToEmptyString()
	emit("mov %%rax, %%rcx")
	emit("pop %%rax # right string")
	emit("push %%rcx")
	emitConvertNilToEmptyString()
	emit("PUSH_8")

	// 3rd arg
	if equal {
		emit("LOAD_NUMBER 1")
	} else {
		emit("LOAD_NUMBER 0")
	}
	emit("PUSH_8")

	call := &IrLowLevelCall{
		symbol:        "iruntime.eqCstrings",
		argsFromStack: 3,
	}
	call.emit()}

// emit []byte(cstring)
func emitConvertStringToSlice(cstring Expr) {
	labelEnd := makeLabel()
	labelThen := makeLabel()
	cstring.emit()
	emit("TEST_IT") // check if string is nil
	emit("jne %s # go to then if not nil", labelThen)
	emit("# if nil ")
	emit("LOAD_EMPTY_SLICE") // emit 0,0,0
	emitEmptyString() // emit ""
	emit("jmp %s", labelEnd)
	emit("%s:", labelThen)

	emit("PUSH_8 # string addr")

	// calc len
	emit("PUSH_8")
	eStrLen := &IrLowLevelCall{
		symbol:        "strlen",
		argsFromStack: 1,
	}
	eStrLen.emit()
	emit("mov %%rax, %%rbx # len")
	emit("mov %%rax, %%rcx # cap")

	emit("POP_8 # string addr")
	emit("%s:", labelEnd)
}

func emitStringConcate(leftCstring Expr, rightCstring Expr) {
	emit("# emitStringConcate")


	emitConvertStringToSlice(leftCstring)
	emit("PUSH_SLICE")

	emitConvertStringToSlice(rightCstring)
	emit("PUSH_SLICE")

	eStrConCate := &IrLowLevelCall{
		symbol:        "iruntime.gostringconcate",
		argsFromStack: 6,
	}
	eStrConCate.emit()
}
