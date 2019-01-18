package main

import "fmt"

var RegsForCall = []string{"rdi", "rsi", "rdx", "rcx", "r8", "r9"}

const INT_SIZE = 8 // not like 8cc

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

	// treat main.main function as a special one
	if f.pkg == "main" && f.fname == "main" {
		return "main"
	}

	// other functions
	return getPackagedFuncName(f.pkg, string(f.fname))
}

func (f *DeclFunc) emitPrologue() {
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
		offset -= INT_SIZE
		param.offset = offset
		emit("push %%%s", RegsForCall[i])
	}

	emit("# Allocating stiack memory for localvars")
	var localarea int
	for _, lvar := range f.localvars {
		if lvar.gtype == nil {
			debugf("%s has nil gtype ", lvar)
		}
		size := lvar.gtype.getSize()
		assert(size != 0, lvar.token(), "size is not zero")
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
		strcttype := variable.gtype.origType.relation.gtype
		field := strcttype.getField(a.fieldname)
		variable.emit()
		emit("add $%d, %%rax", field.offset)
		emit("mov (%%rax), %%rax")
	case G_REL: // struct
		strcttype := variable.gtype.relation.gtype
		field := strcttype.getField(a.fieldname)
		if field.typ == G_ARRAY {
			emit("lea %d(%%rbp), %%rax", variable.offset+field.offset)
		} else {
			emit("mov %d(%%rbp), %%rax", variable.offset+field.offset)
		}
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

func (vr *ExprVariable) emitAddress() {
	if vr.isGlobal {
		emit("lea %s(%%rip), %%rax", vr.varname)
	} else {
		if vr.offset == 0 {
			errorf("offset should not be zero for localvar %s", vr.varname)
		}
		emit("lea %d(%%rbp), %%rax", vr.offset)
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
		errorf("left.rel.expr is nil")
	}
	vr := rel.expr.(*ExprVariable)
	if vr.isGlobal {
		emitGsave(vr.gtype.getSize(), vr.varname)
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
				errorf("rel is not an variable")
			}
			vr.emitAddress()
		case *ExprStructLiteral:
			e := ast.operand.(*ExprStructLiteral)
			assert(e.invisiblevar.offset != 0, nil, "ExprStructLiteral's invisible var has offset")
			assignToStruct(e.invisiblevar, e)
			emit("lea %d(%%rbp), %%rax", e.invisiblevar.offset)
		default:
			errorf("Unknown type: %s", ast.operand)
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
		errorf("unable to handle uop %s", ast.op)
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
		errorf("Unknown binop: %s", ast.op)
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
				errorf("multivalue is not allowed")
			}
			numRight += len(rettypes)
		case *ExprMethodcall:
			rettypes := right.(*ExprMethodcall).getRettypes()
			if singleValueMode && len(rettypes) > 1 {
				errorf("multivalue is not allowed")
			}
			numRight += len(rettypes)
		default:
			numRight++
		}
	}
	if numLeft != numRight {
		errorf("number of exprs does not match")
	}

	for rightIndex, right := range ast.rights {
		left := ast.lefts[rightIndex]
		switch right.(type) {
		case *ExprFuncallOrConversion:
			rettypes := right.(*ExprFuncallOrConversion).getFuncDef().rettypes
			emit("# emitting rhs (funcall)")
			right.emit()
			for i, _ := range rettypes {
				emit("mov %s(%%rip), %%rax", retvals[len(rettypes)-1-i])
				emit("push %%rax")
			}
			for _, left := range ast.lefts {
				emit("pop %%rax")
				emitSave(left)
			}
		case *ExprMethodcall:
			rettypes := right.(*ExprMethodcall).getRettypes()
			emit("# emitting rhs (funcall)")
			right.emit()
			for i, _ := range rettypes {
				emit("mov %s(%%rip), %%rax", retvals[len(rettypes)-1-i])
				emit("push %%rax")
			}
			for _, left := range ast.lefts {
				emit("pop %%rax")
				emitSave(left)
			}
		default:
			gtype := right.getGtype()
			switch {
			case gtype.typ == G_ARRAY:
				assignToLocalArray(left, right)
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
		errorf("Unknown case %T", left)
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
		errorf("TBI %s", e.tok)
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
	var length int
	if f.rng.rangeexpr.getGtype().typ == G_ARRAY {
		length = f.rng.rangeexpr.getGtype().length
	}

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
			initArray(localOffset, fieldtype)
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
		errorf("TBI %s", lhs.token())
	default:
		errorf("unkonwn type %T", lhs)
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
		assignToLocalArray(lit.invisiblevar, arrayLiteral)
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
		errorf("TBI %T %s", rhs, rhs.token())
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

func initArray(headOffset int, arrayType *Gtype) {
	elmSize := arrayType.elementType.getSize()

	for i := 0; i < arrayType.length; i++ {
		emit("mov $0, %%rax")
		localoffset := headOffset + i*elmSize
		emitLsave(elmSize, localoffset)
	}
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
func assignToLocalArray(lhs Expr, rhs Expr) {
	if rel, ok := lhs.(*Relation); ok {
		lhs = rel.expr
	}
	variable, ok := lhs.(*ExprVariable)
	assert(ok, nil, "expect variable in lhs")
	headOffset := variable.offset
	arrayType := lhs.getGtype()
	initArray(headOffset, arrayType)

	if rhs == nil {
		return
	}

	switch rhs.(type) {
	case *Relation:
		rel := rhs.(*Relation)
		arrayVariable, ok := rel.expr.(*ExprVariable)
		assert(ok, nil, "ok")
		elmSize := arrayType.elementType.getSize()
		for i := 0; i < arrayType.length; i++ {
			emit("mov %d(%%rbp), %%rax", arrayVariable.offset+i*elmSize)
			localoffset := headOffset + i*elmSize
			emitLsave(elmSize, localoffset)
		}
	case *ExprStructField:
		strctField := rhs.(*ExprStructField)
		fieldType := strctField.getGtype()
		assert(fieldType.typ == G_ARRAY, nil, "should be array")
		elmSize := arrayType.elementType.getSize()
		for i := 0; i < arrayType.length; i++ {
			emit("mov %d(%%rbp), %%rax", strctField.getOffset()+i*elmSize)
			localoffset := headOffset + i*elmSize
			emitLsave(elmSize, localoffset)
		}

	case *ExprArrayLiteral:
		arrayLiteral := rhs.(*ExprArrayLiteral)
		setValuesToArray(headOffset, arrayType, arrayLiteral)
	default:
		errorf("no supporetd %T", rhs)
	}
}

// for local var
func (decl *DeclVar) emit() {
	gtype := decl.variable.gtype
	switch {
	case gtype.typ == G_ARRAY:
		assignToLocalArray(decl.varname, decl.initval)
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

func (e *ExprIndex) emit() {
	emit("# emit *ExprIndex")
	if e.collection.getGtype().typ == G_ARRAY {
		elmType := e.collection.getGtype().elementType

		e.collection.emit()
		emit("push %%rax") // store address of variable

		e.index.emit()
		emit("mov %%rax, %%rcx") // index

		size := elmType.getSize()
		assert(size > 0, nil, "size > 0")
		emit("mov $%d, %%rax", size) // size of one element
		emit("imul %%rcx, %%rax")    // index * size
		emit("push %%rax")           // store index * size
		emit("pop %%rcx")            // load  index * size
		emit("pop %%rbx")            // load address of variable
		emit("add %%rcx , %%rbx")    // (index * size) + address
		emit("mov (%%rbx), %%rax")   // dereference the content of an emelment
	} else if e.collection.getGtype().typ == G_SLICE {
		elmType := e.collection.getGtype().elementType
		emit("# emit address of the low index")
		e.collection.emit() // eval pointer value
		emit("push %%rax")  // store head address

		e.index.emit() // index
		emit("mov %%rax, %%rcx")

		size := elmType.getSize()
		assert(size > 0, nil, "size > 0")
		emit("mov $%d, %%rax", size) // size of one element
		emit("imul %%rcx, %%rax")    // set e.index * size => %rax
		emit("pop %%rbx")            // load head address
		emit("add %%rax , %%rbx")    // (e.index * size) + head address
		emit("mov (%%rbx), %%rax")   // dereference the content of an emelment

	} else {
		errorf("TBI")
	}
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
	panic("implement me " + e.tok.String())
}

func (e *ExprConversion) emit() {
	e.expr.emit()
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
	//debugf("ast.receiver=%v", ast.receiver)
	//debugf("gtype=%v", gtype)
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

func (funcall *ExprFuncallOrConversion) getFuncDef() *DeclFunc {
	relexpr := funcall.rel.expr
	assertNotNil2(relexpr, funcall.tok, funcall.rel)
	funcref, ok := relexpr.(*ExprFuncRef)
	if !ok {
		errorf("Compiler error: funcref is not *ExprFuncRef but %v", funcref, funcall.fname)
	}
	assertNotNil(funcref.funcdef != nil, nil)
	return funcref.funcdef
}

func emitBuiltinLen(args []Expr) {
	assert(len(args) == 1, nil, "invalid arguments for len()")
	arg := args[0]
	gtype := arg.getGtype()
	switch {
	case gtype.typ == G_ARRAY:
		emit("mov $%d, %%rax", gtype.length)
	case gtype.typ == G_SLICE:
		rel,ok := arg.(*Relation)
		assert(ok, arg.token(), "ok")
		variable ,ok := rel.expr.(*ExprVariable)
		assert(ok, arg.token(), "ok")
		emit("mov %d(%%rbp), %%rax", variable.offset + ptrSize)
	case gtype.typ == G_STRING, gtype.typ == G_REL && gtype.relation.gtype.typ == G_STRING :
		errorf("TBI %s", arg.token())
	case gtype.typ == G_MAP:
		errorf("TBI %s", arg.token())
	default:
		errorf("TBI %s", arg.token())
	}
}

func (funcall *ExprFuncallOrConversion) emit() {
	decl := funcall.getFuncDef() // check existance
	if decl == nil {
		errorf("funcdef not found for funcall %s, rel=%v ", funcall.fname, funcall.rel)
	}

	if decl == builinLen {
		emitBuiltinLen(funcall.args)
		return
	}
	emitCall(getPackagedFuncName(decl.pkg, funcall.fname), funcall.args)
}

func emitCall(fname string, args []Expr) {

	emit("# funcall %s", fname)
	for i, _ := range args {
		emit("push %%%s", RegsForCall[i])
	}

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

	for i, _ := range args {
		j := len(args) - 1 - i
		emit("pop %%%s", RegsForCall[j])
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
		errorf("unkown type %T", e)
	}
	return 0
}

func (decl *DeclVar) emitGlobal() {
	assert(decl.variable.isGlobal, nil, "should be global")
	assertNotNil(decl.variable.gtype != nil, nil)
	emitLabel("%s:", decl.variable.varname)

	if decl.variable.gtype.typ == G_ARRAY {
		arrayliteral, ok := decl.initval.(*ExprArrayLiteral)
		assert(ok, nil, "should be array lieteral")
		elmType := decl.variable.gtype.elementType
		assertNotNil(elmType != nil, nil)
		for _, value := range arrayliteral.values {
			assertNotNil(value != nil, nil)
			size := elmType.getSize()
			if size == 8 {
				if value.getGtype().typ == G_STRING {
					stringLiteral, ok := value.(*ExprStringLiteral)
					assert(ok, nil, "ok")
					emit(".quad .%s", stringLiteral.slabel)
				} else {
					emit(".quad %d", evalIntExpr(value))
				}
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
				emit(".quad %d", val)
			case *ExprConstVariable:
				val = evalIntExpr(decl.initval)
				emit(".quad %d", val)
			case *ExprStringLiteral:
				stringLiteral := decl.initval.(*ExprStringLiteral)
				emit(".quad .%s", stringLiteral.slabel)
				decl.variable.gtype.length = len(stringLiteral.val)
			}
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

	emit("")
	emitComment("GLOBAL VARS")
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
