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
	write([]byte("\n"))
}

var pos *Token // current source position

func setPos(ptok *Token) {
	pos = ptok
}

func writePos() {
	if !emitPosition {
		return
	}
	var spos string
	if pos != nil {
		spos = pos.String()
	}
	var s string
	s = concat3("/*", spos, "*/")
	write([]byte(s))
}

var gasIndentLevel int = 1

func emit(format string, v ...interface{}) {
	writePos()

	for i := 0; i < gasIndentLevel; i++ {
		write([]byte("  "))
	}

	s := Sprintf(format, v...)
	writeln([]byte(s))
}

func emitWithoutIndent(format string, v ...interface{}) {
	writePos()
	s := Sprintf(format, v...)
	writeln([]byte(s))
}

func unwrapRel(e Expr) Expr {
	if rel, ok := e.(*Relation); ok {
		return rel.expr
	}
	return e
}

// Mytype.method -> Mytype#method
func getMethodUniqueName(gtype *Gtype, fname identifier) string {
	assertNotNil(gtype != nil, nil)
	var typename identifier
	if gtype.kind == G_POINTER {
		typename = gtype.origType.relation.name
	} else {
		typename = gtype.relation.name
	}
	s := Sprintf("%s$%s", typename, fname)
	return s
}

// "main","f1" -> "main.f1"
func getFuncSymbol(pkg identifier, fname string) string {
	if eq(string(pkg), "libc") {
		return fname
	}
	s := Sprintf("%s.%s", pkg, fname)
	return s
}

func (f *DeclFunc) getSymbol() string {
	if f.receiver != nil {
		// method
		fname := f.fname
		return getFuncSymbol(f.pkg, getMethodUniqueName(f.receiver.gtype, fname))
	}

	// other functions
	return getFuncSymbol(f.pkg, string(f.fname))
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
	emitNewline()
	emit("# func epilogue")
	// every function has a defer handler
	emit("%s: # defer handler", labelDeferHandler)

	// if the function has a defer statement, jump to there
	if stmtDefer != nil {
		emit("jmp %s", stmtDefer.label)
	}

	emit("LEAVE_AND_RET")
}

func emit_intcast(gtype *Gtype) {
	if gtype.getKind() == G_BYTE {
		emit("CAST_BYTE_TO_INT")
	}
}

func emit_comp_primitive(inst string, binop *ExprBinop) {
	emit("# emit_comp_primitive")
	assert(len(inst) > 0 , binop.token(), "inst shoud not be empty")
	binop.left.emit()
	if binop.left.getGtype().getKind() == G_BYTE {
		emit_intcast(binop.left.getGtype())
	}
	emit("PUSH_8 # left") // left
	binop.right.emit()
	if binop.right.getGtype().getKind() == G_BYTE {
		emit_intcast(binop.right.getGtype())
	}
	emit("PUSH_8 # right") // right
	emit("CMP_FROM_STACK %s", inst)
}

var labelSeq int = 0

func makeLabel() string {
	r := Sprintf(".L%d", labelSeq)
	labelSeq++
	return r
}

func (ast *StmtInc) emit() {
	emitIncrDecl("ADD_NUMBER 1", ast.operand)
}
func (ast *StmtDec) emit() {
	emitIncrDecl("SUB_NUMBER 1", ast.operand)
}

// https://golang.org/ref/spec#IncDecStmt
// As with an assignment, the operand must be addressable or a map index expression.
func emitIncrDecl(inst string, operand Expr) {
	operand.emit()
	emit(inst)

	left := operand
	emitSavePrimitive(left)
}


func (binop *ExprBinop) emitComp() {
	emit("# emitComp")
	assert(binop.left != nil, binop.token(), "should not be nil")

	assert(!binop.left.getGtype().isString(), binop.token(), "should not be string")


	var instruction string
	op := binop.op
	switch op {
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
	default:
		assertNotReached(binop.token())
	}

	emit_comp_primitive(instruction, binop)
}

