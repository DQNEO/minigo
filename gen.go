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

	// put stringLiterals
	for _, ast := range stringLiterals {
		emitLabel(".%s:", ast.slabel)
		emit(".string \"%s\"", ast.val)
	}
}

func emitFuncPrologue(f *AstFuncDecl) {
	emitLabel("# f %s", f.fname)
	emit(".text")
	emit(".globl	%s", f.fname)
	emitLabel("%s:", f.fname)
	emit("push %%rbp")
	emit("mov %%rsp, %%rbp")

	// calc offset
	var offset int
	for i, param := range f.params {
		offset -= 8
		param.offset = offset
		emit("push %%%s", regs[i])
	}

	var localarea int
	for _, lvar := range f.localvars {
		localarea -= 8
		offset -= 8
		lvar.offset = offset
		debugf("set offset %d to lvar %s", lvar.offset, lvar.varname)
	}
	if localarea != 0 {
		emit("# allocate localarea")
		emit("sub $%d, %%rsp", -localarea)
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
		if ast.offset == 0 {
			errorf("offset should not be zero for localvar %s", ast.varname)
		}
		emit("mov %d(%%rbp), %%eax", ast.offset)
	}
}

func (ast *ExprConstVariable) emit() {
	switch ast.val.(type) {
	case *AstIdentExpr:
		e := ast.val.(*AstIdentExpr)
		e.expr.emit()
	default:
		ast.val.emit()
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

func assignLocal(left Expr, right Expr) {
	//assert(!left.isGlobal, "should be a local var")
	right.emit()
	switch left.(type) {
	case *ExprVariable:
		vr := left.(*ExprVariable)
		emit("push %%rax")
		emit("mov %%eax, %d(%%rbp)", vr.offset)
	default:
		errorf("Unexpected type %v", left)
	}
}

func (ast *AstAssignment) emit() {
	assignLocal(ast.left, ast.right)
}

func emitDeclLocalVar(ast *AstVarDecl) {
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
		} else if stmt.constdecl != nil {
			// nothing to do
		}
	}
}

var regs = []string{"rdi", "rsi", "rdx", "rcx", "r8", "r9"}

func (e *ExprIndexAccess) emit() {
	variable := e.variable
	var varname identifier
	switch variable.(type) {
	case *ExprVariable:
		varname = variable.(*ExprVariable).varname
	case *AstIdentExpr:
		varname = variable.(*AstIdentExpr).name
	}
	emit("lea %s(%%rip), %%rax", varname)
	emit("push %%rax")
	e.index.emit()
	emit("mov %%rax, %%rcx")
	size := 4
	emit("mov $%d, %%rax", size)
	emit("imul %%rcx, %%rax")
	emit("push %%rax")
	emit("pop %%rcx")
	emit("pop %%rbx")
	emit("add %%rcx , %%rbx")
	emit("mov (%%rbx), %%rax")
}

func (funcall *ExprFuncall) emit() {
	fname := funcall.fname
	emit("# funcall %s", fname)
	args := funcall.args
	for i, _ := range args {
		emit("push %%%s", regs[i])
	}

	emit("# setting arguments")
	for i, arg := range args {
		debugf("arg = %v", arg)
		arg.emit()
		emit("push %%rax  # argument no %d", i + 1)
	}

	for i, _ := range args {
		j := len(args) - 1 - i
		emit("pop %%%s   # argument no %d", regs[j], j + 1)
	}
	emit("mov $0, %%rax")
	emit("call %s", fname)

	for i, _ := range args {
		j := len(args) - 1 - i
		emit("pop %%%s", regs[j])
	}
}

func emitFuncdef(f *AstFuncDecl) {
	emitFuncPrologue(f)
	emitCompound(f.body)
	emit("mov $0, %%eax") // return 0
	emitFuncEpilogue()
}

func evalIntExpr(e Expr) int {
	switch e.(type) {
	case *ExprNumberLiteral:
		return e.(*ExprNumberLiteral).val
	case *AstIdentExpr:
		return evalIntExpr(e.(*AstIdentExpr).expr)
	case *ExprBinop:
		binop := e.(*ExprBinop)
		switch binop.op {
		case "+":
			return evalIntExpr(binop.left) + evalIntExpr(binop.right)
		case "-":
			return evalIntExpr(binop.left) - evalIntExpr(binop.right)
		case "*":
			return evalIntExpr(binop.left) * evalIntExpr(binop.right)

		}
	default:
		errorf("unkown type %v to eval", e)
	}
	return 0
}

func emitGlobalDeclVar(variable *ExprVariable, initval Expr) {
	assert(variable.isGlobal, "should be global")
	assert(variable.gtype != nil, "variable has gtype")
	emitLabel(".global %s", variable.varname)
	emitLabel("%s:", variable.varname)
	if variable.gtype.typ == G_ARRAY {
		arrayliteral, ok := initval.(*ExprArrayLiteral)
		assert(ok, "should be array lieteral")
		for _, value := range arrayliteral.values {
			debugPrintV(value)
			emit(".long %d", evalIntExpr(value))
		}
	} else {
		if initval == nil {
			// set zero value
			emit(".long %d", 0)
		} else {
			val := evalIntExpr(initval)
			emit(".long %d", val)
		}
	}
}

func generate(f *AstSourceFile) {
	emitDataSection()
	for _, decl := range f.decls {
		if decl.vardecl != nil {
			emitGlobalDeclVar(decl.vardecl.variable, decl.vardecl.initval)
		} else if decl.funcdecl != nil {
			emitFuncdef(decl.funcdecl)
		}
	}
}
