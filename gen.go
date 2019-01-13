package main

import "fmt"

const INT_SIZE = 8 // not like 8cc

func emit(format string, v ...interface{}) {
	fmt.Printf("\t"+format+"\n", v...)
}

func emitComment(format string, v ...interface{}) {
	fmt.Printf("# "+format+"\n", v...)
}

func emitLabel(format string, v ...interface{}) {
	fmt.Printf(format+"\n", v...)
}

func getMethodUniqueName(gtype *Gtype, fname identifier) string {
	assertNotNil(gtype != nil, nil)
	var typename identifier
	if gtype.typ == G_POINTER {
		typename = gtype.ptr.relation.name
	} else {
		typename = gtype.relation.name
	}
	return string(typename) + "_m_" + string(fname)
}

// main.f1 -> main_p_f1
func getPackagedFuncName(pkg identifier, fname string) string {
	if pkg == "libc" {
		return fname
	}

	return fmt.Sprintf("%s_p_%s", pkg, fname)
}

func (f *DeclFunc) getUniqueName() string {
	if f.receiver != nil {
		// method
		return getPackagedFuncName(f.pkg, getMethodUniqueName(f.receiver.gtype, f.fname))
	}
	// treat main.main function as a special one
	if f.pkg == "main" && f.fname == "main" {
		return "main"
	}

	// other functions
	return getPackagedFuncName(f.pkg, string(f.fname))
}

