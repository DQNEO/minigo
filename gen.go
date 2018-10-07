package main

import "fmt"

func emitFuncMainPrologue() {
	fmt.Println("\t.globl	main")
	fmt.Println("main:")
}

func emitFuncMainEpilogue() {
	fmt.Println("\tret")
}

func emitExpr(ast *Ast) {
	if ast.typ == "uop" {
		fmt.Printf("\tmovl	$%d, %%eax\n", ast.operand.ival)
	} else if ast.typ == "binop" {
		fmt.Printf("\tmovl	$%d, %%eax\n", ast.left.operand.ival)
		fmt.Printf("\tmovl	$%d, %%ebx\n", ast.right.operand.ival)
		if ast.op == "+" {
			fmt.Printf("\taddl	%%ebx, %%eax\n")
		} else if ast.op == "-" {
			fmt.Printf("\tsubl	%%ebx, %%eax\n")
		} else if ast.op == "*" {
			fmt.Printf("\timul	%%ebx, %%eax\n")
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
