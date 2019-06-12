// Code generator
// Convention:
//  We SHOULD use the word "emit" for the meaning of "output assembly code",
//  NOT for "load something to %rax".
//  Such usage would make much confusion.

package main

import (
	"fmt"
	"os"
)

const IntSize int = 8 // 64-bit (8 bytes)
const ptrSize int = 8
const sliceWidth int = 3
const interfaceWidth int = 3
const mapWidth int = 3
const sliceSize int = IntSize + ptrSize + ptrSize

func emitNewline() {
	var b []byte = []byte{'\n'}
	os.Stdout.Write(b)
}

func write(s string) {
	var b []byte = []byte(s)
	os.Stdout.Write(b)
}

func writef(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	var b []byte = []byte(s)
	os.Stdout.Write(b)
}

var gasIndentLevel int = 1

func emit(format string, v ...interface{}) {
	var format2 string = format

	write("/* tok */ ")
	for i := 0; i < gasIndentLevel; i++ {
		write("  ")
	}

	frmt := format2 + "\n"
	writef(frmt, v...)
}

func emitWithoutIndent(format string, v ...interface{}) {
	write("/* tok */ ")
	writef(format + "\n", v...)
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
	return string(typename) + "$" + string(fname)
}

// "main","f1" -> "main.f1"
func getFuncSymbol(pkg identifier, fname string) string {
	if pkg == "libc" {
		return fname
	}
	if pkg == "" {
		pkg = ""
	}
	return fmt.Sprintf("%s.%s", pkg, fname)
}

func (f *DeclFunc) getSymbol() string {
	if f.receiver != nil {
		// method
		return getFuncSymbol(f.pkg, getMethodUniqueName(f.receiver.gtype, f.fname))
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

var labelSeq = 0

func makeLabel() string {
	r := fmt.Sprintf(".L%d", labelSeq)
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

	if ast.op == "+" {
		emit("SUM_FROM_STACK")
	} else if ast.op == "-" {
		emit("SUB_FROM_STACK")
	} else if ast.op == "*" {
		emit("IMUL_FROM_STACK")
	} else if ast.op == "%" {
		emit("pop %%rcx")
		emit("pop %%rax")
		emit("mov $0, %%rdx # init %%rdx")
		emit("div %%rcx")
		emit("mov %%rdx, %%rax")
	} else if ast.op == "/" {
		emit("pop %%rcx")
		emit("pop %%rax")
		emit("mov $0, %%rdx # init %%rdx")
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

func (e *ExprConversionToInterface) emit() {
	emit("# ExprConversionToInterface")
	emitConversionToInterface(e.expr)
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
	emit("PUSH_8")
	emitCallMalloc(8)
	emit("PUSH_8")
	emit("STORE_8_INDIRECT_FROM_STACK")
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
	rel, ok := e.(*Relation)
	if ok {
		_, isNil := rel.expr.(*ExprNilLiteral)
		return isNil
	}
	return false
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
	varname := decl.varname
	switch gtype.getKind() {
	case G_ARRAY:
		assignToArray(varname, decl.initval)
	case G_SLICE:
		assignToSlice(varname, decl.initval)
	case G_STRUCT:
		assignToStruct(varname, decl.initval)
	case G_MAP:
		assignToMap(varname, decl.initval)
	case G_INTERFACE:
		assignToInterface(varname, decl.initval)
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
		emit("# LOAD RHS")
		gasIndentLevel++
		rhs.emit()
		gasIndentLevel--
		comment := "initialize " + string(decl.variable.varname)
		emit("# Assign to LHS")
		gasIndentLevel++
		emit("STORE_%d_TO_LOCAL %d # %s",
			decl.variable.getGtype().getSize(), decl.variable.offset, comment)
		gasIndentLevel--
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
		//emit("# Statement: %s", stmt.token())
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

func (ast *StmtShortVarDecl) emit() {
	a := &StmtAssignment{
		tok:    ast.tok,
		lefts:  ast.lefts,
		rights: ast.rights,
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
		TBI(e.token(), "")
	} else {
		// if T is not an interface type,
		// x.(T) asserts that the dynamic type of x is identical to the type T.

		e.expr.emit() // emit interface
		// rax(ptr), rbx(receiverTypeId of method table), rcx(hashed receiverTypeId)
		emit("PUSH_8")
		// @TODO DRY with type switch statement
		typeLabel := symbolTable.getTypeLabel(e.gtype)
		emit("lea .%s(%%rip), %%rax # type: %s", typeLabel, e.gtype.String())

		emit("push %%rcx") // @TODO ????
		emit("PUSH_8")
		emitStringsEqualFromStack(true)

		emit("mov %%rax, %%rbx") // move flag @TODO: this is BUG in slice,map cases
		// @TODO consider big data like slice, struct, etd
		emit("pop %%rax # load ptr")
		emit("TEST_IT")
		labelEnd := makeLabel()
		emit("je %s # jmp if nil", labelEnd)
		emit("LOAD_8_BY_DEREF")
		emitWithoutIndent("%s:", labelEnd)
	}
}

func (ast *StmtExpr) emit() {
	ast.expr.emit()
}

func (e *ExprVaArg) emit() {
	e.expr.emit()
}

func (e *ExprConversion) emit() {
	emit("# ExprConversion.emit()")
	if e.gtype.isString() {
		// s = string(bytes)
		labelEnd := makeLabel()
		e.expr.emit()
		emit("TEST_IT")
		emit("jne %s", labelEnd)
		emitEmptyString()
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

func bool2string(bol bool) string {
	if bol {
		return "true"
	} else {
		return "false"
	}
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
		cnst := e.(*ExprConstVariable)
		constVal, ok := cnst.val.(*Relation)
		if ok && constVal.name == "iota" {
			val, ok := constVal.expr.(*ExprConstVariable)
			if ok && val == eIota {
				return cnst.iotaIndex
			}
		}
		return evalIntExpr(cnst.val)
	default:
		errorft(e.token(), "unkown type %T", e)
	}
	return 0
}

