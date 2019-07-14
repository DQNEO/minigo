package main

func emitEmptyString() {
	emit("LOAD_EMPTY_SLICE")
}

func countStrlen(chars bytes) int {
	var length int
	var isInBackSlash bool

	for _, c := range chars {
		if !isInBackSlash && c == '\\' {
			isInBackSlash = true
			continue
		} else if isInBackSlash && c == '\\' {
			isInBackSlash = false
			length++
		} else {
			isInBackSlash = false
			length++
		}
	}
	return length
}

func (ast *ExprStringLiteral) emit() {
	var length int = countStrlen(ast.val)
	if length == 0 {
		emitEmptyString()
	} else {
		emit("LOAD_STRING_LITERAL .%s", ast.slabel)
		emit("mov $%d, %%rbx", length)
		emit("mov $%d, %%rcx", length)
	}
}

func (e *IrExprStringComparison) token() *Token {
	return e.tok
}

func (binop *IrExprStringComparison) emit() {
	emit("# emitCompareStrings")
	var equal bool
	op := binop.op
	switch op {
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
	// func eq(a []byte, b []byte, eq bool) bool
	// @TODO get params by DeclFunc dynamically
	params = append(params, &ExprVariable{}) // a []byte
	params = append(params, &ExprVariable{}) // b []byte
	params = append(params, &ExprVariable{}) // eq bool
	// eq(left, right, eEqual)
	call := &IrStaticCall{
		tok: binop.token(),
		symbol: "iruntime.eqGostrings",
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
		symbol:        S("iruntime.eqGostrings"),
		argsFromStack: 7,
	}
	call.emit()
}
