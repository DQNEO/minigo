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

func getMethodUniqueName(gtype *Gtype, fname identifier) string {
	assert(gtype != nil, "gtype is not nil")
	var typename identifier
	if gtype.typ == G_POINTER {
		typename = gtype.ptr.relation.name
	} else {
		typename = gtype.relation.name
	}
	return string(typename) + "__xx__" + string(fname)
}

func (f *AstFuncDecl) getUniqueName() string {
	if f.receiver != nil {
		// method
		return getMethodUniqueName(f.receiver.gtype, f.fname)
	}
	// function
	return string(f.fname)
}

func (f *AstFuncDecl) emitPrologue() {
	emitLabel("# func %s", f.getUniqueName())
	emit(".text")
	emit(".globl	%s", f.getUniqueName())
	emitLabel("%s:", f.getUniqueName())
	emit("push %%rbp")
	emit("mov %%rsp, %%rbp")

	// calc offset
	var offset int
	var params []*ExprVariable
	if f.receiver != nil {
		params = []*ExprVariable{f.receiver}
		for _, param := range f.params {
			params = append(params, param)
		}
	} else {
		params = f.params
	}

	for i, param := range params {
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
		debugf("set offset %d to lvar %s, type=%s", lvar.offset, lvar.varname, lvar.gtype)
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

func (a *AstStructFieldAccess) emit() {
	rel, ok := a.strct.(*Relation)
	if !ok {
		errorf("struct is not a variable")
	}
	assert(rel.expr != nil, "rel is a variable")
	variable, ok := rel.expr.(*ExprVariable)
	assert(ok, "rel is a variable")

	switch variable.gtype.typ {
	case G_POINTER: // pointer to struct
		strcttype := variable.gtype.ptr.relation.gtype
		field := strcttype.getField(a.fieldname)
		variable.emit()
		emit("add $%d, %%rax", field.offset)
		emit("mov (%%rax), %%rax")
	case G_REL: // struct
		strcttype := variable.gtype.relation.gtype
		field := strcttype.getField(a.fieldname)
		emit("mov %d(%%rbp), %%rax", variable.offset + field.offset)
	default:
		errorf("internal error: bad gtype %d", variable.gtype.typ)
	}
}

func (ast *ExprVariable) emit() {
	if ast.gtype.typ == G_ARRAY {
		ast.emitAddress()
		return
	}
	if ast.isGlobal {
		emit("mov %s(%%rip), %%rax", ast.varname)
	} else {
		if ast.offset == 0 {
			errorf("offset should not be zero for localvar %s", ast.varname)
		}
		emit("mov %d(%%rbp), %%rax", ast.offset)
	}
}

func (ast *ExprVariable) emitAddress() {
	if ast.isGlobal {
		emit("lea %s(%%rip), %%rax", ast.varname)
	} else {
		if ast.offset == 0 {
			errorf("offset should not be zero for localvar %s", ast.varname)
		}
		emit("lea %d(%%rbp), %%rax", ast.offset)
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

func (ast *AstIncrStmt) emit() {
	emitIncrDecl("add", ast.operand)
}
func (ast *AstDecrStmt) emit() {
	emitIncrDecl("sub", ast.operand)
}

func emitIncrDecl(inst string, operand Expr) {
	rel, ok := operand.(*Relation)
	if !ok {
		errorf("operand should be *Relation")
	}
	vr, ok := rel.expr.(*ExprVariable)
	assert(ok, "operand is a rel")
	vr.emit()
	emit("%s $1, %%rax", inst)
	emitLsave(vr.gtype.getSize(), vr.offset)
}

func (ast *ExprUop) emit() {
	debugf("emitting ExprUop")
	if ast.op == "&" {
		switch ast.operand.(type) {
		case *Relation:
			rel := ast.operand.(*Relation)
			vr, ok := rel.expr.(*ExprVariable)
			if !ok {
				errorf("rel is not an variable")
			}
			vr.emitAddress()
		case *ExprStructLiteral:
			e := ast.operand.(*ExprStructLiteral)
			assert(e.invisiblevar.offset != 0, "ExprStructLiteral's invisible var has offset")
			assignStructLiteral(e.invisiblevar, e)
			emit("lea %d(%%rbp), %%rax", e.invisiblevar.offset)
		default:
			errorf("Unknown type: %s", ast.operand)
		}
	} else if ast.op == "*" {
		// dereferene of a pointer
		debugf("dereferene of a pointer")
		rel, ok := ast.operand.(*Relation)
		debugf("operand:%s", rel)
		vr, ok := rel.expr.(*ExprVariable)
		assert(ok, "operand is a rel")
		vr.emit()
		emit("mov (%%rax), %%rcx")
		emit("mov %%rcx, %%rax")
	} else if ast.op == "!" {
		ast.operand.emit()
		emit("mov $0, %%rcx")
		emit("cmp %%rax, %%rcx")
		emit("sete %%al")
		emit("movzb %%al, %%eax")
	} else {
		errorf("unable to handle uop %s", ast.op)
	}
	debugf("end of emitting ExprUop")

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
	emit("# start AstAssignment")
	done := make(map[int]bool)
	for i, right := range ast.rights {
		switch right.(type) {
		case *ExprStructLiteral: // assign struct literal to var
			rel := ast.lefts[i].(*Relation)
			vr := rel.expr.(*ExprVariable)
			assignStructLiteral(vr, right.(*ExprStructLiteral))
			done[i] = true
		default:
			emit("# emitting rhs")
			right.emit()
			emit("push %%rax")
		}
	}

	for i := len(ast.lefts) - 1; i >= 0; i-- {
		if done[i] {
			continue
		}
		emit("# assigning to lhs")
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
			if vr.isGlobal {
				emitGsave(vr.gtype.getSize(), vr.varname)
			} else {
				emitLsave(vr.gtype.getSize(), vr.offset)
			}
		case *ExprArrayIndex:
			emit("push %%rax") // push RHS value
			// load head address of the array
			// load index
			// multi index * size
			// calc address = head address + offset
			// copy value to the address
			e := left.(*ExprArrayIndex)
			rel, ok := e.array.(*Relation)
			if !ok {
				errorf("should be array variable. array expr is not supported yet")
			}

			vr, ok := rel.expr.(*ExprVariable)
			if !ok {
				errorf("should be array variable. ")
			}
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
			emit("mov %%%s, (%%rbx)", reg) // dereference the content of an emelment
		case *AstStructFieldAccess:
			ast, ok := left.(*AstStructFieldAccess)
			if !ok {
				errorf("left is not AstStructFieldAccess")
			}
			rel := ast.strct.(*Relation)
			vr := rel.expr.(*ExprVariable)
			field := vr.gtype.relation.gtype.getField(ast.fieldname)
			emitLsave(field.getSize(), vr.offset+field.offset)
		default:
			left.dump()
			errorf("Unknown case")
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
		lefts: []Expr{
			f.rng.indexvar,
		},
		rights: []Expr{
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
					array: &Relation{
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
		left:  f.rng.indexvar,             // i
		right: &ExprNumberLiteral{length}, // len(list)
	}
	condition.emit() // i < len(list)
	emit("test %%rax, %%rax")
	emit("je %s  # jump if false", labelEnd)

	f.block.emit()

	indexIncr := &AstIncrStmt{
		operand:f.rng.indexvar,
	}
	indexIncr.emit() // i++
	if f.rng.valuevar != nil {
		assignVar.emit() // v = s[i]
	}
	emit("jmp %s", labelBegin)
	emit("%s: # end loop", labelEnd)
}

func (f *AstForStmt) emitForClause() {
	assert(f.cls != nil, "f.cls must not be nil")
	labelBegin := makeLabel()
	labelEnd := makeLabel()

	if f.cls.init != nil {
		f.cls.init.emit()
	}
	emit("%s: # begin loop ", labelBegin)
	if f.cls.cond != nil {
		f.cls.cond.emit()
		emit("test %%rax, %%rax")
		emit("je %s  # jump if false", labelEnd)
	}
	f.block.emit()
	if f.cls.post != nil {
		f.cls.post.emit()
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

func (stmt *AstReturnStmt) emit() {
	if len(stmt.exprs) == 0 {
		emit("mov $0, %%rax")
	} else {
		stmt.exprs[0].emit()
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

func emitGsave(regSize int, varname identifier) {
	reg := getReg(regSize)
	emit("mov %%%s, %s(%%rip)", reg, varname)
}

func assignStructLiteral(variable *ExprVariable, structliteral *ExprStructLiteral) {
	strcttyp := structliteral.strctname.gtype
	// do assignment for each field
	for _, field := range structliteral.fields {
		field.value.emit()
		fieldtype := strcttyp.getField(field.key)
		localoffset := variable.offset + fieldtype.offset
		regSize := fieldtype.relation.gtype.getSize()
		emitLsave(regSize, localoffset)
	}

}

func (ast *AstVarDecl) emit() {
	if ast.variable.gtype.typ == G_ARRAY && ast.initval != nil {
		// initialize local array
		debugf("initialize local array")
		initvalues, ok := ast.initval.(*ExprArrayLiteral)
		if !ok {
			errorf("error?")
		}
		arraygtype := ast.variable.gtype
		elmType := arraygtype.ptr.relation.gtype
		debugf("gtype:%v", elmType)
		for i, val := range initvalues.values {
			val.emit()
			localoffset := ast.variable.offset + i*elmType.getSize()
			emitLsave(elmType.getSize(), localoffset)
		}
	} else if ast.variable.gtype.relation != nil && ast.variable.gtype.relation.gtype.typ == G_STRUCT && ast.initval != nil {
		// initialize local struct
		debugf("initialize local struct")
		structliteral, ok := ast.initval.(*ExprStructLiteral)
		if !ok {
			errorf("error?")
		}
		assignStructLiteral(ast.variable, structliteral)

	} else {
		debugf("gtype=%v", ast.variable.gtype)
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
	rel,ok := e.array.(*Relation)
	if !ok {
		errorf("array should be a Relation")
	}
	vr, ok := rel.expr.(*ExprVariable)
	if !ok {
		errorf("array should be a variable")
	}
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
	emit("mov (%%rbx), %%rax")   // dereference the content of an emelment
}

func (ast *ExprMethodcall) getUniqueName() string {
	var gtype *Gtype

	switch ast.receiver.(type) {
	case *Relation:
		rel := ast.receiver.(*Relation)
		if vr, ok := rel.expr.(*ExprVariable); ok {
			gtype = vr.gtype
			if gtype.typ == G_REL && gtype.relation.gtype.typ == G_INTERFACE {
				errorf("interface method call is not supported yet. (%s.%s)", gtype.relation.name, ast.fname)
			}
		} else {
			// @TODO must adapt to method chains like foo.Bar().Buz()
			errorf("internal error")
		}
	default:
		errorf("internal error")
	}
	debugf("ast.receiver=%v", ast.receiver)
	debugf("gtype=%v", gtype)
	return getMethodUniqueName(gtype, ast.fname)
}

func (methodCall *ExprMethodcall) emit() {
	args := []Expr{methodCall.receiver}
	for _, arg := range methodCall.args {
		args = append(args, arg)
	}

	name := methodCall.getUniqueName()
	emitCall(name, args)
}

func (funcall *ExprFuncall) emit() {
	emitCall(funcall.fname, funcall.args)
}

func emitCall(fname string, args []Expr) {

	emit("# funcall %s", fname)
	for i, _ := range args {
		emit("push %%%s", regs[i])
	}

	emit("# setting arguments")
	for i, arg := range args {
		debugf("arg[%d] = %v", i, arg)
		arg.emit()
		emit("push %%rax  # argument no %d", i+1)
	}

	for i, _ := range args {
		j := len(args) - 1 - i
		emit("pop %%%s   # argument no %d", regs[j], j+1)
	}
	emit("mov $0, %%rax")
	emit("call %s", fname)

	for i, _ := range args {
		j := len(args) - 1 - i
		emit("pop %%%s", regs[j])
	}
}

func (f *AstFuncDecl) emit() {
	f.emitPrologue()
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
			assert(value != nil, "value is set")
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
			decl.funcdecl.emit()
		}
	}
}
