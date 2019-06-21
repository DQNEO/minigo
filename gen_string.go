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

	emit("POP_TO_ARG_2")
	emit("POP_TO_ARG_1")
	emit("POP_TO_ARG_0")
	emit("FUNCALL iruntime.eqCstrings")
}

func emitStringConcate(left Expr, right Expr) {
	emit("# emitStringConcate")
	left.emit()
	emit("PUSH_8 # left string")

	// get left len
	emit("PUSH_8")
	eStrLen := &IrLowLevelCall{
		symbol:        "strlen",
		argsFromStack: 1,
	}
	eStrLen.emit()
	emit("PUSH_8 # left len")

	right.emit()
	emit("PUSH_8 # right string")

	// get right len
	emit("PUSH_8")
	eStrLen = &IrLowLevelCall{
		symbol:        "strlen",
		argsFromStack: 1,
	}
	eStrLen.emit()

	emit("PUSH_8 # right len")

	emit("pop %%rax # right len")
	emit("pop %%rcx # right string")
	emit("pop %%rbx # left len")
	emit("pop %%rdx # left string")

	emit("push %%rcx # right string")
	emit("push %%rdx # left  string")

	// newSize = strlen(left) + strlen(right) + 1
	emit("add %%rax, %%rbx # len + len")
	emit("add $1, %%rbx # + 1 (null byte)")
	emit("mov %%rbx, %%rax")

	emit("PUSH_8")
	e := &IrLowLevelCall{
		symbol:        "iruntime.malloc",
		argsFromStack: 1,
	}
	e.emit()

	emit("mov %%rax, %%rcx") // malloced
	emit("POP_8")
	emit("push %%rcx") // malloced
	emit("PUSH_8")

	e = &IrLowLevelCall{
		symbol:        "strcat",
		argsFromStack: 2,
	}
	e.emit()

	emit("mov %%rax, %%rcx") // strcatted
	emit("POP_8")
	emit("push %%rcx") // strcatted
	emit("PUSH_8")

	e = &IrLowLevelCall{
		symbol:        "strcat",
		argsFromStack: 2,
	}
	e.emit()
}
