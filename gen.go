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
		assert(lvar.gtype != nil, "lvar has gtype")
		size := lvar.gtype.getSize()
		assert(size != 0, "size is not zero")
		loff := align(size, 8)
		localarea -= loff
		offset -= loff
		lvar.offset = offset
		debugf("set offset %d to lvar %s", lvar.offset, lvar.varname)
	}
	if localarea != 0 {
		emit("sub $%d, %%rsp # allocate localarea", -localarea)
	}

	emit("# end of prologue")
}

func align(n int, m int) int {
	remainder := n % m
	if remainder == 0 {
		return n
	} else {
		return n - remainder + m
	}
}

func emitFuncEpilogue() {
	emit("#")
	emit("leave")
	emit("ret")
}

func (ast *ExprNumberLiteral) emit() {
	emit("mov	$%d, %%rax", ast.val)
}

func (ast *ExprStringLiteral) emit() {
	emit("lea .%s(%%rip), %%rax", ast.slabel)
}

func (ast *ExprVariable) emit() {
	if ast.isGlobal {
		if ast.gtype.typ == G_ARRAY {
			emit("lea %s(%%rip), %%rax", ast.varname)
		} else {
			emit("mov %s(%%rip), %%rax", ast.varname)
		}
	} else {
		if ast.offset == 0 {
			errorf("offset should not be zero for localvar %s", ast.varname)
		}
		if ast.gtype.typ == G_ARRAY {
			emit("lea %d(%%rbp), %%rax", ast.offset)
		} else {
			emit("mov %d(%%rbp), %%rax", ast.offset)
		}
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

func emit_comp(inst string, ast *ExprBinop) {
	ast.left.emit()
	emit("push %%rax")
	ast.right.emit()
	emit("pop %%rcx")
	emit("cmp %%rax, %%rcx")
	emit("%s %%al", inst)
	emit("movzb %%al, %%eax")
}

func (ast *ExprBinop) emit() {
	if ast.op == "<" {
		emit_comp("setl", ast)
		return
	} else if ast.op == ">" {
		emit_comp("setg", ast)
		return
	} else if ast.op == "<=" {
		emit_comp("setle", ast)
		return
	} else if ast.op == ">=" {
		emit_comp("setge", ast)
		return
	} else if ast.op == "!=" {
		errorf("TBD")
	} else if ast.op == "==" {
		emit_comp("sete", ast)
		return
	}
	ast.left.emit()
	emit("push %%rax")
	ast.right.emit()
	emit("push %%rax")
	emit("pop %%rbx")
	emit("pop %%rax")
	if ast.op == "+" {
		emit("add	%%rbx, %%rax")
	} else if ast.op == "-" {
		emit("sub	%%rbx, %%rax")
	} else if ast.op == "*" {
		emit("imul	%%rbx, %%rax")
	} else {
		errorf("Unknown binop: %s", ast.op)
	}
}


func (ast *AstAssignment) emit() {
	ast.right.emit()
	switch ast.left.(type) {
	case *Relation:
		// e.g. x = 1
		// resolve relation
		rel := ast.left.(*Relation)
		vr := rel.expr.(*ExprVariable)
		emitLsave(vr.gtype.getSize(), vr.offset)
	case *ExprArrayIndex:
		emit("push %%rax") // push RHS value
		e := ast.left.(*ExprArrayIndex)
		// load head address of the array
		// load index
		// multi index * size
		// calc address = head address + offset
		// copy value to the address
		vr := e.rel.expr.(*ExprVariable)
		if vr.isGlobal {
			emit("lea %s(%%rip), %%rax", vr.varname)
		} else {
			emit("lea %d(%%rbp), %%rax", vr.offset)
		}
		emit("push %%rax") // store address of variable
		e.index.emit()
		emit("mov %%rax, %%rcx") // index
		elmType := vr.gtype.ptr
		size := elmType.getSize()
		assert(size > 0, "size > 0")
		emit("mov $%d, %%rax", size) // size of one element
		emit("imul %%rcx, %%rax")    // index * size
		emit("push %%rax")           // store index * size
		emit("pop %%rcx")            // load  index * size
		emit("pop %%rbx")            // load address of variable
		emit("add %%rcx , %%rbx")    // (index * size) + address
		emit("pop %%rax")            // load RHS value
		reg := getReg(size)
		emit("mov %%%s, (%%rbx)", reg)   // dereference the content of an emelment

	default:
		errorf("Unexpected type %v", ast.left)
	}
}


func (s *AstIfStmt) emit() {
	emit("# if")
	s.cond.emit()
	emit("test %%rax, %%rax")
	if s.els != nil {
		labelElse := fmt.Sprintf("L%d", labelNo)
		labelNo++
		labelEndif := fmt.Sprintf("L%d", labelNo)
		labelNo++
		emit("je .%s  # jump if 0", labelElse)
		emit("# then block")
		s.then.emit()
		emit("jmp .%s # jump to endif", labelEndif)
		emit("# else block")
		emit(".%s:", labelElse)
		s.els.emit()
		emit("# endif")
		emit(".%s:", labelEndif)
	} else {
		// no else block
		labelEndif := fmt.Sprintf("L%d", labelNo)
		labelNo++
		emit("je .%s  # jump if 0", labelEndif)
		emit("# then block")
		s.then.emit()
		emit("# endif")
		emit(".%s:", labelEndif)
	}
}

func (f *AstForStmt) emit() {
	labelBegin := fmt.Sprintf("L%d", labelNo)
	labelNo++
	labelEnd := fmt.Sprintf("L%d", labelNo)
	labelNo++

	if f.left != nil {
		f.left.emit()
	}
	emit(".%s: # begin loop ", labelBegin)
	if f.middle != nil {
		f.middle.emit()
		emit("test %%rax, %%rax")
		emit("je .%s  # jump if false", labelEnd)
	}
	f.block.emit()
	if f.right != nil {
		f.right.emit()
	}
	emit("jmp .%s", labelBegin)
	emit(".%s: # end loop", labelEnd)
}


func (stmt AstReturnStmt) emit() {
	if stmt.expr == nil {
		emit("mov $0, %%rax")
	} else {
		stmt.expr.emit()
	}
}

func getReg(regSize int) string {
	var reg string
	switch regSize {
	case 1:
		reg = "al"
	case 8:
		reg = "rax"
	default:
		errorf("Unexpected reg size %d", regSize)

	}
	return reg
}

func emitLsave(regSize int, loff int) {
	reg := getReg(regSize)
	emit("mov %%%s, %d(%%rbp)", reg, loff)
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
			val.emit()
			localoffset := ast.variable.offset + i * elmType.getSize()
			emitLsave(elmType.getSize(), localoffset)
		}
	} else {
		if ast.initval == nil {
			// assign zero value
			ast.initval = &ExprNumberLiteral{}
		}
		ast.initval.emit()
		emitLsave(ast.variable.gtype.getSize(), ast.variable.offset)
	}
}

