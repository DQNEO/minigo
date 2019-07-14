// Code generator
// Convention:
//  We SHOULD use the word "emit" for the meaning of "output assembly code",
//  NOT for "load something to %rax".
//  Such usage would make much confusion.

package main

var offset0 int = 0
var offset8 int = 8
var offset16 int = 16

const IntSize int = 8 // 64-bit (8 bytes)
const ptrSize int = 8
const sliceWidth int = 3
const interfaceWidth int = 3
const mapWidth int = 3
const sliceSize int = IntSize + ptrSize + ptrSize

func emitNewline() {
	writePos()
	write(S("\n"))
}

var pos *Token // current source position

func setPos(ptok *Token) {
	pos = ptok
}

func writePos() {
	if !emitPosition {
		return
	}
	var spos bytes
	if pos == nil {
		spos = S("()")
	} else {
		spos = pos.String()
	}
	var b []byte
	b = concat3(S("/*"), spos, S("*/"))
	write(b)
}

var gasIndentLevel int = 1

func emit(format bytes, v ...interface{}) {
	writePos()

	for i := 0; i < gasIndentLevel; i++ {
		write(bytes("  "))
	}

	s := Sprintf(bytes(format), v...)
	writeln(s)
}

func emitWithoutIndent(format bytes, v ...interface{}) {
	writePos()
	s := Sprintf(bytes(format), v...)
	writeln(s)
}

func unwrapRel(e Expr) Expr {
	if rel, ok := e.(*Relation); ok {
		return rel.expr
	}
	return e
}

// Mytype.method -> Mytype#method
func getMethodUniqueName(gtype *Gtype, fname goidentifier) bytes {
	assertNotNil(gtype != nil, nil)
	var typename goidentifier
	if gtype.kind == G_POINTER {
		typename = gtype.origType.relation.name
	} else {
		typename = gtype.relation.name
	}
	return Sprintf(S("%s$%s"), bytes(typename), bytes(fname))
}

// "main","f1" -> "main.f1"
func getFuncSymbol(pkg bytes, fname bytes) bytes {
	if eq(pkg, S("libc")) {
		return fname
	}
	if len(pkg) == 0 {
		pkg = bytes("")
	}
	return Sprintf(S("%s.%s"), pkg, fname)
}

func (f *DeclFunc) getSymbol() bytes {
	if f.receiver != nil {
		// method
		fname := goidentifier(f.fname)
		return getFuncSymbol(bytes(f.pkg), getMethodUniqueName(f.receiver.gtype, fname))
	}

	// other functions
	return getFuncSymbol(bytes(f.pkg), bytes(f.fname))
}

func align(n int, m int) int {
	remainder := n % m
	if remainder == 0 {
		return n
	} else {
		return n - remainder + m
	}
}

func emitFuncEpilogue(labelDeferHandler bytes, stmtDefer *StmtDefer) {
	emitNewline()
	emit(S("# func epilogue"))
	// every function has a defer handler
	emit(S("%s: # defer handler"), labelDeferHandler)

	// if the function has a defer statement, jump to there
	if stmtDefer != nil {
		emit(S("jmp %s"), stmtDefer.label)
	}

	emit(S("LEAVE_AND_RET"))
}

func emit_intcast(gtype *Gtype) {
	if gtype.getKind() == G_BYTE {
		emit(S("CAST_BYTE_TO_INT"))
	}
}

func emit_comp_primitive(inst bytes, binop *ExprBinop) {
	emit(S("# emit_comp_primitive"))
	assert(len(inst) > 0 , binop.token(), S("inst shoud not be empty"))
	binop.left.emit()
	if binop.left.getGtype().getKind() == G_BYTE {
		emit_intcast(binop.left.getGtype())
	}
	emit(S("PUSH_8 # left")) // left
	binop.right.emit()
	if binop.right.getGtype().getKind() == G_BYTE {
		emit_intcast(binop.right.getGtype())
	}
	emit(S("PUSH_8 # right")) // right
	emit(S("CMP_FROM_STACK %s"), inst)
}

var labelSeq int = 0

func makeLabel() bytes {
	r := Sprintf(S(".L%d"), labelSeq)
	labelSeq++
	return r
}

