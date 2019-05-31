// Code generator
// Convention:
//  We SHOULD use the word "emit" for the meaning of "output assembly code",
//  NOT for "load something to %rax".
//  Such usage would make much confusion.

package main

import (
	"fmt"
	"os"
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

var retRegi [14]string = [14]string{
	"rax", "rbx", "rcx", "rdx", "rdi", "rsi", "r8", "r9", "r10", "r11", "r12", "r13", "r14", "r15",
}

var RegsForArguments [12]string = [12]string{"rdi", "rsi", "rdx", "rcx", "r8", "r9", "r10", "r11", "r12", "r13", "r14", "r15"}

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

func (f *DeclFunc) emitPrologue() {
	emitWithoutIndent("%s:", f.getSymbol())
	emit("FUNC_PROLOGUE")

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
		emit("# set params")
	}

	var regIndex int
	for _, param := range params {
		switch param.getGtype().getKind() {
		case G_SLICE, G_INTERFACE, G_MAP:
			offset -= IntSize * 3
			param.offset = offset
			emit("PUSH_ARG_%d # third", regIndex+2)
			emit("PUSH_ARG_%d # second", regIndex+1)
			emit("PUSH_ARG_%d # fist \"%s\" %s", regIndex, param.varname, param.getGtype().String())
			regIndex += sliceWidth
		default:
			offset -= IntSize
			param.offset = offset
			emit("PUSH_ARG_%d # param \"%s\" %s", regIndex, param.varname, param.getGtype().String())
			regIndex += 1
		}
	}

	if len(f.localvars) > 0 {
		emit("# Allocating stack for localvars len=%d", len(f.localvars))
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
		emit("# offset %d variable \"%s\" %s", lvar.offset, lvar.varname, lvar.gtype.String())
	}

	if localarea != 0 {
		emit("sub $%d, %%rsp # total stack size", -localarea)
	}

	emitNewline()
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

func (ast *ExprNumberLiteral) emit() {
	emit("LOAD_NUMBER %d", ast.val)
}

func (ast *ExprStringLiteral) emit() {
	emit("LOAD_STRING_LITERAL .%s", ast.slabel)
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
		if field.kind == G_ARRAY {
			variable.emitAddress(field.offset)
		} else {
			if variable.isGlobal {
				emit("LOAD_8_FROM_GLOBAL %s, %d+%d", variable.varname, field.offset,offset)
			} else {
				emit("LOAD_8_FROM_LOCAL %d+%d+%d", variable.offset, field.offset, offset)
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
	emit("ADD_NUMBER %d", field.offset)
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

func (a *ExprStructField) emit() {
	emit("# LOAD ExprStructField")
	a.calcOffset()
	switch a.strct.getGtype().kind {
	case G_POINTER: // pointer to struct
		strcttype := a.strct.getGtype().origType.relation.gtype
		field := strcttype.getField(a.fieldname)
		a.strct.emit()
		emit("ADD_NUMBER %d", field.offset)
		switch field.getKind() {
		case G_SLICE, G_INTERFACE, G_MAP:
			emit("LOAD_24_BY_DEREF")
		default:
			emit("LOAD_8_BY_DEREF")
		}

	case G_NAMED: // struct
		strcttype := a.strct.getGtype().relation.gtype
		assert(strcttype.size > 0, a.token(), "struct size should be > 0")
		field := strcttype.getField(a.fieldname)
		loadStructField(a.strct, field, 0)
	default:
		errorft(a.token(), "internal error: bad gtype %s", a.strct.getGtype().String())
	}
}

func (ast *ExprVariable) emit() {
	emit("# load variable \"%s\" %s", ast.varname, ast.getGtype().String())
	if ast.isGlobal {
		switch ast.gtype.getKind() {
		case G_INTERFACE:
			emit("LOAD_INTERFACE_FROM_GLOBAL %s", ast.varname)
		case G_SLICE:
			emit("LOAD_SLICE_FROM_GLOBAL %s", ast.varname)
		case G_MAP:
			emit("LOAD_MAP_FROM_GLOBAL %s", ast.varname)
		case G_ARRAY:
			ast.emitAddress(0)
		default:
			if ast.getGtype().getSize() == 1 {
				emit("LOAD_1_FROM_GLOBAL_CAST %s", ast.varname)
			} else {
				emit("LOAD_8_FROM_GLOBAL %s", ast.varname)
			}
		}
	} else {
		if ast.offset == 0 {
			errorft(ast.token(), "offset should not be zero for localvar %s", ast.varname)
		}
		switch ast.gtype.getKind() {
		case G_INTERFACE:
			emit("LOAD_INTERFACE_FROM_LOCAL %d", ast.offset)
		case G_SLICE:
			emit("LOAD_SLICE_FROM_LOCAL %d", ast.offset)
		case G_MAP:
			emit("LOAD_MAP_FROM_LOCAL %d", ast.offset)
		case G_ARRAY:
			ast.emitAddress(0)
		default:
			if ast.getGtype().getSize() == 1 {
				emit("LOAD_1_FROM_LOCAL_CAST %d", ast.offset)
			} else {
				emit("LOAD_8_FROM_LOCAL %d", ast.offset)
			}
		}
	}
}

func (variable *ExprVariable) emitAddress(offset int) {
	if variable.isGlobal {
		emit("LOAD_GLOBAL_ADDR %s, %d", variable.varname, offset)
	} else {
		if variable.offset == 0 {
			errorft(variable.token(), "offset should not be zero for localvar %s", variable.varname)
		}
		emit("LOAD_LOCAL_ADDR %d+%d", variable.offset, offset)
	}
}

func (rel *Relation) emit() {
	if rel.expr == nil {
		errorft(rel.token(), "rel.expr is nil: %s", rel.name)
	}
	rel.expr.emit()
}

func (ast *ExprConstVariable) emit() {
	emit("# *ExprConstVariable.emit() name=%s iotaindex=%d", ast.name, ast.iotaIndex)
	assert(ast.iotaIndex < 10000, ast.token(), "iotaindex is too large")
	assert(ast.val != nil, ast.token(), "const.val for should not be nil:"+string(ast.name))
	rel, ok := ast.val.(*Relation)
	if ok {
		emit("# rel=%s", rel.name)
		cnst, ok := rel.expr.(*ExprConstVariable)
		if ok && cnst == eIota {
			emit("# const is iota")
			// replace the iota expr by a index number
			val := &ExprNumberLiteral{
				val: ast.iotaIndex,
			}
			val.emit()
		} else {
			emit("# Not iota")
			ast.val.emit()
		}
	} else {
		emit("# const is not iota")
		ast.val.emit()
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
	emit("PUSH_PRIMITIVE # left") // left
	binop.right.emit()
	if binop.right.getGtype().getKind() == G_BYTE {
		emit_intcast(binop.right.getGtype())
	}
	emit("PUSH_PRIMITIVE # right") // right
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

// e.g. *x = 1, or *x++
func (uop *ExprUop) emitSave() {
	emit("# *ExprUop.emitSave()")
	assert(uop.op == "*", uop.tok, "uop op should be *")
	emit("PUSH_PRIMITIVE")
	uop.operand.emit()
	emit("PUSH_PRIMITIVE")
	emit("STORE_INDIRECT_FROM_STACK")
}

// e.g. x = 1
func (rel *Relation) emitSave() {
	if rel.expr == nil {
		errorft(rel.token(), "left.rel.expr is nil")
	}
	variable := rel.expr.(*ExprVariable)
	variable.emitOffsetSave(variable.getGtype().getSize(), 0, false)
}

func (variable *ExprVariable) emitOffsetSave(size int, offset int, forceIndirection bool) {
	emit("# ExprVariable.emitOffsetSave(size %d, offset %d)", size, offset)
	assert(0 <= size && size <= 8, variable.token(), fmt.Sprintf("invalid size %d", size))
	if variable.getGtype().kind == G_POINTER && (offset > 0 || forceIndirection) {
		assert(variable.getGtype().kind == G_POINTER, variable.token(), "")
		emit("PUSH_PRIMITIVE")
		variable.emit()
		emit("ADD_NUMBER %d", offset)
		emit("PUSH_PRIMITIVE")

		emit("STORE_INDIRECT_FROM_STACK")
		emit("#")
		return
	}
	if variable.isGlobal {
		emit("STORE_%d_TO_GLOBAL %s %d", size, variable.varname, offset)
	} else {
		emitStoreItToLocal(size, variable.offset+offset, "")
	}
}

func (variable *ExprVariable) emitOffsetLoad(size int, offset int) {
	assert(0 <= size && size <= 8, variable.token(), "invalid size")
	if variable.isGlobal {
		emit("LOAD_%d_FROM_GLOBAL %s %d", size, variable.varname, offset)
	} else {
		emit("LOAD_%d_FROM_LOCAL %d+%d", size,  variable.offset, offset)
	}
}

func (ast *ExprUop) emit() {
	emit("# emitting ExprUop")
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
			ivv := e.invisiblevar
			assignToStruct(ivv, e)

			emitCallMalloc(e.getGtype().getSize()) // => rax
			emit("PUSH_PRIMITIVE")                     // to:ptr addr
			e.invisiblevar.emitAddress(0)
			emit("PUSH_PRIMITIVE") // from:address of invisible var
			emitCopyStructFromStack(e.getGtype())
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
		emit("LOAD_8_BY_DEREF")
	} else if ast.op == "!" {
		ast.operand.emit()
		emit("CMP_EQ_ZERO")
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
	emit("mov $0, %%rax")
	emit("je %s", labelElse)
	emitEmptyString()
	emit("jmp %s", labelEnd)
	emit("%s:", labelElse)
	binop.left.emit()
	emit("%s:", labelEnd)
	emit("PUSH_PRIMITIVE")

	binop.right.emit()
	emit("PUSH_PRIMITIVE")
	emitStringsEqualFromStack(equal)
}

func emitConvertNilToEmptyString(regi string) {
	emit("# emitConvertNilToEmptyString")
	emit("mov %s, %%rax", regi)
	emit("PUSH_PRIMITIVE")
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
	emit("pop %%rax")
	emit("pop %%rcx")
	emitStringsEqual(equal, "%rax", "%rcx")
}

func emitStringsEqual(equal bool, leftReg string, rightReg string) {
	emit("push %s", rightReg) // stash

	emitConvertNilToEmptyString(leftReg)
	emit("mov %s, %%rsi", leftReg)

	emit("pop %%rax # right string")
	emitConvertNilToEmptyString("%rax")

	emit("PUSH_PRIMITIVE")
	emit("POP_TO_ARG_0")
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
	emit("push %%rax # left string")

	emit("PUSH_PRIMITIVE")
	emit("POP_TO_ARG_0")
	emit("mov $0, %%rax")
	emit("call strlen # get left len")

	emit("push %%rax # left len")
	right.emit()
	emit("push %%rax # right string")
	emit("PUSH_PRIMITIVE")
	emit("POP_TO_ARG_0")
	emit("mov $0, %%rax")
	emit("call strlen # get right len")
	emit("push %%rax # right len")

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
	emit("PUSH_PRIMITIVE")
	emit("POP_TO_ARG_0")
	emit("mov $0, %%rax")
	emit("call iruntime.malloc")

	emit("PUSH_PRIMITIVE")
	emit("POP_TO_ARG_0")
	emit("POP_TO_ARG_1")
	emit("mov $0, %%rax")
	emit("call strcat")

	emit("PUSH_PRIMITIVE")
	emit("POP_TO_ARG_0")
	emit("POP_TO_ARG_1")
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
		emit("TEST_IT")
		emit("mov $0, %%rax")
		emit("je %s", labelEnd)
		ast.right.emit()
		emit("TEST_IT")
		emit("mov $0, %%rax")
		emit("je %s", labelEnd)
		emit("mov $1, %%rax")
		emit("%s:", labelEnd)
		return
	case "||":
		labelEnd := makeLabel()
		ast.left.emit()
		emit("TEST_IT")
		emit("mov $1, %%rax")
		emit("jne %s", labelEnd)
		ast.right.emit()
		emit("TEST_IT")
		emit("mov $1, %%rax")
		emit("jne %s", labelEnd)
		emit("mov $0, %%rax")
		emit("%s:", labelEnd)
		return
	}
	ast.left.emit()
	emit("PUSH_PRIMITIVE")
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
	emit("# StmtAssignment")
	// the right hand operand is a single multi-valued expression
	// such as a function call, a channel or map operation, or a type assertion.
	// The number of operands on the left hand side must match the number of values.
	isOnetoOneAssignment := (len(ast.rights) > 1)
	if isOnetoOneAssignment {
		emit("# multi(%d) = multi(%d)", len(ast.lefts), len(ast.rights))
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
			case gtype.getKind() == G_ARRAY:
				assignToArray(left, right)
			case gtype.getKind() == G_SLICE:
				assignToSlice(left, right)
			case gtype.getKind() == G_STRUCT:
				assignToStruct(left, right)
			case gtype.getKind() == G_INTERFACE:
				assignToInterface(left, right)
			default:
				// suppose primitive
				emitAssignPrimitive(left, right)
			}
		}
		return
	} else {
		numLeft := len(ast.lefts)
		emit("# multi(%d) = expr", numLeft)
		// a,b,c = expr
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
			if indexExpr.collection.getGtype().getKind() == G_MAP {
				// map get
				emit("# v, ok = map[k]")
				leftsMayBeTwo = true
			}
			numRight++
		default:
			numRight++
		}

		if leftsMayBeTwo {
			if numLeft > 2 {
				errorft(ast.token(), "number of exprs does not match. numLeft=%d", numLeft)
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
					if left.getGtype().kind == G_SLICE {
						// @TODO: Does this work ?
						emitSave24(left, 0)
					} else if left.getGtype().getKind() == G_INTERFACE {
						// @TODO: Does this work ?
						emitSave24(left, 0)
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
		//emit("# Assign %T %s = %T %s", left, gtype.String(), right, right.getGtype())
		switch {
		case gtype == nil:
			// suppose left is "_"
			right.emit()
		case gtype.getKind() == G_ARRAY:
			assignToArray(left, right)
		case gtype.getKind() == G_SLICE:
			assignToSlice(left, right)
		case gtype.getKind() == G_STRUCT:
			assignToStruct(left, right)
		case gtype.getKind() == G_INTERFACE:
			assignToInterface(left, right)
		case gtype.getKind() == G_MAP:
			assignToMap(left, right)
		default:
			// suppose primitive
			emitAssignPrimitive(left, right)
		}
		if leftsMayBeTwo && len(ast.lefts) == 2 {
			okVariable := ast.lefts[1]
			okRegister := mapOkRegister(right.getGtype().is24Width())
			emit("mov %%%s, %%rax # emit okValue", okRegister)
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
		emit("# %s %s = ", left.(*Relation).name, left.getGtype().String())
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

// save data from stack
func (e *ExprIndex) emitSave24() {
	// load head address of the array
	// load index
	// multi index * size
	// calc address = head address + offset
	// copy value to the address

	collectionType := e.collection.getGtype()
	switch {
	case collectionType.getKind() == G_ARRAY, collectionType.getKind() == G_SLICE, collectionType.getKind() == G_STRING:
		e.collection.emit() // head address
	case collectionType.getKind() == G_MAP:
		e.emitMapSet(true)
		return
	default:
		TBI(e.token(), "unable to handle %s", collectionType)
	}
	emit("PUSH_PRIMITIVE # head address of collection")
	e.index.emit()
	emit("PUSH_PRIMITIVE # index")
	var elmType *Gtype
	if collectionType.isString() {
		elmType = gByte
	} else {
		elmType = collectionType.elementType
	}
	size := elmType.getSize()
	assert(size > 0, nil, "size > 0")
	emit("LOAD_NUMBER %d # elementSize", size)
	emit("PUSH_PRIMITIVE")
	emit("IMUL_FROM_STACK # index * elementSize")
	emit("PUSH_PRIMITIVE # index * elementSize")
	emit("SUM_FROM_STACK # (index * size) + address")
	emit("PUSH_PRIMITIVE")
	emit("STORE_24_INDIRECT_FROM_STACK")
}

func (e *ExprIndex) emitSave() {
	emit("PUSH_PRIMITIVE") // push RHS value

	// load head address of the array
	// load index
	// multi index * size
	// calc address = head address + offset
	// copy value to the address

	collectionType := e.collection.getGtype()
	switch {
	case collectionType.getKind() == G_ARRAY, collectionType.getKind() == G_SLICE, collectionType.getKind() == G_STRING:
		e.collection.emit() // head address
	case collectionType.getKind() == G_MAP:
		e.emitMapSet(false)
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
	if e.strct.getGtype().kind == G_POINTER {
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
	assert(vr.gtype.kind == G_NAMED, e.tok, "expect G_NAMED, but got "+vr.gtype.String())
	field := vr.gtype.relation.gtype.getField(e.fieldname)
	vr.emitOffsetLoad(size, field.offset+offset)
}

// rax: address
// rbx: len
// rcx: cap
func (e *ExprSliceLiteral) emit() {
	emit("# (*ExprSliceLiteral).emit()")
	length := len(e.values)
	//debugf("slice literal %s: underlyingarray size = %d (should be %d)", e.getGtype(), e.gtype.getSize(),  e.gtype.elementType.getSize() * length)
	emitCallMalloc(e.gtype.getSize() * length)
	emit("push %%rax # ptr")
	for i, value := range e.values {
		if e.gtype.elementType.getKind() == G_INTERFACE && value.getGtype().getKind() != G_INTERFACE {
			emitConversionToInterface(value)
		} else {
			value.emit()
		}

		switch e.gtype.elementType.getKind() {
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
	emit("TEST_IT")
	if stmt.els != nil {
		labelElse := makeLabel()
		labelEndif := makeLabel()
		emit("je %s  # jump if 0", labelElse)
		emit("# then block")
		stmt.then.emit()
		emit("jmp %s # jump to endif", labelEndif)
		emit("# else block")
		emit("%s:", labelElse)
		stmt.els.emit()
		emit("# endif")
		emit("%s:", labelEndif)
	} else {
		// no else block
		labelEndif := makeLabel()
		emit("je %s  # jump if 0", labelEndif)
		emit("# then block")
		stmt.then.emit()
		emit("# endif")
		emit("%s:", labelEndif)
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
		emit("PUSH_PRIMITIVE")
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
				if gtype.isNil() {
					emit("mov $0, %%rax")
				} else {
					typeLabel := groot.getTypeLabel(gtype)
					emit("lea .%s(%%rip), %%rax # type: %s", typeLabel, gtype.String())
				}

				emit("pop %%rcx # the subject type")
				emit("push %%rcx # the subject value")
				emitStringsEqual(true, "%rax", "%rcx")
				emit("TEST_IT")
				emit("jne %s # jump if matches", myCaseLabel)
			}
		} else if stmt.cond == nil {
			for _, e := range caseClause.exprs {
				e.emit()
				emit("TEST_IT")
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
				emit("TEST_IT")
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

func (f *StmtFor) emitRangeForList() {
	emitNewline()
	emit("# for range %s", f.rng.rangeexpr.getGtype().String())
	assertNotNil(f.rng.indexvar != nil, f.rng.tok)
	assert(f.rng.rangeexpr.getGtype().kind == G_ARRAY || f.rng.rangeexpr.getGtype().kind == G_SLICE, f.rng.tok, "rangeexpr should be G_ARRAY or G_SLICE, but got "+f.rng.rangeexpr.getGtype().String())

	labelBegin := makeLabel()
	f.labelEndBlock = makeLabel()
	f.labelEndLoop = makeLabel()

	// i = 0
	emit("# init index")
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
	initstmt.emit()

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
	emit("TEST_IT")
	emit("je %s  # if false, go to loop end", f.labelEndLoop)

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

	f.block.emit()
	emit("%s: # end block", f.labelEndBlock)

	// break if i == len(list) - 1
	condition2 := &ExprBinop{
		op:   "==",
		left: f.rng.indexvar, // i
		// @TODO2
		// The range expression x is evaluated once before beginning the loop
		right: &ExprBinop{
			op: "-",
			left: &ExprLen{
				arg: f.rng.rangeexpr, // len(expr)
			},
			right: &ExprNumberLiteral{
				val: 1,
			},
		},
	}
	condition2.emit()
	emit("TEST_IT")
	emit("jne %s  # if this iteration is final, go to loop end", f.labelEndLoop)

	// i++
	indexIncr := &StmtInc{
		operand: f.rng.indexvar,
	}
	indexIncr.emit()

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
		emit("TEST_IT")
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
		if f.rng.rangeexpr.getGtype().getKind() == G_MAP {
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
		if rettype.getKind() == G_INTERFACE && expr.getGtype().getKind() != G_INTERFACE {
			if expr.getGtype() == nil {
				emit("LOAD_EMPTY_INTERFACE")
			} else {
				emitConversionToInterface(expr)
			}
		} else {
			expr.emit()
			if expr.getGtype() == nil && stmt.rettypes[0].kind == G_SLICE {
				emit("LOAD_EMPTY_SLICE")
			}
		}
		stmt.emitDeferAndReturn()
		return
	}
	for i, rettype := range stmt.rettypes {
		expr := stmt.exprs[i]
		expr.emit()
		//		rettype := stmt.rettypes[i]
		if expr.getGtype() == nil && rettype.kind == G_SLICE {
			emit("LOAD_EMPTY_SLICE")
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

func emitStoreItToLocal(size int, loff int, comment string) {
	emit("STORE_%d_TO_LOCAL %d", size, loff)
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

// expect rhs address is in the stack top, lhs is in the second top
func emitCopyStructFromStack(gtype *Gtype) {
	//assert(left.getGtype().getSize() == right.getGtype().getSize(), left.token(),"size does not match")
	emit("pop %%rax") // to
	emit("pop %%rbx") // from
	emit("push %%r15")
	emit("push %%r11")
	emit("mov %%rax, %%r15") // from
	emit("mov %%rbx, %%rax") // to

	emitCopyStructInt(gtype)
}

func emitCopyStructInt(gtype *Gtype) {
	var i int
	for ; i < gtype.getSize(); i += 8 {
		emit("movq %d(%%r15), %%r11", i)
		emit("movq %%r11, %d(%%rax)", i)
	}
	for ; i < gtype.getSize(); i += 4 {
		emit("movl %d(%%r15), %%r11", i)
		emit("movl %%r11, %d(%%rax)", i)
	}
	for ; i < gtype.getSize(); i++ {
		emit("movb %d(%%r15), %%r11", i)
		emit("movb %%r11, %d(%%rax)", i)
	}

	emit("pop %%r11")
	emit("pop %%r15")
}


func assignToStruct(lhs Expr, rhs Expr) {
	emit("# assignToStruct start")

	if rel, ok := lhs.(*Relation); ok {
		lhs = rel.expr
	}
	assert(rhs == nil || (rhs.getGtype().kind == G_NAMED && rhs.getGtype().relation.gtype.kind == G_STRUCT),
		lhs.token(), "rhs should be struct type")
	// initializes with zero values
	emit("# initialize struct with zero values: start")
	for _, fieldtype := range lhs.getGtype().relation.gtype.fields {
		switch {
		case fieldtype.kind == G_ARRAY:
			arrayType := fieldtype
			elementType := arrayType.elementType
			elmSize := arrayType.elementType.getSize()
			switch {
			case elementType.kind == G_NAMED && elementType.relation.gtype.kind == G_STRUCT:
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

		case fieldtype.kind == G_SLICE:
			emit("LOAD_EMPTY_SLICE")
			emit("PUSH_SLICE")
			emitSave24(lhs, fieldtype.offset)
		case fieldtype.kind == G_MAP:
			emit("LOAD_EMPTY_MAP")
			emit("PUSH_MAP")
			emitSave24(lhs, fieldtype.offset)
		case fieldtype.kind == G_NAMED && fieldtype.relation.gtype.kind == G_STRUCT:
			left := &ExprStructField{
				strct:     lhs,
				fieldname: fieldtype.fieldname,
			}
			assignToStruct(left, nil)
		case fieldtype.getKind() == G_INTERFACE:
			emit("LOAD_EMPTY_INTERFACE")
			emit("PUSH_INTERFACE")
			emitSave24(lhs, fieldtype.offset)
		default:
			emit("mov $0, %%rax")
			regSize := fieldtype.getSize()
			assert(0 < regSize && regSize <= 8, lhs.token(), fieldtype.String())
			emitOffsetSave(lhs, regSize, fieldtype.offset)
		}
	}
	emit("# initialize struct with zero values: end")

	if rhs == nil {
		return
	}
	variable := lhs

	strcttyp := rhs.getGtype().Underlying()

	switch rhs.(type) {
	case *Relation:
		emitAddress(lhs)
		emit("PUSH_PRIMITIVE")
		emitAddress(rhs)
		emit("PUSH_PRIMITIVE")
		emitCopyStructFromStack(lhs.getGtype())
	case *ExprUop:
		re := rhs.(*ExprUop)
		if re.op == "*" {
			// copy struct
			emitAddress(lhs)
			emit("PUSH_PRIMITIVE")
			re.operand.emit()
			emit("PUSH_PRIMITIVE")
			emitCopyStructFromStack(lhs.getGtype())
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
			case fieldtype.kind == G_ARRAY:
				initvalues, ok := field.value.(*ExprArrayLiteral)
				assert(ok, nil, "ok")
				arrayType := strcttyp.getField(field.key)
				elementType := arrayType.elementType
				elmSize := elementType.getSize()
				switch {
				case elementType.kind == G_NAMED && elementType.relation.gtype.kind == G_STRUCT:
					left := &ExprStructField{
						strct:     lhs,
						fieldname: fieldtype.fieldname,
					}
					assignToArray(left, field.value)
				default:
					for i, val := range initvalues.values {
						val.emit()
						emitOffsetSave(variable, elmSize, arrayType.offset+i*elmSize)
					}
				}
			case fieldtype.kind == G_SLICE:
				left := &ExprStructField{
					tok:       variable.token(),
					strct:     lhs,
					fieldname: field.key,
				}
				assignToSlice(left, field.value)
			case fieldtype.getKind() == G_MAP:
				left := &ExprStructField{
					tok:       variable.token(),
					strct:     lhs,
					fieldname: field.key,
				}
				assignToMap(left, field.value)
			case fieldtype.getKind() == G_INTERFACE:
				left := &ExprStructField{
					tok:       lhs.token(),
					strct:     lhs,
					fieldname: field.key,
				}
				assignToInterface(left, field.value)
			case fieldtype.kind == G_NAMED && fieldtype.relation.gtype.kind == G_STRUCT:
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

	emit("# assignToStruct end")
}

const sliceOffsetForLen = 8

func emitOffsetSave(lhs Expr, size int, offset int) {
	switch lhs.(type) {
	case *Relation:
		rel := lhs.(*Relation)
		emitOffsetSave(rel.expr, size, offset)
	case *ExprVariable:
		variable := lhs.(*ExprVariable)
		variable.emitOffsetSave(size, offset, false)
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
	emit("# emitOffsetLoad(offset %d)", offset)
	switch lhs.(type) {
	case *Relation:
		rel := lhs.(*Relation)
		emitOffsetLoad(rel.expr, size, offset)
	case *ExprVariable:
		variable := lhs.(*ExprVariable)
		variable.emitOffsetLoad(size, offset)
	case *ExprStructField:
		structfield := lhs.(*ExprStructField)
		structfield.calcOffset()
		fieldType := structfield.getGtype()
		if structfield.strct.getGtype().kind == G_POINTER {
			structfield.strct.emit() // emit address of the struct
			emit("# offset %d + %d = %d", fieldType.offset, offset, fieldType.offset+offset)
			emit("ADD_NUMBER %d+%d", fieldType.offset,offset)
			//reg := getReg(size)
			emit("LOAD_8_BY_DEREF")
		} else {
			emitOffsetLoad(structfield.strct, size, fieldType.offset+offset)
		}
	case *ExprIndex:
		//  e.g. arrayLiteral.values[i].getGtype().getKind()
		indexExpr := lhs.(*ExprIndex)
		loadCollectIndex(indexExpr.collection, indexExpr.index, offset)
	case *ExprMethodcall:
		// @TODO this logic is temporarly. Need to be verified.
		mcall := lhs.(*ExprMethodcall)
		rettypes := mcall.getRettypes()
		assert(len(rettypes) == 1, lhs.token(), "rettype should be single")
		rettype := rettypes[0]
		assert(rettype.getKind() == G_POINTER, lhs.token(), "only pointer is supported")
		mcall.emit()
		emit("ADD_NUMBER %d", offset)
		emit("LOAD_8_BY_DEREF")
	default:
		errorft(lhs.token(), "unkonwn type %T", lhs)
	}
}

// take slice values from stack
func emitSave24(lhs Expr, offset int) {
	assertInterface(lhs)
	//emit("# emitSave24(%T, offset %d)", lhs, offset)
	emit("# emitSave24(?, offset %d)", offset)
	switch lhs.(type) {
	case *Relation:
		rel := lhs.(*Relation)
		emitSave24(rel.expr, offset)
	case *ExprVariable:
		variable := lhs.(*ExprVariable)
		variable.emitSave24(offset)
	case *ExprStructField:
		structfield := lhs.(*ExprStructField)
		fieldType := structfield.getGtype()
		fieldOffset := fieldType.offset
		emit("# fieldOffset=%d (%s)", fieldOffset, fieldType.fieldname)
		emitSave24(structfield.strct, fieldOffset+offset)
	case *ExprIndex:
		indexExpr := lhs.(*ExprIndex)
		indexExpr.emitSave24()
	default:
		errorft(lhs.token(), "unkonwn type %T", lhs)
	}
}

func emitCallMallocDinamicSize(eSize Expr) {
	eSize.emit()
	emit("PUSH_PRIMITIVE")
	emit("POP_TO_ARG_0")
	emit("mov $0, %%rax")
	emit("call iruntime.malloc")
}

func emitCallMalloc(size int) {
	eNumber := &ExprNumberLiteral{
		val: size,
	}
	emitCallMallocDinamicSize(eNumber)
}

func assignToMap(lhs Expr, rhs Expr) {
	emit("# assignToMap")
	if rhs == nil {
		emit("# initialize map with a zero value")
		emit("LOAD_EMPTY_MAP")
		emit("PUSH_MAP")
		emitSave24(lhs, 0)
		return
	}
	switch rhs.(type) {
	case *ExprMapLiteral:
		emit("# map literal")

		lit := rhs.(*ExprMapLiteral)
		lit.emit()
		emit("PUSH_MAP")
	case *Relation, *ExprVariable, *ExprIndex, *ExprStructField, *ExprFuncallOrConversion, *ExprMethodcall:
		rhs.emit()
		emit("PUSH_MAP")
	default:
		TBI(rhs.token(), "unable to handle %T", rhs)
	}
	emitSave24(lhs, 0)
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
	emit("PUSH_PRIMITIVE")
	emitCallMalloc(8)
	emit("PUSH_PRIMITIVE")
	emit("STORE_INDIRECT_FROM_STACK")
	emit("PUSH_PRIMITIVE # addr of dynamicValue") // address

	if receiverType.kind == G_POINTER {
		receiverType = receiverType.origType.relation.gtype
	}
	//assert(receiverType.receiverTypeId > 0,  dynamicValue.token(), "no receiverTypeId")
	emit("LOAD_NUMBER %d # receiverTypeId", receiverType.receiverTypeId)
	emit("PUSH_PRIMITIVE # receiverTypeId")

	gtype := dynamicValue.getGtype()
	label := groot.getTypeLabel(gtype)
	emit("lea .%s, %%rax# dynamicType %s", label, gtype.String())
	emit("PUSH_PRIMITIVE # dynamicType")

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


func assignToInterface(lhs Expr, rhs Expr) {
	emit("# assignToInterface")
	if rhs == nil || isNil(rhs) {
		emit("LOAD_EMPTY_INTERFACE")
		emit("PUSH_INTERFACE")
		emitSave24(lhs, 0)
		return
	}

	assert(rhs.getGtype() != nil, rhs.token(), fmt.Sprintf("rhs gtype is nil:%T", rhs))
	if rhs.getGtype().getKind() == G_INTERFACE {
		rhs.emit()
		emit("PUSH_INTERFACE")
		emitSave24(lhs, 0)
		return
	}

	emitConversionToInterface(rhs)
	emit("PUSH_INTERFACE")
	emitSave24(lhs, 0)
}

func assignToSlice(lhs Expr, rhs Expr) {
	emit("# assignToSlice")
	assertInterface(lhs)
	//assert(rhs == nil || rhs.getGtype().kind == G_SLICE, nil, "should be a slice literal or nil")
	if rhs == nil {
		emit("LOAD_EMPTY_SLICE")
		emit("PUSH_SLICE")
		emitSave24(lhs, 0)
		return
	}

	//	assert(rhs.getGtype().getKind() == G_SLICE, rhs.token(), "rsh should be slice type")

	switch rhs.(type) {
	case *Relation:
		rel := rhs.(*Relation)
		if _, ok := rel.expr.(*ExprNilLiteral); ok {
			emit("LOAD_EMPTY_SLICE")
			emit("PUSH_SLICE")
			emitSave24(lhs, 0)
			return
		}
		rvariable, ok := rel.expr.(*ExprVariable)
		assert(ok, nil, "ok")
		rvariable.emit()
		emit("PUSH_SLICE")
	case *ExprSliceLiteral:
		lit := rhs.(*ExprSliceLiteral)
		lit.emit()
		emit("PUSH_SLICE")
	case *ExprSlice:
		e := rhs.(*ExprSlice)
		e.emit()
		emit("PUSH_SLICE")
	case *ExprConversion:
		// https://golang.org/ref/spec#Conversions
		// Converting a value of a string type to a slice of bytes type
		// yields a slice whose successive elements are the bytes of the string.
		//
		// see also https://blog.golang.org/strings
		conversion := rhs.(*ExprConversion)
		assert(conversion.gtype.kind == G_SLICE, rhs.token(), "must be a slice of bytes")
		assert(conversion.expr.getGtype().kind == G_STRING || conversion.expr.getGtype().relation.gtype.kind == G_STRING, rhs.token(), "must be a string type, but got "+conversion.expr.getGtype().String())
		stringVarname, ok := conversion.expr.(*Relation)
		assert(ok, rhs.token(), "ok")
		stringVariable := stringVarname.expr.(*ExprVariable)
		stringVariable.emit()
		emit("PUSH_PRIMITIVE # ptr")
		strlen := &ExprLen{
			arg: stringVariable,
		}
		strlen.emit()
		emit("PUSH_PRIMITIVE # len")
		emit("PUSH_PRIMITIVE # cap")

	default:
		//emit("# emit rhs of type %T %s", rhs, rhs.getGtype().String())
		rhs.emit() // it should put values to rax,rbx,rcx
		emit("PUSH_SLICE")
	}

	emitSave24(lhs, 0)
}

func (variable *ExprVariable) emitSave24(offset int) {
	emit("# *ExprVariable.emitSave24()")
	emit("pop %%rax # 3rd")
	variable.emitOffsetSave(8, offset+ptrSize+sliceOffsetForLen, false)
	emit("pop %%rax # 2nd")
	variable.emitOffsetSave(8, offset+ptrSize, false)
	emit("pop %%rax # 1st")
	variable.emitOffsetSave(8, offset, true)
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
	assert(rhs == nil || rhs.getGtype().kind == G_ARRAY, nil, "rhs should be array")
	switch {
	case elementType.kind == G_NAMED && elementType.relation.gtype.kind == G_STRUCT:
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
				if elementType.getKind() == G_INTERFACE {
					emit("LOAD_EMPTY_INTERFACE")
					emit("PUSH_INTERFACE")
					emitSave24(lhs, offsetByIndex)
					continue
				} else {
					emit("mov $0, %%rax")
				}
			case *ExprArrayLiteral:
				arrayLiteral := rhs.(*ExprArrayLiteral)
				if elementType.getKind() == G_INTERFACE {
					if i >= len(arrayLiteral.values) {
						// zero value
						emit("LOAD_EMPTY_INTERFACE")
						emit("PUSH_INTERFACE")
						emitSave24(lhs, offsetByIndex)
						continue
					} else if arrayLiteral.values[i].getGtype().getKind() != G_INTERFACE {
						// conversion of dynamic type => interface type
						dynamicValue := arrayLiteral.values[i]
						emitConversionToInterface(dynamicValue)
						emit("LOAD_EMPTY_INTERFACE")
						emit("PUSH_INTERFACE")
						emitSave24(lhs, offsetByIndex)
						continue
					} else {
						arrayLiteral.values[i].emit()
						emitSave24(lhs, offsetByIndex)
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
		emitStoreItToLocal(decl.variable.getGtype().getSize(), decl.variable.offset, comment)
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

func emitCollectIndexSave(array Expr, index Expr, offset int) {
	assert(array.getGtype().kind == G_ARRAY, array.token(), "should be array")
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
	emitNewline()
}

func loadCollectIndex(collection Expr, index Expr, offset int) {
	emit("# loadCollectIndex")
	if collection.getGtype().kind == G_ARRAY {
		elmType := collection.getGtype().elementType
		emit("# collection.emit()")
		collection.emit()  // emit address
		emit("PUSH_PRIMITIVE") // store address of variable

		index.emit()
		emit("mov %%rax, %%rcx") // index

		size := elmType.getSize()
		assert(size > 0, nil, "size > 0")
		emit("mov $%d, %%rax", size) // size of one element
		emit("imul %%rcx, %%rax")    // index * size
		emit("PUSH_PRIMITIVE")           // store index * size
		emit("pop %%rcx")            // load  index * size
		emit("pop %%rbx")            // load address of variable
		emit("add %%rcx , %%rbx")    // (index * size) + address
		if offset > 0 {
			emit("add $%d,  %%rbx", offset)
		}
		emit("mov %%rbx, %%rax")
		if collection.getGtype().elementType.getKind() == G_INTERFACE {
			emit("LOAD_24_BY_DEREF")
		} else {
			emit("LOAD_8_BY_DEREF")
		}
		return
	} else if collection.getGtype().kind == G_SLICE {
		elmType := collection.getGtype().elementType
		emit("# emit address of the low index")
		collection.emit()  // eval pointer value
		emit("PUSH_PRIMITIVE") // store head address

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
		emit("mov %%rbx, %%rax")

		primType := collection.getGtype().elementType.getKind()
		if primType == G_INTERFACE || primType == G_MAP || primType == G_SLICE {
			emit("LOAD_24_BY_DEREF")
		} else {
			// dereference the content of an emelment
			if size == 1 {
				emit("LOAD_1_BY_DEREF")
			} else {
				emit("LOAD_8_BY_DEREF")
			}
		}
	} else if collection.getGtype().getKind() == G_MAP {
		loadMapIndexExpr(collection, index)
	} else if collection.getGtype().getKind() == G_STRING {
		// https://golang.org/ref/spec#Index_expressions
		// For a of string type:
		//
		// a constant index must be in range if the string a is also constant
		// if x is out of range at run time, a run-time panic occurs
		// a[x] is the non-constant byte value at index x and the type of a[x] is byte
		// a[x] may not be assigned to
		emit("# load head address of the string")
		collection.emit()  // emit address
		emit("PUSH_PRIMITIVE")
		index.emit()
		emit("PUSH_PRIMITIVE")
		emit("SUM_FROM_STACK")
		emit("ADD_NUMBER %d", offset)
		emit("LOAD_8_BY_DEREF")
	} else {
		TBI(collection.token(), "unable to handle %s", collection.getGtype())
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

func (e *ExprSlice) emitSubString() {
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
	emit("PUSH_PRIMITIVE")
	e.low.emit()
	emit("PUSH_PRIMITIVE")
	emit("SUM_FROM_STACK")
	emit("PUSH_PRIMITIVE")

	emitCallMallocDinamicSize(eMemSize)
	emit("PUSH_PRIMITIVE")

	eNewStrlen.emit()
	emit("PUSH_PRIMITIVE")

	emit("POP_TO_ARG_2")
	emit("POP_TO_ARG_1")
	emit("POP_TO_ARG_0")

	emit("FUNCALL iruntime.strcopy")
}

func (e *ExprSlice) emit() {
	if e.collection.getGtype().isString() {
		e.emitSubString()
	} else {
		e.emitSlice()
	}
}

func (e *ExprSlice) emitSlice() {
	elmType := e.collection.getGtype().elementType
	size := elmType.getSize()
	assert(size > 0, nil, "size > 0")

	emit("# assign to a slice")
	emit("#   emit address of the array")
	e.collection.emit()
	emit("PUSH_PRIMITIVE # head of the array")
	e.low.emit()
	emit("PUSH_PRIMITIVE # low index")
	emit("LOAD_NUMBER %d", size)
	emit("PUSH_PRIMITIVE")
	emit("IMUL_FROM_STACK")
	emit("PUSH_PRIMITIVE")
	emit("SUM_FROM_STACK")
	emit("PUSH_PRIMITIVE")

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
	emit("PUSH_PRIMITIVE")

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

	emit("PUSH_PRIMITIVE")
	emit("POP_SLICE")
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
		emit("PUSH_PRIMITIVE")
		// @TODO DRY with type switch statement
		typeLabel := groot.getTypeLabel(e.gtype)
		emit("lea .%s(%%rip), %%rax # type: %s", typeLabel, e.gtype.String())
		emitStringsEqual(true, "%rax", "%rcx")

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
	labelStart := makeLabel() + "_defer"
	labelEnd := makeLabel() + "_defer"
	ast.label = labelStart

	emit("jmp %s", labelEnd)
	emit("%s: # defer start", labelStart)

	for i := 0; i < len(retRegi); i++ {
		emit("push %%%s", retRegi[i])
	}

	ast.expr.emit()

	for i := len(retRegi) - 1; i >= 0; i-- {
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
		emit("PUSH_PRIMITIVE")
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
		emit("PUSH_PRIMITIVE")

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
		emit("PUSH_PRIMITIVE")

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
		emit("lea %s, %%rdi", slabel)
		emit("mov $0, %%rax")
		emit("call %s", ".panic")

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

func (ircall *IrStaticCall) emit(args []Expr) {
	// nothing to do
	emit("# emitCall %s", ircall.symbol)

	var numRegs int
	var param *ExprVariable
	var collectVariadicArgs bool // gather variadic args into a slice
	var variadicArgs []Expr
	var arg Expr
	var argIndex int
	for argIndex, arg = range args {
		var fromGtype string = ""
		if arg.getGtype() != nil {
			emit("# get fromGtype")
			fromGtype = arg.getGtype().String()
		}
		emit("# from %s", fromGtype)
		if argIndex < len(ircall.callee.params) {
			param = ircall.callee.params[argIndex]
			if param.isVariadic {
				if _, ok := arg.(*ExprVaArg); !ok {
					collectVariadicArgs = true
				}
			}
		}

		if collectVariadicArgs {
			variadicArgs = append(variadicArgs, arg)
			continue
		}

		var doConvertToInterface bool

		// do not convert receiver
		if !ircall.isMethodCall || argIndex != 0 {
			if param != nil && ircall.symbol != "printf" {
				emit("# has a corresponding param")

				var fromGtype *Gtype
				if arg.getGtype() != nil {
					fromGtype = arg.getGtype()
					emit("# fromGtype:%s", fromGtype.String())
				}

				var toGtype *Gtype
				if param.getGtype() != nil {
					toGtype = param.getGtype()
					emit("# toGtype:%s", toGtype.String())
				}

				if toGtype != nil && toGtype.getKind() == G_INTERFACE && fromGtype != nil && fromGtype.getKind() != G_INTERFACE {
					doConvertToInterface = true
				}
			}
		}

		if ircall.symbol == ".println" {
			doConvertToInterface = false
		}

		emit("# arg %d, doConvertToInterface=%s, collectVariadicArgs=%s",
			argIndex, bool2string(doConvertToInterface), bool2string(collectVariadicArgs))

		if doConvertToInterface {
			emit("# doConvertToInterface !!!")
			emitConversionToInterface(arg)
		} else {
			arg.emit()
		}

		var primType GTYPE_KIND = 0
		if arg.getGtype() != nil {
			primType = arg.getGtype().getKind()
		}
		var width int
		if doConvertToInterface || primType == G_INTERFACE {
			emit("PUSH_INTERFACE")
			width = interfaceWidth
		} else if primType == G_SLICE {
			emit("PUSH_SLICE")
			width = sliceWidth
		} else if primType == G_MAP {
			emit("PUSH_MAP")
			width = mapWidth
		} else {
			emit("PUSH_PRIMITIVE")
			width = 1
		}
		numRegs += width
	}

	// check if callee has a variadic
	// https://golang.org/ref/spec#Passing_arguments_to_..._parameters
	// If f is invoked with no actual arguments for p, the value passed to p is nil.
	if !collectVariadicArgs {
		if argIndex+1 < len(ircall.callee.params) {
			param = ircall.callee.params[argIndex+1]
			if param.isVariadic {
				collectVariadicArgs = true
			}
		}
	}

	if collectVariadicArgs {
		emit("# collectVariadicArgs = true")
		lenArgs := len(variadicArgs)
		if lenArgs == 0 {
			emit("LOAD_EMPTY_SLICE")
			emit("PUSH_SLICE")
		} else {
			// var a []interface{}
			for vargIndex, varg := range variadicArgs {
				emit("# emit variadic arg")
				if vargIndex == 0 {
					emit("# make an empty slice to append")
					emit("LOAD_EMPTY_SLICE")
					emit("PUSH_SLICE")
				}
				// conversion : var ifc = x
				if varg.getGtype().getKind() == G_INTERFACE {
					varg.emit()
				} else {
					emitConversionToInterface(varg)
				}
				emit("PUSH_INTERFACE")
				emit("# calling append24")
				emit("POP_TO_ARG_5 # ifc_c")
				emit("POP_TO_ARG_4 # ifc_b")
				emit("POP_TO_ARG_3 # ifc_a")
				emit("POP_TO_ARG_2 # cap")
				emit("POP_TO_ARG_1 # len")
				emit("POP_TO_ARG_0 # ptr")
				emit("mov $0, %%rax")
				emit("call iruntime.append24")
				emit("PUSH_SLICE")
			}
		}
		numRegs += 3
	}

	for i := numRegs - 1; i >= 0; i-- {
		if i >= len(RegsForArguments) {
			errorft(args[0].token(), "too many arguments")
		}
		emit("POP_TO_ARG_%d", i)
	}

	emit("mov $0, %%rax")
	emit("call %s", ircall.symbol)
	emitNewline()
}

func emitRuntimeArgs() {
	emitWithoutIndent(".runtime_args:")
	emit("push %%rbp")
	emit("mov %%rsp, %%rbp")

	emit("# set argv, argc, argc")
	emit("mov runtimeArgv(%%rip), %%rax # ptr")
	emit("mov runtimeArgc(%%rip), %%rbx # len")
	emit("mov runtimeArgc(%%rip), %%rcx # cap")

	emitFuncEpilogue(".runtime_args_noop_handler", nil)
}

func emitMainFunc(importOS bool) {
	fname := "main"
	emit(".global	%s", fname)
	emitWithoutIndent("%s:", fname)
	emit("push %%rbp")
	emit("mov %%rsp, %%rbp")

	emit("mov %%rsi, runtimeArgv(%%rip)")
	emit("mov %%rdi, runtimeArgc(%%rip)")
	emit("mov $0, %%rsi")
	emit("mov $0, %%rdi")

	// init runtime
	emit("# init runtime")
	emit("mov $0, %%rax")
	emit("call iruntime.init")

	// init imported packages
	if importOS {
		emit("# init os")
		emit("mov $0, %%rax")
		emit("call os.init")
	}

	emitNewline()
	emit("mov $0, %%rax")
	emit("call main.main")
	emitFuncEpilogue("noop_handler", nil)
}

func emitMakeSliceFunc() {
	// makeSlice
	emitWithoutIndent("%s:", "iruntime.makeSlice")
	emit("push %%rbp")
	emit("mov %%rsp, %%rbp")
	emitNewline()
	emit("PUSH_ARG_2") // -8
	emit("PUSH_ARG_1") // -16
	emit("PUSH_ARG_0") // -24

	emit("mov -16(%%rbp), %%rax # newcap")
	emit("mov -8(%%rbp), %%rcx # unit")
	emit("imul %%rcx, %%rax")
	emit("add $16, %%rax") // pure buffer

	emit("PUSH_PRIMITIVE")
	emit("POP_TO_ARG_0")
	emit("mov $0, %%rax")
	emit("call iruntime.malloc")

	emit("mov -24(%%rbp), %%rbx # newlen")
	emit("mov -16(%%rbp), %%rcx # newcap")

	emit("leave")
	emit("ret")
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

// gloabal var which should be initialized with zeros
// https://en.wikipedia.org/wiki/.bss
func (decl *DeclVar) emitBss() {
	emit(".data")
	// https://sourceware.org/binutils/docs-2.30/as/Lcomm.html#Lcomm
	emit(".lcomm %s, %d", decl.variable.varname, decl.variable.getGtype().getSize())
}

func (decl *DeclVar) emitData() {
	ptok := decl.token()
	gtype := decl.variable.gtype
	right := decl.initval

	emit("# emitData()")
	emit(".data 0")
	emitWithoutIndent("%s: # gtype=%s", decl.variable.varname, gtype.String())
	emit("# right.gtype = %s", right.getGtype().String())
	doEmitData(ptok, right.getGtype(), right, "", 0)
}

func (e *ExprStructLiteral) lookup(fieldname identifier) Expr {
	for _, field := range e.fields {
		if field.key == fieldname {
			return field.value
		}
	}

	return nil
}

func doEmitData(ptok *Token /* left type */, gtype *Gtype, value /* nullable */ Expr, containerName string, depth int) {
	emit("# doEmitData: containerName=%s, depth=%d", containerName, depth)
	primType := gtype.getKind()
	if primType == G_ARRAY {
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
				doEmitData(ptok, elmType, nil, selector, depth)
			} else {
				value := arrayliteral.values[i]
				assertNotNil(value != nil, nil)
				size := elmType.getSize()
				if size == 8 {
					if value.getGtype().kind == G_STRING {
						stringLiteral, ok := value.(*ExprStringLiteral)
						assert(ok, nil, "ok")
						emit(".quad .%s", stringLiteral.slabel)
					} else {
						switch value.(type) {
						case *ExprUop:
							uop := value.(*ExprUop)
							rel, ok := uop.operand.(*Relation)
							assert(ok, uop.token(), "only variable is allowed")
							emit(".quad %s # %s %s", rel.name, value.getGtype().String(), selector)
						case *Relation:
							assert(false, value.token(), "variable here is not allowed")
						default:
							emit(".quad %d # %s %s", evalIntExpr(value), value.getGtype().String(), selector)
						}
					}
				} else if size == 1 {
					emit(".byte %d", evalIntExpr(value))
				} else {
					doEmitData(ptok, gtype.elementType, value, selector, depth)
				}
			}
		}
		emit(".quad 0 # nil terminator")

	} else if primType == G_SLICE {
		switch value.(type) {
		case nil:
			return
		case *ExprSliceLiteral:
			// initialize a hidden array
			lit := value.(*ExprSliceLiteral)
			arrayLiteral := &ExprArrayLiteral{
				gtype:  lit.invisiblevar.gtype,
				values: lit.values,
			}

			emitDataAddr(arrayLiteral, depth)               // emit underlying array
			emit(".quad %d", lit.invisiblevar.gtype.length) // len
			emit(".quad %d", lit.invisiblevar.gtype.length) // cap
		default:
			TBI(ptok, "unable to handle gtype %s", gtype.String())
		}
	} else if primType == G_MAP || primType == G_INTERFACE {
		// @TODO
		emit(".quad 0")
		emit(".quad 0")
		emit(".quad 0")
	} else if primType == G_BOOL {
		if value == nil {
			// zero value
			emit(".quad %d # %s %s", 0, gtype.String(), containerName)
			return
		}
		val := evalIntExpr(value)
		emit(".quad %d # %s %s", val, gtype.String(), containerName)
	} else if primType == G_STRUCT {
		containerName = containerName + "." + string(gtype.relation.name)
		gtype.relation.gtype.calcStructOffset()
		for _, field := range gtype.relation.gtype.fields {
			emit("# padding=%d", field.padding)
			switch field.padding {
			case 1:
				emit(".byte 0 # padding")
			case 4:
				emit(".double 0 # padding")
			case 8:
				emit(".quad 0 # padding")
			default:
			}
			emit("# field:offesr=%d, fieldname=%s", field.offset, field.fieldname)
			if value == nil {
				doEmitData(ptok, field, nil, containerName+"."+string(field.fieldname), depth)
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
			doEmitData(ptok, gtype, value, containerName+"."+string(field.fieldname), depth)
		}
	} else {
		var val int
		switch value.(type) {
		case nil:
			emit(".quad %d # %s %s zero value", 0, gtype.String(), containerName)
		case *ExprNumberLiteral:
			val = value.(*ExprNumberLiteral).val
			emit(".quad %d # %s %s", val, gtype.String(), containerName)
		case *ExprConstVariable:
			cnst := value.(*ExprConstVariable)
			val = evalIntExpr(cnst)
			emit(".quad %d # %s ", val, gtype.String())
		case *ExprVariable:
			vr := value.(*ExprVariable)
			val = evalIntExpr(vr)
			emit(".quad %d # %s ", val, gtype.String())
		case *ExprBinop:
			val = evalIntExpr(value)
			emit(".quad %d # %s ", val, gtype.String())
		case *ExprStringLiteral:
			stringLiteral := value.(*ExprStringLiteral)
			emit(".quad .%s", stringLiteral.slabel)
		case *Relation:
			rel := value.(*Relation)
			doEmitData(ptok, gtype, rel.expr, "rel", depth)
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
				emitDataAddr(operand, depth)
			}
		default:
			TBI(ptok, "unable to handle %d", primType)
		}
	}
}

// this logic is stolen from 8cc.
func emitDataAddr(operand Expr, depth int) {
	emit(".data %d", depth+1)
	label := makeLabel()
	emit("%s:", label)
	doEmitData(nil, operand.getGtype(), operand, "", depth+1)
	emit(".data %d", depth)
	emit(".quad %s", label)
}

func (decl *DeclVar) emitGlobal() {
	emitWithoutIndent("# emitGlobal for %s", decl.variable.varname)
	assertNotNil(decl.variable.gtype != nil, nil)

	if decl.initval == nil {
		decl.emitBss()
	} else {
		decl.emitData()
	}
}

func makeDynamicTypeLabel(id int) string {
	return fmt.Sprintf("DynamicTypeId%d", id)
}

func (root *IrRoot) getTypeLabel(gtype *Gtype) string {
	dynamicTypeId := get_index(gtype.String(), root.uniquedDTypes)
	if dynamicTypeId == -1 {
		errorft(nil, "type %s not found in uniquedDTypes", gtype.String())
	}
	return makeDynamicTypeLabel(dynamicTypeId)
}

// builtin string
var builtinStringKey1 string = "SfmtDumpInterface"
var builtinStringValue1 string = "# interface = {ptr:%p,receiverTypeId:%d,dtype:'%s'}\\n"
var builtinStringKey2 string = "SfmtDumpSlice"
var builtinStringValue2 string = "# slice = {underlying:%p,len:%d,cap:%d}\\n"

func (root *IrRoot) emitSpecialStrings() {
	// https://sourceware.org/binutils/docs-2.30/as/Data.html#Data
	emit(".data 0")
	emit("# special strings")

	// emit builtin string
	emitWithoutIndent(".%s:", builtinStringKey1)
	emit(".string \"%s\"", builtinStringValue1)
	emitWithoutIndent(".%s:", builtinStringKey2)
	emit(".string \"%s\"", builtinStringValue2)

	// empty string
	eEmptyString.slabel = "empty"
	emitWithoutIndent(".%s:", eEmptyString.slabel)
	emit(".string \"%s\"", eEmptyString.val)
}

func (root *IrRoot) emitDynamicTypes() {
	emitNewline()
	emit("# Dynamic Types")
	for dynamicTypeId, gs := range root.uniquedDTypes {
		label := makeDynamicTypeLabel(dynamicTypeId)
		emitWithoutIndent(".%s:", label)
		emit(".string \"%s\"", gs)
	}
}

func (root *IrRoot) emitMethodTable() {
	emit("# Method table")

	emitWithoutIndent("%s:", "receiverTypes")
	emit(".quad 0 # receiverTypeId:0")
	for i := 1; i <= len(root.methodTable); i++ {
		emit(".quad receiverType%d # receiverTypeId:%d", i, i)
	}

	var shortMethodNames []string

	for i := 1; i <= len(root.methodTable); i++ {
		emitWithoutIndent("receiverType%d:", i)
		mt := root.methodTable
		methods, ok := mt[i]
		if !ok {
			debugf("methods not found in methodTable %d", i)
			continue
		}
		for _, methodNameFull := range methods {
			splitted := strings.Split(methodNameFull, "$")
			shortMethodName := splitted[1]
			emit(".quad .M%s # key", shortMethodName)
			emit(".quad %s # method", methodNameFull)
			if !in_array(shortMethodName, shortMethodNames) {
				shortMethodNames = append(shortMethodNames, shortMethodName)
			}
		}
	}

	emit("# METHOD NAMES")
	for _, shortMethodName := range shortMethodNames {
		emitWithoutIndent(".M%s:", shortMethodName)
		emit(".string \"%s\"", shortMethodName)
	}

}

// generate code
func (root *IrRoot) emit() {
	groot = root

	emitMacroDefinitions()

	emit(".data 0")
	root.emitSpecialStrings()
	root.emitDynamicTypes()
	root.emitMethodTable()

	emitWithoutIndent(".text")
	emitRuntimeArgs()
	emitMainFunc(root.importOS)
	emitMakeSliceFunc()

	// emit packages
	for _, pkg := range root.packages {
		emitWithoutIndent("#--------------------------------------------------------")
		emitWithoutIndent("# package %s", pkg.name)
		emitWithoutIndent("# string literals")
		emitWithoutIndent(".data 0")
		for _, ast := range pkg.stringLiterals {
			emitWithoutIndent(".%s:", ast.slabel)
			// https://sourceware.org/binutils/docs-2.30/as/String.html#String
			// the assembler marks the end of each string with a 0 byte.
			emit(".string \"%s\"", ast.val)
		}

		for _, vardecl := range pkg.vars {
			emitNewline()
			vardecl.emit()
		}
		emitNewline()

		emitWithoutIndent(".text")
		for _, funcdecl := range pkg.funcs {
			funcdecl.emit()
			emitNewline()
		}

	}

}