func (ast *ExprBinop) emit() {
	if eq(ast.op , "+") && ast.left.getGtype().isString() {
		TBI(ast.token(), "concat strings")
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
		emit("TEST_IT")
		emit("LOAD_NUMBER 0")
		emit("je %s", labelEnd)
		ast.right.emit()
		emit("TEST_IT")
		emit("LOAD_NUMBER 0")
		emit("je %s", labelEnd)
		emit("LOAD_NUMBER 1")
		emit("%s:", labelEnd)
		return
	case "||":
		labelEnd := makeLabel()
		ast.left.emit()
		emit("TEST_IT")
		emit("LOAD_NUMBER 1")
		emit("jne %s", labelEnd)
		ast.right.emit()
		emit("TEST_IT")
		emit("LOAD_NUMBER 1")
		emit("jne %s", labelEnd)
		emit("LOAD_NUMBER 0")
		emit("%s:", labelEnd)
		return
	}
	ast.left.emit()
	emit("PUSH_8")
	ast.right.emit()
	emit("PUSH_8")

	op = string(ast.op)
	switch op {
	case "+":
		emit("SUM_FROM_STACK")
	case "-":
		emit("SUB_FROM_STACK")
	case "*":
		emit("IMUL_FROM_STACK")
	case "%":
		emit("pop %%rcx")
		emit("pop %%rax")
		emit("mov $0, %%rdx # init %%rdx")
		emit("div %%rcx")
		emit("mov %%rdx, %%rax")
	case"/":
		emit("pop %%rcx")
		emit("pop %%rax")
		emit("mov $0, %%rdx # init %%rdx")
		emit("div %%rcx")
	default:
		errorft(ast.token(), "Unknown binop: %s", op)
	}
}

func isUnderScore(e Expr) bool {
	rel, ok := e.(*Relation)
	if !ok {
		return false
	}
	return eq(string(rel.name), "_")
}

// expect rhs address is in the stack top, lhs is in the second top
func emitCopyStructFromStack(size int) {
	emit("pop %%rbx") // to
	emit("pop %%rax") // from

	var i int
	for ; i < size; i += 8 {
		emit("movq %d(%%rbx), %%rcx", i)
		emit("movq %%rcx, %d(%%rax)", i)
	}
	for ; i < size; i += 4 {
		emit("movl %d(%%rbx), %%rcx", i)
		emit("movl %%rcx, %d(%%rax)", i)
	}
	for ; i < size; i++ {
		emit("movb %d(%%rbx), %%rcx", i)
		emit("movb %%rcx, %d(%%rax)", i)
	}
}

const sliceOffsetForLen = 8

func emitCallMallocDinamicSize(eSize Expr) {
	eSize.emit()
	emit("PUSH_8")
	emit("POP_TO_ARG_0")
	emit("FUNCALL iruntime.malloc")
}

func emitCallMalloc(size int) {
	eNumber := &ExprNumberLiteral{
		val: size,
	}
	emitCallMallocDinamicSize(eNumber)
}

func (e *IrExprConversionToInterface) emit() {
	emit("# IrExprConversionToInterface")
	emitConversionToInterface(e.arg)
}

func emitConversionToInterface(dynamicValue Expr) {
	receiverType := dynamicValue.getGtype()
	if receiverType == nil {
		emit("# receiverType is nil. emit nil for interface")
		emit("LOAD_EMPTY_INTERFACE")
		return
	}

	emit("# emitConversionToInterface from %s", dynamicValue.getGtype().String())
	dynamicValue.emit()
	if dynamicValue.getGtype().is24WidthType() {
		emit("PUSH_24")
		emitCallMalloc(24)
		emit("PUSH_8")
		emit("STORE_24_INDIRECT_FROM_STACK")
	} else {
		emit("PUSH_8")
		emitCallMalloc(8)
		emit("PUSH_8")
		emit("STORE_8_INDIRECT_FROM_STACK")
	}

	emit("PUSH_8 # addr of dynamicValue") // address

	if receiverType.kind == G_POINTER {
		receiverType = receiverType.origType.relation.gtype
	}
	//assert(receiverType.receiverTypeId > 0,  dynamicValue.token(), "no receiverTypeId")
	emit("LOAD_NUMBER %d # receiverTypeId", receiverType.receiverTypeId)
	emit("PUSH_8 # receiverTypeId")

	gtype := dynamicValue.getGtype()
	label := symbolTable.getTypeLabel(gtype)
	emit("lea .%s, %%rax# dynamicType %s", label, gtype.String())
	emit("PUSH_8 # dynamicType")

	emit("POP_INTERFACE")
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
	emit("# DeclVar \"%s\"", decl.variable.varname)
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
		emit("# Statement")
		gasIndentLevel++
		stmt.emit()
		gasIndentLevel--
	}
}

