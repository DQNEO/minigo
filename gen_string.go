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

func (e *IrExprStringComparison) token() *Token {
	return e.tok
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

	// func eq(a string, b string, eq bool) bool
	// @TODO get params by DeclFunc dynamically
	params = append(params, &ExprVariable{}) // a []byte
	params = append(params, &ExprVariable{}) // b []byte
	params = append(params, &ExprVariable{}) // eq bool
	// eq(left, right, eFlag)
	call := &IrStaticCall{
		tok: binop.token(),
		symbol: "iruntime.cmpStrings",
		isMethodCall:false,
		args: args,
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
		symbol:        "iruntime.cmpStrings",
		argsFromStack: 7,
	}
	call.emit()
}
