package main

import "fmt"

var RegsForCall = [...]string{"rdi", "rsi", "rdx", "rcx", "r8", "r9"}

const INT_SIZE = 8 // not like 8cc

var hiddenArrayId = 1;

func getHidddenArrayId() int {
	r := hiddenArrayId
	hiddenArrayId++
	return r
}

func emit(format string, v ...interface{}) {
	fmt.Printf("\t"+format+"\n", v...)
}

func emitComment(format string, v ...interface{}) {
	fmt.Printf("/* "+format+" */\n", v...)
}

func emitLabel(format string, v ...interface{}) {
	fmt.Printf(format+"\n", v...)
}

func getMethodUniqueName(gtype *Gtype, fname identifier) string {
	assertNotNil(gtype != nil, nil)
	var typename identifier
	if gtype.typ == G_POINTER {
		typename = gtype.origType.relation.name
	} else {
		typename = gtype.relation.name
	}
	return string(typename) + "." + string(fname)
}

// main.f1 -> main.f1
func getPackagedFuncName(pkg identifier, fname string) string {
	if pkg == "libc" {
		return fname
	}

	return fmt.Sprintf("%s.%s", pkg, fname)
}

func (f *DeclFunc) getUniqueName() string {
	if f.receiver != nil {
		// method
		return getPackagedFuncName(f.pkg, getMethodUniqueName(f.receiver.gtype, f.fname))
	}

	// main.main => main
	if f.isMainMain {
		return "main"
	}

	// other functions
	return getPackagedFuncName(f.pkg, string(f.fname))
}

func (f *DeclFunc) emitPrologue() {
	if f.isMainMain {
		emit("# main.main")
	}
	uniquName := f.getUniqueName()
	//emitComment("func %s.%s()", f.pkg, f.fname)
	if f.getUniqueName() == "main" {
		emit(".global	%s", f.getUniqueName())
	}
	emitLabel("%s:", uniquName)
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
		if i == 0 {
			emit("# Allocating stack for params")
		}
		offset -= INT_SIZE
		param.offset = offset
		emit("push %%%s", RegsForCall[i])
	}

	var localarea int
	for i, lvar := range f.localvars {
		if i == 0 {
			emit("# Allocating stack for localvars")
		}
		if lvar.gtype == nil {
			debugf("%s has nil gtype ", lvar)
		}
		size := lvar.gtype.getSize()
		assert(size != 0, lvar.token(), "size should  not be zero:" + lvar.gtype.String())
		loff := align(size, 8)
		localarea -= loff
		offset -= loff
		lvar.offset = offset
		emit("# offset %d for variable \"%s\" of %s", offset, lvar.varname, lvar.gtype)
		//debugf("set offset %d to lvar %s, type=%s", lvar.offset, lvar.varname, lvar.gtype)
	}
	if localarea != 0 {
		emit("sub $%d, %%rsp # total stack size", -localarea)
	}

	emit("")
	if f.isMainMain {
		if importOS {
			emit("# set Args")
			emit("mov %%rsi, Args(%%rip)")       // set pointer (**argv)
			emit("mov %%rdi, 8+Args(%%rip)")     // set len (argc)
			emit("mov %%rdi, 16+Args(%%rip)")     // set cap (argc)
		}
	}
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
	emit("")
	emit("leave")
	emit("ret")
	emit("")
}

func (ast *ExprNumberLiteral) emit() {
	emit("mov	$%d, %%rax", ast.val)
}

func (ast *ExprStringLiteral) emit() {
	emit("lea .%s(%%rip), %%rax", ast.slabel)
}

func loadStructField(strct Expr, field *Gtype, offset int) {
	emit("# loadStructField")
	switch strct.(type) {
	case *Relation:
		rel := strct.(*Relation)
		assertNotNil(rel.expr != nil, nil)
		variable, ok := rel.expr.(*ExprVariable)
		assert(ok, nil, "rel is a variable")
		if field.typ == G_ARRAY {
			emit("lea %d(%%rbp), %%rax", variable.offset+field.offset)
		} else {
			if variable.isGlobal {
				emit("mov %s+%d(%%rip), %%rax # ", variable.varname, field.offset + offset)
			} else {
				emit("mov %d(%%rbp), %%rax", variable.offset + field.offset + offset)
			}
		}
	case *ExprStructField: // strct.field.field
		a := strct.(*ExprStructField)
		strcttype := a.strct.getGtype().relation.gtype
		assert(strcttype.size > 0, a.token(), "struct size should be > 0")
		field2 := strcttype.getField(a.fieldname)
		loadStructField(a.strct, field2, offset + field.offset)
	case *ExprIndex: // array[1].field
		indexExpr := strct.(*ExprIndex)
		loadCollectIndex(indexExpr.collection, indexExpr.index, offset + field.offset)
	default:
		// funcall().field
		// methodcall().field
		// *ptr.field
		// (MyStruct{}).field
		// (&MyStruct{}).field
		TBI(strct.token(), "unable to handle %T", strct)
	}

}