func (e *ExprIndex) emit() {
	emit("# emit *ExprIndex")
	e.emitOffsetLoad(0)
}

func (e *ExprNilLiteral) emit() {
	emit("LOAD_NUMBER 0 # nil literal")
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
	emit("LOAD_NUMBER 1 # funcref") // emit 1 for now.  @FIXME
}

func (e ExprArrayLiteral) emit() {
	errorft(e.token(), "DO NOT EMIT")
}

// https://golang.org/ref/spec#Type_assertions
func (e *ExprTypeAssertion) emit() {
	assert(e.expr.getGtype().getKind() == G_INTERFACE, e.token(), "expr must be an Interface type")
	if e.gtype.getKind() == G_INTERFACE {
		TBI(e.token(), "type assertion")
	} else {
		// if T is not an interface type,
		// x.(T) asserts that the dynamic type of x is identical to the type T.

		e.expr.emit() // emit interface
		// rax(ptr), rbx(receiverTypeId of method table), rcx(gtype as astring)
		emit("PUSH_8 # push dynamic data")

		emit("push %%rcx # push dynamic type addr")
		emitCompareDynamicTypeFromStack(e.gtype)

		// move ok value
		if e.gtype.is24WidthType() {
			emit("mov %%rax, %%rdx")
		} else {
			emit("mov %%rax, %%rbx")
		}
		emit("pop %%rax # load dynamic data")
		emit("TEST_IT")
		labelEnd := makeLabel()
		emit("je %s # exit if nil", labelEnd)
		if e.gtype.is24WidthType() {
			emit("LOAD_24_BY_DEREF")
		} else {
			emit("LOAD_8_BY_DEREF")
		}
		emitWithoutIndent("%s:", labelEnd)
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
	emit("# IrExprConversion.emit()")
	if  e.arg.getGtype().isBytesSlice() && e.toGtype.isString() {
		// cstring(bytes)
		emit("# convert slice to string")
		labelEnd := makeLabel()
		e.arg.emit()    // load slice
		emit("TEST_IT") // check if ptr is nil
		emit("jne %s # exit if not nil", labelEnd)
		emit("# if nil then")
		emitEmptyString()
		emit("%s:", labelEnd)
	} else {
		e.arg.emit()
	}
}

func (e *ExprStructLiteral) emit() {
	errorft(e.token(), "This cannot be emitted alone")
}

func (e *ExprTypeSwitchGuard) emit() {
	e.expr.emit()
}

func bool2string(bol bool) string {
	if bol {
		return "true"
	} else {
		return "false"
	}
}

func (f *DeclFunc) emit() {
	f.prologue.emit()
	f.body.emit()
	emit("mov $0, %%rax")
	emitFuncEpilogue(f.labelDeferHandler, f.stmtDefer)
}

func evalIntExpr(e Expr) int {
	e = unwrapRel(e)

	switch e.(type) {
	case nil:
		errorf("e is nil")
	case *ExprNumberLiteral:
		return e.(*ExprNumberLiteral).val
	case *ExprVariable:
		errorft(e.token(), "variable cannot be inteppreted at compile time :%#v", e)
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
		errorft(e.token(), "unkown type %T", e)
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
