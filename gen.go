package main

import (
	"fmt"
	"strings"
)

/**
  Intel® 64 and IA-32 Architectures Software Developer’s Manual
  Combined Volumes: 1, 2A, 2B, 2C, 2D, 3A, 3B, 3C, 3D and 4

  3.4.1.1 General-Purpose Registers in 64-Bit Mode

  In 64-bit mode, there are 16 general purpose registers and the default operand size is 32 bits.
  However, general-purpose registers are able to work with either 32-bit or 64-bit operands.
  If a 32-bit operand size is specified: EAX, EBX, ECX, EDX, EDI, ESI, EBP, ESP, R8D - R15D are available.
  If a 64-bit operand size is specified: RAX, RBX, RCX, RDX, RDI, RSI, RBP, RSP, R8-R15 are available.
  R8D-R15D/R8-R15 represent eight new general-purpose registers.
  All of these registers can be accessed at the byte, word, dword, and qword level.
  REX prefixes are used to generate 64-bit operand sizes or to reference registers R8-R15.
*/

var retRegi = [14]string{
	"rax", "rbx", "rcx", "rdx", "rdi", "rsi", "r8", "r9", "r10", "r11", "r12", "r13", "r14", "r15",
}

var RegsForCall = [...]string{"rdi", "rsi", "rdx", "rcx", "r8", "r9", "r10", "r11", "r12", "r13", "r14", "r15"}

const IntSize int = 8 // 64-bit (8 bytes)
const ptrSize int = 8
const sliceWidth int = 3
const sliceSize int = IntSize + ptrSize + ptrSize

var hiddenArrayId = 1

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

// Mytype.method -> Mytype#method

func getMethodUniqueName(gtype *Gtype, fname identifier) string {
	assertNotNil(gtype != nil, nil)
	var typename identifier
	if gtype.typ == G_POINTER {
		typename = gtype.origType.relation.name
	} else {
		typename = gtype.relation.name
	}
	return string(typename) + "$" + string(fname)
}

// "main","f1" -> "main.f1"
func getFuncSymbol(pkg identifier, fname string) string {
	if pkg == "libc" {
		return fname
	}

	return fmt.Sprintf("%s.%s", pkg, fname)
}

func (f *DeclFunc) getUniqueName() string {
	if f.receiver != nil {
		// method
		return getFuncSymbol(f.pkg, getMethodUniqueName(f.receiver.gtype, f.fname))
	}

	// other functions
	return getFuncSymbol(f.pkg, string(f.fname))
}

