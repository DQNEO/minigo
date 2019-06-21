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

	labelElse := makeLabel()
	labelEnd := makeLabel()

	binop.left.emit()

	// convert nil to the empty string
	emit("CMP_EQ_ZERO")
	emit("TEST_IT")
	emit("LOAD_NUMBER 0")
	emit("je %s", labelElse)
	emitEmptyString()
	emit("jmp %s", labelEnd)
	emit("%s:", labelElse)
	binop.left.emit()
	emit("%s:", labelEnd)
	emit("PUSH_8")

	binop.right.emit()
	emit("PUSH_8")
	emitCStringsEqualFromStack(equal)
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

	emit("PUSH_8")
	emit("POP_TO_ARG_0")
	emit("FUNCALL strlen # get left len")

	emit("PUSH_8 # left len")
	right.emit()
	emit("PUSH_8 # right string")
	emit("PUSH_8")
	emit("POP_TO_ARG_0")
	emit("FUNCALL strlen # get right len")
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
	emit("POP_TO_ARG_0")
	emit("FUNCALL iruntime.malloc")

	emit("PUSH_8")
	emit("POP_TO_ARG_0")
	emit("POP_TO_ARG_1")
	emit("FUNCALL strcat")

	emit("PUSH_8")
	emit("POP_TO_ARG_0")
	emit("POP_TO_ARG_1")
	emit("FUNCALL strcat")
}
