package main

var eEmptyString = ExprStringLiteral{
	val: gostring(""),
}

func emitEmptyString() {
	eEmpty := &eEmptyString
	eEmpty.emit()
	emit("mov $0, %%rbx")
	emit("mov $0, %%rcx")
}

func (ast *ExprStringLiteral) emit() {
	emit("LOAD_STRING_LITERAL .%s", ast.slabel)
	emit("mov $%d, %%rbx", len(ast.val))
	emit("mov $%d, %%rcx", len(ast.val))
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
	// func eqGostrings(a []byte, b []byte, eq bool) bool
	// @TODO get params by DeclFunc dynamically
	params = append(params, &ExprVariable{}) // a []byte
	params = append(params, &ExprVariable{}) // b []byte
	params = append(params, &ExprVariable{}) // eq bool
	// eqGostrings(left, right, eEqual)
	call := &IrStaticCall{
		tok: binop.token(),
		symbol: S("iruntime.eqGostrings"),
		isMethodCall:false,
		args: args,
		callee: &DeclFunc{
			params: params,
		},
	}
	call.emit()
}

func emitGoStringsEqualFromStack() {
	emit("LOAD_NUMBER 1")
	emit("PUSH_8")

	call := &IrLowLevelCall{
		symbol:        "iruntime.eqGostrings",
		argsFromStack: 7,
	}
	call.emit()
}

// emit []byte(cstring)
func emitConvertCstringFromStackToSlice() {
	labelEnd := makeLabel()
	labelThen := makeLabel()

	emit("POP_8 # restore string")
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
// emit []byte(cstring)
func emitConvertCstringToSlice(cstring Expr) {
	cstring.emit()

	if gString.is24WidthType() {
		return
	}

	emit("PUSH_8")

	emitConvertCstringFromStackToSlice()
}

func emitStringConcate(leftCstring Expr, rightCstring Expr) {
	emit("# emitStringConcate")


	emitConvertCstringToSlice(leftCstring)
	emit("PUSH_SLICE")

	emitConvertCstringToSlice(rightCstring)
	emit("PUSH_SLICE")

	eStrConCate := &IrLowLevelCall{
		symbol:        "iruntime.strcat",
		argsFromStack: 6,
	}
	eStrConCate.emit()
}
