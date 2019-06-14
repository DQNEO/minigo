// gen_load handles loading of expressions
package main

import "fmt"

func (ast *ExprNumberLiteral) emit() {
	emit("LOAD_NUMBER %d", ast.val)
}

func (ast *ExprStringLiteral) emit() {
	emit("LOAD_STRING_LITERAL .%s", ast.slabel)
}

func loadStructField(strct Expr, field *Gtype, offset int) {
	strct = unwrapRel(strct)
	emit("# loadStructField")
	switch strct.(type) {
	case *ExprVariable:
		variable := strct.(*ExprVariable)
		if field.getKind() == G_ARRAY {
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
		strcttype := a.strct.getGtype().Underlying()
		assert(strcttype.size > 0, a.token(), "struct size should be > 0")
		field2 := strcttype.getField(a.fieldname)
		loadStructField(a.strct, field2, offset+field.offset)
	case *ExprIndex: // array[1].field
		indexExpr := strct.(*ExprIndex)
		indexExpr.emitOffsetLoad(offset+field.offset)
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

func (a *ExprStructField) emit() {
	emit("# LOAD ExprStructField")
	switch a.strct.getGtype().getKind() {
	case G_POINTER: // pointer to struct
		strcttype := a.strct.getGtype().origType.relation.gtype
		field := strcttype.getField(a.fieldname)
		a.strct.emit()
		emit("ADD_NUMBER %d", field.offset)
		switch field.is24WidthType() {
		case true:
			emit("LOAD_24_BY_DEREF")
		default:
			emit("LOAD_8_BY_DEREF")
		}

	case G_STRUCT:
		strcttype := a.strct.getGtype().relation.gtype
		assert(strcttype.size > 0, a.token(), "struct size should be > 0")
		field := strcttype.getField(a.fieldname)
		loadStructField(a.strct, field, 0)
	default:
		errorft(a.token(), "internal error: bad gtype %s", a.strct.getGtype().String())
	}
}


func (e *ExprStructField) emitOffsetLoad(size int, offset int) {
	strct := e.strct
	strct = unwrapRel(strct)
	vr, ok := strct.(*ExprVariable)
	assert(ok, e.tok, "should be *ExprVariable")
	assert(vr.gtype.kind == G_NAMED, e.tok, "expect G_NAMED, but got "+vr.gtype.String())
	field := vr.gtype.relation.gtype.getField(e.fieldname)
	vr.emitOffsetLoad(size, field.offset+offset)
}

func (ast *ExprVariable) emit() {
	emit("# load variable \"%s\" %s", ast.varname, ast.getGtype().String())
	if ast.isGlobal {
		if ast.gtype.getKind() == G_ARRAY {
			ast.emitAddress(0)
		} else if ast.getGtype().is24WidthType() {
			emit("LOAD_24_FROM_GLOBAL %s", ast.varname)
		} else if ast.getGtype().getSize() == 1 {
			emit("LOAD_1_FROM_GLOBAL_CAST %s", ast.varname)
		} else {
			emit("LOAD_8_FROM_GLOBAL %s", ast.varname)
		}

	} else {
		if ast.offset == 0 {
			errorft(ast.token(), "offset should not be zero for localvar %s", ast.varname)
		}
		if ast.gtype.getKind() == G_ARRAY {
			ast.emitAddress(0)
		} else if ast.gtype.is24WidthType() {
			emit("LOAD_24_FROM_LOCAL %d", ast.offset)
		} else if ast.getGtype().getSize() == 1 {
			emit("LOAD_1_FROM_LOCAL_CAST %d", ast.offset)
		} else {
			emit("LOAD_8_FROM_LOCAL %d", ast.offset)
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
	assert(rel.expr != nil, rel.token(), fmt.Sprintf("rel.expr is nil: %s", rel.name))
	rel.expr.emit()
}

func (ast *ExprConstVariable) emit() {
	emit("# *ExprConstVariable.emit() name=%s iotaindex=%d", ast.name, ast.iotaIndex)
	assert(ast.iotaIndex < 10000, ast.token(), "iotaindex is too large")
	assert(ast.val != nil, ast.token(), "const.val for should not be nil:"+string(ast.name))
	if ast.hasIotaValue() {
		emit("# const is iota")
		val := &ExprNumberLiteral{
			val: ast.iotaIndex,
		}
		val.emit()
	} else {
		emit("# const is not iota")
		ast.val.emit()
	}
}

func (ast *ExprUop) emit() {
	operand := unwrapRel(ast.operand)
	ast.operand = operand
	emit("# emitting ExprUop")
	if ast.op == "&" {
		switch ast.operand.(type) {
		case *ExprVariable:
			vr := ast.operand.(*ExprVariable)
			vr.emitAddress(0)
		case *ExprStructLiteral:
			e := ast.operand.(*ExprStructLiteral)
			assert(e.invisiblevar.offset != 0, nil, "ExprStructLiteral's invisible var has offset")
			ivv := e.invisiblevar
			assignToStruct(ivv, e)

			emitCallMalloc(e.getGtype().getSize())
			emit("PUSH_8")                     // to:ptr addr
			e.invisiblevar.emitAddress(0)
			emit("PUSH_8") // from:address of invisible var
			emitCopyStructFromStack(e.getGtype().getSize())
			// emit address
		case *ExprStructField:
			e := ast.operand.(*ExprStructField)
			e.emitAddress()
		default:
			errorft(ast.token(), "Unknown type: %T", ast.operand)
		}
	} else if ast.op == "*" {
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

func (variable *ExprVariable) emitOffsetLoad(size int, offset int) {
	assert(0 <= size && size <= 8, variable.token(), "invalid size")
	if variable.isGlobal {
		emit("LOAD_%d_FROM_GLOBAL %s %d", size, variable.varname, offset)
	} else {
		emit("LOAD_%d_FROM_LOCAL %d+%d", size,  variable.offset, offset)
	}
}

// rax: address
// rbx: len
// rcx: cap
func (e *ExprSliceLiteral) emit() {
	emit("# (*ExprSliceLiteral).emit()")
	length := len(e.values)
	//debugf("slice literal %s: underlyingarray size = %d (should be %d)", e.getGtype(), e.gtype.getSize(),  e.gtype.elementType.getSize() * length)
	emitCallMalloc(e.gtype.getSize() * length)
	emit("PUSH_8 # ptr")
	for i, value := range e.values {
		if e.gtype.elementType.getKind() == G_INTERFACE && value.getGtype().getKind() != G_INTERFACE {
			emitConversionToInterface(value)
		} else {
			value.emit()
		}

		emit("pop %%r10 # ptr")

		if e.gtype.elementType.is24WidthType() {
			emit("mov %%rax, %d+%d(%%r10)", IntSize*3*i,0)
			emit("mov %%rbx, %d+%d(%%r10)", IntSize*3*i,8)
			emit("mov %%rcx, %d+%d(%%r10)", IntSize*3*i,16)
		} else if e.gtype.elementType.getSize() <= 8 {
			emit("mov %%rax, %d(%%r10)", IntSize*i)
		} else {
			TBI(e.token(), "")
		}
		emit("push %%r10 # ptr")
	}

	emit("pop %%rax # ptr")
	emit("mov $%d, %%rbx # len", length)
	emit("mov $%d, %%rcx # cap", length)
}

func emitAddress(e Expr) {
	e = unwrapRel(e)
	switch e.(type) {
	case *ExprVariable:
		e.(*ExprVariable).emitAddress(0)
	default:
		TBI(e.token(), "")
	}
}

func emitOffsetLoad(lhs Expr, size int, offset int) {
	lhs = unwrapRel(lhs)
	emit("# emitOffsetLoad(offset %d)", offset)
	switch lhs.(type) {
	case *ExprVariable:
		variable := lhs.(*ExprVariable)
		variable.emitOffsetLoad(size, offset)
	case *ExprStructField:
		structfield := lhs.(*ExprStructField)
		fieldType := structfield.getGtype()
		if structfield.strct.getGtype().getKind() == G_POINTER {
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
		indexExpr.emitOffsetLoad(offset)
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

func loadArrayOrSliceIndex(collection Expr, index Expr, offset int) {
	elmType := collection.getGtype().elementType
	elmSize := elmType.getSize()
	assert(elmSize > 0, nil, "elmSize > 0")

	collection.emit()
	emit("PUSH_8 # head")

	index.emit()
	emit("IMUL_NUMBER %d", elmSize)
	emit("PUSH_8 # index * elmSize")

	emit("SUM_FROM_STACK # (index * elmSize) + head")
	emit("ADD_NUMBER %d", offset)

	// dereference the content of an emelment
	if elmType.is24WidthType() {
		emit("LOAD_24_BY_DEREF")
	} else if elmSize == 1 {
			emit("LOAD_1_BY_DEREF")
	} else {
		emit("LOAD_8_BY_DEREF")
	}
}

func (e *ExprIndex) emitOffsetLoad(offset int) {
	emit("# ExprIndex.emitOffsetLoad")
	collection := e.collection
	index := e.index
	switch collection.getGtype().getKind() {
	case G_ARRAY, G_SLICE:
		loadArrayOrSliceIndex(collection, index, offset)
		return
	case G_MAP:
		loadMapIndexExpr(e)
	case G_STRING:
		// https://golang.org/ref/spec#Index_expressions
		// For a of string type:
		//
		// a constant index must be in range if the string a is also constant
		// if x is out of range at run time, a run-time panic occurs
		// a[x] is the non-constant byte value at index x and the type of a[x] is byte
		// a[x] may not be assigned to
		emit("# load head address of the string")
		collection.emit() // emit address
		emit("PUSH_8")
		index.emit()
		emit("PUSH_8")
		emit("SUM_FROM_STACK")
		emit("ADD_NUMBER %d", offset)
		emit("LOAD_8_BY_DEREF")
	default:
		TBI(collection.token(), "unable to handle %s", collection.getGtype())
	}
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
	emit("PUSH_8")
	e.low.emit()
	emit("PUSH_8")
	emit("SUM_FROM_STACK")
	emit("PUSH_8")

	emitCallMallocDinamicSize(eMemSize)
	emit("PUSH_8")

	eNewStrlen.emit()
	emit("PUSH_8")

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
	emit("PUSH_8 # head of the array")
	e.low.emit()
	emit("PUSH_8 # low index")
	emit("LOAD_NUMBER %d", size)
	emit("PUSH_8")
	emit("IMUL_FROM_STACK")
	emit("PUSH_8")
	emit("SUM_FROM_STACK")
	emit("PUSH_8")

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
	emit("PUSH_8")

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

	emit("PUSH_8")
	emit("POP_SLICE")
}