func (ast *StmtInc) emit() {
	emitIncrDecl(S("ADD_NUMBER 1"), ast.operand)
}
func (ast *StmtDec) emit() {
	emitIncrDecl(S("SUB_NUMBER 1"), ast.operand)
}

// https://golang.org/ref/spec#IncDecStmt
// As with an assignment, the operand must be addressable or a map index expression.
func emitIncrDecl(inst bytes, operand Expr) {
	operand.emit()
	emit(inst)

	left := operand
	emitSavePrimitive(left)
}


func (binop *ExprBinop) emitComp() {
	emit(S("# emitComp"))
	assert(binop.left != nil, binop.token(), S("should not be nil"))

	assert(!binop.left.getGtype().isString(), binop.token(), S("should not be clike string"))


	var instruction bytes
	op := string(binop.op)
	switch op {
	case "<":
		instruction = S("setl")
	case ">":
		instruction = S("setg")
	case "<=":
		instruction = S("setle")
	case ">=":
		instruction = S("setge")
	case "!=":
		instruction = S("setne")
	case "==":
		instruction = S("sete")
	default:
		assertNotReached(binop.token())
	}

	emit_comp_primitive(instruction, binop)
}

func (ast *ExprBinop) emit() {
	if eq(ast.op , bytes("+")) && ast.left.getGtype().isString() {
		TBI(ast.token(), S("concat strings"))
		return
	}
	op := string(ast.op)
	switch op {
	case "<", ">", "<=", ">=", "!=", "==":
		ast.emitComp()
		return
	case "&&":
		labelEnd := makeLabel()
		ast.left.emit()
		emit(S("TEST_IT"))
		emit(S("LOAD_NUMBER 0"))
		emit(S("je %s"), labelEnd)
		ast.right.emit()
		emit(S("TEST_IT"))
		emit(S("LOAD_NUMBER 0"))
		emit(S("je %s"), labelEnd)
		emit(S("LOAD_NUMBER 1"))
		emit(S("%s:"), labelEnd)
		return
	case "||":
		labelEnd := makeLabel()
		ast.left.emit()
		emit(S("TEST_IT"))
		emit(S("LOAD_NUMBER 1"))
		emit(S("jne %s"), labelEnd)
		ast.right.emit()
		emit(S("TEST_IT"))
		emit(S("LOAD_NUMBER 1"))
		emit(S("jne %s"), labelEnd)
		emit(S("LOAD_NUMBER 0"))
		emit(S("%s:"), labelEnd)
		return
	}
	ast.left.emit()
	emit(S("PUSH_8"))
	ast.right.emit()
	emit(S("PUSH_8"))

	op = string(ast.op)
	switch op {
	case "+":
		emit(S("SUM_FROM_STACK"))
	case "-":
		emit(S("SUB_FROM_STACK"))
	case "*":
		emit(S("IMUL_FROM_STACK"))
	case "%":
		emit(S("pop %%rcx"))
		emit(S("pop %%rax"))
		emit(S("mov $0, %%rdx # init %%rdx"))
		emit(S("div %%rcx"))
		emit(S("mov %%rdx, %%rax"))
	case"/":
		emit(S("pop %%rcx"))
		emit(S("pop %%rax"))
		emit(S("mov $0, %%rdx # init %%rdx"))
		emit(S("div %%rcx"))
	default:
		errorft(ast.token(), S("Unknown binop: %s"), op)
	}
}

func isUnderScore(e Expr) bool {
	rel, ok := e.(*Relation)
	if !ok {
		return false
	}
	return eq(bytes(rel.name), bytes("_"))
}

// expect rhs address is in the stack top, lhs is in the second top
func emitCopyStructFromStack(size int) {
	emit(S("pop %%rbx")) // to
	emit(S("pop %%rax")) // from

	var i int
	for ; i < size; i += 8 {
		emit(S("movq %d(%%rbx), %%rcx"), i)
		emit(S("movq %%rcx, %d(%%rax)"), i)
	}
	for ; i < size; i += 4 {
		emit(S("movl %d(%%rbx), %%rcx"), i)
		emit(S("movl %%rcx, %d(%%rax)"), i)
	}
	for ; i < size; i++ {
		emit(S("movb %d(%%rbx), %%rcx"), i)
		emit(S("movb %%rcx, %d(%%rax)"), i)
	}
}

