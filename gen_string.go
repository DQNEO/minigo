package main

var eEmptyString = ExprStringLiteral{
	val: gostring(""),
}

func emitEmptyString() {
	eEmpty := &eEmptyString
	eEmpty.emit()
	emit2("mov $0, %%rbx")
	emit2("mov $0, %%rcx")
}

func (ast *ExprStringLiteral) emit() {
	emit2("LOAD_STRING_LITERAL .%s", ast.slabel)
	var length int = len(ast.val)
	emit2("mov $%d, %%rbx", length)
	emit2("mov $%d, %%rcx", length)
}

func (e *IrExprStringComparison) token() *Token {
	return e.tok
}

func (binop *IrExprStringComparison) emit() {
	emit2("# emitCompareStrings")
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
	emit2("LOAD_NUMBER 1")
	emit2("PUSH_8")

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

	emit2("POP_8 # restore string")
	emit2("TEST_IT") // check if string is nil
	emit2("jne %s # go to then if not nil", labelThen)
	emit2("# if nil ")
	emit2("LOAD_EMPTY_SLICE") // emit 0,0,0
	emitEmptyString() // emit ""
	emit2("jmp %s", labelEnd)
	emit2("%s:", labelThen)

	emit2("PUSH_8 # string addr")

	// calc len
	emit2("PUSH_8")
	eStrLen := &IrLowLevelCall{
		symbol:        "strlen",
		argsFromStack: 1,
	}
	eStrLen.emit()
	emit2("mov %%rax, %%rbx # len")
	emit2("mov %%rax, %%rcx # cap")

	emit2("POP_8 # string addr")
	emit2("%s:", labelEnd)
}
// emit []byte(cstring)
func emitConvertCstringToSlice(cstring Expr) {
	cstring.emit()

	if gString.is24WidthType() {
		return
	}

	emit2("PUSH_8")

	emitConvertCstringFromStackToSlice()
}

func emitStringConcate(leftCstring Expr, rightCstring Expr) {
	emit2("# emitStringConcate")


	emitConvertCstringToSlice(leftCstring)
	emit2("PUSH_SLICE")

	emitConvertCstringToSlice(rightCstring)
	emit2("PUSH_SLICE")

	eStrConCate := &IrLowLevelCall{
		symbol:        "iruntime.strcat",
		argsFromStack: 6,
	}
	eStrConCate.emit()
}