func (f *DeclFunc) emitPrologue() {
	emitComment("FUNCTION %s", f.getUniqueName())
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
		if lvar.gtype == nil {
			debugf("%s has nil gtype ", lvar)
		}
		size := lvar.gtype.getSize()
		assert(size != 0, nil, "size is not zero")
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

func (a *ExprStructField) emit() {
	rel, ok := a.strct.(*Relation)
	if !ok {
		errorf("struct is not a variable")
	}
	assertNotNil(rel.expr != nil, nil)
	variable, ok := rel.expr.(*ExprVariable)
	assert(ok, nil, "rel is a variable")

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
		emit("mov %d(%%rbp), %%rax", variable.offset+field.offset)
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
	assertNotNil(ast.val != nil, nil)
	rel, ok := ast.val.(*Relation)
	if ok && rel.expr == eIota {
		// replace the iota expr by a index number
		val := &ExprNumberLiteral{
			val: ast.iotaIndex,
		}
		val.emit()
	} else {
		ast.val.emit()
	}
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

var labelSeq = 0

func makeLabel() string {
	r := fmt.Sprintf(".L%d", labelSeq)
	labelSeq++
	return r
}

func (ast *StmtInc) emit() {
	emitIncrDecl("add", ast.operand)
}
func (ast *StmtDec) emit() {
	emitIncrDecl("sub", ast.operand)
}

func emitIncrDecl(inst string, operand Expr) {
	rel, ok := operand.(*Relation)
	if !ok {
		errorf("operand should be *Relation")
	}
	vr, ok := rel.expr.(*ExprVariable)
	assert(ok, nil, "operand is a rel")
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
			assert(e.invisiblevar.offset != 0, nil, "ExprStructLiteral's invisible var has offset")
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
		assert(ok, nil, "operand is a rel")
		vr.emit()
		emit("mov (%%rax), %%rcx")
		emit("mov %%rcx, %%rax")
	} else if ast.op == "!" {
		ast.operand.emit()
		emit("mov $0, %%rcx")
		emit("cmp %%rax, %%rcx")
		emit("sete %%al")
		emit("movzb %%al, %%eax")
	} else if ast.op == "-" {
		// delegate to biop
		// -(x) -> (-1) * (x)
		binop := &ExprBinop{
			op:    "*",
			left:  &ExprNumberLiteral{val: -1},
			right: ast.operand,
		}
		binop.emit()
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

func (ast *StmtAssignment) emit() {
	emit("# start StmtAssignment")
	numLeft := len(ast.lefts)
	numRight := 0
	for _, right := range ast.rights {
		switch right.(type) {
		case *ExprFuncall:
			funcdef := right.(*ExprFuncall).getFuncDef()
			numRight += len(funcdef.rettypes)
		case *ExprMethodcall:
			rettypes := right.(*ExprMethodcall).getRettypes()
			numRight += len(rettypes)
		default:
			numRight++
		}
	}
	if numLeft != numRight {
		errorf("number of exprs does not match")
	}

	var done map[int]bool
	done = make(map[int]bool) // @FIXME this is not correct any more
	for i, right := range ast.rights {
		switch right.(type) {
		case *ExprStructLiteral: // assign struct literal to var
			rel := ast.lefts[i].(*Relation)
			vr := rel.expr.(*ExprVariable)
			assignStructLiteral(vr, right.(*ExprStructLiteral))
			done[i] = true // @FIXME this is not correct any more
		case *ExprFuncall:
			funcdef := right.(*ExprFuncall).getFuncDef()
			emit("# emitting rhs (funcall)")
			right.emit()
			for i, _ := range funcdef.rettypes {
				emit("mov %s(%%rip), %%rax", retvals[i])
				emit("push %%rax")
			}
		case *ExprMethodcall:
			rettypes := right.(*ExprMethodcall).getRettypes()
			emit("# emitting rhs (funcall)")
			right.emit()
			for i, _ := range rettypes {
				emit("mov %s(%%rip), %%rax", retvals[i])
				emit("push %%rax")
			}
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
		case *ExprIndex:
			emit("push %%rax") // push RHS value
			// load head address of the array
			// load index
			// multi index * size
			// calc address = head address + offset
			// copy value to the address
			e := left.(*ExprIndex)
			rel, ok := e.collection.(*Relation)
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
			assert(size > 0, nil, "size > 0")
			emit("mov $%d, %%rax", size) // size of one element
			emit("imul %%rcx, %%rax")    // index * size
			emit("push %%rax")           // store index * size
			emit("pop %%rcx")            // load  index * size
			emit("pop %%rbx")            // load address of variable
			emit("add %%rcx , %%rbx")    // (index * size) + address
			emit("pop %%rax")            // load RHS value
			reg := getReg(size)
			emit("mov %%%s, (%%rbx)", reg) // dereference the content of an emelment
		case *ExprStructField:
			ast, ok := left.(*ExprStructField)
			if !ok {
				errorf("left is not ExprStructField")
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

func (s *StmtIf) emit() {
	emit("# if")
	if s.simplestmt != nil {
		s.simplestmt.emit()
	}
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

func (f *StmtFor) emitRange() {
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

	initstmt := &StmtAssignment{
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
	var assignVar *StmtAssignment
	if f.rng.valuevar != nil {
		assignVar = &StmtAssignment{
			lefts: []Expr{
				f.rng.valuevar,
			},
			rights: []Expr{
				&ExprIndex{
					collection: &Relation{
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
		left:  f.rng.indexvar,                  // i
		right: &ExprNumberLiteral{val: length}, // len(list)
	}
	condition.emit() // i < len(list)
	emit("test %%rax, %%rax")
	emit("je %s  # jump if false", labelEnd)

	f.block.emit()

	indexIncr := &StmtInc{
		operand: f.rng.indexvar,
	}
	indexIncr.emit() // i++
	if f.rng.valuevar != nil {
		assignVar.emit() // v = s[i]
	}
	emit("jmp %s", labelBegin)
	emit("%s: # end loop", labelEnd)
}

func (f *StmtFor) emitForClause() {
	assertNotNil(f.cls != nil, nil)
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

func (f *StmtFor) emit() {
	if f.rng != nil {
		f.emitRange()
		return
	}
	f.emitForClause()
}

func (stmt *StmtReturn) emit() {
	if len(stmt.exprs) == 0 {
		emit("mov $0, %%rax")
		return
	}

	if len(stmt.exprs) > 7 {
		errorf("TBI")
	}

	for i, expr := range stmt.exprs {
		expr.emit()
		emit("mov %%rax, %s(%%rip)", retvals[i])
	}
	emit("leave")
	emit("ret")
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

// for local var
func (decl *DeclVar) emit() {
	if decl.variable.gtype.typ == G_ARRAY && decl.initval != nil {
		// initialize local array
		debugf("initialize local array")
		initvalues, ok := decl.initval.(*ExprArrayLiteral)
		if !ok {
			errorf("error?")
		}
		arraygtype := decl.variable.gtype
		elmType := arraygtype.ptr.relation.gtype
		debugf("gtype:%v", elmType)
		for i, val := range initvalues.values {
			val.emit()
			localoffset := decl.variable.offset + i*elmType.getSize()
			emitLsave(elmType.getSize(), localoffset)
		}
	} else if decl.variable.gtype.relation != nil && decl.variable.gtype.relation.gtype.typ == G_STRUCT && decl.initval != nil {
		// initialize local struct
		debugf("initialize local struct")
		structliteral, ok := decl.initval.(*ExprStructLiteral)
		if !ok {
			errorf("error?")
		}
		assignStructLiteral(decl.variable, structliteral)

	} else {
		debugf("gtype=%v", decl.variable.gtype)
		if decl.initval == nil {
			// assign zero value
			decl.initval = &ExprNumberLiteral{}
		}
		decl.initval.emit()
		emitLsave(decl.variable.gtype.getSize(), decl.variable.offset)
	}
}

func (decl *DeclType) emit() {
	// nothing to do
}

func (decl *DeclConst) emit() {
	// nothing to do
}

func (ast *StmtSatementList) emit() {
	for _, stmt := range ast.stmts {
		stmt.emit()
	}
}

var regs = []string{"rdi", "rsi", "rdx", "rcx", "r8", "r9"}

func (e *ExprIndex) emit() {
	emit("# emit *ExprIndex")
	rel, ok := e.collection.(*Relation)
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
	assert(size > 0, nil, "size > 0")
	emit("mov $%d, %%rax", size) // size of one element
	emit("imul %%rcx, %%rax")    // index * size
	emit("push %%rax")           // store index * size
	emit("pop %%rcx")            // load  index * size
	emit("pop %%rbx")            // load address of variable
	emit("add %%rcx , %%rbx")    // (index * size) + address
	emit("mov (%%rbx), %%rax")   // dereference the content of an emelment
}

func (e *ExprNilLiteral) emit() {
	emit("mov $0, %%rax")
}

func (ast *StmtShortVarDecl) emit() {
	a := &StmtAssignment{
		lefts:  ast.lefts,
		rights: ast.rights,
	}
	a.emit()
}

func (f *ExprFuncRef) emit() {
	emit("mov $1, %%rax") // emit 1 for now.  @FIXME
}

func (e *ExprSliced) emit() {
	errorf("TBD")
}

func (e ExprArrayLiteral) emit() {
	errorf("DO NOT EMIT")
}

func (e *ExprTypeAssertion) emit() {
	panic("implement me")
}

func (ast *StmtContinue) emit() {
	panic("implement me")
}

func (ast *StmtBreak) emit() {
	panic("implement me")
}

func (ast *StmtExpr) emit() {
	ast.expr.emit()
}

func (ast *StmtDefer) emit() {
	panic("implement me")
}

func (e *ExprVaArg) emit() {
	panic("implement me")
}

func (e *ExprConversion) emit() {
	panic("implement me")
}

func (e *ExprStructLiteral) emit() {
	errorf("This cannot be emitted alone")
}

func (e *ExprTypeSwitchGuard) emit() {
	panic("implement me")
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

func (methodCall *ExprMethodcall) getPkgName() identifier {
	origType := methodCall.getOrigType()
	if origType.typ == G_INTERFACE {
		errorf("TBI")
	} else {
		funcref, ok := origType.methods[methodCall.fname]
		if !ok {
			errorf("method %s is not found in type %s", methodCall.fname, methodCall.receiver.getGtype())
		}
		return funcref.funcdef.pkg
	}
	return ""
}

func (methodCall *ExprMethodcall) getOrigType() *Gtype {
	gtype := methodCall.receiver.getGtype()
	assertNotNil(gtype != nil, methodCall.tok)
	assert(gtype.typ == G_REL || gtype.typ == G_POINTER || gtype.typ == G_INTERFACE, methodCall.tok, "method must be an interface or belong to a named type")
	var typeToBeloing *Gtype
	if gtype.typ == G_POINTER {
		typeToBeloing = gtype.ptr
	} else {
		typeToBeloing = gtype
	}
	assert(typeToBeloing.typ == G_REL, methodCall.tok, "method must belong to a named type")
	debugf("typeToBeloing = %s", typeToBeloing)
	origType := typeToBeloing.relation.gtype
	debugf("origType = %v", origType)
	return origType
}


func (methodCall *ExprMethodcall) getRettypes() []*Gtype {
	origType := methodCall.getOrigType()
	if origType.typ == G_INTERFACE {
		return origType.imethods[methodCall.fname].rettypes
	} else {
		funcref, ok := origType.methods[methodCall.fname]
		if !ok {
			errorf("method %s is not found in type %s", methodCall.fname, methodCall.receiver.getGtype())
		}
		return funcref.funcdef.rettypes
	}
}

func (methodCall *ExprMethodcall) emit() {

	args := []Expr{methodCall.receiver}
	for _, arg := range methodCall.args {
		args = append(args, arg)
	}

	pkgname := methodCall.getPkgName()
	name := methodCall.getUniqueName()
	emitCall(getPackagedFuncName(pkgname, name), args)
}

func (funcall *ExprFuncall) getFuncDef() *DeclFunc {
	relexpr := funcall.rel.expr
	assertNotNil(relexpr != nil, funcall.tok)
	funcref, ok := relexpr.(*ExprFuncRef)
	if !ok {
		errorf("Compiler error: funcref is not *ExprFuncRef but %v", funcref, funcall.fname)
	}
	assertNotNil(funcref.funcdef != nil, nil)
	return funcref.funcdef
}

func (funcall *ExprFuncall) emit() {
	decl := funcall.getFuncDef() // check existance
	if decl == nil {
		errorf("funcdef not found for funcall %s, rel=%v ", funcall.fname, funcall.rel)
	}

	emitCall(getPackagedFuncName(decl.pkg, funcall.fname), funcall.args)
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

func (f *DeclFunc) emit() {
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

func (decl *DeclVar) emitGlobal() {
	assert(decl.variable.isGlobal, nil, "should be global")
	assertNotNil(decl.variable.gtype != nil, nil)
	emitLabel(".global %s", decl.variable.varname)
	emitLabel("%s:", decl.variable.varname)

	if decl.variable.gtype.typ == G_ARRAY {
		arrayliteral, ok := decl.initval.(*ExprArrayLiteral)
		assert(ok, nil, "should be array lieteral")
		elmType := decl.variable.gtype.ptr
		assertNotNil(elmType != nil, nil)
		for _, value := range arrayliteral.values {
			assertNotNil(value != nil, nil)
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
		if decl.initval == nil {
			// set zero value
			emit(".quad %d", 0)
		} else {
			var val int
			switch decl.initval.(type) {
			case *ExprNumberLiteral:
				val = decl.initval.(*ExprNumberLiteral).val
			case *ExprConstVariable:
				val = evalIntExpr(decl.initval)
			}
			emit(".quad %d", val)
		}
	}
}

type IrRoot struct {
	vars           []*DeclVar
	funcs          []*DeclFunc
	stringLiterals []*ExprStringLiteral
}

var retvals = []string{"rt1", "rt2", "rt3", "rt4", "rt5", "rt6"}

func (root *IrRoot) emit() {

	// generate code
	emit(".data")

	emitComment("STRING LITERALS")
	for id, ast := range root.stringLiterals {
		ast.slabel = fmt.Sprintf("S%d", id)
		emitLabel(".%s:", ast.slabel)
		emit(".string \"%s\"", ast.val)
	}

	emitComment("GLOBAL RETVALS")
	for _, name := range retvals {
		emitLabel(".global %s", name)
		emitLabel("%s:", name)
		emit(".quad 0")
	}

	emitComment("GLOBAL VARS")
	for _, vardecl := range root.vars {
		vardecl.emitGlobal()
	}

	emitComment("FUNCTIONS")
	for _, funcdecl := range root.funcs {
		funcdecl.emit()
	}
}