const sliceOffsetForLen = 8

func emitCallMallocDinamicSize(eSize Expr) {
	eSize.emit()
	emit(S("PUSH_8"))
	emit(S("POP_TO_ARG_0"))
	emit(S("FUNCALL iruntime.malloc"))
}

func emitCallMalloc(size int) {
	eNumber := &ExprNumberLiteral{
		val: size,
	}
	emitCallMallocDinamicSize(eNumber)
}

func (e *IrExprConversionToInterface) emit() {
	emit(S("# IrExprConversionToInterface"))
	emitConversionToInterface(e.arg)
}

func emitConversionToInterface(dynamicValue Expr) {
	receiverType := dynamicValue.getGtype()
	if receiverType == nil {
		emit(S("# receiverType is nil. emit nil for interface"))
		emit(S("LOAD_EMPTY_INTERFACE"))
		return
	}

	emit(S("# emitConversionToInterface from %s"), dynamicValue.getGtype().String())
	dynamicValue.emit()
	if dynamicValue.getGtype().is24WidthType() {
		emit(S("PUSH_24"))
		emitCallMalloc(24)
		emit(S("PUSH_8"))
		emit(S("STORE_24_INDIRECT_FROM_STACK"))
	} else {
		emit(S("PUSH_8"))
		emitCallMalloc(8)
		emit(S("PUSH_8"))
		emit(S("STORE_8_INDIRECT_FROM_STACK"))
	}

	emit(S("PUSH_8 # addr of dynamicValue")) // address

	if receiverType.kind == G_POINTER {
		receiverType = receiverType.origType.relation.gtype
	}
	//assert(receiverType.receiverTypeId > 0,  dynamicValue.token(), S("no receiverTypeId"))
	emit(S("LOAD_NUMBER %d # receiverTypeId"), receiverType.receiverTypeId)
	emit(S("PUSH_8 # receiverTypeId"))

	gtype := dynamicValue.getGtype()
	label := symbolTable.getTypeLabel(gtype)
	emit(S("lea .%s, %%rax# dynamicType %s"), label, gtype.String())
	emit(S("PUSH_8 # dynamicType"))

	emit(S("POP_INTERFACE"))
	emitNewline()
}

func isNil(e Expr) bool {
	e = unwrapRel(e)
	_, isNil := e.(*ExprNilLiteral)
	return isNil
}

func (decl *DeclVar) emit() {
	if decl.variable.isGlobal {
		decl.emitGlobal()
	} else {
		decl.emitLocal()
	}
}

func (decl *DeclVar) emitLocal() {
	emit(S("# DeclVar \"%s\""), decl.variable.varname)
	gtype := decl.variable.gtype
	variable := decl.variable
	rhs := decl.initval
	switch gtype.getKind() {
	case G_ARRAY:
		assignToArray(variable, rhs)
	case G_SLICE, G_STRING:
		assignToSlice(variable, rhs)
	case G_STRUCT:
		assignToStruct(variable, rhs)
	case G_INTERFACE:
		assignToInterface(variable, rhs)
	default:
		emitAssignPrimitive(variable, rhs)
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
		setPos(ast.token())
		emit(S("# Statement"))
		gasIndentLevel++
		stmt.emit()
		gasIndentLevel--
	}
}

func (e *ExprIndex) emit() {
	emit(S("# emit *ExprIndex"))
	e.emitOffsetLoad(0)
}

func (e *ExprNilLiteral) emit() {
	emit(S("LOAD_NUMBER 0 # nil literal"))
}

func (s *StmtShortVarDecl) emit() {
	// this emitter cannot be removed due to lack of for.cls.init conversion
	a := &StmtAssignment{
		tok:    s.tok,
		lefts:  s.lefts,
		rights: s.rights,
	}
	a.emit()
}

func (f *ExprFuncRef) emit() {
	emit(S("LOAD_NUMBER 1 # funcref")) // emit 1 for now.  @FIXME
}

