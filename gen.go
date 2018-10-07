package main

import "fmt"

func emit(format string, v ...interface{})  {
	fmt.Printf("\t" + format + "\n", v...)
}

func emitLabel(format string, v ...interface{})  {
	fmt.Printf(format + "\n", v...)
}


func emitFuncMainPrologue() {
	emit(".globl	main")
	emitLabel("main:")
}

func emitFuncMainEpilogue() {
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
	emitFuncMainPrologue()
	emitExpr(expr)
	emitFuncMainEpilogue()
}