func (a *ExprStructField) emit() {
	emit("# ExprStructField.emit()")
	switch a.strct.getGtype().typ {
	case G_POINTER: // pointer to struct
		strcttype := a.strct.getGtype().origType.relation.gtype
		field := strcttype.getField(a.fieldname)
		a.strct.emit()
		emit("add $%d, %%rax", field.offset)
		emit("mov (%%rax), %%rax")
	case G_REL: // struct
		strcttype := a.strct.getGtype().relation.gtype
		assert(strcttype.size > 0, a.token(), "struct size should be > 0")
		field := strcttype.getField(a.fieldname)
		loadStructField(a.strct, field, 0)
	default:
		errorft(a.token(), "internal error: bad gtype %s", a.strct.getGtype())
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
			errorft(ast.token(),"offset should not be zero for localvar %s", ast.varname)
		}
		emit("mov %d(%%rbp), %%rax", ast.offset)
	}
}

func (vr *ExprVariable) emitAddress() {
	if vr.isGlobal {
		emit("lea %s(%%rip), %%rax", vr.varname)
	} else {
		if vr.offset == 0 {
			errorft(vr.token(), "offset should not be zero for localvar %s", vr.varname)
		}
		emit("lea %d(%%rbp), %%rax", vr.offset)
	}
}

