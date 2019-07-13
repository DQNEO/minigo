package main

var eEmptyString = ExprStringLiteral{
	val: gostring(""),
}

func emitEmptyString() {
	eEmpty := &eEmptyString
	emit(S("LOAD_STRING_LITERAL .%s"), eEmpty.slabel)
	emit(S("mov $0, %%rbx"))
	emit(S("mov $0, %%rcx"))
}

func countStrlen(chars gostring) int {
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
	emit(S("LOAD_STRING_LITERAL .%s"), ast.slabel)
	var length int = countStrlen(ast.val)
	emit(S("mov $%d, %%rbx"), length)
	emit(S("mov $%d, %%rcx"), length)
}

func (e *IrExprStringComparison) token() *Token {
	return e.tok
}

func (binop *IrExprStringComparison) emit() {
	emit(S("# emitCompareStrings"))
	var equal bool
	switch switchexpr(binop.op) {
	case "<":
		TBI(binop.token(), S(""))
	case ">":
		TBI(binop.token(), S(""))
	case "<=":
		TBI(binop.token(), S(""))
	case ">=":
		TBI(binop.token(), S(""))
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
	emit(S("LOAD_NUMBER 1"))
	emit(S("PUSH_8"))

	call := &IrLowLevelCall{
		symbol:        S("iruntime.eqGostrings"),
		argsFromStack: 7,
	}
	call.emit()
}
