// Code generator
// Convention:
//  We SHOULD use the word "emit" for the meaning of "output assembly code",
//  NOT for "load something to %rax".
//  Such usage would make much confusion.

package main

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
	var spos gostring
	if pos == nil {
		spos = S("()")
	} else {
		spos = pos.GoString()
	}
	var b []byte
	b = concat3(S("/*"), spos, S("*/"))
	write(b)
}

var gasIndentLevel int = 1

func emit2(format string, v ...interface{}) {
	var frmt gostring = gostring(format)
	writePos()

	for i := 0; i < gasIndentLevel; i++ {
		write(gostring("  "))
	}

	s := GoSprintf2(frmt, v...)
	writeln(s)
}

func emit(format string, v ...interface{}) {
	var frmt gostring = gostring(format)
	writePos()

	for i := 0; i < gasIndentLevel; i++ {
		write(gostring("  "))
	}

	s := GoSprintf(frmt, v...)
	writeln(s)
}

func emitWithoutIndent2(format string, v ...interface{}) {
	writePos()
	s := GoSprintf2(gostring(format), v...)
	writeln(s)
}

func unwrapRel(e Expr) Expr {
	if rel, ok := e.(*Relation); ok {
		return rel.expr
	}
	return e
}

// Mytype.method -> Mytype#method
func getMethodUniqueName(gtype *Gtype, fname identifier) gostring {
	assertNotNil(gtype != nil, nil)
	var typename identifier
	if gtype.kind == G_POINTER {
		typename = gtype.origType.relation.name
	} else {
		typename = gtype.relation.name
	}
	return GoSprintf2(S("%s$%s"), gostring(typename), gostring(fname))
}

// "main","f1" -> "main.f1"
func getFuncSymbol(pkg gostring, fname gostring) gostring {
	if eq(pkg, "libc") {
		return fname
	}
	if len(pkg) == 0 {
		pkg = gostring("")
	}
	return GoSprintf2(S("%s.%s"), pkg, fname)
}

func (f *DeclFunc) getSymbol() gostring {
	if f.receiver != nil {
		// method
		return getFuncSymbol(gostring(f.pkg), getMethodUniqueName(f.receiver.gtype, f.fname))
	}

	// other functions
	return getFuncSymbol(gostring(f.pkg), gostring(f.fname))
}

func align(n int, m int) int {
	remainder := n % m
	if remainder == 0 {
		return n
	} else {
		return n - remainder + m
	}
}

func emitFuncEpilogue(labelDeferHandler gostring, stmtDefer *StmtDefer) {
	emitNewline()
	emit2("# func epilogue")
	// every function has a defer handler
	emit2("%s: # defer handler", labelDeferHandler)

	// if the function has a defer statement, jump to there
	if stmtDefer != nil {
		emit2("jmp %s", stmtDefer.label)
	}

	emit2("LEAVE_AND_RET")
}

func emit_intcast(gtype *Gtype) {
	if gtype.getKind() == G_BYTE {
		emit2("CAST_BYTE_TO_INT")
	}
}

func emit_comp_primitive(inst gostring, binop *ExprBinop) {
	emit2("# emit_comp_primitive")
	assert(len(inst) > 0 , binop.token(), "inst shoud not be empty")
	binop.left.emit()
	if binop.left.getGtype().getKind() == G_BYTE {
		emit_intcast(binop.left.getGtype())
	}
	emit2("PUSH_8 # left") // left
	binop.right.emit()
	if binop.right.getGtype().getKind() == G_BYTE {
		emit_intcast(binop.right.getGtype())
	}
	emit2("PUSH_8 # right") // right
	emit2("CMP_FROM_STACK %s", inst)
}

var labelSeq int = 0

func makeLabel() gostring {
	r := GoSprintf2(S(".L%d"), labelSeq)
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
	emit2(inst)

	left := operand
	emitSavePrimitive(left)
}


