// gen_load handles loading of expressions
package main

func (ast *ExprNumberLiteral) emit() {
	emit(S("LOAD_NUMBER %d"), ast.val)
}

func loadStructField(strct Expr, field *Gtype, offset int) {
	strct = unwrapRel(strct)
	emit(S("# loadStructField"))
	switch strct.(type) {
	case *ExprVariable:
		variable := strct.(*ExprVariable)
		if field.getKind() == G_ARRAY {
			variable.emitAddress(field.offset)
		} else {
			if variable.isGlobal {
				emit(S("LOAD_8_FROM_GLOBAL %s, %d+%d"), gostring(variable.varname), field.offset, offset)
			} else {
				emit(S("LOAD_8_FROM_LOCAL %d+%d+%d"), variable.offset, field.offset, offset)
			}
		}
	case *ExprStructField: // strct.field.field
		a := strct.(*ExprStructField)
		strcttype := a.strct.getGtype().Underlying()
		assert(strcttype.size > 0, a.token(), S("struct size should be > 0"))
		field2 := strcttype.getField(a.fieldname)
		loadStructField(a.strct, field2, offset+field.offset)
	case *ExprIndex: // array[1].field
		indexExpr := strct.(*ExprIndex)
		indexExpr.emitOffsetLoad(offset + field.offset)
	default:
		// funcall().field
		// methodcall().field
		// *ptr.field
		// (MyStruct{}).field
		// (&MyStruct{}).field
		TBI(strct.token(), S("unable to handle %T"), strct)
	}

}

func (a *ExprStructField) emitAddress() {
	strcttype := a.strct.getGtype().origType.relation.gtype
	field := strcttype.getField(a.fieldname)
	a.strct.emit()
	emit(S("ADD_NUMBER %d"), field.offset)
}

func (a *ExprStructField) emit() {
	emit(S("# LOAD ExprStructField"))
	switch a.strct.getGtype().getKind() {
	case G_POINTER: // pointer to struct
		strcttype := a.strct.getGtype().origType.relation.gtype
		field := strcttype.getField(a.fieldname)
		a.strct.emit()
		emit(S("ADD_NUMBER %d"), field.offset)
		switch field.is24WidthType() {
		case true:
			emit(S("LOAD_24_BY_DEREF"))
		default:
			emit(S("LOAD_8_BY_DEREF"))
		}

	case G_STRUCT:
		strcttype := a.strct.getGtype().relation.gtype
		assert(strcttype.size > 0, a.token(), S("struct size should be > 0"))
		field := strcttype.getField(a.fieldname)
		loadStructField(a.strct, field, 0)
	default:
		errorft(a.token(), S("internal error: bad gtype %s"), a.strct.getGtype().String())
	}
}

func (e *ExprStructField) emitOffsetLoad(size int, offset int) {
	strct := e.strct
	strct = unwrapRel(strct)
	vr, ok := strct.(*ExprVariable)
	assert(ok, e.tok, S("should be *ExprVariable"))
	assert(vr.gtype.kind == G_NAMED, e.tok, S("expect G_NAMED, but got "), vr.gtype.String())
	field := vr.gtype.relation.gtype.getField(e.fieldname)
	vr.emitOffsetLoad(size, field.offset+offset)
}

func (ast *ExprVariable) emit() {
	emit(S("# load variable \"%s\" %s"), gostring(ast.varname), ast.getGtype().String())
	if ast.isGlobal {
		if ast.gtype.getKind() == G_ARRAY {
			ast.emitAddress(0)
		} else if ast.getGtype().is24WidthType() {
			emit(S("LOAD_24_FROM_GLOBAL %s"), gostring(ast.varname))
		} else if ast.getGtype().getSize() == 1 {
			emit(S("LOAD_1_FROM_GLOBAL_CAST %s"), gostring(ast.varname))
		} else {
			emit(S("LOAD_8_FROM_GLOBAL %s"), gostring(ast.varname))
		}

	} else {
		if ast.offset == 0 {
			errorft(ast.token(), S("offset should not be zero for localvar %s"), ast.varname)
		}
		if ast.gtype.getKind() == G_ARRAY {
			ast.emitAddress(0)
		} else if ast.gtype.is24WidthType() {
			emit(S("LOAD_24_FROM_LOCAL %d"), ast.offset)
		} else if ast.getGtype().getSize() == 1 {
			emit(S("LOAD_1_FROM_LOCAL_CAST %d"), ast.offset)
		} else {
			emit(S("LOAD_8_FROM_LOCAL %d"), ast.offset)
		}
	}
}

func (variable *ExprVariable) emitAddress(offset int) {
	if variable.isGlobal {
		emit(S("LOAD_GLOBAL_ADDR %s, %d"), gostring(variable.varname), offset)
	} else {
		if variable.offset == 0 {
			errorft(variable.token(), S("offset should not be zero for localvar %s"), variable.varname)
		}
		emit(S("LOAD_LOCAL_ADDR %d+%d"), variable.offset, offset)
	}
}

