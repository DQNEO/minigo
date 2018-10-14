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
		emit(".string \"%s\"", ast.val)
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

func (ast *ExprNumberLiteral) emit() {
	emit("movl	$%d, %%eax", ast.val)
}

func (ast *ExprStringLiteral) emit() {
	emit("lea .%s(%%rip), %%rax", ast.slabel)
}

func (ast *ExprVariable) emit() {
	if ast.isGlobal {
		emit("mov %s(%%rip), %%eax", ast.varname)
	} else {
		emit("mov %d(%%rbp), %%eax", ast.offset)
	}
}

func (ast *ExprBinop) emit() {
	ast.left.emit()
	emit("push %%rax")
	ast.right.emit()
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

func assignLocal(left *ExprVariable, right Expr) {
	assert(!left.isGlobal, "should be a local var")
	right.emit()
	emit("push %%rax")
	emit("mov %%eax, %d(%%rbp)", left.offset)
}

func (ast *AstAssignment) emit() {
	assignLocal(ast.left, ast.right)
}

func emitDeclLocalVar(ast *AstDeclVar) {
	if ast.initval == nil {
		// assign zero value
		ast.initval = &ExprNumberLiteral{}
	}

	assignLocal(ast.variable, ast.initval)
}

func emitCompound(ast *AstCompountStmt) {
	for _, stmt := range ast.stmts {
		if stmt.expr != nil {
			stmt.expr.emit()
		} else if stmt.assignment != nil {
			stmt.assignment.emit()
		} else if stmt.declvar != nil {
			emitDeclLocalVar(stmt.declvar)
		}
	}
}

var regs = []string{"rdi", "rsi", "rdx", "rcx", "r8", "r9"}

func (funcall *ExprFuncall) emit() {
	fname := funcall.fname
	emit("# funcall %s", fname)
	args := funcall.args
	for i, _ := range args {
		emit("push %%%s", regs[i])
	}

	emit("# setting arguments")
	for _, arg := range args {
		arg.emit()
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

func emitGlobalDeclVar(declvar *AstDeclVar) {
	variable := declvar.variable
	assert(variable.isGlobal, "should be global")
	emitLabel(".global %s", variable.varname)
	emitLabel("%s:", variable.varname)
	if declvar.initval == nil {
		// set zero value
		emit(".long %d", 0)
	} else {
		ival, ok := declvar.initval.(*ExprNumberLiteral)
		if !ok {
			errorf("only number can be assign to global variables")
		}
		emit(".long %d", ival.val)
	}
}

func generate(a *AstFile) {
	emitDataSection()
	for _, declvar := range a.decls {
		emitGlobalDeclVar(declvar)
	}
	for _, f := range a.funcdefs {
		emitFuncdef(f)
	}
}