func (f *DeclFunc) emitPrologue() {
	uniquName := f.getUniqueName()
	emitLabel("%s:", uniquName)
	emit("push %%rbp")
	emit("mov %%rsp, %%rbp")

	var params []*ExprVariable

	// prepend receiver to params
	if f.receiver != nil {
		params = []*ExprVariable{f.receiver}
		for _, param := range f.params {
			params = append(params, param)
		}
	} else {
		params = f.params
	}

	// offset for params and local variables
	var offset int

	if len(params) > 0 {
		emit("# Allocating stack for params")
	}

	var regIndex int
	for _, param := range params {
		switch param.getGtype().getPrimType() {
		case G_SLICE, G_INTERFACE, G_MAP:
			offset -= IntSize * 3
			param.offset = offset
			emit("push %%%s # c", RegsForCall[regIndex+2])
			emit("push %%%s # b", RegsForCall[regIndex+1])
			emit("push %%%s # a", RegsForCall[regIndex])
			regIndex += sliceWidth
		default:
			offset -= IntSize
			param.offset = offset
			emit("push %%%s", RegsForCall[regIndex])
			regIndex += 1
		}
	}

	if len(f.localvars) > 0 {
		emit("# Allocating stack for localvars")
	}

	var localarea int
	for _, lvar := range f.localvars {
		if lvar.gtype == nil {
			debugf("%s has nil gtype ", lvar)
		}
		size := lvar.gtype.getSize()
		assert(size != 0, lvar.token(), "size should  not be zero:"+lvar.gtype.String())
		loff := align(size, 8)
		localarea -= loff
		offset -= loff
		lvar.offset = offset
		//debugf("set offset %d to lvar %s, type=%s", lvar.offset, lvar.varname, lvar.gtype)
	}

	for i := len(f.localvars) - 1; i >= 0; i-- {
		lvar := f.localvars[i]
		emit("# offset %d for variable \"%s\" of %s", lvar.offset, lvar.varname, lvar.gtype)
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

func emitFuncEpilogue(labelDeferHandler string, stmtDefer *StmtDefer) {
	emit("# func epilogue")
	// every function has a defer handler
	emit("%s: # defer handler", labelDeferHandler)

	// if the function has a defer statement, jump to there
	if stmtDefer != nil {
		emit("jmp %s", stmtDefer.label)
	}

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
		loadStructField(variable, field, offset)
	case *ExprVariable:
		variable := strct.(*ExprVariable)
		if field.typ == G_ARRAY {
			variable.emitAddress(field.offset)
		} else {
			if variable.isGlobal {
				emit("mov %s+%d(%%rip), %%rax # ", variable.varname, field.offset+offset)
			} else {
				emit("mov %d(%%rbp), %%rax", variable.offset+field.offset+offset)
			}
		}
	case *ExprStructField: // strct.field.field
		a := strct.(*ExprStructField)
		strcttype := a.strct.getGtype().relation.gtype
		assert(strcttype.size > 0, a.token(), "struct size should be > 0")
		field2 := strcttype.getField(a.fieldname)
		loadStructField(a.strct, field2, offset+field.offset)
	case *ExprIndex: // array[1].field
		indexExpr := strct.(*ExprIndex)
		loadCollectIndex(indexExpr.collection, indexExpr.index, offset+field.offset)
	default:
		// funcall().field
		// methodcall().field
		// *ptr.field
		// (MyStruct{}).field
		// (&MyStruct{}).field
		TBI(strct.token(), "unable to handle %T", strct)
	}

}

func (a *ExprStructField) emitAddress() {
	strcttype := a.strct.getGtype().origType.relation.gtype
	field := strcttype.getField(a.fieldname)
	a.strct.emit()
	emit("add $%d, %%rax", field.offset)
}

func (a *ExprStructField) emit() {
	emit("# ExprStructField.emit()")
	switch a.strct.getGtype().typ {
	case G_POINTER: // pointer to struct
		strcttype := a.strct.getGtype().origType.relation.gtype
		// very dirty hack
		if strcttype.size == undefinedSize {
			strcttype.calcStructOffset()
		}
		field := strcttype.getField(a.fieldname)
		a.strct.emit()
		emit("add $%d, %%rax", field.offset)
		emit("mov %%rax, %%rdx")
		switch field.getPrimType() {
		case G_SLICE, G_INTERFACE, G_MAP:
			emit("mov (%%rdx), %%rax")
			emit("mov %d(%%rdx), %%rbx", ptrSize)
			emit("mov %d(%%rdx), %%rcx", ptrSize+ptrSize)
		default:
			emit("mov (%%rdx), %%rax")
		}

	case G_REL: // struct
		strcttype := a.strct.getGtype().relation.gtype
		assert(strcttype.size > 0, a.token(), "struct size should be > 0")
		field := strcttype.getField(a.fieldname)
		loadStructField(a.strct, field, 0)
	default:
		errorft(a.token(), "internal error: bad gtype %s", a.strct.getGtype())
	}
}

func getLoadInst(size int) string {
	var inst string
	if size == 1 {
		inst = "movsbq"
	} else {
		inst = "mov"
	}

	return inst
}

func (ast *ExprVariable) emit() {
	emit("# emit variable \"%s\" of type %s", ast.varname, ast.getGtype())
	if ast.gtype.typ == G_ARRAY {
		ast.emitAddress(0)
		return
	} else if ast.gtype.getPrimType() == G_INTERFACE {
		if ast.isGlobal {
			emit("mov %s+%d(%%rip), %%rcx", ast.varname, ptrSize+ptrSize)
			emit("mov %s+%d(%%rip), %%rbx", ast.varname, ptrSize)
			emit("mov %s(%%rip), %%rax", ast.varname)
		} else {
			emit("mov %d(%%rbp), %%rcx", ast.offset+ptrSize+ptrSize)
			emit("mov %d(%%rbp), %%rbx", ast.offset+ptrSize)
			emit("mov %d(%%rbp), %%rax", ast.offset)
		}
		return
	}

	if ast.isGlobal {
		switch {
		case ast.getGtype().typ == G_SLICE:
			emit("#   emit slice variable")
			emit("mov %s(%%rip), %%rax # ptr", ast.varname)
			emit("mov %s+%d(%%rip), %%rbx # len", ast.varname, ptrSize)
			emit("mov %s+%d(%%rip), %%rcx # cap", ast.varname, ptrSize+IntSize)
		case ast.getGtype().typ == G_MAP:
			emit("#   emit map variable")
			emit("mov %s(%%rip), %%rax # ptr", ast.varname)
			emit("mov %s+%d(%%rip), %%rbx # len", ast.varname, ptrSize)
			emit("mov %s+%d(%%rip), %%rcx # cap", ast.varname, ptrSize+IntSize)
		default:
			inst := getLoadInst(ast.getGtype().getSize())
			emit("%s %s(%%rip), %%rax", inst, ast.varname)
		}
	} else {
		if ast.offset == 0 {
			errorft(ast.token(), "offset should not be zero for localvar %s", ast.varname)
		}
		switch {
		case ast.getGtype().typ == G_SLICE:
			emit("#   emit slice variable")
			emit("mov %d(%%rbp), %%rax # ptr", ast.offset)
			emit("mov %d(%%rbp), %%rbx # len", ast.offset+ptrSize)
			emit("mov %d(%%rbp), %%rcx # cap", ast.offset+ptrSize+IntSize)
		case ast.getGtype().typ == G_MAP:
			emit("#   emit map variable")
			emit("mov %d(%%rbp), %%rax # ptr", ast.offset)
			emit("mov %d(%%rbp), %%rbx # len", ast.offset+ptrSize)
			emit("mov %d(%%rbp), %%rcx # cap", ast.offset+ptrSize+IntSize)
		default:
			inst := getLoadInst(ast.getGtype().getSize())
			emit("%s %d(%%rbp), %%rax", inst, ast.offset)
		}
	}
}

func (variable *ExprVariable) emitAddress(offset int) {
	if variable.isGlobal {
		emit("lea %s+%d(%%rip), %%rax", variable.varname, offset)
	} else {
		if variable.offset == 0 {
			errorft(variable.token(), "offset should not be zero for localvar %s", variable.varname)
		}
		emit("lea %d(%%rbp), %%rax", variable.offset+offset)
	}
}

func (rel *Relation) emit() {
	if rel.expr == nil {
		errorft(rel.token(), "rel.expr is nil: %s", rel.name)
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

func emit_intcast(gtype *Gtype) {
	if gtype.getPrimType() == G_BYTE {
		emit("movzbq %%al, %%rax")
	}
}

func emit_comp_primitive(inst string, binop *ExprBinop) {
	emit("# emit_comp_primitive")
	binop.left.emit()
	if binop.left.getGtype().getPrimType() == G_BYTE {
		emit_intcast(binop.left.getGtype())
	}
	emit("push %%rax")
	binop.right.emit()
	if binop.right.getGtype().getPrimType() == G_BYTE {
		emit_intcast(binop.right.getGtype())
	}
	emit("pop %%rcx")
	emit("cmp %%rax, %%rcx") // right, left
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
	emit("# *ExprUop.emitSave()")
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
	variable := rel.expr.(*ExprVariable)
	variable.emitOffsetSave(variable.getGtype().getSize(), 0)
}

func (variable *ExprVariable) emitOffsetSave(size int, offset int) {
	emit("# ExprVariable.emitOffsetSave")
	assert(0 <= size && size <= 8, variable.token(), fmt.Sprintf("invalid size %d", size))
	if variable.getGtype().typ == G_POINTER && offset > 0 {
		assert(variable.getGtype().typ == G_POINTER, variable.token(), "")
		emit("push %%rax # save the value")
		variable.emit()
		emit("mov %%rax, %%rcx")
		emit("add $%d, %%rcx", offset)
		emit("pop %%rax")
		emit("mov %%rax, (%%rcx)")
		return
	}
	if variable.isGlobal {
		emitGsave(size, variable.varname, offset)
	} else {
		emitLsave(size, variable.offset+offset)
	}
}

func (variable *ExprVariable) emitOffsetLoad(size int, offset int) {
	assert(0 <= size && size <= 8, variable.token(), "invalid size")
	if variable.isGlobal {
		emitGload(size, variable.varname, offset)
	} else {
		emitLload(size, variable.offset+offset)
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
			vr.emitAddress(0)
		case *ExprStructLiteral:
			e := ast.operand.(*ExprStructLiteral)
			assert(e.invisiblevar.offset != 0, nil, "ExprStructLiteral's invisible var has offset")
			assignToStruct(e.invisiblevar, e)

			emitCallMalloc(e.getGtype().getSize()) // => rax
			emit("push %%rax")                     // to:ptr addr
			// @TODO handle global vars
			emit("lea %d(%%rbp), %%rax", e.invisiblevar.offset)
			emit("push %%rax") // from:address of invisible var
			emitCopyStructFromStack(e.getGtype())
			emit("pop %%rax") // from
			emit("pop %%rax") // to
			// emit address
		case *ExprStructField:
			e := ast.operand.(*ExprStructField)
			e.emitAddress()
		default:
			errorft(ast.token(), "Unknown type: %T", ast.operand)
		}
	} else if ast.op == "*" {
		// dereferene of a pointer
		//debugf("dereferene of a pointer")
		//rel, ok := ast.operand.(*Relation)
		//debugf("operand:%s", rel)
		//vr, ok := rel.expr.(*ExprVariable)
		//assert(ok, nil, "operand is a rel")
		ast.operand.emit()
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
		left := &ExprNumberLiteral{val: -1}
		binop := &ExprBinop{
			op:    "*",
			left:  left,
			right: ast.operand,
		}
		binop.emit()
	} else {
		errorft(ast.token(), "unable to handle uop %s", ast.op)
	}
	//debugf("end of emitting ExprUop")

}

func (binop *ExprBinop) emitCompareStrings() {
	var equal bool
	switch binop.op {
	case "<":
		TBI(binop.token(), "")
	case ">":
		TBI(binop.token(), "")
	case "<=":
		TBI(binop.token(), "")
	case ">=":
		TBI(binop.token(), "")
	case "!=":
		equal = false
	case "==":
		equal = true
	}

	binop.left.emit()
	emit("push %%rax")
	binop.right.emit()
	emit("pop %%rcx")
	emit("# rax = right, rcx = left")
	emitStringsEqual(equal, "%rcx", "%rax")
}

// call strcmp
func emitStringsEqual(equal bool, leftReg string, rightReg string) {
	emit("mov %s, %%rsi", leftReg)
	emit("mov %s, %%rdi", rightReg)
	emit("mov $0, %%rax")
	emit("call strcmp")
	emit("cmp $0, %%rax") // retval == 0
	if equal {
		emit("sete %%al")
	} else {
		emit("setne %%al")
	}
	emit("movzb %%al, %%eax")
}

func (binop *ExprBinop) emitComp() {
	if binop.left.getGtype().isString() {
		binop.emitCompareStrings()
		return
	}

	var instruction string
	switch binop.op {
	case "<":
		instruction = "setl"
	case ">":
		instruction = "setg"
	case "<=":
		instruction = "setle"
	case ">=":
		instruction = "setge"
	case "!=":
		instruction = "setne"
	case "==":
		instruction = "sete"
	}

	emit_comp_primitive(instruction, binop)
}

func emitStringConcate(left Expr, right Expr) {
	emit("# emitStringConcate")
	// newSize = strlen(left) + strlen(right) + 1
	binop := &ExprBinop{
		op: "+",
		left: &ExprLen{
			arg: left,
		},
		right: &ExprLen{
			arg: right,
		},
	}
	binop2 := &ExprBinop{
		op: "+",
		left: binop,
		right: &ExprNumberLiteral{
				val:1,
		},
	}

	// strcat(newstring, left)
	emitCallMallocDinamicSize(binop2) 	// malloc(newSize)
	emit("push %%rax")
	left.emit()
	emit("push %%rax")
	emit("pop %%rsi")
	emit("pop %%rdi")
	emit("mov $0, %%rax")
	emit("call strcat")

	emit("push %%rax")
	right.emit()
	emit("push %%rax")
	emit("pop %%rsi")
	emit("pop %%rdi")
	emit("mov $0, %%rax")
	emit("call strcat")
}

func (ast *ExprBinop) emit() {
	if ast.op == "+" && ast.left.getGtype().isString() {
		emitStringConcate(ast.left, ast.right)
		return
	}
	switch ast.op {
	case "<", ">", "<=", ">=", "!=", "==":
		ast.emitComp()
		return
	case "&&":
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
	case "||":
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

func isUnderScore(e Expr) bool {
	rel, ok := e.(*Relation)
	if !ok {
		return false
	}
	return rel.name == "_"
}

// https://golang.org/ref/spec#Assignments
// A tuple assignment assigns the individual elements of a multi-valued operation to a list of variables.
// There are two forms.
//
// In the first,
// the right hand operand is a single multi-valued expression such as a function call, a channel or map operation, or a type assertion.
// The number of operands on the left hand side must match the number of values.
// For instance, if f is a function returning two values,
//
//	x, y = f()
//
// assigns the first value to x and the second to y.
//
// In the second form,
// the number of operands on the left must equal the number of expressions on the right,
// each of which must be single-valued, and the nth expression on the right is assigned to the nth operand on the left:
//
//  one, two, three = '一', '二', '三'
//
func (ast *StmtAssignment) emit() {
	emit("")
	emit("# StmtAssignment")
	// the right hand operand is a single multi-valued expression
	// such as a function call, a channel or map operation, or a type assertion.
	// The number of operands on the left hand side must match the number of values.
	isOnetoOneAssignment := (len(ast.rights) > 1)
	if isOnetoOneAssignment {
		// a,b,c = expr1,expr2,expr3
		if len(ast.lefts) != len(ast.rights) {
			errorft(ast.token(), "number of exprs does not match")
		}

		for rightIndex, right := range ast.rights {
			left := ast.lefts[rightIndex]
			switch right.(type) {
			case *ExprFuncallOrConversion, *ExprMethodcall:
				rettypes := getRettypes(right)
				assert(len(rettypes) == 1, ast.token(), "return values should be one")
			}
			gtype := left.getGtype()
			switch {
			case gtype.typ == G_ARRAY:
				assignToArray(left, right)
			case gtype.typ == G_SLICE:
				assignToSlice(left, right)
			case gtype.typ == G_REL && gtype.relation.gtype.typ == G_STRUCT:
				assignToStruct(left, right)
			case gtype.getPrimType() == G_INTERFACE:
				assignToInterface(left, right)
			default:
				// suppose primitive
				emitAssignPrimitive(left, right)
			}
		}
		return
	} else {
		// a,b,c = expr
		numLeft := len(ast.lefts)
		numRight := 0
		right := ast.rights[0]

		var leftsMayBeTwo bool // a(,b) := expr // map index or type assertion
		switch right.(type) {
		case *ExprFuncallOrConversion, *ExprMethodcall:
			rettypes := getRettypes(right)
			if isOnetoOneAssignment && len(rettypes) > 1 {
				errorft(ast.token(), "multivalue is not allowed")
			}
			numRight += len(rettypes)
		case *ExprTypeAssertion:
			leftsMayBeTwo = true
			numRight++
		case *ExprIndex:
			indexExpr := right.(*ExprIndex)
			if indexExpr.collection.getGtype().getPrimType() == G_MAP {
				// map get
				leftsMayBeTwo = true
			}
			numRight++
		default:
			numRight++
		}

		if leftsMayBeTwo {
			if numLeft > 2 {
				errorft(ast.token(), "number of exprs does not match")
			}
		} else {
			if numLeft != numRight {
				errorft(ast.token(), "number of exprs does not match: %d <=> %d", numLeft, numRight)
			}
		}

		left := ast.lefts[0]
		switch right.(type) {
		case *ExprFuncallOrConversion, *ExprMethodcall:
			rettypes := getRettypes(right)
			if len(rettypes) > 1 {
				// a,b,c = f()
				emit("# a,b,c = f()")
				right.emit()
				var retRegiLen int
				for _, rettype := range rettypes {
					retSize := rettype.getSize()
					if retSize < 8 {
						retSize = 8
					}
					retRegiLen += retSize / 8
				}
				emit("# retRegiLen=%d\n", retRegiLen)
				for i := retRegiLen - 1; i >= 0; i-- {
					emit("push %%%s # %d", retRegi[i], i)
				}
				for _, left := range ast.lefts {
					if isUnderScore(left) {
						continue
					}
					assert(left.getGtype() != nil, left.token(), "should not be nil")
					if left.getGtype().typ == G_SLICE {
						// @TODO: Does this work ?
						emitSave3Elements(left, 0)
					} else if left.getGtype().getPrimType() == G_INTERFACE {
						// @TODO: Does this work ?
						emitSaveInterface(left, 0)
					} else {
						emit("pop %%rax")
						emitSave(left)
					}
				}
				return
			}
		}

		gtype := left.getGtype()
		if _, ok := left.(*Relation); ok {
			emit("# \"%s\" = ", left.(*Relation).name)
		}
		emit("# Assign %T %s = %T %s", left, gtype, right, right.getGtype())
		switch {
		case gtype == nil:
			// suppose left is "_"
			right.emit()
		case gtype.typ == G_ARRAY:
			assignToArray(left, right)
		case gtype.typ == G_SLICE:
			assignToSlice(left, right)
		case gtype.typ == G_REL && gtype.relation.gtype.typ == G_STRUCT:
			assignToStruct(left, right)
		case gtype.getPrimType() == G_INTERFACE:
			assignToInterface(left, right)
		case gtype.getPrimType() == G_MAP:
			assignToMap(left, right)
		default:
			// suppose primitive
			emitAssignPrimitive(left, right)
		}
		if leftsMayBeTwo && len(ast.lefts) == 2 {
			okVariable := ast.lefts[1]
			// @TODO consider big data like slice, struct, etd
			emit("mov %%rbx, %%rax") // ok
			emitSave(okVariable)
		}
		return
	}

}

func emitAssignPrimitive(left Expr, right Expr) {
	assert(left.getGtype().getSize() <= 8, left.token(), fmt.Sprintf("invalid type for lhs: %s", left.getGtype()))
	assert(right != nil || right.getGtype().getSize() <= 8, right.token(), fmt.Sprintf("invalid type for rhs: %s", right.getGtype()))
	right.emit()   //   expr => %rax
	emitSave(left) //   %rax => memory
}

// Each left-hand side operand must be addressable,
// a map index expression,
// or (for = assignments only) the blank identifier.
func emitSave(left Expr) {
	switch left.(type) {
	case *Relation:
		emit("# %s %s = ", left.(*Relation).name, left.getGtype())
		left.(*Relation).emitSave()
	case *ExprIndex:
		left.(*ExprIndex).emitSave()
	case *ExprStructField:
		left.(*ExprStructField).emitSave()
	case *ExprUop:
		left.(*ExprUop).emitSave()
	default:
		left.dump()
		errorft(left.token(), "Unknown case %T", left)
	}
}

// m[k] = v
// append key and value to the tail of map data, and increment its length
func (e *ExprIndex) emitMapSet() {

	labelAppend := makeLabel()
	labelSave := makeLabel()

	// map get to check if exists
	e.emit()
	// jusdge update or append
	emit("cmp $0, %%rcx")
	emit("setne %%al")
	emit("movzb %%al, %%eax")
	emit("test %%rax, %%rax")
	emit("je %s  # jump to append if not found", labelAppend)

	// update
	emit("push %%rcx") // push address of the key
	emit("jmp %s", labelSave)

	// append
	emit("%s: # append to a map ", labelAppend)
	e.collection.emit() // emit pointer address to %rax
	emit("push %%rax # stash head address of mapData")

	// emit len of the map
	elen := &ExprLen{
		arg: e.collection,
	}
	elen.emit()
	emit("imul $%d, %%rax", 2*8) // distance from head to tail
	emit("pop %%rcx")            // head
	emit("add %%rax, %%rcx")     // now rcx is the tail address
	emit("push %%rcx")

	// map len++
	elen.emit()
	emit("add $1, %%rax")
	emitOffsetSave(e.collection, IntSize, ptrSize) // update map len

	// Save key and value
	emit("%s: # end loop", labelSave)
	e.index.emit()
	emit("push %%rax") // index value

	mapType := e.collection.getGtype().getSource()
	mapKeyType := mapType.mapKey

	if mapKeyType.isString() {
		emit("pop %%rcx")          // index value
		emit("pop %%rax")          // map tail address
		emit("mov %%rcx, (%%rax)") // save indexvalue to malloced area
		emit("push %%rax")         // push map tail
	} else {
		// malloc(8)
		emit("mov $%d, %%rdi", 8) // malloc 8 bytes
		emit("mov $0, %%rax")
		emit("call .malloc")
		// %%rax : malloced address
		// stack : [map tail address, index value]
		emit("pop %%rcx")            // index value
		emit("mov %%rcx, (%%rax)")   // save indexvalue to malloced area
		emit("pop %%rcx")            // map tail address
		emit("mov %%rax, (%%rcx) #") // save index address to the tail
		emit("push %%rcx")           // push map tail
	}

	// save value

	// malloc(8)
	emit("mov $%d, %%rdi", 8) // malloc 8 bytes
	emit("mov $0, %%rax")
	emit("call .malloc")

	emit("pop %%rcx")           // map tail address
	emit("mov %%rax, 8(%%rcx)") // set malloced address to tail+8

	emit("pop %%rcx") // rhs value

	// save value
	emit("mov %%rcx, (%%rax)") // save value address to the malloced area
}

// save data from stack
func (e *ExprIndex) emitSaveInterface() {
	// load head address of the array
	// load index
	// multi index * size
	// calc address = head address + offset
	// copy value to the address

	collectionType := e.collection.getGtype()
	switch {
	case collectionType.getPrimType() == G_ARRAY, collectionType.getPrimType() == G_SLICE, collectionType.getPrimType() == G_STRING:
		e.collection.emit() // head address
	case collectionType.getPrimType() == G_MAP:
		e.emitMapSet()
		return
	default:
		TBI(e.token(), "unable to handle %s", collectionType)
	}

	emit("push %%rax # stash head address of collection")
	e.index.emit()
	emit("mov %%rax, %%rcx") // index
	var elmType *Gtype
	if collectionType.isString() {
		elmType = gByte
	} else {
		elmType = collectionType.elementType
	}
	size := elmType.getSize()
	assert(size > 0, nil, "size > 0")
	emit("mov $%d, %%rax # size of one element", size)
	emit("imul %%rcx, %%rax # index * size")
	emit("push %%rax # store index * size")
	emit("pop %%rcx # load index * size")
	emit("pop %%rax # load address of variable")
	emit("add %%rcx , %%rax # (index * size) + address")

	emit("mov %%rax, %%rbx")
	emit("pop %%rax # load RHS value(c)")
	emit("mov %%rax, 16(%%rbx)")
	emit("pop %%rax # load RHS value(b)")
	emit("mov %%rax, 8(%%rbx)")
	emit("pop %%rax # load RHS value(a)")
	emit("mov %%rax, (%%rbx)")
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
	case collectionType.getPrimType() == G_ARRAY, collectionType.getPrimType() == G_SLICE, collectionType.getPrimType() == G_STRING:
		e.collection.emit() // head address
	case collectionType.getPrimType() == G_MAP:
		e.emitMapSet()
		return
	default:
		TBI(e.token(), "unable to handle %s", collectionType)
	}

	emit("push %%rax # stash head address of collection")
	e.index.emit()
	emit("mov %%rax, %%rcx") // index
	var elmType *Gtype
	if collectionType.isString() {
		elmType = gByte
	} else {
		elmType = collectionType.elementType
	}
	size := elmType.getSize()
	assert(size > 0, nil, "size > 0")
	emit("mov $%d, %%rax # size of one element", size)
	emit("imul %%rcx, %%rax # index * size")
	emit("push %%rax # store index * size")
	emit("pop %%rcx # load index * size")
	emit("pop %%rax # load address of variable")
	emit("add %%rcx , %%rax # (index * size) + address")

	emit("mov %%rax, %%rbx")
	emit("pop %%rax # load RHS value")
	reg := getReg(size)
	emit("mov %%%s, (%%rbx) # finally save value to an element", reg)
}

func (e *ExprStructField) emitSave() {
	fieldType := e.getGtype()
	if e.strct.getGtype().typ == G_POINTER {
		emit("push %%rax # store rhs")
		// structptr.field = x
		e.strct.emit() // emit address
		emit("add $%d, %%rax", fieldType.offset)
		emit("mov %%rax, %%rbx")
		emit("pop %%rax # load rhs")
		emit("mov %%rax, (%%rbx)")
	} else {
		emitOffsetSave(e.strct, 8, fieldType.offset)
	}
}

func (e *ExprStructField) emitOffsetLoad(size int, offset int) {
	rel, ok := e.strct.(*Relation)
	assert(ok, e.tok, "should be *Relation")
	vr, ok := rel.expr.(*ExprVariable)
	assert(ok, e.tok, "should be *ExprVariable")
	assert(vr.gtype.typ == G_REL, e.tok, "expect G_REL, but got "+vr.gtype.String())
	field := vr.gtype.relation.gtype.getField(e.fieldname)
	vr.emitOffsetLoad(size, field.offset+offset)
}

// rax: address
// rbx: len
// rcx: cap
func (e *ExprSliceLiteral) emit() {
	emit("# (*ExprSliceLiteral).emit()")
	length := len(e.values)
	//emitCallMalloc(e.gtype.elementType.getSize() * length)
	//debugf("slice literal %s: underlyingarray size = %d (should be %d)", e.getGtype(), e.gtype.getSize(),  e.gtype.elementType.getSize() * length)
	emitCallMalloc(e.gtype.getSize()) // why does this work??
	emit("push %%rax # ptr")
	for i, value := range e.values {
		if e.gtype.elementType.getPrimType() == G_INTERFACE && value.getGtype().getPrimType() != G_INTERFACE {
			emitConversionToInterface(value)
		} else {
			value.emit()
		}

		switch e.gtype.elementType.getPrimType() {
		case G_BYTE, G_INT, G_POINTER, G_STRING:
			emit("pop %%rbx # ptr")
			emit("mov %%rax, %d(%%rbx)", IntSize*i)
			emit("push %%rbx # ptr")
		case G_INTERFACE, G_SLICE, G_MAP:
			emit("pop %%rdx # ptr")
			emit("mov %%rax, %d(%%rdx)", IntSize*3*i)
			emit("mov %%rbx, %d(%%rdx)", IntSize*3*i+ptrSize)
			emit("mov %%rcx, %d(%%rdx)", IntSize*3*i+ptrSize+ptrSize)
			emit("push %%rdx # ptr")
		default:
			TBI(e.token(), "")
		}
	}
	emit("pop %%rax # ptr")
	emit("mov $%d, %%rbx # len", length)
	emit("mov $%d, %%rcx # cap", length)
}

func (stmt *StmtIf) emit() {
	emit("# if")
	if stmt.simplestmt != nil {
		stmt.simplestmt.emit()
	}
	stmt.cond.emit()
	emit("test %%rax, %%rax")
	if stmt.els != nil {
		labelElse := makeLabel()
		labelEndif := makeLabel()
		emit("je %stmt  # jump if 0", labelElse)
		emit("# then block")
		stmt.then.emit()
		emit("jmp %stmt # jump to endif", labelEndif)
		emit("# else block")
		emit("%stmt:", labelElse)
		stmt.els.emit()
		emit("# endif")
		emit("%stmt:", labelEndif)
	} else {
		// no else block
		labelEndif := makeLabel()
		emit("je %stmt  # jump if 0", labelEndif)
		emit("# then block")
		stmt.then.emit()
		emit("# endif")
		emit("%stmt:", labelEndif)
	}
}

func (stmt *StmtSwitch) emit() {

	emit("#")
	emit("# switch statement")
	labelEnd := makeLabel()
	var labels []string

	// switch (expr) {
	if stmt.cond != nil {
		emit("# the subject expression")
		stmt.cond.emit()
		emit("push %%rax")
		emit("#")
	} else {
		// switch {
		emit("# no condition")
	}

	// case exp1,exp2,..:
	//     stmt1;
	//     stmt2;
	//     ...
	for i, caseClause := range stmt.cases {
		emit("# case %d", i)
		myCaseLabel := makeLabel()
		labels = append(labels, myCaseLabel)
		if stmt.isTypeSwitch {
			// compare type
			for _, gtype := range caseClause.gtypes {
				typeLabel := groot.getTypeLabel(gtype)
				emit("lea .%s(%%rip), %%rax # type: %s", typeLabel, gtype)
				emit("pop %%rcx # the subject type")
				emit("push %%rcx # the subject value")
				emitStringsEqual(true, "%rax", "%rcx")
				emit("test %%rax, %%rax")
				emit("jne %s # jump if matches", myCaseLabel)
			}
		} else if stmt.cond == nil {
			for _, e := range caseClause.exprs {
				e.emit()
				emit("test %%rax, %%rax")
				emit("jne %s # jump if matches", myCaseLabel)
				emit("push %%rcx # the subject value")
			}
		} else {
			for _, e := range caseClause.exprs {
				e.emit()
				emit("pop %%rcx # the subject value")
				if e.getGtype().isString() {
					emit("push %%rcx")
					emitStringsEqual(true, "%rax", "%rcx")
					emit("pop %%rcx")
				} else {
					emit("cmp %%rax, %%rcx") // right, left
					emit("sete %%al")
					emit("movzb %%al, %%eax")
				}
				emit("test %%rax, %%rax")
				emit("jne %s # jump if matches", myCaseLabel)
				emit("push %%rcx # the subject value")
			}
		}
	}

	var defaultLabel string
	if stmt.dflt == nil {
		emit("jmp %s", labelEnd)
	} else {
		emit("# default")
		defaultLabel = makeLabel()
		emit("jmp %s", defaultLabel)
	}

	emit("pop %%rax # destroy the subject value")
	emit("#")
	for i, caseClause := range stmt.cases {
		emit("# case stmts")
		emit("%s:", labels[i])
		caseClause.compound.emit()
		emit("jmp %s", labelEnd)
	}

	if stmt.dflt != nil {
		emit("%s:", defaultLabel)
		stmt.dflt.emit()
	}

	emit("%s: # end of switch", labelEnd)
}

func (f *StmtFor) emitRangeForMap() {
	emit("# for range %T", f.rng.rangeexpr.getGtype())
	assertNotNil(f.rng.indexvar != nil, f.rng.tok)
	labelBegin := makeLabel()
	f.labelEndBlock = makeLabel()
	f.labelEndLoop = makeLabel()

	mapCounter := &Relation{
		name: "",
		expr: f.rng.invisibleMapCounter,
	}
	// counter = 0
	initstmt := &StmtAssignment{
		lefts: []Expr{
			mapCounter,
		},
		rights: []Expr{
			&ExprNumberLiteral{
				val: 0,
			},
		},
	}
	emit("# init index")
	initstmt.emit()

	emit("%s: # begin loop ", labelBegin)

	// counter < len(list)
	condition := &ExprBinop{
		op:   "<",
		left: mapCounter, // i
		// @TODO
		// The range expression x is evaluated once before beginning the loop
		right: &ExprLen{
			arg: f.rng.rangeexpr, // len(expr)
		},
	}
	condition.emit()
	emit("test %%rax, %%rax")
	emit("je %s  # if false, exit loop", f.labelEndLoop)

	// set key and value
	mapCounter.emit()
	emit("imul $16, %%rax")
	emit("push %%rax")
	f.rng.rangeexpr.emit() // emit address of map data head
	mapType := f.rng.rangeexpr.getGtype()
	mapKeyType := mapType.mapKey

	emit("pop %%rcx")
	emit("add %%rax, %%rcx")
	emit("mov (%%rcx), %%rax")
	if !mapKeyType.isString() {
		emit("mov (%%rax), %%rax")
	}
	f.rng.indexvar.emitSave()

	if f.rng.valuevar != nil {
		emit("# Setting valuevar")
		emit("## rangeexpr.emit()")
		f.rng.rangeexpr.emit()
		emit("mov %%rax, %%rcx")

		emit("## mapCounter.emit()")
		mapCounter.emit()

		//assert(f.rng.valuevar.getGtype().getSize() <= 8, f.rng.token(), "invalid size")
		emit("## eval value")
		emit("imul $16, %%rax")
		emit("add $8, %%rax")
		emit("add %%rax, %%rcx")
		emit("mov (%%rcx), %%rdx")

		switch f.rng.valuevar.getGtype().getPrimType() {
		case G_SLICE, G_MAP:
			emit("mov (%%rdx), %%rax")
			emit("mov 8(%%rdx), %%rbx")
			emit("mov 16(%%rdx), %%rcx")
			emit("push %%rax")
			emit("push %%rbx")
			emit("push %%rcx")
			emitSave3Elements(f.rng.valuevar, 0)
		default:
			emit("mov (%%rdx), %%rax")
			f.rng.valuevar.emitSave()
		}

	}

	f.block.emit()
	emit("%s: # end block", f.labelEndBlock)

	// counter++
	indexIncr := &StmtInc{
		operand: mapCounter,
	}
	indexIncr.emit()

	emit("jmp %s", labelBegin)
	emit("%s: # end loop", f.labelEndLoop)
}

func (f *StmtFor) emitRangeForList() {
	emit("# for range %T", f.rng.rangeexpr.getGtype())
	assertNotNil(f.rng.indexvar != nil, f.rng.tok)
	assert(f.rng.rangeexpr.getGtype().typ == G_ARRAY || f.rng.rangeexpr.getGtype().typ == G_SLICE, f.rng.tok, "rangeexpr should be G_ARRAY or G_SLICE")

	labelBegin := makeLabel()
	f.labelEndBlock = makeLabel()
	f.labelEndLoop = makeLabel()

	// check if 0 == len(list)
	conditionEmpty := &ExprBinop{
		op: "==",
		left: &ExprNumberLiteral{
			val: 0,
		},
		right: &ExprLen{
			arg: f.rng.rangeexpr, // len(expr)
		},
	}
	conditionEmpty.emit()
	emit("test %%rax, %%rax")
	emit("jne %s  # if true, go to loop end", f.labelEndLoop)

	// i = 0
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
	initstmt.emit()

	// v = s[i]
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
		assignVar.emit()
	}

	emit("%s: # begin loop ", labelBegin)

	// i < len(list)
	condition := &ExprBinop{
		op:   "<",
		left: f.rng.indexvar, // i
		// @TODO
		// The range expression x is evaluated once before beginning the loop
		right: &ExprLen{
			arg: f.rng.rangeexpr, // len(expr)
		},
	}
	condition.emit()
	emit("test %%rax, %%rax")
	emit("je %s  # if false, go to loop end", f.labelEndLoop)

	f.block.emit()
	emit("%s: # end block", f.labelEndBlock)

	// i++
	indexIncr := &StmtInc{
		operand: f.rng.indexvar,
	}
	indexIncr.emit()

	// v = s[i]
	if f.rng.valuevar != nil {
		assignVar.emit()
	}
	emit("jmp %s", labelBegin)
	emit("%s: # end loop", f.labelEndLoop)
}

func (f *StmtFor) emitForClause() {
	assertNotNil(f.cls != nil, nil)
	labelBegin := makeLabel()
	f.labelEndBlock = makeLabel()
	f.labelEndLoop = makeLabel()

	if f.cls.init != nil {
		f.cls.init.emit()
	}
	emit("%s: # begin loop ", labelBegin)
	if f.cls.cond != nil {
		f.cls.cond.emit()
		emit("test %%rax, %%rax")
		emit("je %s  # jump if false", f.labelEndLoop)
	}
	f.block.emit()
	emit("%s: # end block", f.labelEndBlock)
	if f.cls.post != nil {
		f.cls.post.emit()
	}
	emit("jmp %s", labelBegin)
	emit("%s: # end loop", f.labelEndLoop)
}

func (f *StmtFor) emit() {
	if f.rng != nil {
		if f.rng.rangeexpr.getGtype().typ == G_MAP {
			f.emitRangeForMap()
		} else {
			f.emitRangeForList()
		}
		return
	}
	f.emitForClause()
}

func (stmt *StmtReturn) emitDeferAndReturn() {
	if stmt.labelDeferHandler != "" {
		emit("# defer and return")
		emit("jmp %s", stmt.labelDeferHandler)
	}
}

func (stmt *StmtReturn) emit() {
	if len(stmt.exprs) == 0 {
		// return void
		emit("mov $0, %%rax")
		stmt.emitDeferAndReturn()
		return
	}

	if len(stmt.exprs) > 7 {
		TBI(stmt.token(), "too many number of arguments")
	}

	var retRegiIndex int
	if len(stmt.exprs) == 1 {
		expr := stmt.exprs[0]
		rettype := stmt.rettypes[0]
		if rettype.getPrimType() == G_INTERFACE && expr.getGtype().getPrimType() != G_INTERFACE {
			if expr.getGtype() == nil {
				emit("mov $0, %%rax")
				emit("mov $0, %%rbx")
				emit("mov $0, %%rcx")
			} else {
				emitConversionToInterface(expr)
			}
		} else {
			expr.emit()
			if expr.getGtype() == nil && stmt.rettypes[0].typ == G_SLICE {
				emit("mov $0, %%rbx")
				emit("mov $0, %%rcx")
			}
		}
		stmt.emitDeferAndReturn()
		return
	}
	for i, rettype := range stmt.rettypes {
		expr := stmt.exprs[i]
		expr.emit()
		//		rettype := stmt.rettypes[i]
		if expr.getGtype() == nil && rettype.typ == G_SLICE {
			emit("mov $0, %%rbx")
			emit("mov $0, %%rcx")
		}
		size := rettype.getSize()
		if size < 8 {
			size = 8
		}
		var num64bit int = size / 8 // @TODO odd size
		for j := 0; j < num64bit; j++ {
			emit("push %%%s", retRegi[num64bit-1-j])
			retRegiIndex++
		}
	}
	for i := 0; i < retRegiIndex; i++ {
		emit("pop %%%s", retRegi[retRegiIndex-1-i])
	}

	stmt.emitDeferAndReturn()
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
	emit("mov %%%s, %s+%d(%%rip)", reg, varname, offset)
}

func emitLload(regSize int, loff int) {
	reg := getReg(regSize)
	emit("mov %d(%%rbp), %%%s", loff, reg)
}

func emitGload(regSize int, varname identifier, offset int) {
	reg := getReg(regSize)
	emit("mov %s+%d(%%rip), %%%s", varname, offset, reg)
}

func emitAddress(e Expr) {
	switch e.(type) {
	case *Relation:
		emitAddress(e.(*Relation).expr)
	case *ExprVariable:
		e.(*ExprVariable).emitAddress(0)
	default:
		TBI(e.token(), "")
	}
}

// expect lhs address is in the stack top, rhs is in the second top
func emitCopyStructFromStack(gtype *Gtype) {
	//assert(left.getGtype().getSize() == right.getGtype().getSize(), left.token(),"size does not match")
	emit("pop %%rax") // from
	emit("pop %%rbx") // to
	emit("push %%rcx")
	emit("push %%r11")
	emit("mov %%rax, %%rcx") // from
	emit("mov %%rbx, %%rax") // to

	var i int
	for i = i; i < gtype.getSize(); i += 8 {
		emit("movq %d(%%rcx), %%r11", i)
		emit("movq %%r11, %d(%%rax)", i)
	}
	for i = i; i < gtype.getSize(); i += 4 {
		emit("movl %d(%%rcx), %%r11", i)
		emit("movl %%r11, %d(%%rax)", i)
	}
	for i = i; i < gtype.getSize(); i++ {
		emit("movb %d(%%rcx), %%r11", i)
		emit("movb %%r11, %d(%%rax)", i)
	}

	emit("pop %%r11")
	emit("pop %%rcx")

	// recover stack
	emit("push %%rax") // to
	emit("push %%rcx") // from
}

// expect rhs address is in the stack top
func emitCopyStruct(left Expr) {
	//assert(left.getGtype().getSize() == right.getGtype().getSize(), left.token(),"size does not match")
	emit("pop %%rax")
	emit("push %%rcx")
	emit("push %%r11")
	emit("mov %%rax, %%rcx")
	emitAddress(left)

	var i int
	for i = i; i < left.getGtype().getSize(); i += 8 {
		emit("movq %d(%%rcx), %%r11", i)
		emit("movq %%r11, %d(%%rax)", i)
	}
	for i = i; i < left.getGtype().getSize(); i += 4 {
		emit("movl %d(%%rcx), %%r11", i)
		emit("movl %%r11, %d(%%rax)", i)
	}
	for i = i; i < left.getGtype().getSize(); i++ {
		emit("movb %d(%%rcx), %%r11", i)
		emit("movb %%r11, %d(%%rax)", i)
	}

	emit("pop %%r11")
	emit("pop %%rcx")
}

func assignToStruct(lhs Expr, rhs Expr) {
	emit("# assignToStruct")
	if rel, ok := lhs.(*Relation); ok {
		lhs = rel.expr
	}
	assert(rhs == nil || (rhs.getGtype().typ == G_REL && rhs.getGtype().relation.gtype.typ == G_STRUCT),
		lhs.token(), "rhs should be struct type")
	// initializes with zero values
	for _, fieldtype := range lhs.getGtype().relation.gtype.fields {
		//debugf("%#v", fieldtype)
		switch {
		case fieldtype.typ == G_ARRAY:
			arrayType := fieldtype
			elementType := arrayType.elementType
			elmSize := arrayType.elementType.getSize()
			switch {
			case elementType.typ == G_REL && elementType.relation.gtype.typ == G_STRUCT:
				left := &ExprStructField{
					strct:     lhs,
					fieldname: fieldtype.fieldname,
				}
				assignToArray(left, nil)
			default:
				assert(0 <= elmSize && elmSize <= 8, lhs.token(), "invalid size")
				for i := 0; i < arrayType.length; i++ {
					emit("mov $0, %%rax")
					emitOffsetSave(lhs, elmSize, fieldtype.offset+i*elmSize)
				}
			}

		case fieldtype.typ == G_SLICE:
			emit("# initialize slice with a zero value")
			emit("push $0")
			emit("push $0")
			emit("push $0")
			emitSave3Elements(lhs, fieldtype.offset)
		case fieldtype.typ == G_MAP:
			emit("# initialize slice with a zero value")
			emit("push $0")
			emit("push $0")
			emit("push $0")
			emitSave3Elements(lhs, fieldtype.offset)
		case fieldtype.typ == G_REL && fieldtype.relation.gtype.typ == G_STRUCT:
			left := &ExprStructField{
				strct:     lhs,
				fieldname: fieldtype.fieldname,
			}
			assignToStruct(left, nil)
		case fieldtype.getPrimType() == G_INTERFACE:
			emit("push $0")
			emit("push $0")
			emit("push $0")
			emitSaveInterface(lhs, fieldtype.offset)
		default:
			emit("mov $0, %%rax")
			regSize := fieldtype.getSize()
			assert(0 < regSize && regSize <= 8, lhs.token(), fieldtype.String())
			emitOffsetSave(lhs, regSize, fieldtype.offset)
		}
	}

	if rhs == nil {
		return
	}
	variable := lhs

	strcttyp := rhs.getGtype().getSource()

	switch rhs.(type) {
	case *Relation:
		emitAddress(rhs)
		emit("push %%rax")
		emitCopyStruct(lhs)
	case *ExprUop:
		re := rhs.(*ExprUop)
		if re.op == "*" {
			// copy struct
			re.operand.emit()
			emit("push %%rax")
			emitCopyStruct(lhs)
		} else {
			TBI(rhs.token(), "")
		}
	case *ExprStructLiteral:
		structliteral, ok := rhs.(*ExprStructLiteral)
		assert(ok || rhs == nil, rhs.token(), fmt.Sprintf("invalid rhs: %T", rhs))

		// do assignment for each field
		for _, field := range structliteral.fields {
			emit("# .%s", field.key)
			fieldtype := strcttyp.getField(field.key)

			switch {
			case fieldtype.typ == G_ARRAY:
				initvalues, ok := field.value.(*ExprArrayLiteral)
				assert(ok, nil, "ok")
				fieldtype := strcttyp.getField(field.key)
				arrayType := fieldtype
				elementType := arrayType.elementType
				elmSize := elementType.getSize()
				switch {
				case elementType.typ == G_REL && elementType.relation.gtype.typ == G_STRUCT:
					left := &ExprStructField{
						strct:     lhs,
						fieldname: fieldtype.fieldname,
					}
					assignToArray(left, field.value)
				default:
					for i, val := range initvalues.values {
						val.emit()
						emitOffsetSave(variable, elmSize, fieldtype.offset+i*elmSize)
					}
				}
			case fieldtype.typ == G_SLICE:
				left := &ExprStructField{
					tok:       variable.token(),
					strct:     lhs,
					fieldname: field.key,
				}
				assignToSlice(left, field.value)
			case fieldtype.getPrimType() == G_MAP:
				left := &ExprStructField{
					tok:       variable.token(),
					strct:     lhs,
					fieldname: field.key,
				}
				assignToMap(left, field.value)
			case fieldtype.getPrimType() == G_INTERFACE:
				left := &ExprStructField{
					tok:       lhs.token(),
					strct:     lhs,
					fieldname: field.key,
				}
				assignToInterface(left, field.value)
			case fieldtype.typ == G_REL && fieldtype.relation.gtype.typ == G_STRUCT:
				left := &ExprStructField{
					tok:       variable.token(),
					strct:     lhs,
					fieldname: field.key,
				}
				assignToStruct(left, field.value)
			default:
				field.value.emit()

				regSize := fieldtype.getSize()
				assert(0 < regSize && regSize <= 8, variable.token(), fieldtype.String())
				emitOffsetSave(variable, regSize, fieldtype.offset)
			}
		}
	default:
		TBI(rhs.token(), "")
	}

}

const sliceOffsetForLen = 8

func emitOffsetSave(lhs Expr, size int, offset int) {
	switch lhs.(type) {
	case *Relation:
		rel := lhs.(*Relation)
		emitOffsetSave(rel.expr, size, offset)
	case *ExprVariable:
		variable := lhs.(*ExprVariable)
		variable.emitOffsetSave(size, offset)
	case *ExprStructField:
		structfield := lhs.(*ExprStructField)
		fieldType := structfield.getGtype()
		emitOffsetSave(structfield.strct, size, fieldType.offset+offset)
	case *ExprIndex:
		indexExpr := lhs.(*ExprIndex)
		emitCollectIndexSave(indexExpr.collection, indexExpr.index, offset)

	default:
		errorft(lhs.token(), "unkonwn type %T", lhs)
	}
}

func emitOffsetLoad(lhs Expr, size int, offset int) {
	switch lhs.(type) {
	case *Relation:
		rel := lhs.(*Relation)
		emitOffsetLoad(rel.expr, size, offset)
	case *ExprVariable:
		variable := lhs.(*ExprVariable)
		variable.emitOffsetLoad(size, offset)
	case *ExprStructField:
		structfield := lhs.(*ExprStructField)
		fieldType := structfield.getGtype()
		if structfield.strct.getGtype().typ == G_POINTER {
			structfield.strct.emit() // emit address of the struct
			emit("add $%d, %%rax", fieldType.offset+offset)
			reg := getReg(size)
			emit("mov (%%rax), %%%s", reg)
		} else {
			emitOffsetLoad(structfield.strct, size, fieldType.offset+offset)
		}
	case *ExprIndex:
		//  e.g. arrayLiteral.values[i].getGtype().getPrimType()
		indexExpr := lhs.(*ExprIndex)
		loadCollectIndex(indexExpr.collection, indexExpr.index, offset)
	case *ExprMethodcall:
		// @TODO this logic is temporarly. Need to be verified.
		mcall := lhs.(*ExprMethodcall)
		rettypes := mcall.getRettypes()
		assert(len(rettypes) == 1, lhs.token(), "rettype should be single")
		rettype := rettypes[0]
		assert(rettype.getPrimType() == G_POINTER, lhs.token(), "only pointer is supported")
		mcall.emit()
		emit("# START DEBUG")
		emit("add $%d, %%rax", offset)
		emit("mov (%%rax), %%rax")
	default:
		errorft(lhs.token(), "unkonwn type %T", lhs)
	}
}

func emitSaveInterface(lhs Expr, offset int) {
	switch lhs.(type) {
	case *Relation:
		rel := lhs.(*Relation)
		emitSaveInterface(rel.expr, offset)
	case *ExprVariable:
		variable := lhs.(*ExprVariable)
		variable.saveInterface(offset)
	case *ExprStructField:
		structfield := lhs.(*ExprStructField)
		fieldType := structfield.getGtype()
		emitSaveInterface(structfield.strct, fieldType.offset+offset)
	case *ExprIndex:
		indexExpr := lhs.(*ExprIndex)
		indexExpr.emitSaveInterface()
	default:
		errorft(lhs.token(), "unkonwn type %T", lhs)
	}
}

// take slice values from stack
func emitSave3Elements(lhs Expr, offset int) {
	switch lhs.(type) {
	case *Relation:
		rel := lhs.(*Relation)
		emitSave3Elements(rel.expr, offset)
	case *ExprVariable:
		variable := lhs.(*ExprVariable)
		variable.saveSlice(offset)
	case *ExprStructField:
		structfield := lhs.(*ExprStructField)
		fieldType := structfield.getGtype()
		emitSave3Elements(structfield.strct, fieldType.offset+offset)
	case *ExprIndex:
		indexExpr := lhs.(*ExprIndex)
		indexExpr.emitSave()
	default:
		errorft(lhs.token(), "unkonwn type %T", lhs)
	}
}

func emitCallMallocDinamicSize(eSize Expr) {
	eSize.emit()
	emit("mov %%rax, %%rdi")
	emit("mov $0, %%rax")
	emit("call .malloc")
}

func emitCallMalloc(size int) {
	emit("mov $%d, %%rdi", size)
	emit("mov $0, %%rax")
	emit("call .malloc")
}

// push addr, len, cap
func (lit *ExprMapLiteral) emitPush() {
	length := len(lit.elements)

	// allocaated address of the map head
	var size int
	if length == 0 {
		size = ptrSize * 1024
	} else {
		size = length * ptrSize * 1024
	}
	emitCallMalloc(size)
	emit("push %%rax") // map head

	mapType := lit.getGtype()
	mapKeyType := mapType.mapKey

	for i, element := range lit.elements {
		// alloc key
		if mapKeyType.isString() {
			element.key.emit()
		} else {
			element.key.emit()
			emit("push %%rax") // value of key
			// call malloc for key
			emitCallMalloc(8)
			emit("pop %%rcx")          // value of key
			emit("mov %%rcx, (%%rax)") // save key to heap
		}

		emit("pop %%rbx")                     // map head
		emit("mov %%rax, %d(%%rbx) #", i*2*8) // save key address
		emit("push %%rbx")                    // map head

		if element.value.getGtype().getSize() <= 8 {
			element.value.emit()
			emit("push %%rax") // value of value
			// call malloc
			emitCallMalloc(8)
			emit("pop %%rcx")          // value of value
			emit("mov %%rcx, (%%rax)") // save value to heap
		} else {
			switch element.value.getGtype().getPrimType() {
			case G_MAP, G_SLICE, G_INTERFACE:
				// rax,rbx,rcx
				element.value.emit()
				emit("push %%rax") // ptr
				emitCallMalloc(8 * 3)
				emit("pop %%rdx") // ptr
				emit("mov %%rdx, (%%rax)")
				emit("mov %%rbx, %d(%%rax)", 8*1)
				emit("mov %%rcx, %d(%%rax)", 8*2)

			default:
				TBI(element.value.token(), "unable to handle %s", element.value.getGtype())
			}
		}

		emit("pop %%rbx") // map head
		emit("mov %%rax, %d(%%rbx) #", i*2*8+8)
		emit("push %%rbx")
	}

	emit("pop %%rax")
	emit("push %%rax")       // address (head of the heap)
	emit("push $%d", length) // len
	emit("push $%d", length) // cap
}

func assignToMap(lhs Expr, rhs Expr) {
	emit("# assignToMap")
	if rhs == nil {
		emit("# initialize map with a zero value")
		emit("push $0")
		emit("push $0")
		emit("push $0")
		emitSave3Elements(lhs, 0)
		return
	}
	switch rhs.(type) {
	case *ExprMapLiteral:
		emit("# map literal")

		lit := rhs.(*ExprMapLiteral)
		lit.emitPush()
	case *Relation, *ExprVariable, *ExprIndex, *ExprStructField, *ExprFuncallOrConversion:
		rhs.emit()
		emit("push %%rax")
		emit("push %%rbx")
		emit("push %%rcx")
	default:
		TBI(rhs.token(), "unable to handle %T", rhs)
	}
	emitSave3Elements(lhs, 0)
}

func (e *ExprConversionToInterface) emit() {
	emit("# ExprConversionToInterface")
	emitConversionToInterface(e.expr)
}

func emitConversionToInterface(dynamicValue Expr) {
	emit("# emitConversionToInterface")
	dynamicValue.emit()
	emit("push %%rax")
	emitCallMalloc(8)
	emit("pop %%rcx")          // dynamicValue
	emit("mov %%rcx, (%%rax)") // store value to heap
	emit("push %%rax # addr of dynamicValue")         // address

	receiverType := dynamicValue.getGtype()
	if receiverType == nil {
		errorft(dynamicValue.token(), "type is nil:%s", dynamicValue)
	}
	if receiverType.typ == G_POINTER {
		receiverType = receiverType.origType.relation.gtype
	}
	//assert(receiverType.typeId > 0,  dynamicValue.token(), "no typeId")
	emit("mov $%d, %%rax # receiverTypeId", receiverType.typeId)
	emit("push %%rax # receiverTypeId")

	gtype := dynamicValue.getGtype()
	//debugf("dynamic type:%s", gtype)
	dynamicTypeId, ok := groot.hashedTypes[gtype.String()]
	if !ok {
		//debugf("types:%#v", groot.hashedTypes)
		//debugf("gtype.origType.relation.pkg:%s", gtype.origType.relation.pkg)
		errorft(dynamicValue.token(), "type %s not found for %s", gtype, dynamicValue)
	}
	label := fmt.Sprintf("DT%d", dynamicTypeId)
	emit("lea .%s, %%rax# dynamicType %s", label, gtype.String())
	emit("mov %%rax, %%rcx # dynamicType")
	emit("pop %%rbx # receiverTypeId")
	emit("pop %%rax # addr of dynamicValue")
	emit("")
}

func isNil(e Expr) bool {
	rel, ok := e.(*Relation)
	if ok {
		_, isNil := rel.expr.(*ExprNilLiteral)
		return isNil
	}
	return false
}

func assignToInterface(lhs Expr, rhs Expr) {
	emit("# assignToInterface")
	if rhs == nil || isNil(rhs) {
		emit("# initialize interface with a zero value")
		emit("push $0")
		emit("push $0")
		emit("push $0")
		emitSaveInterface(lhs, 0)
		return
	}

	assert(rhs.getGtype() != nil, rhs.token(), fmt.Sprintf("rhs gtype is nil:%T", rhs))
	if rhs.getGtype().getPrimType() == G_INTERFACE {
		rhs.emit()
		emit("push %%rax")
		emit("push %%rbx")
		emit("push %%rcx")
		emitSaveInterface(lhs, 0)
		return
	}

	emitConversionToInterface(rhs)
	emit("push %%rax")
	emit("push %%rbx")
	emit("push %%rcx")
	emitSaveInterface(lhs, 0)
}

func assignToSlice(lhs Expr, rhs Expr) {
	emit("# assignToSlice")
	//assert(rhs == nil || rhs.getGtype().typ == G_SLICE, nil, "should be a slice literal or nil")
	if rhs == nil {
		emit("# initialize slice with a zero value")
		emit("push $0")
		emit("push $0")
		emit("push $0")
		emitSave3Elements(lhs, 0)
		return
	}

	//	assert(rhs.getGtype().getPrimType() == G_SLICE, rhs.token(), "rsh should be slice type")

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
		lit.emit()
		emit("push %%rax")
		emit("push %%rbx")
		emit("push %%rcx")
	case *ExprSlice:
		e := rhs.(*ExprSlice)
		e.emitToStack()
	case *ExprConversion:
		// https://golang.org/ref/spec#Conversions
		// Converting a value of a string type to a slice of bytes type
		// yields a slice whose successive elements are the bytes of the string.
		//
		// see also https://blog.golang.org/strings
		conversion := rhs.(*ExprConversion)
		assert(conversion.gtype.typ == G_SLICE, rhs.token(), "must be a slice of bytes")
		assert(conversion.expr.getGtype().typ == G_STRING || conversion.expr.getGtype().relation.gtype.typ == G_STRING, rhs.token(), "must be a string type, but got "+conversion.expr.getGtype().String())
		stringVarname, ok := conversion.expr.(*Relation)
		assert(ok, rhs.token(), "ok")
		stringVariable := stringVarname.expr.(*ExprVariable)
		stringVariable.emit()
		emit("push %%rax")
		strlen := &ExprLen{
			arg: stringVariable,
		}
		strlen.emit()
		emit("push %%rax # len")
		emit("push %%rax # cap")

	default:
		emit("# emit rhs of type %T %s", rhs, rhs.getGtype())
		rhs.emit() // it should put values to rax,rbx,rcx
		emit("push %%rax")
		emit("push %%rbx")
		emit("push %%rcx")
	}

	emitSave3Elements(lhs, 0)
}

func (variable *ExprVariable) saveSlice(offset int) {
	emit("# *ExprVariable.saveSlice()")
	emit("pop %%rax")
	variable.emitOffsetSave(8, offset+ptrSize+sliceOffsetForLen)
	emit("pop %%rax")
	variable.emitOffsetSave(8, offset+ptrSize)
	emit("pop %%rax")
	variable.emitOffsetSave(8, offset)
}

func (variable *ExprVariable) saveInterface(offset int) {
	emit("# *ExprVariable.saveInterface()")
	emit("pop %%rax # dynamic type id")
	variable.emitOffsetSave(8, offset+ptrSize+ptrSize)
	emit("pop %%rax # reciverTypeId")
	variable.emitOffsetSave(8, offset+ptrSize)
	emit("pop %%rax # ptr")
	variable.emitOffsetSave(8, offset)
}

// copy each element
func assignToArray(lhs Expr, rhs Expr) {
	emit("# assignToArray")
	if rel, ok := lhs.(*Relation); ok {
		lhs = rel.expr
	}

	arrayType := lhs.getGtype()
	elementType := arrayType.elementType
	elmSize := elementType.getSize()
	assert(rhs == nil || rhs.getGtype().typ == G_ARRAY, nil, "rhs should be array")
	switch {
	case elementType.typ == G_REL && elementType.relation.gtype.typ == G_STRUCT:
		//TBI
		for i := 0; i < arrayType.length; i++ {
			left := &ExprIndex{
				collection: lhs,
				index:      &ExprNumberLiteral{val: i},
			}
			if rhs == nil {
				assignToStruct(left, nil)
				continue
			}
			arrayLiteral, ok := rhs.(*ExprArrayLiteral)
			assert(ok, nil, "ok")
			assignToStruct(left, arrayLiteral.values[i])
		}
		return
	default: // prrimitive type or interface
		for i := 0; i < arrayType.length; i++ {
			offsetByIndex := i * elmSize
			switch rhs.(type) {
			case nil:
				// assign zero values
				if elementType.getPrimType() == G_INTERFACE {
					emit("push $0")
					emit("push $0")
					emit("push $0")
					emitSaveInterface(lhs, offsetByIndex)
					continue
				} else {
					emit("mov $0, %%rax")
				}
			case *ExprArrayLiteral:
				arrayLiteral := rhs.(*ExprArrayLiteral)
				if elementType.getPrimType() == G_INTERFACE {
					if i >= len(arrayLiteral.values) {
						// zero value
						emit("push $0")
						emit("push $0")
						emit("push $0")
						emitSaveInterface(lhs, offsetByIndex)
						continue
					} else if arrayLiteral.values[i].getGtype().getPrimType() != G_INTERFACE {
						// conversion of dynamic type => interface type
						dynamicValue := arrayLiteral.values[i]
						emitConversionToInterface(dynamicValue)
						emit("push %%rax")
						emit("push %%rbx")
						emit("push %%rcx")
						emitSaveInterface(lhs, offsetByIndex)
						continue
					} else {
						arrayLiteral.values[i].emit()
						emitSaveInterface(lhs, offsetByIndex)
						continue
					}
				}

				if i >= len(arrayLiteral.values) {
					// zero value
					emit("mov $0, %%rax")
				} else {
					val := arrayLiteral.values[i]
					val.emit()
				}
			case *Relation:
				rel := rhs.(*Relation)
				arrayVariable, ok := rel.expr.(*ExprVariable)
				assert(ok, nil, "ok")
				arrayVariable.emitOffsetLoad(elmSize, offsetByIndex)
			case *ExprStructField:
				strctField := rhs.(*ExprStructField)
				strctField.emitOffsetLoad(elmSize, offsetByIndex)
			default:
				TBI(rhs.token(), "no supporetd %T", rhs)
			}

			emitOffsetSave(lhs, elmSize, offsetByIndex)
		}
	}
}

// for local var
func (decl *DeclVar) emit() {
	emit("")
	emit("# DeclVar %s", decl.variable.varname)
	gtype := decl.variable.gtype
	switch {
	case gtype.typ == G_ARRAY:
		assignToArray(decl.varname, decl.initval)
	case gtype.typ == G_SLICE:
		assignToSlice(decl.varname, decl.initval)
	case gtype.typ == G_REL && gtype.relation.gtype.typ == G_STRUCT:
		assignToStruct(decl.varname, decl.initval)
	case gtype.getPrimType() == G_MAP:
		assignToMap(decl.varname, decl.initval)
	case gtype.getPrimType() == G_INTERFACE:
		assignToInterface(decl.varname, decl.initval)
	default:
		assert(decl.variable.getGtype().getSize() <= 8, decl.token(), "invalid type:"+gtype.String())
		// primitive types like int,bool,byte
		rhs := decl.initval
		if rhs == nil {
			if gtype.isString() {
				rhs = &eEmptyString
			} else {
				// assign zero value
				rhs = &ExprNumberLiteral{}
			}
		}
		rhs.emit()
		emitLsave(decl.variable.getGtype().getSize(), decl.variable.offset)
	}
}

var eEmptyString = ExprStringLiteral{}

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

func emitCollectIndexSave(array Expr, index Expr, offset int) {
	assert(array.getGtype().typ == G_ARRAY, array.token(), "should be array")
	elmType := array.getGtype().elementType
	emit("push %%rax # STACK 1 : the value") // stash value

	emit("# array.emit()")
	array.emit()                 // emit address
	emit("push %%rax # STACK 2") // store address of variable

	index.emit()
	emit("mov %%rax, %%rcx") // index

	size := elmType.getSize()
	assert(size > 0, nil, "size > 0")
	emit("mov $%d, %%rax    # size of one element", size)
	emit("imul %%rcx, %%rax # index * size")
	emit("push %%rax        # STACK 3 : store index * size")
	emit("pop %%rcx         # STACK 3: load  index * size")
	emit("pop %%rbx         # STACK 2 : load address of variable")
	emit("add %%rcx , %%rbx # (index * size) + address")
	if offset > 0 {
		emit("add $%d,  %%rbx # offset", offset)
	}
	emit("pop %%rax # STACK 1: restore the value")
	emit("mov %%rax, (%%rbx) # save the value")
	emit("")
}

func loadCollectIndex(array Expr, index Expr, offset int) {
	emit("# loadCollectIndex")
	if array.getGtype().typ == G_ARRAY {
		elmType := array.getGtype().elementType
		emit("# array.emit()")
		array.emit()       // emit address
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
		if array.getGtype().elementType.getPrimType() == G_INTERFACE {
			emit("# emit the element of interface type")
			emit("mov %%rbx, %%rdx")
			emit("mov (%%rdx), %%rax")
			emit("mov 8(%%rdx), %%rbx")
			emit("mov 16(%%rdx), %%rcx")
		} else {
			emit("# emit the element of primitive type")
			emit("mov (%%rbx), %%rax")
		}
	} else if array.getGtype().typ == G_SLICE {
		elmType := array.getGtype().elementType
		emit("# emit address of the low index")
		array.emit()       // eval pointer value
		emit("push %%rax") // store head address

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

		primType := array.getGtype().elementType.getPrimType()
		if primType == G_INTERFACE || primType == G_MAP || primType == G_SLICE {
			emit("# emit the element of interface type")
			emit("mov %%rbx, %%rdx")
			emit("mov (%%rdx), %%rax")
			emit("mov 8(%%rdx), %%rbx")
			emit("mov 16(%%rdx), %%rcx")
		} else {
			// dereference the content of an emelment
			inst := getLoadInst(size)
			emit("# emit the element of primitive type")
			emit("%s (%%rbx), %%rax", inst)
		}
	} else if array.getGtype().getPrimType() == G_MAP {
		// e.g. x[key]
		emit("# emit map index expr")
		emit("# r10: map header address")
		emit("# r11: map len")
		emit("# r12: specified index value")
		emit("# r13: loop counter")

		// rax: found value (zero if not found)
		// rcx: ok (found: address of the index,  not found:0)
		_map := array
		emit("# emit mapData head address")
		_map.emit()
		emit("mov %%rax, %%r10 # copy head address")
		emitOffsetLoad(_map, IntSize, IntSize)
		emit("mov %%rax, %%r11 # copy len ")
		index.emit()
		emit("mov %%rax, %%r12 # index value")
		emitMapGet(array.getGtype(), true)
	} else if array.getGtype().getPrimType() == G_STRING {
		// https://golang.org/ref/spec#Index_expressions
		// For a of string type:
		//
		// a constant index must be in range if the string a is also constant
		// if x is out of range at run time, a run-time panic occurs
		// a[x] is the non-constant byte value at index x and the type of a[x] is byte
		// a[x] may not be assigned to

		emit("# load head address of the string")
		array.emit()       // emit address
		emit("push %%rax") // store address of variable
		index.emit()
		emit("mov %%rax, %%rcx")  // load  index * 1
		emit("pop %%rbx")         // load address of variable
		emit("add %%rcx , %%rbx") // (index * size) + address
		if offset > 0 {
			emit("add $%d,  %%rbx", offset)
		}
		emit("mov (%%rbx), %%rax") // dereference the content of an emelment	} else {
	} else {
		TBI(array.token(), "unable to handle %s", array.getGtype())
	}
}

func emitMapGet(mapType *Gtype, deref bool) {
	if mapType.typ == G_REL {
		// @TODO handle infinite chain of relations
		mapType = mapType.relation.gtype
	}

	emit("# emitMapGet")
	emit("mov $0, %%r13 # init loop counter") // i = 0

	labelBegin := makeLabel()
	labelEnd := makeLabel()
	emit("%s: # begin loop ", labelBegin)

	labelIncr := makeLabel()
	// break if i < len
	emit("cmp %%r11, %%r13") // len > i
	emit("setl %%al")
	emit("movzb %%al, %%eax")
	emit("test %%rax, %%rax")
	emit("mov $0, %%rax") // key not found. set zero value.
	emit("mov $0, %%rcx") // ok = false
	emit("je %s  # jump if false", labelEnd)

	// check if index value matches
	emit("mov %%r13, %%rax")   // i
	emit("imul $16, %%rax")    // i * 16
	emit("mov %%r10, %%rcx")   // head
	emit("add %%rax, %%rcx")   // head + i * 16
	emit("mov (%%rcx), %%rax") // emit index address

	mapKeyType := mapType.mapKey
	assert(mapKeyType != nil, nil, "key typ should not be nil:"+mapType.String())
	if !mapKeyType.isString() {
		emit("mov (%%rax), %%rax") // dereference
	}
	if mapKeyType.isString() {
		emit("push %%r13")
		emit("push %%r11")
		emit("push %%r10")
		emit("push %%rcx")
		emitStringsEqual(true, "%r12", "%rax")
		emit("pop %%rcx")
		emit("pop %%r10")
		emit("pop %%r11")
		emit("pop %%r13")
	} else {
		// primitive comparison
		emit("cmp %%r12, %%rax # compare specifiedvalue vs indexvalue")
		emit("sete %%al")
		emit("movzb %%al, %%eax")
	}

	emit("test %%rax, %%rax")
	emit("je %s  # jump if false", labelIncr)

	emit("# Value found!")
	emit("mov 8(%%rcx), %%rax # set the found value address")
	if deref {
		emit("mov (%%rax), %%rax # dereference")
	}
	emit("jmp %s", labelEnd)

	emit("%s: # incr", labelIncr)
	emit("add $1, %%r13") // i++
	emit("jmp %s", labelBegin)

	emit("%s: # end loop", labelEnd)
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
	if e.collection.getGtype().isString() {
		// s[n:m]
		// new strlen: m - n
		var high Expr
		if e.high == nil {
			high = &ExprLen{
				tok: e.token(),
				arg: e.collection,
			}
		} else {
			high = e.high
		}
		eNewStrlen := &ExprBinop{
			tok:   e.token(),
			op:    "-",
			left:  high,
			right: e.low,
		}
		// mem size = strlen + 1
		eMemSize := &ExprBinop{
			tok:  e.token(),
			op:   "+",
			left: eNewStrlen,
			right: &ExprNumberLiteral{
				val: 1,
			},
		}

		// src address + low
		e.collection.emit()
		emit("push %%rax # src address")
		e.low.emit()
		emit("pop %%rbx")
		emit("add %%rax, %%rbx")
		emit("push %%rbx")

		emitCallMallocDinamicSize(eMemSize)
		emit("push %%rax # dst address")

		eNewStrlen.emit()
		emit("push %%rax # strlen")

		emit("pop %%%s", RegsForCall[2])
		emit("pop %%%s", RegsForCall[1])
		emit("pop %%%s", RegsForCall[0])
		emit("mov $0, %%rax")
		emit("call .strcopy")
	} else {
		e.emitToStack()
		emit("pop %%rcx")
		emit("pop %%rbx")
		emit("pop %%rax")
	}
}

func (e *ExprSlice) emitToStack() {
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

	if e.high == nil {
		e.high = &ExprNumberLiteral{
			val: e.collection.getGtype().length,
		}
	}
	calcLen := &ExprBinop{
		op:    "-",
		left:  e.high,
		right: e.low,
	}
	calcLen.emit()
	emit("push %%rax")

	emit("#   calc and set cap")
	var max Expr
	if e.max != nil {
		max = e.max
	} else {
		max = &ExprCap{
			tok: e.token(),
			arg: e.collection,
		}
	}
	calcCap := &ExprBinop{
		op:    "-",
		left:  max,
		right: e.low,
	}

	calcCap.emit()

	emit("push %%rax")
}

func (e ExprArrayLiteral) emit() {
	errorft(e.token(), "DO NOT EMIT")
}

// https://golang.org/ref/spec#Type_assertions
func (e *ExprTypeAssertion) emit() {
	assert(e.expr.getGtype().getPrimType() == G_INTERFACE, e.token(), "expr must be an Interface type")
	if e.gtype.getPrimType() == G_INTERFACE {
		TBI(e.token(), "")
	} else {
		// if T is not an interface type,
		// x.(T) asserts that the dynamic type of x is identical to the type T.

		e.expr.emit() // emit interface
		// rax(ptr), rbx(typeId of method table), rcx(hashed typeId)
		emit("push %%rax")
		// @TODO DRY with type switch statement
		typeLabel := groot.getTypeLabel(e.gtype)
		emit("lea .%s(%%rip), %%rax # type: %s", typeLabel, e.gtype)
		emitStringsEqual(true, "%rax", "%rcx")

		emit("mov %%rax, %%rbx") // move flag
		// @TODO consider big data like slice, struct, etd
		emit("pop %%rax")          // load ptr
		emit("mov (%%rax), %%rax") // deref
	}
}

func (ast *StmtContinue) emit() {
	assert(ast.stmtFor.labelEndBlock != "", ast.token(), "labelEndLoop should not be empty")
	emit("jmp %s # continue", ast.stmtFor.labelEndBlock)
}

func (ast *StmtBreak) emit() {
	assert(ast.stmtFor.labelEndLoop != "", ast.token(), "labelEndLoop should not be empty")
	emit("jmp %s # break", ast.stmtFor.labelEndLoop)
}

func (ast *StmtExpr) emit() {
	ast.expr.emit()
}

func (ast *StmtDefer) emit() {
	emit("# defer")
	/*
	// arguments should be evaluated immediately
	var args []Expr
	switch ast.expr.(type) {
	case *ExprMethodcall:
		call := ast.expr.(*ExprMethodcall)
		args = call.args
	case *ExprFuncallOrConversion:
		call := ast.expr.(*ExprFuncallOrConversion)
		args = call.args
	default:
		errorft(ast.token(), "defer should be a funcall")
	}
	*/
	labelStart := makeLabel()  + "_defer"
	labelEnd := makeLabel()  + "_defer"
	ast.label = labelStart

	emit("jmp %s", labelEnd)
	emit("%s: # defer start", labelStart)

	for i:=0 ; i < len(retRegi) ; i++ {
		emit("push %%%s", retRegi[i])
	}

	ast.expr.emit()

	for i:= len(retRegi) - 1 ; i >= 0 ; i-- {
		emit("pop %%%s", retRegi[i])
	}

	emit("leave")
	emit("ret")
	emit("%s: # defer end", labelEnd)

}

func (e *ExprVaArg) emit() {
	e.expr.emit()
}

func (e *ExprConversion) emit() {
	emitComment("ExprConversion.emit()")
	if e.gtype.isString() {
		// s = string(bytes)
		labelEnd := makeLabel()
		e.expr.emit()
		emit("push %%rax")
		emit("test %%rax, %%rax")
		emit("pop %%rax")
		emit("jne %s", labelEnd)
		emit("lea .%s(%%rip), %%rax # set empty strinf", eEmptyString.slabel)
		emit("%s:", labelEnd)
	} else {
		e.expr.emit()
	}
}

func (e *ExprStructLiteral) emit() {
	errorft(e.token(), "This cannot be emitted alone")
}

func (e *ExprTypeSwitchGuard) emit() {
	e.expr.emit()
	emit("mov %%rcx, %%rax # copy type id")
}

func (e *ExprMapLiteral) emit() {
	e.emitPush()
	emit("pop %%rcx")
	emit("pop %%rbx")
	emit("pop %%rax")
}

func (ast *ExprMethodcall) getUniqueName() string {
	gtype := ast.receiver.getGtype()
	return getMethodUniqueName(gtype, ast.fname)
}

func (methodCall *ExprMethodcall) getOrigType() *Gtype {
	gtype := methodCall.receiver.getGtype()
	assertNotNil(methodCall.receiver != nil, methodCall.token())
	assertNotNil(gtype != nil, methodCall.tok)
	assert(gtype.typ == G_REL || gtype.typ == G_POINTER || gtype.typ == G_INTERFACE, methodCall.tok, "method must be an interface or belong to a named type")
	var typeToBeloing *Gtype
	if gtype.typ == G_POINTER {
		typeToBeloing = gtype.origType
		assert(typeToBeloing != nil, methodCall.token(), "shoudl not be nil:"+gtype.String())
	} else {
		typeToBeloing = gtype
	}
	assert(typeToBeloing.typ == G_REL, methodCall.tok, "method must belong to a named type")
	origType := typeToBeloing.relation.gtype
	//debugf("origType = %v", origType)
	assert(typeToBeloing.relation.gtype != nil, methodCall.token(), fmt.Sprintf("origType should not be nil:%#v", typeToBeloing.relation))
	return origType
}

func getRettypes(call Expr) []*Gtype {
	switch call.(type) {
	case *ExprFuncallOrConversion:
		return call.(*ExprFuncallOrConversion).getRettypes()
	case *ExprMethodcall:
		return call.(*ExprMethodcall).getRettypes()
	}
	errorf("no reach here")
	return nil
}

func (funcall *ExprFuncallOrConversion) getRettypes() []*Gtype {
	if funcall.rel.gtype != nil {
		// Conversion
		return []*Gtype{funcall.rel.gtype}
	}

	return funcall.getFuncDef().rettypes
}

func (methodCall *ExprMethodcall) getRettypes() []*Gtype {
	origType := methodCall.getOrigType()
	if origType == nil {
		errorft(methodCall.token(), "origType should not be nil")
	}
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

type IrInterfaceMethodCall struct {
	receiver   Expr
	methodName identifier
}

func (call *IrInterfaceMethodCall) emit(args []Expr) {
	emit("# emit interface method call \"%s\"", call.methodName)
	if true {
		mapType := &Gtype{
			typ: G_MAP,
			mapKey: &Gtype{
				typ: G_STRING,
			},
			mapValue: &Gtype{
				typ: G_STRING,
			},
		}
		emit("# emit receiverTypeId of %s", call.receiver.getGtype())
		emitOffsetLoad(call.receiver, ptrSize, ptrSize)
		emit("imul $8, %%rax")
		emit("push %%rax")
		emit("lea receiverTypes(%%rip), %%rax")
		emit("pop %%rcx")
		emit("add %%rcx, %%rax")
		emit("# find method %s", call.methodName)
		emit("mov (%%rax), %%r10") // address of receiverType

		emit("mov $128, %%rax")  // copy len
		emit("mov %%rax, %%r11") // copy len

		emit("lea .M%s, %%rax", call.methodName) // index value
		emit("mov %%rax, %%r12")                 // index value
		emitMapGet(mapType, false)
	}

	emit("push %%rax")

	emit("# setting arguments %v", args)

	receiver := args[0]
	emit("mov $0, %%rax")
	receiverType := receiver.getGtype()
	assert(receiverType.getPrimType() == G_INTERFACE, nil, "should be interface")

	// dereference: convert an interface value to a concrete value
	receiver.emit()

	emit("mov (%%rax), %%rax")

	emit("push %%rax  # receiver")

	otherArgs := args[1:]
	for i, arg := range otherArgs {
		if _, ok := arg.(*ExprVaArg); ok {
			// skip VaArg for now
			emit("mov $0, %%rax")
		} else {
			arg.emit()
		}
		emit("push %%rax  # argument no %d", i+2)
	}

	for i, _ := range args {
		j := len(args) - 1 - i
		emit("pop %%%s   # argument no %d", RegsForCall[j], j+1)
	}

	emit("pop %%rax")
	emit("call *%%rax")
}

func (methodCall *ExprMethodcall) emitInterfaceMethodCall() {
	args := []Expr{methodCall.receiver}
	for _, arg := range methodCall.args {
		args = append(args, arg)
	}
	call := &IrInterfaceMethodCall{
		receiver:   methodCall.receiver,
		methodName: methodCall.fname,
	}
	call.emit(args)
}

func (methodCall *ExprMethodcall) emit() {
	origType := methodCall.getOrigType()
	if origType.typ == G_INTERFACE {
		methodCall.emitInterfaceMethodCall()
		return
	}

	args := []Expr{methodCall.receiver}
	for _, arg := range methodCall.args {
		args = append(args, arg)
	}

	funcref, ok := origType.methods[methodCall.fname]
	if !ok {
		errorft(methodCall.token(), "method %s is not found in type %s", methodCall.fname, methodCall.receiver.getGtype())
	}
	pkgname := funcref.funcdef.pkg
	name := methodCall.getUniqueName()
	var staticCall *IrStaticCall = &IrStaticCall{
		symbol: getFuncSymbol(pkgname, name),
		callee: funcref.funcdef,
	}
	staticCall.emit(args)
}

func (funcall *ExprFuncallOrConversion) getFuncDef() *DeclFunc {
	relexpr := funcall.rel.expr
	assert(relexpr != nil, funcall.token(), "relexpr should NOT be nil")
	funcref, ok := relexpr.(*ExprFuncRef)
	if !ok {
		errorft(funcall.token(), "Compiler error: funcref is not *ExprFuncRef but %v", funcref, funcall.fname)
	}
	assertNotNil(funcref.funcdef != nil, nil)
	return funcref.funcdef
}

func (e *ExprLen) emit() {
	emit("# emit len()")
	arg := e.arg
	gtype := arg.getGtype()
	assert(gtype != nil, e.token(), "gtype should not be  nil:\n" + fmt.Sprintf("%#v", arg))

	switch {
	case gtype.typ == G_ARRAY:
		emit("mov $%d, %%rax", gtype.length)
	case gtype.typ == G_SLICE:
		emit("# len(slice)")
		switch arg.(type) {
		case *Relation:
			emit("# Relation")
			emitOffsetLoad(arg, 8, ptrSize)
		case *ExprStructField:
			emit("# ExprStructField")
			emitOffsetLoad(arg, 8, ptrSize)
		case *ExprIndex:
			emitOffsetLoad(arg, 8, ptrSize)
		case *ExprSliceLiteral:
			emit("# ExprSliceLiteral")
			_arg := arg.(*ExprSliceLiteral)
			length := len(_arg.values)
			emit("mov $%d, %%rax", length)
		case *ExprSlice:
			sliceExpr := arg.(*ExprSlice)
			uop := &ExprBinop{
				op:    "-",
				left:  sliceExpr.high,
				right: sliceExpr.low,
			}
			uop.emit()
		default:
			TBI(arg.token(), "unable to handle %T", arg)
		}
	case gtype.getPrimType() == G_MAP:
		switch arg.(type) {
		case *Relation:
			emit("# Relation")
			emitOffsetLoad(arg, 8, ptrSize)
		case *ExprStructField:
			emit("# ExprStructField")
			emitOffsetLoad(arg, 8, ptrSize)
		case *ExprMapLiteral:
			TBI(arg.token(), "unable to handle %T", arg)
		default:
			TBI(arg.token(), "unable to handle %T", arg)
		}
	case gtype.getPrimType() == G_STRING:
		arg.emit()
		emit("mov %%rax, %%rdi")
		emit("mov $0, %%rax")
		emit("call strlen")
	default:
		TBI(arg.token(), "unable to handle %s", gtype)
	}
}

func (e *ExprCap) emit() {
	emit("# emit cap()")
	arg := e.arg
	gtype := arg.getGtype()
	switch {
	case gtype.typ == G_ARRAY:
		emit("mov $%d, %%rax", gtype.length)
	case gtype.typ == G_SLICE:
		switch arg.(type) {
		case *Relation:
			emit("# Relation")
			emitOffsetLoad(arg, 8, ptrSize*2)
		case *ExprStructField:
			emit("# ExprStructField")
			emitOffsetLoad(arg, 8, ptrSize*2)
		case *ExprIndex:
			emitOffsetLoad(arg, 8, ptrSize*2)
		case *ExprSliceLiteral:
			emit("# ExprSliceLiteral")
			_arg := arg.(*ExprSliceLiteral)
			length := len(_arg.values)
			emit("mov $%d, %%rax", length)
		case *ExprSlice:
			sliceExpr := arg.(*ExprSlice)
			if sliceExpr.collection.getGtype().typ == G_ARRAY {
				cp := &ExprBinop{
					tok: e.tok,
					op:  "-",
					left: &ExprLen{
						tok: e.tok,
						arg: sliceExpr.collection,
					},
					right: sliceExpr.low,
				}
				cp.emit()
			} else {
				TBI(arg.token(), "unable to handle %T", arg)
			}
		default:
			TBI(arg.token(), "unable to handle %T", arg)
		}
	case gtype.getPrimType() == G_MAP:
		TBI(arg.token(), "unable to handle %T", arg)
	case gtype.getPrimType() == G_STRING:
		TBI(arg.token(), "unable to handle %T", arg)
	default:
		TBI(arg.token(), "unable to handle %s", gtype)
	}
}

func (funcall *ExprFuncallOrConversion) emit() {
	if funcall.rel.expr == nil && funcall.rel.gtype != nil {
		// Conversion
		conversion := &ExprConversion{
			tok:   funcall.token(),
			gtype: funcall.rel.gtype,
			expr:  funcall.args[0],
		}
		conversion.emit()
		return
	}

	assert(funcall.rel.expr != nil && funcall.rel.gtype == nil, funcall.token(), "this is conversion")
	assert(funcall.getFuncDef() != nil, funcall.token(), "funcdef is nil")
	decl := funcall.getFuncDef()

	// check if it's a builtin function
	switch decl {
	case builinLen:
		assert(len(funcall.args) == 1, funcall.token(), "invalid arguments for len()")
		arg := funcall.args[0]
		exprLen := &ExprLen{
			tok: arg.token(),
			arg: arg,
		}
		exprLen.emit()
	case builinCap:
		arg := funcall.args[0]
		e := &ExprCap{
			tok: arg.token(),
			arg: arg,
		}
		e.emit()
	case builtinAppend:
		assert(len(funcall.args) == 2, funcall.token(), "append() should take 2 argments")
		slice := funcall.args[0]
		valueToAppend := funcall.args[1]
		emit("# append(%s, %s)", slice.getGtype(), valueToAppend.getGtype())
		var staticCall *IrStaticCall = &IrStaticCall{
			callee: decl,
		}
		switch slice.getGtype().elementType.getSize() {
		case 1:
			staticCall.symbol = getFuncSymbol("", "append1")
			staticCall.emit(funcall.args)
		case 8:
			staticCall.symbol = getFuncSymbol("", "append8")
			staticCall.emit(funcall.args)
		case 24:
			if slice.getGtype().elementType.getPrimType() == G_INTERFACE && valueToAppend.getGtype().getPrimType() != G_INTERFACE {
				eConvertion := &ExprConversionToInterface{
					tok: valueToAppend.token(),
					expr: valueToAppend,
				}
				funcall.args[1] = eConvertion
			}
			staticCall.symbol = getFuncSymbol("", "append24")
			staticCall.emit(funcall.args)
		default:
			TBI(slice.token(), "")
		}
	case builtinDumpInterface:
		arg := funcall.args[0]

		emit("lea .%s, %%rax", builtinStringKey1)
		emit("push %%rax")
		arg.emit()
		emit("push %%rax  # interface ptr")
		emit("push %%rbx  # interface receverTypeId")
		emit("push %%rcx  # interface dynamicTypeId")

		numRegs := 4
		for i := numRegs - 1; i >= 0; i-- {
			emit("pop %%%s   # RegsForCall[%d]", RegsForCall[i], i)
		}

		emit("mov $0, %%rax")
		emit("call %s", "printf")
		emit("")
	case builtinAsComment:
		arg := funcall.args[0]
		if stringLiteral, ok := arg.(*ExprStringLiteral); ok {
			emit("# %s", stringLiteral.val)
		}
	default:
		var staticCall *IrStaticCall = &IrStaticCall{
			symbol: getFuncSymbol(decl.pkg, funcall.fname),
			callee: decl,
		}
		staticCall.emit(funcall.args)
	}
}

type IrStaticCall struct {
	symbol string
	callee *DeclFunc
}

func (ircall *IrStaticCall) emit(args []Expr) {
	// nothing to do
	emit("")
	emit("# emitCall %s", ircall.symbol)

	var numRegs int
	var param *ExprVariable
	var collectVariadicArgs bool // gather variadic args into a slice
	var variadicArgs []Expr
	var arg Expr
	var i int
	for i, arg = range args {
		if i < len(ircall.callee.params) {
			param = ircall.callee.params[i]
			if param.isVariadic {
				if ircall.symbol != "fmt.Printf" {
					// ignore fmt.Printf variadic
					if _, ok := arg.(*ExprVaArg); !ok {
						collectVariadicArgs = true
					}
				}
			}
		}

		if collectVariadicArgs {
			variadicArgs = append(variadicArgs, arg)
			continue
		}

		emit("# arg %d, collectVariadicArgs=%v", i, collectVariadicArgs)
		if param != nil {
			emit("# %s <- %T", param.getGtype(), arg.getGtype())
		}
		arg.emit()

		var primType GTYPE_TYPE = 0
		if arg.getGtype() != nil {
			primType = arg.getGtype().getPrimType()
		}
		emit("#")
		if primType == G_SLICE || primType == G_INTERFACE || primType == G_MAP {
			emit("push %%rax  # argument 1/3")
			emit("push %%rbx  # argument 2/3")
			emit("push %%rcx  # argument 3/3")
			numRegs += sliceWidth
		} else {
			emit("push %%rax  # argument primitive")
			numRegs += 1
		}
	}

	// check if callee has a variadic
	// https://golang.org/ref/spec#Passing_arguments_to_..._parameters
	// If f is invoked with no actual arguments for p, the value passed to p is nil.
	if !collectVariadicArgs {
		if i+1 < len(ircall.callee.params) {
			param = ircall.callee.params[i+1]
			if param.isVariadic {
				collectVariadicArgs = true
			}
		}
	}

	if collectVariadicArgs {
		emit("# passing variadic args")
		lenArgs := len(variadicArgs)
		if lenArgs == 0 {
			// pass an empty slice
			emit("push $0")
			emit("push $0")
			emit("push $0")
		} else {
			// var a []interface{}
			for i, varg := range variadicArgs {
				if i == 0 {
					// make an empty slice
					emit("push $0")
					emit("push $0")
					emit("push $0")
				}
				// conversion : var ifc = x
				emitConversionToInterface(varg)
				emit("push %%rax")
				emit("push %%rbx")
				emit("push %%rcx")
				emit("# calling append24")
				emit("pop %%r9  # ifc_c")
				emit("pop %%r8  # ifc_b")
				emit("pop %%rcx # ifc_a")
				emit("pop %%rdx # cap")
				emit("pop %%rsi # len")
				emit("pop %%rdi # ptr")
				emit("mov $0, %%rax")
				emit("call .append24")
				emit("push %%rax # slice.ptr")
				emit("push %%rbx # slice.len")
				emit("push %%rcx # slice.cap")
			}
		}
		numRegs += 3
	}

	for i := numRegs - 1; i >= 0; i-- {
		if i >= len(RegsForCall) {
			errorft(args[0].token(), "too many arguments")
		}
		emit("pop %%%s   # RegsForCall[%d]", RegsForCall[i], i)
	}

	emit("mov $0, %%rax")
	if ircall.symbol == "fmt.Printf" {
		// replace fmt.Printf by libc's printf
		emit("call printf")
	} else {
		emit("call %s", ircall.symbol)
	}
	emit("")
}

func emitRuntimeArgs() {
	emitLabel(".runtime_args:")
	emit("push %%rbp")
	emit("mov %%rsp, %%rbp")

	emit("# set argv, argc, argc")
	emit("mov runtimeArgv(%%rip), %%rax # ptr")
	emit("mov runtimeArgc(%%rip), %%rbx # len")
	emit("mov runtimeArgc(%%rip), %%rcx # cap")

	emitFuncEpilogue(".runtime_args_noop_handler",nil)
}

func emitMainFunc(importOS bool) {
	fname := "main"
	emit(".global	%s", fname)
	emitLabel("%s:", fname)
	emit("push %%rbp")
	emit("mov %%rsp, %%rbp")

	emit("mov %%rsi, runtimeArgv(%%rip)")
	emit("mov %%rdi, runtimeArgc(%%rip)")
	emit("mov $0, %%rsi")
	emit("mov $0, %%rdi")

	// init runtime
	emit("# init runtime")
	emit("mov $0, %%rax")
	emit("call .init")

	// init imported packages
	if importOS {
		emit("# init os")
		emit("mov $0, %%rax")
		emit("call os.init")
	}

	emit("")
	emit("mov $0, %%rax")
	emit("call main.main")
	emitFuncEpilogue("noop_handler", nil,)
}

func (f *DeclFunc) emit() {
	f.emitPrologue()
	f.body.emit()
	emit("mov $0, %%rax")
	emitFuncEpilogue(f.labelDeferHandler, f.stmtDefer)
}

func evalIntExpr(e Expr) int {
	switch e.(type) {
	case nil:
		errorf("e is nil")
	case *ExprNumberLiteral:
		return e.(*ExprNumberLiteral).val
	case *ExprVariable:
		errorft(e.token(), "variable cannot be inteppreted at compile time :%#v", e)
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

func emitGlobalDeclInit(ptok *Token /* left type */, gtype *Gtype, value /* nullable */ Expr, containerName string) {
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
						switch value.(type) {
						case *Relation:
							rel := value.(*Relation)
							vr,ok := rel.expr.(*ExprVariable)
							if !ok {
								errorft(value.token(), "cannot compile")
							}
							emit(".quad %s # %s %s", vr.varname, value.getGtype(), selector)
						default:
							emit(".quad %d # %s %s", evalIntExpr(value), value.getGtype(), selector)
						}
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
		case nil:
			return
		case *ExprSliceLiteral:
			// initialize a hidden array
			lit := value.(*ExprSliceLiteral)
			lit.invisiblevar.varname = identifier(fmt.Sprintf("$hiddenArray$%d", getHidddenArrayId()))
			emit(".quad %s", lit.invisiblevar.varname)      // address of the hidden array
			emit(".quad %d", lit.invisiblevar.gtype.length) // len
			emit(".quad %d", lit.invisiblevar.gtype.length) // cap
			arrayLiteral := &ExprArrayLiteral{
				gtype:  lit.invisiblevar.gtype,
				values: lit.values,
			}
			arrayDecl := &DeclVar{
				tok:      ptok,
				variable: lit.invisiblevar,
				initval:  arrayLiteral,
			}
			arrayDecl.emitGlobal()

		default:
			TBI(ptok, "unable to handle T=%s, value=%#v", gtype, value)
		}
	} else if gtype.typ == G_MAP {
		// @TODO
		emit(".quad 0")
		emit(".quad 0")
		emit(".quad 0")
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
			emit("# field:%s", field.fieldname)
			if value == nil {
				emitGlobalDeclInit(ptok, field, nil, containerName+"."+string(field.fieldname))
				continue
			}
			structLiteral, ok := value.(*ExprStructLiteral)
			assert(ok, nil, "ok:"+containerName)
			value := structLiteral.lookup(field.fieldname)
			if value == nil {
				// zero value
				//continue
			}
			gtype := field
			emitGlobalDeclInit(ptok, gtype, value, containerName+"."+string(field.fieldname))
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
		case *ExprBinop:
			val = evalIntExpr(value)
			emit(".quad %d # %s ", val, gtype)
		case *ExprStringLiteral:
			stringLiteral := value.(*ExprStringLiteral)
			emit(".quad .%s", stringLiteral.slabel)
		case *Relation:
			rel := value.(*Relation)
			emit(".quad 0 # (TBI) rel:%s", rel.name)
		case *ExprUop:
			uop := value.(*ExprUop)
			assert(uop.op == "&", ptok, "only uop & is allowed")
			operand := uop.operand
			rel, ok := operand.(*Relation)
			if ok {
				assert(ok, value.token(), "operand should be *Relation")
				vr, ok := rel.expr.(*ExprVariable)
				assert(ok, value.token(), "operand should be a variable")
				assert(vr.isGlobal, value.token(), "operand should be a global variable")
				emit(".quad %s", vr.varname)
			} else {
				// var gv = &Struct{_}
				entityLabel :=  makeLabel() + "_global_entity"
				emit(".quad %s", entityLabel)
				entity := &GlobalInternalEntity{
					token: ptok,
					label: entityLabel,
					expr: operand,
				}
				globalInternalEntities = append(globalInternalEntities, entity)
			}
		default:
			TBI(ptok, "unable to handle %T", value)
		}
	}
}

var globalInternalEntities []*GlobalInternalEntity
type GlobalInternalEntity struct {
	token *Token
	label string
	expr Expr
}

func (e *GlobalInternalEntity) emit() {
	emitLabel(e.label + ":")
	emitGlobalDeclInit(e.token, e.expr.getGtype(), e.expr, "")
}

func (decl *DeclVar) emitGlobal() {
	emitLabel("# emitGlobal for %s" , decl.variable.varname)
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
	emit("# initval=%#v", right)
	emitGlobalDeclInit(ptok, right.getGtype(), right, "")
}

type IrRoot struct {
	vars           []*DeclVar
	funcs          []*DeclFunc
	stringLiterals []*ExprStringLiteral
	methodTable    map[int][]string
	hashedTypes    map[string]int
	importOS       bool
}

var groot *IrRoot

func (root *IrRoot) getTypeLabel(gtype *Gtype) string {
	dynamicTypeId := root.hashedTypes[gtype.String()]
	label := fmt.Sprintf("DT%d", dynamicTypeId)
	return label
}

// builtin string
var builtinStringKey1 string = "SfmtDumpInterface"
var builtinStringValue1 string = "# interface = {ptr:%p,receiverTypeId:%d,dtype:'%s'}\\n"

func (root *IrRoot) emit() {
	groot = root
	// generate code
	emit(".data")

	emit("")
	emitComment("STRING LITERALS")

	// emit builtin string
	emitLabel(".%s:", builtinStringKey1)
	emit(".string \"%s\"", builtinStringValue1)

	// empty string
	eEmptyString.slabel = fmt.Sprintf("S%d", 0)
	emitLabel(".%s:", eEmptyString.slabel)
	emit(".string \"%s\"", eEmptyString.val)

	for id, ast := range root.stringLiterals {
		ast.slabel = fmt.Sprintf("S%d", id+1)
		emitLabel(".%s:", ast.slabel)
		// https://sourceware.org/binutils/docs-2.30/as/String.html#String
		// the assembler marks the end of each string with a 0 byte.
		emit(".string \"%s\"", ast.val)
	}

	emit("")
	emitComment("Dynamic Types")
	var dynamicTypeId int // 0 means nil
	for hashedType, _ := range root.hashedTypes {
		dynamicTypeId++
		root.hashedTypes[hashedType] = dynamicTypeId
		label := fmt.Sprintf("DT%d", dynamicTypeId)
		emitLabel(".%s:", label)
		emit(".string \"%s\"", hashedType)
	}

	emitComment("Method table")

	emitLabel("%s:", "receiverTypes")
	emit(".quad 0 # receiverTypeId:0")
	for i := 1; i <= len(root.methodTable); i++ {
		emit(".quad receiverType%d # receiverTypeId:%d, %s", i, i, root.methodTable[i])
	}

	var shortMethodNames map[string]string = map[string]string{}

	for i := 1; i <= len(root.methodTable); i++ {
		emitLabel("receiverType%d:", i)
		for _, methodNameFull := range root.methodTable[i] {
			splitted := strings.Split(methodNameFull, "$")
			shortMethodName := splitted[1]
			emit(".quad .M%s # key", shortMethodName)
			emit(".quad %s # method", methodNameFull)
			shortMethodNames[shortMethodName] = shortMethodName
		}
	}

	emitComment("METHOD NAMES")
	for shortMethodName := range shortMethodNames {
		emitLabel(".M%s:", shortMethodName)
		emit(".string \"%s\"", shortMethodName)
	}

	emitComment("GLOBAL VARS")
	emit("")
	for _, vardecl := range root.vars {
		vardecl.emitGlobal()
	}

	// @TODO do this infinitly
	for _, entity := range globalInternalEntities {
		entity.emit()
	}


	emitComment("FUNCTIONS")
	emit(".text")
	for _, funcdecl := range root.funcs {
		funcdecl.emit()
	}

	emitRuntimeArgs()
	emitMainFunc(root.importOS)
}