func (rel *Relation) emit() {
	assert(rel.expr != nil, rel.token(), S("rel.expr is nil"))
	rel.expr.emit()
}

func (ast *ExprConstVariable) emit() {
	emit(S("# *ExprConstVariable.emit() name=%s iotaindex=%d"), gostring(ast.name), ast.iotaIndex)
	assert(ast.iotaIndex < 10000, ast.token(), S("iotaindex is too large"))
	assert(ast.val != nil, ast.token(), S("const.val for should not be nil:%s"), gostring(ast.name))
	if ast.hasIotaValue() {
		emit(S("# const is iota"))
		val := &ExprNumberLiteral{
			val: ast.iotaIndex,
		}
		val.emit()
	} else {
		emit(S("# const is not iota"))
		ast.val.emit()
	}
}

func (ast *ExprUop) emit() {
	operand := unwrapRel(ast.operand)
	ast.operand = operand
	emit(S("# emitting ExprUop"))
	op := ast.op
	switch cstring(op) {

	case "&" :
		switch ast.operand.(type) {
		case *ExprVariable:
			vr := ast.operand.(*ExprVariable)
			vr.emitAddress(0)
		case *ExprStructLiteral:
			e := ast.operand.(*ExprStructLiteral)
			assert(e.invisiblevar.offset != 0, nil, S("ExprStructLiteral's invisible var has offset"))
			ivv := e.invisiblevar
			assignToStruct(ivv, e)

			emitCallMalloc(e.getGtype().getSize())
			emit(S("PUSH_8")) // to:ptr addr
			e.invisiblevar.emitAddress(0)
			emit(S("PUSH_8")) // from:address of invisible var
			emitCopyStructFromStack(e.getGtype().getSize())
			// emit address
		case *ExprStructField:
			e := ast.operand.(*ExprStructField)
			e.emitAddress()
		case *ExprIndex:
			e := ast.operand.(*ExprIndex)
			emitAddress(e)
		default:
			errorft(ast.token(), S("Unknown type: %T"), ast.operand)
		}
	case "*":
		ast.operand.emit()
		emit(S("LOAD_8_BY_DEREF"))
	case "!":
		ast.operand.emit()
		emit(S("CMP_EQ_ZERO"))
	case "-":
		// delegate to biop
		// -(x) -> (-1) * (x)
		left := &ExprNumberLiteral{val: -1}
		binop := &ExprBinop{
			op:    gostring("*"),
			left:  left,
			right: ast.operand,
		}
		binop.emit()
	default:
		errorft(ast.token(), S("unable to handle uop %s"), ast.op)
	}
	//debugf(S("end of emitting ExprUop"))

}

func (variable *ExprVariable) emitOffsetLoad(size int, offset int) {
	assert(0 <= size && size <= 8, variable.token(), S("invalid size"))
	if variable.isGlobal {
		emit(S("LOAD_%d_FROM_GLOBAL %s %d"), size, gostring(variable.varname), offset)
	} else {
		emit(S("LOAD_%d_FROM_LOCAL %d+%d"), size, variable.offset, offset)
	}
}

// rax: address
// rbx: len
// rcx: cap
func (e *ExprSliceLiteral) emit() {
	emit(S("# (*ExprSliceLiteral).emit()"))
	var length int = len(e.values)
	//debugf(S("slice literal %s: underlyingarray size = %d (should be %d)"), e.getGtype(), e.gtype.getSize(),  e.gtype.elementType.getSize() * length)
	emitCallMalloc(e.gtype.getSize() * length)
	emit(S("PUSH_8 # ptr"))
	for i, value := range e.values {
		if e.gtype.elementType.getKind() == G_INTERFACE && value.getGtype().getKind() != G_INTERFACE {
			emitConversionToInterface(value)
		} else {
			value.emit()
		}

		emit(S("pop %%r10 # ptr"))

		var baseOffset int = IntSize*3*i
		if e.gtype.elementType.is24WidthType() {
			emit(S("mov %%rax, %d+%d(%%r10)"), baseOffset, offset0)
			emit(S("mov %%rbx, %d+%d(%%r10)"), baseOffset, offset8)
			emit(S("mov %%rcx, %d+%d(%%r10)"), baseOffset, offset16)
		} else if e.gtype.elementType.getSize() <= 8 {
			var offset int = IntSize*i
			emit(S("mov %%rax, %d(%%r10)"), offset)
		} else {
			TBI(e.token(), S(""))
		}
		emit(S("push %%r10 # ptr"))
	}

	emit(S("pop %%rax # ptr"))
	emit(S("mov $%d, %%rbx # len"), length)
	emit(S("mov $%d, %%rcx # cap"), length)
}

func emitAddress(e Expr) {
	e = unwrapRel(e)
	switch e.(type) {
	case *ExprVariable:
		e.(*ExprVariable).emitAddress(0)
	case *ExprIndex:
		e.(*ExprIndex).emitAddress()
	default:
		TBI(e.token(), S(""))
	}
}