func (rel *Relation) emit() {
	if rel.expr == nil {
		errorft(rel.token(),"rel.expr is nil: %s", rel.name)
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

// https://golang.org/ref/spec#IncDecStmt
// As with an assignment, the operand must be addressable or a map index expression.
func emitIncrDecl(inst string, operand Expr) {
	operand.emit()
	emit("%s $1, %%rax", inst)

	left := operand
	emitSave(left)
}

// e.g. *x = 1, or *x++
func (uop *ExprUop) emitSave() {
	assert(uop.op == "*", uop.tok, "uop op should be *")
	emit("push %%rax")
	uop.operand.emit()
	emit("mov %%rax, %%rcx")
	emit("pop %%rax")
	reg := getReg(uop.operand.getGtype().getSize())
	emit("mov %%%s, (%%rcx)", reg)

}

// e.g. x = 1
func (rel *Relation) emitSave() {
	if rel.expr == nil {
		errorft(rel.token(), "left.rel.expr is nil")
	}
	vr := rel.expr.(*ExprVariable)
	assert(0 <= vr.gtype.getSize() && vr.gtype.getSize() <= 8, rel.token(), "invalid size")
	if vr.isGlobal {
		emitGsave(vr.gtype.getSize(), vr.varname, 0)
	} else {
		emitLsave(vr.gtype.getSize(), vr.offset)
	}
}

func (ast *ExprUop) emit() {
	//debugf("emitting ExprUop")
	if ast.op == "&" {
		switch ast.operand.(type) {
		case *Relation:
			rel := ast.operand.(*Relation)
			vr, ok := rel.expr.(*ExprVariable)
			if !ok {
				errorft(ast.token(), "rel is not an variable")
			}
			vr.emitAddress()
		case *ExprStructLiteral:
			e := ast.operand.(*ExprStructLiteral)
			assert(e.invisiblevar.offset != 0, nil, "ExprStructLiteral's invisible var has offset")
			assignToStruct(e.invisiblevar, e)
			emit("lea %d(%%rbp), %%rax", e.invisiblevar.offset)
		default:
			errorft(ast.token(), "Unknown type: %s", ast.operand)
		}
	} else if ast.op == "*" {
		// dereferene of a pointer
		//debugf("dereferene of a pointer")
		rel, ok := ast.operand.(*Relation)
		//debugf("operand:%s", rel)
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
		errorft(ast.token(), "unable to handle uop %s", ast.op)
	}
	//debugf("end of emitting ExprUop")

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
		errorft(ast.token(), "Unknown binop: %s", ast.op)
	}
}

func (ast *StmtAssignment) emit() {
	emit("# start StmtAssignment")
	// the right hand operand is a single multi-valued expression
	// such as a function call, a channel or map operation, or a type assertion.
	// The number of operands on the left hand side must match the number of values.
	singleValueMode := (len(ast.rights) > 1)

	numLeft := len(ast.lefts)
	numRight := 0
	for _, right := range ast.rights {
		switch right.(type) {
		case *ExprFuncallOrConversion:
			rettypes := right.(*ExprFuncallOrConversion).getFuncDef().rettypes
			if singleValueMode && len(rettypes) > 1 {
				errorft(ast.token(), "multivalue is not allowed")
			}
			numRight += len(rettypes)
		case *ExprMethodcall:
			rettypes := right.(*ExprMethodcall).getRettypes()
			if singleValueMode && len(rettypes) > 1 {
				errorft(ast.token(), "multivalue is not allowed")
			}
			numRight += len(rettypes)
		default:
			numRight++
		}
	}
	if numLeft != numRight {
		errorft(ast.token(), "number of exprs does not match")
	}

	for rightIndex, right := range ast.rights {
		left := ast.lefts[rightIndex]
		switch right.(type) {
		case *ExprFuncallOrConversion:
			rettypes := right.(*ExprFuncallOrConversion).getFuncDef().rettypes
			emit("# emitting rhs (funcall)")
			if len(rettypes) == 1 {
				right.emit()
				emitSave(left)
			} else {
				right.emit()
				for i, _ := range rettypes {
					emit("mov %s(%%rip), %%rax", retvals[len(rettypes)-1-i])
					emit("push %%rax")
				}
				for _, left := range ast.lefts {
					emit("pop %%rax")
					emitSave(left)
				}
			}
		case *ExprMethodcall:
			rettypes := right.(*ExprMethodcall).getRettypes()
			emit("# emitting rhs (method)")
			if len(rettypes) == 1 {
				right.emit()
				emitSave(left)
			} else {
				right.emit()
				for i, _ := range rettypes {
					emit("mov %s(%%rip), %%rax", retvals[len(rettypes)-1-i])
					emit("push %%rax")
				}
				for _, left := range ast.lefts {
					emit("pop %%rax")
					emitSave(left)
				}
			}
		default:
			gtype := right.getGtype()
			switch {
			case gtype.typ == G_ARRAY:
				assignToArray(left, right)
			case gtype.typ == G_SLICE:
				assignToSlice(left, right)
			case gtype.typ == G_REL && gtype.relation.gtype.typ == G_STRUCT:
				assignToStruct(left, right)
			default:
				// suppose primitive
				emit("# evaluating rhs")
				right.emit()
				emit("# saving it to lhs")
				emitSave(left)
			}
		}
	}

}

// Each left-hand side operand must be addressable,
// a map index expression,
// or (for = assignments only) the blank identifier.
func emitSave(left Expr) {
	switch left.(type) {
	case *Relation:
		left.(*Relation).emitSave()
	case *ExprIndex:
		left.(*ExprIndex).emitSave()
	case *ExprStructField:
		left.(*ExprStructField).emitLsave()
	case *ExprUop:
		left.(*ExprUop).emitSave()
	default:
		left.dump()
		errorft(left.token(), "Unknown case %T", left)
	}
}

func (e *ExprIndex) emitSave() {
	emit("push %%rax") // push RHS value

	// load head address of the array
	// load index
	// multi index * size
	// calc address = head address + offset
	// copy value to the address

	collectionType := e.collection.getGtype()
	switch {
	case collectionType.typ == G_ARRAY, collectionType.typ == G_SLICE:
		e.collection.emit() // head address
	default:
		TBI(e.token(), "unable to handle %s", collectionType)
	}

	emit("push %%rax") // stash head address of collection
	e.index.emit()
	emit("mov %%rax, %%rcx") // index
	elmType := collectionType.elementType
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
}

func (e *ExprStructField) getOffset() int {
	var vr *ExprVariable
	if rel, ok := e.strct.(*Relation); ok {
		vr, ok = rel.expr.(*ExprVariable)
		assert(ok, e.tok, "should be *ExprVariable")
	} else {
		var ok bool
		vr, ok = e.strct.(*ExprVariable)
		assert(ok, e.tok, "should be *ExprVariable")
	}
	assert(vr.gtype.typ == G_REL, e.tok, "expect G_REL|G_POINTER , but got "+vr.gtype.String())
	field := vr.gtype.relation.gtype.getField(e.fieldname)
	return vr.offset + field.offset
}

func (e *ExprStructField) emitLsave() {
	rel, ok := e.strct.(*Relation)
	assert(ok, e.tok, "should be *Relation")
	vr, ok := rel.expr.(*ExprVariable)
	assert(ok, e.tok, "should be *ExprVariable")
	assert(vr.gtype.typ == G_REL || vr.gtype.typ == G_POINTER, e.tok, "expect G_REL|G_POINTER , but got "+vr.gtype.String())
	if vr.gtype.typ == G_REL {
		field := vr.gtype.relation.gtype.getField(e.fieldname)
		emitLsave(field.getSize(), vr.offset+field.offset)
	} else if vr.gtype.typ == G_POINTER {
		field := vr.gtype.origType.relation.gtype.getField(e.fieldname)
		emit("push %%rax # rhs")
		emit("# assigning to a struct pointer field")
		vr.emit()
		emit("mov %%rax, %%rcx")
		emit("add $%d, %%rcx", field.offset)
		emit("pop %%rax  # rhs")
		reg := getReg(field.getSize())
		emit("mov %%%s, (%%rcx)", reg)
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
	assertNotNil(f.rng.indexvar != nil, f.rng.tok)
	assert(f.rng.rangeexpr.getGtype().typ == G_ARRAY || f.rng.rangeexpr.getGtype().typ == G_SLICE, f.rng.tok, "rangeexpr should be G_ARRAY or G_SLICE")

	emit("# for range")

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
					collection: f.rng.rangeexpr,
					index:      f.rng.indexvar,
				},
			},
		}
		assignVar.emit() // v = s[i]
	}

	emit("%s: # begin loop ", labelBegin)

	condition := &ExprBinop{
		op:    "<",
		left:  f.rng.indexvar,                  // i
		// @TODO
		// The range expression x is evaluated once before beginning the loop
		right: &ExprLen{
			arg: f.rng.rangeexpr,
		},
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
		TBI(stmt.token(), "too many number of arguments")
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

func emitGsave(regSize int, varname identifier, offset int) {
	reg := getReg(regSize)
	if offset != 0 {
		emit("mov %%%s, %s+%d(%%rip)", reg, varname, offset)
	} else {
		emit("mov %%%s, %s(%%rip)", reg, varname)
	}
}

func assignToStruct(lhs Expr, rhs Expr) {
	if rel, ok := lhs.(*Relation); ok {
		lhs = rel.expr
	}
	variable, ok := lhs.(*ExprVariable)
	assert(ok, nil, "lhs should be a variable")
	structliteral, ok := rhs.(*ExprStructLiteral)
	assert(ok || rhs == nil, nil, "invalid rhs")

	// initializes with zero values
	for _, fieldtype := range variable.gtype.relation.gtype.fields {
		//debugf("%#v", fieldtype)
		localOffset := variable.offset + fieldtype.offset
		switch fieldtype.typ {
		case G_ARRAY:
			//initArray(localOffset, fieldtype)
		case G_SLICE:
			initLocalSlice(localOffset)
		default:
			emit("mov $0, %%rax")
			regSize := fieldtype.getSize()
			assert(0 < regSize && regSize <= 8, variable.tok, fieldtype.String())
			emitLsave(regSize, localOffset)
		}
	}

	if structliteral == nil {
		return
	}

	strcttyp := structliteral.strctname.gtype
	// do assignment for each field
	for _, field := range structliteral.fields {
		fieldtype := strcttyp.getField(field.key)
		localOffset := variable.offset + fieldtype.offset

		switch {
		case fieldtype.typ == G_ARRAY:
			initvalues, ok := field.value.(*ExprArrayLiteral)
			assert(ok, nil, "ok")
			fieldtype := strcttyp.getField(field.key)
			setValuesToArray(localOffset, fieldtype, initvalues)
		case fieldtype.typ == G_SLICE:
			left := &ExprStructField{
				tok:       variable.tok,
				strct:     lhs,
				fieldname: field.key,
			}
			assignToSlice(left, field.value)

		default:
			field.value.emit()

			regSize := fieldtype.getSize()
			assert(0 < regSize && regSize <= 8, structliteral.tok, fieldtype.String())
			emitLsave(regSize, localOffset)
		}
	}

}

func initLocalSlice(offset int) {
	emit("# initialize slice with a zero value")
	emit("push $0")
	emit("push $0")
	emit("push $0")
	saveSlice(offset)
}

const sliceOffsetForLen = 8

func assignToSlice(lhs Expr, rhs Expr) {
	if rel, ok := lhs.(*Relation); ok {
		lhs = rel.expr
	}
	var targetOffset int
	switch lhs.(type) {
	case *ExprVariable:
		targetOffset = lhs.(*ExprVariable).offset
	case *ExprStructField:
		targetOffset = lhs.(*ExprStructField).getOffset()
	case *ExprIndex:
		TBI(lhs.token(), "Unable to assign to *ExprIndex")
	default:
		errorft(lhs.token(), "unkonwn type %T", lhs)
	}

	//assert(rhs == nil || rhs.getGtype().typ == G_SLICE, nil, "should be a slice literal or nil")
	if rhs == nil {
		initLocalSlice(targetOffset)
		return
	}

	switch rhs.(type) {
	case *Relation:
		rel := rhs.(*Relation)
		if _, ok := rel.expr.(*ExprNilLiteral); ok {
			// already initialied above
			return
		}
		rvariable, ok := rel.expr.(*ExprVariable)
		assert(ok, nil, "ok")
		// copy address
		rvariable.emit()
		emit("push %%rax")

		// copy len
		emit("mov %d(%%rbp), %%rax", rvariable.offset+ptrSize)
		emit("push %%rax")

		// copy cap
		emit("mov %d(%%rbp), %%rax", rvariable.offset+ptrSize+sliceOffsetForLen)
		emit("push %%rax")
	case *ExprSliceLiteral:
		lit := rhs.(*ExprSliceLiteral)
		// initialize a hidden array
		arrayLiteral := &ExprArrayLiteral{
			gtype:  lit.invisiblevar.gtype,
			values: lit.values,
		}
		assignToArray(lit.invisiblevar, arrayLiteral)
		lit.invisiblevar.emitAddress()
		emit("push %%rax")
		emit("push $%d", lit.invisiblevar.gtype.length) // len
		emit("push $%d", lit.invisiblevar.gtype.length) // cap
	case *ExprSlice:
		e := rhs.(*ExprSlice)
		emit("# assign to a slice")
		emit("#   emit address of the array")
		e.collection.emit()
		emit("push %%rax") // head of the array
		emit("#   emit low index")
		e.low.emit()
		emit("mov %%rax, %%rcx") // low index
		elmType := e.collection.getGtype().elementType
		size := elmType.getSize()
		assert(size > 0, nil, "size > 0")
		emit("mov $%d, %%rax", size) // size of one element
		emit("imul %%rcx, %%rax")    // index * size
		emit("pop %%rcx")            // head of the array
		emit("add %%rcx , %%rax")    // (index * size) + address
		emit("push %%rax")

		emit("#   calc and set len")
		calcLen := &ExprBinop{
			op:    "-",
			left:  e.high,
			right: e.low,
		}
		calcLen.emit()
		emit("push %%rax")

		emit("#   calc and set cap")
		calcCap := &ExprBinop{
			op: "-",
			left: &ExprNumberLiteral{
				val: e.collection.getGtype().length,
			},
			right: e.low,
		}
		calcCap.emit()
		emit("push %%rax")
	case *ExprConversion:
		// https://golang.org/ref/spec#Conversions
		// Converting a value of a string type to a slice of bytes type
		// yields a slice whose successive elements are the bytes of the string.
		//
		// see also https://blog.golang.org/strings
		conversion := rhs.(*ExprConversion)
		assert(conversion.gtype.typ == G_SLICE, rhs.token(), "must be a slice of bytes")
		assert(conversion.expr.getGtype().typ == G_STRING ||conversion.expr.getGtype().relation.gtype.typ == G_STRING , rhs.token(), "must be a string type, but got " + conversion.expr.getGtype().String())
		stringVarname,ok := conversion.expr.(*Relation)
		assert(ok, rhs.token(), "ok")
		stringVariable := stringVarname.expr.(*ExprVariable)
		stringVariable.emit()
		emit("push %%rax")
		strlen := stringVariable.getGtype().length
		emit("push $%d", strlen) // len
		emit("push $%d", strlen) // cap

	default:
		TBI(rhs.token(), "unable to handle %T", rhs)
	}

	saveSlice(targetOffset)
}

func saveSlice(targetOffset int) {
	emit("pop %%rax")
	emit("mov %%rax, %d(%%rbp)", targetOffset+ptrSize+sliceOffsetForLen)
	emit("pop %%rax")
	emit("mov %%rax, %d(%%rbp)", targetOffset+ptrSize)
	emit("pop %%rax")
	emit("mov %%rax, %d(%%rbp)", targetOffset)
}


func setValuesToArray(headOffset int, arrayType *Gtype, arrayLiteral *ExprArrayLiteral) {
	elmSize := arrayType.elementType.getSize()
	for i, val := range arrayLiteral.values {
		val.emit()
		localoffset := headOffset + i*elmSize
		emitLsave(elmSize, localoffset)
	}
}

// copy each element
func assignToArray(lhs Expr, rhs Expr) {
	if rel, ok := lhs.(*Relation); ok {
		lhs = rel.expr
	}
	variable, ok := lhs.(*ExprVariable)
	assert(ok, nil, "expect variable in lhs")
	headOffset := variable.offset

	arrayType := lhs.getGtype()
	elmSize := arrayType.elementType.getSize()
	for i := 0; i < arrayType.length; i++ {
		emit("mov $0, %%rax")
		if variable.isGlobal {
				emitGsave(elmSize, variable.varname, i*elmSize)
		} else {
				localoffset := variable.offset + i*elmSize
				emitLsave(elmSize, localoffset)
		}
	}

	if rhs == nil {
		return
	}

	switch rhs.(type) {
	case *Relation:
		rel := rhs.(*Relation)
		arrayVariable, ok := rel.expr.(*ExprVariable)
		assert(ok, nil, "ok")
		for i := 0; i < arrayType.length; i++ {
			emit("mov %d(%%rbp), %%rax", arrayVariable.offset+ i*elmSize)
			if variable.isGlobal {
				emitGsave(elmSize, variable.varname, i*elmSize)
			} else {
				localoffset := variable.offset + i*elmSize
				emitLsave(elmSize, localoffset)
			}
		}
	case *ExprStructField:
		strctField := rhs.(*ExprStructField)
		fieldType := strctField.getGtype()
		assert(fieldType.typ == G_ARRAY, nil, "should be array")
		for i := 0; i < arrayType.length; i++ {
			emit("mov %d(%%rbp), %%rax", strctField.getOffset()+ i*elmSize)
			if variable.isGlobal {
				emitGsave(elmSize, variable.varname, i*elmSize)
			} else {
				localoffset := variable.offset + i*elmSize
				emitLsave(elmSize, localoffset)
			}
		}

	case *ExprArrayLiteral:
		arrayLiteral := rhs.(*ExprArrayLiteral)
		//setValuesToArray(headOffset, arrayType, arrayLiteral)
		for i, val := range arrayLiteral.values {
			val.emit()
			if variable.isGlobal {
				emitGsave(elmSize, variable.varname, i*elmSize)
			} else {
				localoffset := headOffset + i*elmSize
				emitLsave(elmSize, localoffset)
			}
		}
	default:
		errorft(rhs.token(), "no supporetd %T", rhs)
	}
}

// for local var
func (decl *DeclVar) emit() {
	gtype := decl.variable.gtype
	switch {
	case gtype.typ == G_ARRAY:
		assignToArray(decl.varname, decl.initval)
	case gtype.typ == G_SLICE:
		assignToSlice(decl.varname, decl.initval)
	case gtype.typ == G_REL && gtype.relation.gtype.typ == G_STRUCT:
		assignToStruct(decl.varname, decl.initval)
	default:
		// primitive types like int,bool,byte
		rhs := decl.initval
		if rhs == nil {
			// assign zero value
			rhs = &ExprNumberLiteral{}
		}
		rhs.emit()
		emitLsave(decl.variable.getGtype().getSize(), decl.variable.offset)
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

func loadCollectIndex(array Expr, index Expr, offset int) {
	if array.getGtype().typ == G_ARRAY {
		elmType := array.getGtype().elementType
		array.emit()
		emit("push %%rax") // store address of variable

		index.emit()
		emit("mov %%rax, %%rcx") // index

		size := elmType.getSize()
		assert(size > 0, nil, "size > 0")
		emit("mov $%d, %%rax", size) // size of one element
		emit("imul %%rcx, %%rax")    // index * size
		emit("push %%rax")           // store index * size
		emit("pop %%rcx")            // load  index * size
		emit("pop %%rbx")            // load address of variable
		emit("add %%rcx , %%rbx")    // (index * size) + address
		if offset > 0 {
			emit("add $%d,  %%rbx", offset)
		}
		emit("mov (%%rbx), %%rax")   // dereference the content of an emelment
	} else if array.getGtype().typ == G_SLICE {
		elmType := array.getGtype().elementType
		emit("# emit address of the low index")
		array.emit() // eval pointer value
		emit("push %%rax")  // store head address

		index.emit() // index
		emit("mov %%rax, %%rcx")

		size := elmType.getSize()
		assert(size > 0, nil, "size > 0")
		emit("mov $%d, %%rax", size) // size of one element
		emit("imul %%rcx, %%rax")    // set e.index * size => %rax
		emit("pop %%rbx")            // load head address
		emit("add %%rax , %%rbx")    // (e.index * size) + head address
		if offset > 0 {
			emit("add $%d,  %%rbx", offset)
		}
		emit("mov (%%rbx), %%rax")   // dereference the content of an emelment

	} else {
		TBI(array.token(), "unable to handle %s", array.getGtype())
	}

}

func (e *ExprIndex) emit() {
	emit("# emit *ExprIndex")
	loadCollectIndex(e.collection, e.index, 0)
}

func (e *ExprNilLiteral) emit() {
	emit("mov $0, %%rax # nil literal")
}

func (ast *StmtShortVarDecl) emit() {
	a := &StmtAssignment{
		tok:    ast.tok,
		lefts:  ast.lefts,
		rights: ast.rights,
	}
	a.emit()
}

func (f *ExprFuncRef) emit() {
	emit("mov $1, %%rax") // emit 1 for now.  @FIXME
}

func (e *ExprSlice) emit() {
	TBI(e.token(), "")
}

func (e ExprArrayLiteral) emit() {
	errorft(e.token(),"DO NOT EMIT")
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
	panic("implement me " + e.tok.String())
}

func (e *ExprConversion) emit() {
	e.expr.emit()
}

func (e *ExprStructLiteral) emit() {
	errorft(e.token(),"This cannot be emitted alone")
}

func (e *ExprTypeSwitchGuard) emit() {
	TBI(e.token(), "")
}

func (e *ExprMapLiteral) emit() {
	TBI(e.token(), "")
}


func (ast *ExprMethodcall) getUniqueName() string {
	var gtype *Gtype

	switch ast.receiver.(type) {
	case *Relation:
		rel := ast.receiver.(*Relation)
		if vr, ok := rel.expr.(*ExprVariable); ok {
			gtype = vr.gtype
			if gtype.typ == G_REL && gtype.relation.gtype.typ == G_INTERFACE {
				TBI(ast.token(), "interface method call is not supported yet. (%s.%s)", gtype.relation.name, ast.fname)
			}
		} else {
			// @TODO must adapt to method chains like foo.Bar().Buz()
			TBI(ast.token(), "")
		}
	default:
		TBI(ast.token(), "unable to handle %T", ast.receiver)
	}
	//debugf("ast.receiver=%v", ast.receiver)
	//debugf("gtype=%v", gtype)
	return getMethodUniqueName(gtype, ast.fname)
}

func (methodCall *ExprMethodcall) getPkgName() identifier {
	origType := methodCall.getOrigType()
	if origType.typ == G_INTERFACE {
		TBI(methodCall.token(), "G_INTERFACE is not supported yet")
	} else {
		funcref, ok := origType.methods[methodCall.fname]
		if !ok {
			errorft(methodCall.token(), "method %s is not found in type %s", methodCall.fname, methodCall.receiver.getGtype())
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
		typeToBeloing = gtype.origType
	} else {
		typeToBeloing = gtype
	}
	assert(typeToBeloing.typ == G_REL, methodCall.tok, "method must belong to a named type")
	//debugf("typeToBeloing = %s", typeToBeloing)
	origType := typeToBeloing.relation.gtype
	//debugf("origType = %v", origType)
	return origType
}

func (methodCall *ExprMethodcall) getRettypes() []*Gtype {
	origType := methodCall.getOrigType()
	if origType.typ == G_INTERFACE {
		return origType.imethods[methodCall.fname].rettypes
	} else {
		funcref, ok := origType.methods[methodCall.fname]
		if !ok {
			errorft(methodCall.token(), "method %s is not found in type %s", methodCall.fname, methodCall.receiver.getGtype())
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

func (funcall *ExprFuncallOrConversion) getFuncDef() *DeclFunc {
	relexpr := funcall.rel.expr
	assertNotNil2(relexpr, funcall.tok, funcall.rel)
	funcref, ok := relexpr.(*ExprFuncRef)
	if !ok {
		errorft(funcall.token(),"Compiler error: funcref is not *ExprFuncRef but %v", funcref, funcall.fname)
	}
	assertNotNil(funcref.funcdef != nil, nil)
	return funcref.funcdef
}

type ExprLen struct {
	tok *Token
	arg Expr
}

func (e *ExprLen) token() *Token {
	panic("implement me")
}

func (e *ExprLen) dump() {
	panic("implement me")
}

func (e *ExprLen) getGtype() *Gtype {
	return gInt
}

func (e *ExprLen) emit() {
	emit("# emit len()")
	arg := e.arg
	gtype := arg.getGtype()
	switch {
	case gtype.typ == G_ARRAY:
		emit("mov $%d, %%rax", gtype.length)
	case gtype.typ == G_SLICE:
		var headOffset int
		switch arg.(type) {
		case *Relation:
			rel := arg.(*Relation)
			variable, ok := rel.expr.(*ExprVariable)
			assert(ok, arg.token(), "ok")
			if variable.isGlobal {
				emit("mov %d+%s(%%rip), %%rax", ptrSize, variable.varname)
			} else {
				headOffset = variable.offset
				emit("mov %d(%%rbp), %%rax", headOffset + ptrSize)
			}
		case *ExprStructField:
			headOffset = arg.(*ExprStructField).getOffset()
			emit("mov %d(%%rbp), %%rax", headOffset + ptrSize)
		case *ExprSliceLiteral:
			length := len(arg.(*ExprSliceLiteral).values)
			emit("mov $%d, %%rax", length)
		case *ExprSlice:
			sliceExpr := arg.(*ExprSlice)
			uop := &ExprBinop{
				op:"-",
				left: sliceExpr.high,
				right: sliceExpr.low,
			}
			uop.emit()
		default:
			TBI(arg.token(), "unable to handle %T", arg)
		}
	case gtype.typ == G_STRING, gtype.typ == G_REL && gtype.relation.gtype.typ == G_STRING :
		TBI(arg.token(), "unable to handle %s", gtype)
	case gtype.typ == G_MAP:
		TBI(arg.token(), "unable to handle %s", gtype)
	default:
		TBI(arg.token(), "unable to handle %s", gtype)
	}
}

func (funcall *ExprFuncallOrConversion) emit() {
	decl := funcall.getFuncDef() // check existance
	if decl == nil {
		errorft(funcall.token(), "funcdef not found for funcall %s, rel=%v ", funcall.fname, funcall.rel)
	}

	// len()
	if decl == builinLen {
		assert(len(funcall.args) == 1, nil, "invalid arguments for len()")
		arg := funcall.args[0]
		exprLen := &ExprLen{
			tok:arg.token(),
			arg: arg,
		}
		exprLen.emit()
		return
	}
	emitCall(getPackagedFuncName(decl.pkg, funcall.fname), funcall.args)
}

func emitCall(fname string, args []Expr) {

	emit("# funcall %s", fname)
	/*
	for i, _ := range args {
		emit("push %%%s", RegsForCall[i])
	}
	*/
	emit("# setting arguments")
	for i, arg := range args {
		//debugf("arg[%d] = %v", i, arg)
		if _, ok := arg.(*ExprVaArg); ok {
			// skip VaArg for now
			emit("mov $0, %%rax")
		} else {
			arg.emit()
		}
		emit("push %%rax  # argument no %d", i+1)
	}

	for i, _ := range args {
		j := len(args) - 1 - i
		emit("pop %%%s   # argument no %d", RegsForCall[j], j+1)
	}

	emit("mov $0, %%rax")
	emit("call %s", fname)

	/*
	for i, _ := range args {
		j := len(args) - 1 - i
		emit("pop %%%s", RegsForCall[j])
	}
	*/
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
		errorft(e.token(), "variable cannot be inteppreted at compile time")
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
	case *ExprConstVariable:
		return evalIntExpr(e.(*ExprConstVariable).val)
	default:
		errorft(e.token(), "unkown type %T", e)
	}
	return 0
}

// gloabal var which should be initialized with zeros
// https://en.wikipedia.org/wiki/.bss
func (decl *DeclVar) emitBss() {
	// https://sourceware.org/binutils/docs-2.30/as/Lcomm.html#Lcomm
	emit(".lcomm %s, %d", decl.variable.varname, decl.variable.getGtype().getSize())
}

func (e *ExprStructLiteral) lookup(fieldname identifier) Expr {
	for _, field := range e.fields {
		if field.key == fieldname {
			return field.value
		}
	}

	return nil
}

func emitGlobalDeclInit(ptok *Token, /* left type */ gtype *Gtype , value /* nullable */ Expr, containerName string) {
	if gtype.typ == G_ARRAY {
		arrayliteral, ok := value.(*ExprArrayLiteral)
		var values []Expr
		if ok {
			values = arrayliteral.values
		}
		assert(ok || arrayliteral == nil, ptok, fmt.Sprintf("*ExprArrayLiteral expected, but got %T", value))
		elmType := gtype.elementType
		assertNotNil(elmType != nil, nil)
		for i := 0; i < gtype.length; i++ {
			selector := fmt.Sprintf("%s[%d]", containerName, i)
			if i >= len(values) {
				// zero value
				emitGlobalDeclInit(ptok, elmType, nil, selector)
			} else {
				value := arrayliteral.values[i]
				assertNotNil(value != nil, nil)
				size := elmType.getSize()
				if size == 8 {
					if value.getGtype().typ == G_STRING {
						stringLiteral, ok := value.(*ExprStringLiteral)
						assert(ok, nil, "ok")
						emit(".quad .%s # %s", stringLiteral.slabel)
					} else {
						emit(".quad %d # %s %s", evalIntExpr(value),  value.getGtype(), selector)
					}
				} else if size == 1 {
					emit(".byte %d", evalIntExpr(value))
				} else {
					emitGlobalDeclInit(ptok, gtype.elementType, value, selector)
				}
			}
		}
	} else if gtype.typ == G_SLICE {
		switch value.(type) {
		case *ExprSliceLiteral:
			// initialize a hidden array
			lit := value.(*ExprSliceLiteral)
			lit.invisiblevar.varname = identifier(fmt.Sprintf("$hiddenArray$%d", getHidddenArrayId()))
			emit(".quad %s", lit.invisiblevar.varname) // address of the hidden array
			emit(".quad %d", lit.invisiblevar.gtype.length) // len
			emit(".quad %d", lit.invisiblevar.gtype.length) // cap
			arrayLiteral := &ExprArrayLiteral{
				gtype:  lit.invisiblevar.gtype,
				values: lit.values,
			}
			arrayDecl := &DeclVar{
				tok:ptok,
				variable:lit.invisiblevar,
				initval:arrayLiteral,
			}
			arrayDecl.emitGlobal()


		default:
			TBI(ptok,"unable to handle %s", gtype)
		}
	} else if gtype == gBool || (gtype.typ == G_REL && gtype.relation.gtype == gBool) {
		if value == nil {
			// zero value
			emit(".quad %d # %s %s", 0, gtype, containerName)
			return
		}
		val := evalIntExpr(value)
		emit(".quad %d # %s %s", val, gtype, containerName)
	} else if gtype.typ == G_REL && gtype.relation.gtype.typ == G_STRUCT {
		containerName = containerName + "." + string(gtype.relation.name)
		gtype.relation.gtype.calcStructOffset()
		for _, field := range gtype.relation.gtype.fields {
			if value == nil {
				emitGlobalDeclInit(ptok, field, nil, containerName + "." + string(field.fieldname))
				continue
			}
			structLiteral, ok := value.(*ExprStructLiteral)
			assert(ok, nil, "ok:" +containerName)
			value := structLiteral.lookup(field.fieldname)
			if value == nil {
				// zero value
				//continue
			}
			gtype := field
			emitGlobalDeclInit(ptok, gtype, value, containerName + "." + string(field.fieldname))
		}
	} else {
		var val int
		switch value.(type) {
		case nil:
			emit(".quad %d # %s %s zero value", 0, gtype, containerName)
		case *ExprNumberLiteral:
			val = value.(*ExprNumberLiteral).val
			emit(".quad %d # %s %s", val, gtype, containerName)
		case *ExprConstVariable:
			val = evalIntExpr(value)
			emit(".quad %d # %s ", val, gtype)
		case *ExprStringLiteral:
			stringLiteral := value.(*ExprStringLiteral)
			emit(".quad .%s", stringLiteral.slabel)
			value.getGtype().length = len(stringLiteral.val)
		case *Relation:
			emit(".quad 0 # TBI") // TBI
		default:
			TBI(ptok, "unable to handle %T", value)
		}
	}
}

func (decl *DeclVar) emitGlobal() {
	assert(decl.variable.isGlobal, nil, "should be global")
	assertNotNil(decl.variable.gtype != nil, nil)

	if decl.initval == nil {
		decl.emitBss()
		return
	}

	ptok := decl.token()
	gtype := decl.variable.gtype
	right := decl.initval

	emitLabel("%s: # %s", decl.variable.varname, gtype)
	emitGlobalDeclInit(ptok, gtype, right, "")
}

type IrRoot struct {
	vars           []*DeclVar
	funcs          []*DeclFunc
	stringLiterals []*ExprStringLiteral
}

var retvals = [...]string{"rt1", "rt2", "rt3", "rt4", "rt5", "rt6"}

func (root *IrRoot) emit() {

	// generate code
	emit(".data")

	emit("")
	emitComment("STRING LITERALS")
	for id, ast := range root.stringLiterals {
		ast.slabel = fmt.Sprintf("S%d", id)
		emitLabel(".%s:", ast.slabel)
		// https://sourceware.org/binutils/docs-2.30/as/String.html#String
		// the assembler marks the end of each string with a 0 byte.
		emit(".string \"%s\"", ast.val)
	}

	emit("")
	emitComment("GLOBAL RETVALS")
	for _, name := range retvals {
		emitLabel("%s:", name)
		emit(".quad 0")
	}

	emitComment("GLOBAL VARS")
	emit("")
	for _, vardecl := range root.vars {
		vardecl.emitGlobal()
	}

	emit("")
	emitComment("FUNCTIONS")
	emit(".text")
	for _, funcdecl := range root.funcs {
		funcdecl.emit()
	}
}
