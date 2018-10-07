package main

import "fmt"

func emit(format string, v ...interface{})  {
	fmt.Printf("\t" + format + "\n", v...)
}

func emitLabel(format string, v ...interface{})  {
	fmt.Printf(format + "\n", v...)
}

func emitDataSection() {
	emit(".data")

	// put dummy label
	emitLabel(".L0:")
	emit(".string \"%%d\\n\"")
}

func emitFuncMainPrologue() {
	emit(".text")
	emit(".globl	main")
	emitLabel("main:")
	emit("push %%rbp")
	emit("mov %%rsp, %%rbp")
}

func emitFuncMainEpilogue() {
	emit("leave")
	emit("ret")
}

func emitExpr(ast *Ast) {
	if ast.typ == "uop" {
		emit("movl	$%d, %%eax", ast.operand.ival)
	} else if ast.typ == "binop" {
		emit("movl	$%d, %%eax", ast.left.operand.ival)
		emit("movl	$%d, %%ebx", ast.right.operand.ival)
		if ast.op == "+" {
			emit("addl	%%ebx, %%eax")
		} else if ast.op == "-" {
			emit("subl	%%ebx, %%eax")
		} else if ast.op == "*" {
			emit("imul	%%ebx, %%eax")
		}
	} else {
		panic(fmt.Sprintf("unexpected ast type %s", ast.typ))
	}
}

func generate(expr *Ast) {
	emitDataSection()
	emitFuncMainPrologue()

	// call printf("%d\n", expr)
	emit("push %%rdi")
	emit("push %%rsi")

	emitExpr(expr)
	emit("push %%rax")

	emit("lea .L0(%%rip), %%rdi") // first argument
	emit("pop %%rsi") // second argument
	emit("mov $0, %%rax")
	emit("call printf")
	emit("pop %%rsi")
	emit("pop %%rdi")

	emit("mov $0, %%eax") // return 0
	emitFuncMainEpilogue()
}
