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

func makeLabel() string {
	r := fmt.Sprintf(".L%d", labelSeq)
	labelSeq++
	return r
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
		emit_comp("setne", ast)
		return
	} else if ast.op == "==" {
		emit_comp("sete", ast)
		return
	} else if ast.op == "&&" {
		labelEnd := makeLabel()
		ast.left.emit()
		emit("test %%rax, %%rax")
		emit("mov $0, %%rax")
		emit("je %s", labelEnd)
		ast.right.emit()
		emit("test %%rax, %%rax")
		emit("mov $0, %%rax")
		emit("je %s", labelEnd)
		emit("mov $1, %%rax")
		emit("%s:", labelEnd)
		return
	} else if ast.op == "||" {
		labelEnd := makeLabel()
		ast.left.emit()
		emit("test %%rax, %%rax")
		emit("mov $1, %%rax")
		emit("jne %s", labelEnd)
		ast.right.emit()
		emit("test %%rax, %%rax")
		emit("mov $1, %%rax")
		emit("jne %s", labelEnd)
		emit("mov $0, %%rax")
		emit("%s:", labelEnd)
		return
	}
	ast.left.emit()
	emit("push %%rax")
	ast.right.emit()
	emit("mov %%rax, %%rcx")
	emit("pop %%rax")
	if ast.op == "+" {
		emit("add	%%rcx, %%rax")
	} else if ast.op == "-" {
		emit("sub	%%rcx, %%rax")
	} else if ast.op == "*" {
		emit("imul	%%rcx, %%rax")
	} else if ast.op == "%" {
		emit("mov $0, %%rdx # init %%rdx")
		emit("div %%rcx")
		emit("mov %%rdx, %%rax")
	} else if ast.op == "/" {
		emit("mov $0, %%rdx # init %%rdx")
		emit("mov $0, %%rdx")
		emit("div %%rcx")
	} else {
		errorf("Unknown binop: %s", ast.op)
	}
}


func (ast *AstAssignment) emit() {
	for _, right := range ast.rights {
		right.emit()
		emit("push %%rax")
	}


	for i := len(ast.lefts) - 1; i >= 0; i-- {
		emit("pop %%rax")
		left := ast.lefts[i]

	switch left.(type) {
	case *Relation:
		// e.g. x = 1
		// resolve relation
		rel := left.(*Relation)
		if rel.expr == nil {
			errorf("left.rel.expr is nil")
		}
		vr := rel.expr.(*ExprVariable)
		emitLsave(vr.gtype.getSize(), vr.offset)
	case *ExprArrayIndex:
		emit("push %%rax") // push RHS value
		e := left.(*ExprArrayIndex)
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
		errorf("Unexpected type %v", ast.lefts)
	}
	}
}


func (s *AstIfStmt) emit() {
	emit("# if")
	s.cond.emit()
	emit("test %%rax, %%rax")
	if s.els != nil {
		labelElse := makeLabel()
		labelEndif := makeLabel()
		emit("je %s  # jump if 0", labelElse)
		emit("# then block")
		s.then.emit()
		emit("jmp %s # jump to endif", labelEndif)
		emit("# else block")
		emit("%s:", labelElse)
		s.els.emit()
		emit("# endif")
		emit("%s:", labelEndif)
	} else {
		// no else block
		labelEndif := makeLabel()
		emit("je %s  # jump if 0", labelEndif)
		emit("# then block")
		s.then.emit()
		emit("# endif")
		emit("%s:", labelEndif)
	}
}

func (f *AstForStmt) emitRange() {
	if f.rng.indexvar == nil {
		errorf("indexVar is nil")
	}

	emit("# for range")
	var length int

	if rel, ok := f.rng.rangeexpr.(*Relation); ok {
		if variable, ok := rel.expr.(*ExprVariable); ok {
			if variable.gtype.typ != G_ARRAY {
				panic("variable should be an array")
			}
			emit("# range expr is %v", variable)
			length = variable.gtype.length
			emit("# length = %d", length)
		} else {
			panic("rel should be a variable")
		}
	} else {
		panic("range expression should be a variable")
	}

	labelBegin := makeLabel()
	labelEnd := makeLabel()

	initstmt := &AstAssignment{
		lefts:[]Expr{
			f.rng.indexvar,
		},
		rights:[]Expr{
			&ExprNumberLiteral{
				val: 0,
			},
		},
	}
	emit("# init index")
	initstmt.emit() // i=0
	var assignVar *AstAssignment
	if f.rng.valuevar != nil {
		assignVar = &AstAssignment{
			lefts: []Expr{
				f.rng.valuevar,
			},
			rights: []Expr{
				&ExprArrayIndex{
					rel: &Relation{
						expr: f.rng.rangeexpr.(*Relation).expr.(*ExprVariable),
					},
					index: f.rng.indexvar,
				},
			},
		}
		assignVar.emit() // v = s[i]
	}

	emit("%s: # begin loop ", labelBegin)
	condition := &ExprBinop{
		op:    "<",
		left:  f.rng.indexvar, // i
		right: &ExprNumberLiteral{length}, // len(list)
	}
	condition.emit() // i < len(list)
	emit("test %%rax, %%rax")
	emit("je %s  # jump if false", labelEnd)

	f.block.emit()

	indexIncr := &AstAssignment{
		lefts:[]Expr{
			f.rng.indexvar, // i =
		},
		rights:[]Expr{
			&ExprBinop{ // @TODO replace by a unary operator
				op: "+",
				left: f.rng.indexvar,
				right: &ExprNumberLiteral{1},
			},
		},
	}
	indexIncr.emit() // i = i + 1
	if f.rng.valuevar != nil {
		assignVar.emit() // v = s[i]
	}
	emit("jmp %s", labelBegin)
	emit("%s: # end loop", labelEnd)
}

func (f *AstForStmt) emitForClause() {
	assert(f.cls != nil , "f.cls must not be nil")
	labelBegin := makeLabel()
	labelEnd := makeLabel()

	if f.cls.initstmt != nil {
		f.cls.initstmt.emit()
	}
	emit("%s: # begin loop ", labelBegin)
	if f.cls.condition != nil {
		f.cls.condition.emit()
		emit("test %%rax, %%rax")
		emit("je %s  # jump if false", labelEnd)
	}
	f.block.emit()
	if f.cls.poststmt != nil {
		f.cls.poststmt.emit()
	}
	emit("jmp %s", labelBegin)
	emit("%s: # end loop", labelEnd)
}

func (f *AstForStmt) emit() {
	if f.rng != nil {
		f.emitRange()
		return
	}
	f.emitForClause()
}


func (stmt AstReturnStmt) emit() {
	if stmt.expr == nil {
		emit("mov $0, %%rax")
	} else {
		stmt.expr.emit()
		emit("leave")
		emit("ret")
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

func (ast *AstVarDecl) emit() {
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

func (decl *AstConstDecl) emit() {
	// nothing to do
}

func (decl *AstTypeDecl) emit() {
	errorf("TBD")
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