func emitOffsetLoad(lhs Expr, size int, offset int) {
	lhs = unwrapRel(lhs)
	emit(S("# emitOffsetLoad(offset %d)"), offset)
	switch lhs.(type) {
	case *ExprVariable:
		variable := lhs.(*ExprVariable)
		variable.emitOffsetLoad(size, offset)
	case *ExprStructField:
		structfield := lhs.(*ExprStructField)
		fieldType := structfield.getGtype()
		if structfield.strct.getGtype().getKind() == G_POINTER {
			structfield.strct.emit() // emit address of the struct
			var sum int = fieldType.offset+offset
			emit(S("# offset %d + %d = %d"), fieldType.offset, offset, sum)
			emit(S("ADD_NUMBER %d+%d"), fieldType.offset, offset)
			//reg := getReg(size)
			emit(S("LOAD_8_BY_DEREF"))
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
		assert(len(rettypes) == 1, lhs.token(), S("rettype should be single"))
		rettype := rettypes[0]
		assert(rettype.getKind() == G_POINTER, lhs.token(), S("only pointer is supported"))
		mcall.emit()
		emit(S("ADD_NUMBER %d"), offset)
		emit(S("LOAD_8_BY_DEREF"))
	default:
		errorft(lhs.token(), S("unkonwn type %T"), lhs)
	}
}


func (e *ExprIndex) emitAddressOfArrayOrSliceIndex() {
	collection := e.collection
	index := e.index
	elmType := collection.getGtype().Underlying().elementType
	assert(elmType != nil, collection.token(), S("elmType should not be nil"))
	elmSize := elmType.getSize()
	assert(elmSize > 0, nil, S("elmSize > 0"))

	collection.emit()
	emit(S("PUSH_8 # head"))

	index.emit()
	emit(S("IMUL_NUMBER %d"), elmSize)
	emit(S("PUSH_8 # index * elmSize"))

	emit(S("SUM_FROM_STACK # (index * elmSize) + head"))
}

func (e *ExprIndex) loadArrayOrSliceIndex(offset int) {
	elmType := e.collection.getGtype().Underlying().elementType
	elmSize := elmType.getSize()

	e.emitAddressOfArrayOrSliceIndex()
	emit(S("ADD_NUMBER %d"), offset)

	// dereference the content of an emelment
	if elmType.is24WidthType() {
		emit(S("LOAD_24_BY_DEREF"))
	} else if elmSize == 1 {
		emit(S("LOAD_1_BY_DEREF"))
	} else {
		emit(S("LOAD_8_BY_DEREF"))
	}
}

func (e *ExprIndex) emitAddress() {
	switch e.collection.getGtype().getKind() {
	case G_ARRAY, G_SLICE:
		e.emitAddressOfArrayOrSliceIndex()
	default:
		TBI(e.collection.token(), S(""))
	}
}

func (e *ExprIndex) emitOffsetLoad(offset int) {
	emit(S("# ExprIndex.emitOffsetLoad"))
	switch e.collection.getGtype().getKind() {
	case G_ARRAY, G_SLICE:
		e.loadArrayOrSliceIndex(offset)
		return
	case G_MAP:
		loadMapIndexExpr(e)
	default:
		TBI(e.collection.token(), S("unable to handle %s"), e.collection.getGtype())
	}
}

func (e *ExprSlice) emit() {
	if e.collection.getGtype().isString() {
		TBI(e.token(), S(""))
	} else {
		e.emitSlice()
	}
}

func (e *ExprSlice) emitSlice() {
	elmType := e.collection.getGtype().Underlying().elementType
	assert(elmType != nil, e.token(),S("type should not be nil:T %s"), e.collection.getGtype().String())
	size := elmType.getSize()
	assert(size > 0, nil, S("size > 0"))

	emit(S("# assign to a slice"))
	emit(S("#   emit address of the array"))
	e.collection.emit()
	emit(S("PUSH_8 # head of the array"))
	e.low.emit()
	emit(S("PUSH_8 # low index"))
	emit(S("LOAD_NUMBER %d"), size)
	emit(S("PUSH_8"))
	emit(S("IMUL_FROM_STACK"))
	emit(S("PUSH_8"))
	emit(S("SUM_FROM_STACK"))
	emit(S("PUSH_8"))

	emit(S("#   calc and set len"))

	if e.high == nil {
		e.high = &ExprLen{
			tok:e.token(),
			arg: e.collection,
		}
	}
	calcLen := &ExprBinop{
		op:    gostring("-"),
		left:  e.high,
		right: e.low,
	}
	calcLen.emit()
	emit(S("PUSH_8"))

	emit(S("#   calc and set cap"))
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
		op:    gostring("-"),
		left:  max,
		right: e.low,
	}

	calcCap.emit()

	emit(S("PUSH_8"))
	emit(S("POP_SLICE"))
}
