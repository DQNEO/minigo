package main

import "fmt"

func emit(format string, v ...interface{}) {
	fmt.Printf("\t"+format+"\n", v...)
}

func emitLabel(format string, v ...interface{}) {
	fmt.Printf(format+"\n", v...)
}

func emitDataSection() {
	emit(".data")

	// put strings
	for _, ast := range strings {
		emitLabel(".%s:", ast.slabel)
		emit(".string \"%s\"", ast.sval)
	}
}

func emitFuncPrologue(f *AstFuncDef) {
	emitLabel("# f %s", f.fname)
	emit(".text")
	emit(".globl	%s", f.fname)
	emitLabel("%s:", f.fname)
	emit("push %%rbp")
	emit("mov %%rsp, %%rbp")

	if len(f.localvars) > 0 {
		emit("# allocate stack area")
		stack_size := 8 * len(f.localvars)
		emit("sub $%d, %%rsp", stack_size)
	}

	emit("# end of prologue")
}

func emitFuncEpilogue() {
	emit("#")
	emit("leave")
	emit("ret")
}

func emitInt(ast *Ast) {
	emit("movl	$%d, %%eax", ast.ival)
}

func emitString(ast *Ast) {
	emit("lea .%s(%%rip), %%rax", ast.slabel)
}

func emitBinop(ast *Ast) {
	emitExpr(ast.left)
	emit("push %%rax")
	emitExpr(ast.right)
	emit("push %%rax")
	emit("pop %%rbx")
	emit("pop %%rax")
	if ast.op == "+" {
		emit("addl	%%ebx, %%eax")
	} else if ast.op == "-" {
		emit("subl	%%ebx, %%eax")
	} else if ast.op == "*" {
		emit("imul	%%ebx, %%eax")
	}
}

func emitLocalVariable(ast *Ast) {
	emit("mov %d(%%rbp), %%eax", ast.offset)
}

func emitAssignment(ast *Ast) {
	emitExpr(ast.right)
	emit("push %%rax")
	emit("mov %%eax, %d(%%rbp)", ast.left.offset)
}

func emitCompound(ast *AstCompountStmt) {
	for _, stmt := range ast.stmts {
		emitExpr(stmt)
	}
}

func emitExpr(ast *Ast) {
	switch ast.typ {
	case "int":
		emitInt(ast)
	case "binop":
		emitBinop(ast)
	case "lvar":
		emitLocalVariable(ast)
	case "assign":
		emitAssignment(ast)
	case "string":
		emitString(ast)
	case "funcall":
		emitFuncall(ast)
	case "decl":
		;
	default:
		panic(fmt.Sprintf("unexpected ast type %s", ast.typ))
	}
}

var regs = []string{"rdi", "rsi", "rdx", "rcx", "r8", "r9"}

func emitFuncall(funcall *Ast) {
	fname := funcall.fname
	emit("# funcall %s", fname)
	args := funcall.args
	for i, _ := range args {
		emit("push %%%s", regs[i])
	}

	emit("# setting arguments")
	for _, arg := range args {
		emitExpr(arg)
		emit("push %%rax")
	}

	for i, _ := range args {
		j := len(args) - 1 - i
		emit("pop %%%s", regs[j])
	}
	emit("mov $0, %%rax")
	emit("call %s", fname)

	for i, _ := range args {
		j := len(args) - 1 - i
		emit("pop %%%s", regs[j])
	}
}

func emitFuncdef(f *AstFuncDef) {
	// calc offset
	var offset int
	for _, lvar := range f.localvars {
		offset -= 8
		lvar.offset = offset
	}

	emitFuncPrologue(f)
	emitCompound(f.body)
	emit("mov $0, %%eax") // return 0
	emitFuncEpilogue()
}

func generate(a *AstFile) {
	emitDataSection()
	for _, f := range a.funcdefs {
		emitFuncdef(f)
	}
}
