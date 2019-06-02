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

func emitOut(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	var b []byte = []byte(s)
	os.Stdout.Write(b)
}

var gasIndentLevel int = 1

func emit(format string, v ...interface{}) {
	var format2 string = format

	for i := 0; i < gasIndentLevel; i++ {
		format2 = "  " + format2
	}

	frmt := format2+"\n"
	emitOut(frmt, v...)
}

func emitWithoutIndent(format string, v ...interface{}) {
	frmt := format+"\n"
	emitOut(frmt, v...)
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

func (structfield *ExprStructField) calcOffset() {
	fieldType := structfield.getGtype()
	if fieldType.offset != undefinedSize {
		return
	}

	structType := structfield.strct.getGtype()
	switch structType.getKind() {
	case G_POINTER:
		origType := structType.origType.relation.gtype
		if origType.size == undefinedSize {
			origType.calcStructOffset()
		}
	case G_STRUCT:
		structType.calcStructOffset()
	default:
		errorf("invalid case")
	}

	if fieldType.offset == undefinedSize {
		errorf("filed type %s [named %s] offset should not be minus.", fieldType.String(), structfield.fieldname)
	}
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
	emitSave(left)
}

func (binop *ExprBinop) emitCompareStrings() {
	emit("# emitCompareStrings")
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

	labelElse := makeLabel()
	labelEnd := makeLabel()

	binop.left.emit()

	// convert nil to the empty string
	emit("CMP_EQ_ZERO")
	emit("TEST_IT")
	emit("LOAD_NUMBER 0")
	emit("je %s", labelElse)
	emitEmptyString()
	emit("jmp %s", labelEnd)
	emit("%s:", labelElse)
	binop.left.emit()
	emit("%s:", labelEnd)
	emit("PUSH_8")

	binop.right.emit()
	emit("PUSH_8")
	emitStringsEqualFromStack(equal)
}

func emitConvertNilToEmptyString() {
	emit("# emitConvertNilToEmptyString")

	emit("PUSH_8")
	emit("# convert nil to an empty string")
	emit("TEST_IT")
	emit("pop %%rax")
	labelEnd := makeLabel()
	emit("jne %s # jump if not nil", labelEnd)
	emit("# if nil then")
	emitEmptyString()
	emit("%s:", labelEnd)
}

// call strcmp
func emitStringsEqualFromStack(equal bool) {
	emit("pop %%rax") // left

	emitConvertNilToEmptyString()
	emit("mov %%rax, %%rcx")
	emit("pop %%rax # right string")
	emit("push %%rcx")
	emitConvertNilToEmptyString()

	emit("PUSH_8")
	emit("POP_TO_ARG_0")
	emit("POP_TO_ARG_1")
	emit("FUNCALL strcmp")
	if equal {
		emit("CMP_EQ_ZERO") // retval == 0
	} else {
		emit("CMP_NE_ZERO") // retval != 0
	}
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

func emitStringConcate(left Expr, right Expr) {
	emit("# emitStringConcate")
	left.emit()
	emit("PUSH_8 # left string")

	emit("PUSH_8")
	emit("POP_TO_ARG_0")
	emit("FUNCALL strlen # get left len")

	emit("PUSH_8 # left len")
	right.emit()
	emit("PUSH_8 # right string")
	emit("PUSH_8")
	emit("POP_TO_ARG_0")
	emit("FUNCALL strlen # get right len")
	emit("PUSH_8 # right len")

	emit("pop %%rax # right len")
	emit("pop %%rcx # right string")
	emit("pop %%rbx # left len")
	emit("pop %%rdx # left string")

	emit("push %%rcx # right string")
	emit("push %%rdx # left  string")

	// newSize = strlen(left) + strlen(right) + 1
	emit("add %%rax, %%rbx # len + len")
	emit("add $1, %%rbx # + 1 (null byte)")
	emit("mov %%rbx, %%rax")
	emit("PUSH_8")
	emit("POP_TO_ARG_0")
	emit("FUNCALL iruntime.malloc")

	emit("PUSH_8")
	emit("POP_TO_ARG_0")
	emit("POP_TO_ARG_1")
	emit("FUNCALL strcat")

	emit("PUSH_8")
	emit("POP_TO_ARG_0")
	emit("POP_TO_ARG_1")
	emit("FUNCALL strcat")
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
	label := groot.getTypeLabel(gtype)
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
	switch {
	case gtype.kind == G_ARRAY:
		assignToArray(varname, decl.initval)
	case gtype.kind == G_SLICE:
		assignToSlice(varname, decl.initval)
	case gtype.kind == G_NAMED && gtype.relation.gtype.kind == G_STRUCT:
		assignToStruct(varname, decl.initval)
	case gtype.getKind() == G_MAP:
		assignToMap(varname, decl.initval)
	case gtype.getKind() == G_INTERFACE:
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

var eEmptyString = ExprStringLiteral{
	val: "",
}

func (decl *DeclType) emit() {
	// nothing to do
}

func (decl *DeclConst) emit() {
	// nothing to do
}

func (ast *StmtSatementList) emit() {
	for _, stmt := range ast.stmts {
		emit("# Statement")
		gasIndentLevel++
		stmt.emit()
		gasIndentLevel--
	}
}

func emitEmptyString() {
	eEmpty := &eEmptyString
	eEmpty.emit()
}

func (e *ExprIndex) emit() {
	emit("# emit *ExprIndex")
	loadCollectIndex(e.collection, e.index, 0)
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
		typeLabel := groot.getTypeLabel(e.gtype)
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

func (ast *ExprMethodcall) getUniqueName() string {
	gtype := ast.receiver.getGtype()
	return getMethodUniqueName(gtype, ast.fname)
}

func (methodCall *ExprMethodcall) getOrigType() *Gtype {
	gtype := methodCall.receiver.getGtype()
	assertNotNil(methodCall.receiver != nil, methodCall.token())
	assertNotNil(gtype != nil, methodCall.tok)
	assert(gtype.kind == G_NAMED || gtype.kind == G_POINTER || gtype.kind == G_INTERFACE, methodCall.tok, "method must be an interface or belong to a named type")
	var typeToBeloing *Gtype
	if gtype.kind == G_POINTER {
		typeToBeloing = gtype.origType
		assert(typeToBeloing != nil, methodCall.token(), "shoudl not be nil:"+gtype.String())
	} else {
		typeToBeloing = gtype
	}
	assert(typeToBeloing.kind == G_NAMED, methodCall.tok, "method must belong to a named type")
	origType := typeToBeloing.relation.gtype
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
	if origType.kind == G_INTERFACE {
		return origType.imethods[methodCall.fname].rettypes
	} else {
		funcref, ok := origType.methods[methodCall.fname]
		if !ok {
			errorft(methodCall.token(), "method %s is not found in type %s", methodCall.fname, methodCall.receiver.getGtype().String())
		}
		return funcref.funcdef.rettypes
	}
}

type IrInterfaceMethodCall struct {
	receiver   Expr
	methodName identifier
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
	if origType.kind == G_INTERFACE {
		methodCall.emitInterfaceMethodCall()
		return
	}

	args := []Expr{methodCall.receiver}
	for _, arg := range methodCall.args {
		args = append(args, arg)
	}

	funcref, ok := origType.methods[methodCall.fname]
	if !ok {
		errorft(methodCall.token(), "method %s is not found in type %s", methodCall.fname, methodCall.receiver.getGtype().String())
	}
	pkgname := funcref.funcdef.pkg
	name := methodCall.getUniqueName()
	var staticCall *IrStaticCall = &IrStaticCall{
		symbol:       getFuncSymbol(pkgname, name),
		callee:       funcref.funcdef,
		isMethodCall: true,
	}
	staticCall.emit(args)
}

func (funcall *ExprFuncallOrConversion) getFuncDef() *DeclFunc {
	relexpr := funcall.rel.expr
	assert(relexpr != nil, funcall.token(), fmt.Sprintf("relexpr should NOT be nil for %s", funcall.fname))
	funcref, ok := relexpr.(*ExprFuncRef)
	if !ok {
		errorft(funcall.token(), "Compiler error: funcref is not *ExprFuncRef (%s)", funcall.fname)
	}
	assertNotNil(funcref.funcdef != nil, nil)
	return funcref.funcdef
}

func (e *ExprLen) emit() {
	emit("# emit len()")
	arg := e.arg
	gtype := arg.getGtype()
	assert(gtype != nil, e.token(), "gtype should not be  nil:\n"+fmt.Sprintf("%#v", arg))

	switch {
	case gtype.kind == G_ARRAY:
		emit("LOAD_NUMBER %d", gtype.length)
	case gtype.kind == G_SLICE:
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
			emit("LOAD_NUMBER %d", length)
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
	case gtype.getKind() == G_MAP:
		emit("# emit len(map)")
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
	case gtype.getKind() == G_STRING:
		arg.emit()
		emit("PUSH_8")
		emit("POP_TO_ARG_0")
		emit("FUNCALL strlen")
	default:
		TBI(arg.token(), "unable to handle %s", gtype)
	}
}

func (e *ExprCap) emit() {
	emit("# emit cap()")
	arg := e.arg
	gtype := arg.getGtype()
	switch {
	case gtype.kind == G_ARRAY:
		emit("LOAD_NUMBER %d", gtype.length)
	case gtype.kind == G_SLICE:
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
			emit("LOAD_NUMBER %d", length)
		case *ExprSlice:
			sliceExpr := arg.(*ExprSlice)
			if sliceExpr.collection.getGtype().kind == G_ARRAY {
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
	case gtype.getKind() == G_MAP:
		TBI(arg.token(), "unable to handle %T", arg)
	case gtype.getKind() == G_STRING:
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
	case builtinLen:
		assert(len(funcall.args) == 1, funcall.token(), "invalid arguments for len()")
		arg := funcall.args[0]
		exprLen := &ExprLen{
			tok: arg.token(),
			arg: arg,
		}
		exprLen.emit()
	case builtinCap:
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
		emit("# append(%s, %s)", slice.getGtype().String(), valueToAppend.getGtype().String())
		var staticCall *IrStaticCall = &IrStaticCall{
			callee: decl,
		}
		switch slice.getGtype().elementType.getSize() {
		case 1:
			staticCall.symbol = getFuncSymbol("iruntime", "append1")
			staticCall.emit(funcall.args)
		case 8:
			staticCall.symbol = getFuncSymbol("iruntime", "append8")
			staticCall.emit(funcall.args)
		case 24:
			if slice.getGtype().elementType.getKind() == G_INTERFACE && valueToAppend.getGtype().getKind() != G_INTERFACE {
				eConvertion := &ExprConversionToInterface{
					tok:  valueToAppend.token(),
					expr: valueToAppend,
				}
				funcall.args[1] = eConvertion
			}
			staticCall.symbol = getFuncSymbol("iruntime", "append24")
			staticCall.emit(funcall.args)
		default:
			TBI(slice.token(), "")
		}
	case builtinMakeSlice:
		assert(len(funcall.args) == 3, funcall.token(), "append() should take 3 argments")
		var staticCall *IrStaticCall = &IrStaticCall{
			callee: decl,
		}
		staticCall.symbol = getFuncSymbol("iruntime", "makeSlice")
		staticCall.emit(funcall.args)
	case builtinDumpSlice:
		arg := funcall.args[0]

		emit("lea .%s, %%rax", builtinStringKey2)
		emit("PUSH_8")

		arg.emit()
		emit("PUSH_SLICE")

		numRegs := 4
		for i := numRegs - 1; i >= 0; i-- {
			emit("POP_TO_ARG_%d", i)
		}

		emit("FUNCALL %s", "printf")
		emitNewline()
	case builtinDumpInterface:
		arg := funcall.args[0]

		emit("lea .%s, %%rax", builtinStringKey1)
		emit("PUSH_8")

		arg.emit()
		emit("PUSH_INTERFACE")

		numRegs := 4
		for i := numRegs - 1; i >= 0; i-- {
			emit("POP_TO_ARG_%d", i)
		}

		emit("FUNCALL %s", "printf")
		emitNewline()
	case builtinAssertInterface:
		emit("# builtinAssertInterface")
		labelEnd := makeLabel()
		arg := funcall.args[0]
		arg.emit() // rax=ptr, rbx=receverTypeId, rcx=dynamicTypeId

		// (ptr != nil && rcx == nil) => Error

		emit("CMP_NE_ZERO")
		emit("TEST_IT")
		emit("je %s", labelEnd)

		emit("mov %%rcx, %%rax")

		emit("CMP_EQ_ZERO")
		emit("TEST_IT")
		emit("je %s", labelEnd)

		slabel := makeLabel()
		emit(".data 0")
		emitWithoutIndent("%s:", slabel)
		emit(".string \"%s\"", "assertInterface failed")
		emit(".text")
		emit("lea %s, %%rax", slabel)
		emit("PUSH_8")
		emit("POP_TO_ARG_0")
		emit("FUNCALL %s", ".panic")

		emitWithoutIndent("%s:", labelEnd)
		emitNewline()

	case builtinAsComment:
		arg := funcall.args[0]
		if stringLiteral, ok := arg.(*ExprStringLiteral); ok {
			emitWithoutIndent("# %s", stringLiteral.val)
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
	// https://sourceware.org/binutils/docs-2.30/as/Symbol-Intro.html#Symbol-Intro
	// A symbol is one or more characters chosen from the set of all letters (both upper and lower case), digits and the three characters ‘_.$’.
	symbol       string
	callee       *DeclFunc
	isMethodCall bool
}

func bool2string(bol bool) string {
	if bol {
		return "true"
	} else {
		return "false"
	}
}

func emitMakeSliceFunc() {
	// makeSlice
	emitWithoutIndent("%s:", "iruntime.makeSlice")
	emit("FUNC_PROLOGUE")
	emitNewline()

	emit("PUSH_ARG_2") // -8
	emit("PUSH_ARG_1") // -16
	emit("PUSH_ARG_0") // -24

	emit("mov -16(%%rbp), %%rax # newcap")
	emit("mov -8(%%rbp), %%rcx # unit")
	emit("imul %%rcx, %%rax")
	emit("ADD_NUMBER 16 # pure buffer")

	emit("PUSH_8")
	emit("POP_TO_ARG_0")
	emit("FUNCALL iruntime.malloc")

	emit("mov -24(%%rbp), %%rbx # newlen")
	emit("mov -16(%%rbp), %%rcx # newcap")

	emit("LEAVE_AND_RET")
	emitNewline()
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