func (stmt *AstStmt) emit() {
	if stmt.expr != nil {
		stmt.expr.emit()
	} else if stmt.assignment != nil {
		stmt.assignment.emit()
	} else if stmt.ifstmt != nil {
		stmt.ifstmt.emit()
	} else if stmt.forstmt != nil {
		stmt.forstmt.emit()
	} else if stmt.rtrnstmt != nil {
		stmt.rtrnstmt.emit()
	} else if stmt.declvar != nil {
		emitDeclLocalVar(stmt.declvar)
	} else if stmt.constdecl != nil {
		// nothing to do
	} else if stmt.compound != nil {
		stmt.compound.emit()
	} else {
		errorf("Unknown statement: %s", stmt)
	}
}

func (ast *AstCompountStmt) emit() {
	for _, stmt := range ast.stmts {
		stmt.emit()
	}
}

var regs = []string{"rdi", "rsi", "rdx", "rcx", "r8", "r9"}

func (e *ExprArrayIndex) emit() {
	emit("# emit *ExprArrayIndex")
	vr := e.rel.expr.(*ExprVariable)
	if vr.isGlobal {
		emit("lea %s(%%rip), %%rax", vr.varname)
	} else {
		emit("lea %d(%%rbp), %%rax", vr.offset)
	}
	emit("push %%rax") // store address of variable
	e.index.emit()
	emit("mov %%rax, %%rcx") // index
	elmType := vr.gtype.ptr
	size :=  elmType.getSize()
	assert(size > 0, "size > 0")
	emit("mov $%d, %%rax", size) // size of one element
	emit("imul %%rcx, %%rax")    // index * size
	emit("push %%rax")           // store index * size
	emit("pop %%rcx")            // load  index * size
	emit("pop %%rbx")            // load address of variable
	emit("add %%rcx , %%rbx")    // (index * size) + address
	emit("mov (%%rbx), %%rax")   // dereference the content of an emelment
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
	f.body.emit()
	emit("mov $0, %%rax")
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
