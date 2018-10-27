package main

import "fmt"

const INT_SIZE = 8 // not like 8cc

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
		offset -= INT_SIZE
		param.offset = offset
		emit("push %%%s", regs[i])
	}

	var localarea int
	for _, lvar := range f.localvars {
		if lvar.gtype == nil {
			errorf("lvar has no gtype: %s", lvar)
		}
		assert(lvar.gtype != nil, "lvar has gtype")
		area := calcAreaOfVar(lvar.gtype)
		localarea -= area
		offset -= area
		lvar.offset = offset
		debugf("set offset %d to lvar %s", lvar.offset, lvar.varname)
	}
	if localarea != 0 {
		emit("# allocate localarea")
		emit("sub $%d, %%rsp", -localarea)
	}

	emit("# end of prologue")
}

func calcAreaOfVar(gtype *Gtype) int {
	assert(gtype != nil, "gtype exists")
	if gtype.typ == G_ARRAY {
		return gtype.length * calcAreaOfVar(gtype.ptr)
	} else {
		return INT_SIZE
	}
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
		if ast.gtype.typ == G_ARRAY {
			emit("lea %s(%%rip), %%eax", ast.varname)
		} else {
			emit("mov %s(%%rip), %%eax", ast.varname)
		}
	} else {
		if ast.offset == 0 {
			errorf("offset should not be zero for localvar %s", ast.varname)
		}
		emit("mov %d(%%rbp), %%eax", ast.offset)
	}
}

func (rel *Relation) emit() {
	if rel.expr == nil {
		errorf("rel.expr is nil: %s", rel.name)
	}
	rel.expr.emit()
}

func (ast *ExprConstVariable) emit() {
	ast.val.emit()
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

func assignLocal(left Expr, loff int ,right Expr) {
	right.emit()
	switch left.(type) {
	case *ExprVariable:
		vr := left.(*ExprVariable)
		emit("mov %%eax, %d(%%rbp)", vr.offset + loff)
	case *Relation:
		rel := left.(*Relation)
		assert(rel.expr != nil, "rel links to a variable")
		vr := rel.expr.(*ExprVariable)
		emit("mov %%eax, %d(%%rbp)", vr.offset + loff)
	default:
		errorf("Unexpected type %v", left)
	}
}

func (ast *AstAssignment) emit() {
	assignLocal(ast.left, 0, ast.right)
}

func emitDeclLocalVar(ast *AstVarDecl) {
	if ast.variable.gtype.typ == G_ARRAY &&   ast.initval != nil {
		// initialize local array
		debugf("initialize local array")
		initvalues,ok := ast.initval.(*ExprArrayLiteral)
		if !ok {
			errorf("error?")
		}
		arraygtype :=  ast.variable.gtype
		elmType := arraygtype.ptr.relation.gtype
		debugf("gtype:%v", elmType)
		for i, val := range initvalues.values {
			assignLocal(ast.variable, i * elmType.size, val)
		}
		return
	}

	if ast.initval == nil {
		// assign zero value
		ast.initval = &ExprNumberLiteral{}
	}

	assignLocal(ast.variable,0, ast.initval)
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

func (e *ExprArrayIndex) emit() {
	emit("# emit *ExprArrayIndex")
	variable := e.variable
	var varname identifier
	var vr *ExprVariable
	switch variable.(type) {
	case *ExprVariable:
		vr = variable.(*ExprVariable)
	case *Relation:
		ex := variable.(*Relation).expr
		switch ex.(type) {
		case *ExprVariable:
			vr = ex.(*ExprVariable)
		}
	}
	if vr.isGlobal {
		varname = vr.varname
		emit("lea %s(%%rip), %%rax", varname)
	} else {
		emit("lea %d(%%rbp), %%rax", vr.offset)
	}
	emit("push %%rax")
	e.index.emit()
	emit("mov %%rax, %%rcx")
	elmType := vr.gtype.ptr
	gtype := elmType
	for {
		if gtype.typ == G_REL {
			debugf("gtype is rel. name is %s ", gtype.relname)
			gtype = gtype.relation.gtype
		}else {
			break
		}
	}

	size :=  gtype.size
	if size == 0 {
		errorf("size 0 %v", gtype)
	}
	assert(size > 0, "size > 0")
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
	case nil:
		errorf("e is nil")
	case *ExprNumberLiteral:
		return e.(*ExprNumberLiteral).val
	case *ExprVariable:
		errorf("variable cannot be inteppreted at compile time")
	case *Relation:
		return evalIntExpr(e.(*Relation).expr)
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
		errorf("unkown type %v", e)
	}
	return 0
}

func (gtype *Gtype) getSize() int {
	if gtype.typ == G_REL {
		return gtype.relation.gtype.getSize()
	} else {
		return gtype.size
	}
}

func emitGlobalDeclVar(variable *ExprVariable, initval Expr) {
	assert(variable.isGlobal, "should be global")
	assert(variable.gtype != nil, "variable has gtype")
	emitLabel(".global %s", variable.varname)
	emitLabel("%s:", variable.varname)
	if variable.gtype.typ == G_ARRAY {
		arrayliteral, ok := initval.(*ExprArrayLiteral)
		assert(ok, "should be array lieteral")
		elmType := variable.gtype.ptr
		assert(elmType != nil, "elm is not nil")
		for _, value := range arrayliteral.values {
			assert(value !=nil, "value is set")
			size := elmType.getSize()
			if size == 8 {
				emit(".quad %d", evalIntExpr(value))
			} else if size == 1 {
				emit(".byte %d", evalIntExpr(value))
			} else {
				errorf("Unexpected size %d", size)
			}
		}
	} else {
		if initval == nil {
			// set zero value
			emit(".quad %d", 0)
		} else {
			var val int
			switch initval.(type) {
			case *ExprNumberLiteral:
				val = initval.(*ExprNumberLiteral).val
			case *ExprConstVariable:
				val = evalIntExpr(initval)
			}
			emit(".quad %d", val)
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