func (binop *ExprBinop) emitComp() {
	emit2("# emitComp")
	if binop.left.getGtype().isString() {
		e := &IrExprStringComparison{
			tok: binop.token(),
			op: binop.op,
			cstringLeft: binop.left,
			cstringRight: binop.right,
		}
		e.emit()
		return
	}

	var instruction gostring
	op := binop.op
	switch cstring(op) {
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
	if eqGostrings(ast.op , gostring("+")) && ast.left.getGtype().isString() {
		emitStringConcate(ast.left, ast.right)
		return
	}
	switch cstring(ast.op) {
	case "<", ">", "<=", ">=", "!=", "==":
		ast.emitComp()
		return
	case "&&":
		labelEnd := makeLabel()
		ast.left.emit()
		emit2("TEST_IT")
		emit2("LOAD_NUMBER 0")
		emit2("je %s", labelEnd)
		ast.right.emit()
		emit2("TEST_IT")
		emit2("LOAD_NUMBER 0")
		emit2("je %s", labelEnd)
		emit2("LOAD_NUMBER 1")
		emit2("%s:", labelEnd)
		return
	case "||":
		labelEnd := makeLabel()
		ast.left.emit()
		emit2("TEST_IT")
		emit2("LOAD_NUMBER 1")
		emit2("jne %s", labelEnd)
		ast.right.emit()
		emit2("TEST_IT")
		emit2("LOAD_NUMBER 1")
		emit2("jne %s", labelEnd)
		emit2("LOAD_NUMBER 0")
		emit2("%s:", labelEnd)
		return
	}
	ast.left.emit()
	emit2("PUSH_8")
	ast.right.emit()
	emit2("PUSH_8")

	op := ast.op
	switch cstring(op) {
	case "+":
		emit2("SUM_FROM_STACK")
	case "-":
		emit2("SUB_FROM_STACK")
	case "*":
		emit2("IMUL_FROM_STACK")
	case "%":
		emit2("pop %%rcx")
		emit2("pop %%rax")
		emit2("mov $0, %%rdx # init %%rdx")
		emit2("div %%rcx")
		emit2("mov %%rdx, %%rax")
	case"/":
		emit2("pop %%rcx")
		emit2("pop %%rax")
		emit2("mov $0, %%rdx # init %%rdx")
		emit2("div %%rcx")
	default:
		errorft(ast.token(), "Unknown binop: %s", op)
	}
}

func isUnderScore(e Expr) bool {
	rel, ok := e.(*Relation)
	if !ok {
		return false
	}
	return rel.name == "_"
}

// expect rhs address is in the stack top, lhs is in the second top
func emitCopyStructFromStack(size int) {
	emit2("pop %%rbx") // to
	emit2("pop %%rax") // from

	var i int
	for ; i < size; i += 8 {
		emit2("movq %d(%%rbx), %%rcx", i)
		emit2("movq %%rcx, %d(%%rax)", i)
	}
	for ; i < size; i += 4 {
		emit2("movl %d(%%rbx), %%rcx", i)
		emit2("movl %%rcx, %d(%%rax)", i)
	}
	for ; i < size; i++ {
		emit2("movb %d(%%rbx), %%rcx", i)
		emit2("movb %%rcx, %d(%%rax)", i)
	}
}

const sliceOffsetForLen = 8

func emitCallMallocDinamicSize(eSize Expr) {
	eSize.emit()
	emit2("PUSH_8")
	emit2("POP_TO_ARG_0")
	emit2("FUNCALL iruntime.malloc")
}

func emitCallMalloc(size int) {
	eNumber := &ExprNumberLiteral{
		val: size,
	}
	emitCallMallocDinamicSize(eNumber)
}

func (e *IrExprConversionToInterface) emit() {
	emit2("# IrExprConversionToInterface")
	emitConversionToInterface(e.arg)
}

func emitConversionToInterface(dynamicValue Expr) {
	receiverType := dynamicValue.getGtype()
	if receiverType == nil {
		emit2("# receiverType is nil. emit nil for interface")
		emit2("LOAD_EMPTY_INTERFACE")
		return
	}

	emit2("# emitConversionToInterface from %s", dynamicValue.getGtype().String())
	dynamicValue.emit()
	if dynamicValue.getGtype().is24WidthType() {
		emit2("PUSH_24")
		emitCallMalloc(24)
		emit2("PUSH_8")
		emit2("STORE_24_INDIRECT_FROM_STACK")
	} else {
		emit2("PUSH_8")
		emitCallMalloc(8)
		emit2("PUSH_8")
		emit2("STORE_8_INDIRECT_FROM_STACK")
	}

	emit2("PUSH_8 # addr of dynamicValue") // address

	if receiverType.kind == G_POINTER {
		receiverType = receiverType.origType.relation.gtype
	}
	//assert(receiverType.receiverTypeId > 0,  dynamicValue.token(), "no receiverTypeId")
	emit2("LOAD_NUMBER %d # receiverTypeId", receiverType.receiverTypeId)
	emit2("PUSH_8 # receiverTypeId")

	gtype := dynamicValue.getGtype()
	label := symbolTable.getTypeLabel(gtype)
	emit2("lea .%s, %%rax# dynamicType %s", label, gtype.String())
	emit2("PUSH_8 # dynamicType")

	emit2("POP_INTERFACE")
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
	emit2("# DeclVar \"%s\"", decl.variable.varname)
	gtype := decl.variable.gtype
	variable := decl.variable
	rhs := decl.initval
	switch gtype.getKind() {
	case G_ARRAY:
		assignToArray(variable, rhs)
	case G_SLICE:
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
		emit2("# Statement")
		gasIndentLevel++
		stmt.emit()
		gasIndentLevel--
	}
}

func (e *ExprIndex) emit() {
	emit2("# emit *ExprIndex")
	e.emitOffsetLoad(0)
}

func (e *ExprNilLiteral) emit() {
	emit2("LOAD_NUMBER 0 # nil literal")
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
	emit2("LOAD_NUMBER 1 # funcref") // emit 1 for now.  @FIXME
}

func (e ExprArrayLiteral) emit() {
	errorft(e.token(), "DO NOT EMIT")
}

// https://golang.org/ref/spec#Type_assertions
func (e *ExprTypeAssertion) emit() {
	assert(e.expr.getGtype().getKind() == G_INTERFACE, e.token(), "expr must be an Interface type")
	if e.gtype.getKind() == G_INTERFACE {
		TBI(e.token(), "")
	} else {
		// if T is not an interface type,
		// x.(T) asserts that the dynamic type of x is identical to the type T.

		e.expr.emit() // emit interface
		// rax(ptr), rbx(receiverTypeId of method table), rcx(gtype as astring)
		emit2("PUSH_8 # push dynamic data")

		emit2("push %%rcx # push dynamic type addr")
		emitCompareDynamicTypeFromStack(e.gtype)

		// move ok value
		if e.gtype.is24WidthType() {
			emit2("mov %%rax, %%rdx")
		} else {
			emit2("mov %%rax, %%rbx")
		}
		emit2("pop %%rax # load dynamic data")
		emit2("TEST_IT")
		labelEnd := makeLabel()
		emit2("je %s # exit if nil", labelEnd)
		if e.gtype.is24WidthType() {
			emit2("LOAD_24_BY_DEREF")
		} else {
			emit2("LOAD_8_BY_DEREF")
		}
		emitWithoutIndent2("%s:", labelEnd)
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
	emit2("# IrExprConversion.emit()")
	if  e.arg.getGtype().isBytesSlice() && e.toGtype.isString() {
		emit2("# convert slice to string")
		// string(bytes)
		labelEnd := makeLabel()
		e.arg.emit() // load slice
		emit2("TEST_IT") // check if ptr is nil
		emit2("jne %s # exit if not nil", labelEnd)
		emit2("# if nil then")
		emitEmptyString()
		emit2("%s:", labelEnd)
	} else if e.arg.getGtype().isString() && e.toGtype.isBytesSlice() {
		//  []byte(cstring)
		cstring := e.arg
		emitConvertCstringToSlice(cstring)
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
	emit2("mov $0, %%rax")
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
		switch cstring(binop.op) {
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