func (e ExprArrayLiteral) emit() {
	errorft(e.token(), S("DO NOT EMIT"))
}

// https://golang.org/ref/spec#Type_assertions
func (e *ExprTypeAssertion) emit() {
	assert(e.expr.getGtype().getKind() == G_INTERFACE, e.token(), S("expr must be an Interface type"))
	if e.gtype.getKind() == G_INTERFACE {
		TBI(e.token(), S("type assertion"))
	} else {
		// if T is not an interface type,
		// x.(T) asserts that the dynamic type of x is identical to the type T.

		e.expr.emit() // emit interface
		// rax(ptr), rbx(receiverTypeId of method table), rcx(gtype as astring)
		emit(S("PUSH_8 # push dynamic data"))

		emit(S("push %%rcx # push dynamic type addr"))
		emitCompareDynamicTypeFromStack(e.gtype)

		// move ok value
		if e.gtype.is24WidthType() {
			emit(S("mov %%rax, %%rdx"))
		} else {
			emit(S("mov %%rax, %%rbx"))
		}
		emit(S("pop %%rax # load dynamic data"))
		emit(S("TEST_IT"))
		labelEnd := makeLabel()
		emit(S("je %s # exit if nil"), labelEnd)
		if e.gtype.is24WidthType() {
			emit(S("LOAD_24_BY_DEREF"))
		} else {
			emit(S("LOAD_8_BY_DEREF"))
		}
		emitWithoutIndent(S("%s:"), labelEnd)
	}
}

func (ast *StmtExpr) emit() {
	setPos(ast.token())
	ast.expr.emit()
}

func (e *ExprVaArg) emit() {
	e.expr.emit()
}

func (e *IrExprConversion) emit() {
	emit(S("# IrExprConversion.emit()"))
	if  e.arg.getGtype().isBytesSlice() && e.toGtype.isString() {
		// cstring(bytes)
		emit(S("# convert slice to string"))
		labelEnd := makeLabel()
		e.arg.emit()    // load slice
		emit(S("TEST_IT")) // check if ptr is nil
		emit(S("jne %s # exit if not nil"), labelEnd)
		emit(S("# if nil then"))
		emitEmptyString()
		emit(S("%s:"), labelEnd)
	} else {
		e.arg.emit()
	}
}

func (e *ExprStructLiteral) emit() {
	errorft(e.token(), S("This cannot be emitted alone"))
}

func (e *ExprTypeSwitchGuard) emit() {
	e.expr.emit()
}

func bool2string(bol bool) bytes {
	if bol {
		return S("true")
	} else {
		return S("false")
	}
}

func (f *DeclFunc) emit() {
	f.prologue.emit()
	f.body.emit()
	emit(S("mov $0, %%rax"))
	emitFuncEpilogue(f.labelDeferHandler, f.stmtDefer)
}

func evalIntExpr(e Expr) int {
	e = unwrapRel(e)

	switch e.(type) {
	case nil:
		errorf(S("e is nil"))
	case *ExprNumberLiteral:
		return e.(*ExprNumberLiteral).val
	case *ExprVariable:
		errorft(e.token(), S("variable cannot be inteppreted at compile time :%#v"), e)
	case *ExprBinop:
		binop := e.(*ExprBinop)
		op := string(binop.op)
		switch op {
		case "+":
			return evalIntExpr(binop.left) + evalIntExpr(binop.right)
		case "-":
			return evalIntExpr(binop.left) - evalIntExpr(binop.right)
		case "*":
			return evalIntExpr(binop.left) * evalIntExpr(binop.right)

		}
	case *ExprConstVariable:
		cnst := e.(*ExprConstVariable)
		if cnst.hasIotaValue() {
			return cnst.iotaIndex
		}
		return evalIntExpr(cnst.val)
	default:
		errorft(e.token(), S("unkown type %T"), e)
	}
	return 0
}

func (cnst *ExprConstVariable) hasIotaValue() bool {
	rel, ok := cnst.val.(*Relation)
	if !ok {
		return false
	}

	val := rel.expr.(*ExprConstVariable)
	return val == eIota
}
