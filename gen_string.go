package main

func emitEmptyString() {
	emit("LOAD_EMPTY_SLICE")
}

func countStrlen(chars []byte) int {
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

func (e *IrStringConcat) emit() {
	emit("# IrExprStringComparison")
	var args []Expr
	args = append(args, e.left)
	args = append(args, e.right)

	var params []*ExprVariable

	var dummyVariable = &ExprVariable{
		isVariadic: false,
	}
	params = append(params, dummyVariable) // 1st arg
	params = append(params, dummyVariable) // 2nd arg

	// left + right
	call := &IrCall{
		tok:          e.token(),
		symbol:       getFuncSymbol(IRuntimePath, "concat"),
		args:         args,
		callee: &DeclFunc{
			params: params,
		},
	}
	call.emit()

}

func (binop *IrExprStringComparison) emit() {
	emit("# IrExprStringComparison")
	var equal bool
	switch binop.op {
	case "!=":
		equal = false
	case "==":
		equal = true
	default:
		TBI(binop.token(), "")
	}

	// 3rd arg
	var eFlag Expr
	if equal {
		eFlag = eTrue
	} else {
		eFlag = eFalse
	}

	var args []Expr
	args = append(args, binop.left)
	args = append(args, binop.right)
	args = append(args, eFlag)

	var params []*ExprVariable

	var dummyVariable = &ExprVariable{
		isVariadic: false,
	}
	// func eq(a string, b string, eq bool) bool
	// @TODO get params by DeclFunc dynamically
	params = append(params, dummyVariable) // a
	params = append(params, dummyVariable) // b
	params = append(params, dummyVariable) // eq
	// eq(left, right, eFlag)
	call := &IrCall{
		tok:          binop.token(),
		symbol:       getFuncSymbol(IRuntimePath, "cmpStrings"),
		args:         args,
		callee: &DeclFunc{
			params: params,
		},
	}
	call.emit()
}

func emitGoStringsEqualFromStack() {
	eTrue.emit()
	emit("PUSH_8")

	call := &IrLowLevelCall{
		symbol:        getFuncSymbol(IRuntimePath, "cmpStrings"),
		argsFromStack: 7,
	}
	call.emit()
}
