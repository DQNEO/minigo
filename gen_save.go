package main

func emitPop(gtype *Gtype) {
	if gtype.is24WidthType() {
		emit(S("POP_24"))
	} else {
		emit(S("POP_8"))
	}
}

func emitOffsetSave(lhs Expr, offset int) {
	if lhs.getGtype().is24WidthType() {
		emitSave24(lhs, offset)
	} else {
		emitSavePrimitive(lhs)
	}
}

// Each left-hand side operand must be addressable,
// a map index expression,
// or (for = assignments only) the blank identifier.
func emitSavePrimitive(lhs Expr) {
	lhs = unwrapRel(lhs)
	switch lhs.(type) {
	case nil:
		return
	case *ExprVariable:
		emitOffsetSavePrimitive(lhs, lhs.getGtype().getSize(), 0)
	case *ExprIndex:
		emitOffsetSavePrimitive(lhs, lhs.getGtype().getSize(), 0)
	case *ExprStructField:
		lhs.(*ExprStructField).emitSavePrimitive()
	case *ExprUop:
		lhs.(*ExprUop).emitSavePrimitive()
	default:
		lhs.dump()
		errorft(lhs.token(), S("Unknown case %T"), lhs)
	}
}

func emitOffsetSavePrimitive(lhs Expr, size int, offset int) {
	lhs = unwrapRel(lhs)
	switch lhs.(type) {
	case *ExprVariable:
		variable := lhs.(*ExprVariable)
		variable.emitOffsetSavePrimitive(size, offset, false)
	case *ExprIndex:
		indexExpr := lhs.(*ExprIndex)
		indexExpr.emitOffsetSavePrimitive(offset)
	case *ExprStructField:
		structfield := lhs.(*ExprStructField)
		fieldType := structfield.getGtype()
		emitOffsetSavePrimitive(structfield.strct, size, fieldType.offset+offset)
	case *ExprUop:
		errorft(lhs.token(), S("unkonwn type %T"), lhs)
	default:
		errorft(lhs.token(), S("unkonwn type %T"), lhs)
	}
}

// e.g. *x = 1, or *x++
func (uop *ExprUop) emitSavePrimitive() {
	emit(S("# *ExprUop.emitSavePrimitive()"))
	assert(eq(uop.op , bytes("*")), uop.tok, S("uop op should be *"))
	emit(S("PUSH_8 # what"))
	uop.operand.emit()
	emit(S("PUSH_8 # where"))
	emit(S("STORE_8_INDIRECT_FROM_STACK"))
}

// x = 1
func (variable *ExprVariable) emitOffsetSavePrimitive(size int, offset int, forceIndirection bool) {
	assert(0 <= size && size <= 8, variable.token(), S("invalid size"))
	if variable.getGtype().getKind() == G_POINTER && (offset > 0 || forceIndirection) {
		assert(variable.getGtype().getKind() == G_POINTER, variable.token(), S(""))
		emit(S("PUSH_8 # what"))
		variable.emit()
		emit(S("ADD_NUMBER %d"), offset)
		emit(S("PUSH_8 # where"))
		emit(S("STORE_8_INDIRECT_FROM_STACK # %s"), bytes(variable.varname))
		return
	}
	if variable.isGlobal {
		emit(S("STORE_%d_TO_GLOBAL %s %d # %s "), size, bytes(variable.varname), offset, bytes(variable.varname))
	} else {
		emit(S("STORE_%d_TO_LOCAL %d+%d # %s"), size, variable.offset, offset, bytes(variable.varname))
	}
}

// save data from stack
func (e *ExprIndex) emitSave24() {
	// load head address of the array
	// load index
	// multi index * size
	// calc address = head address + offset
	// copy value to the address
	emit(S("PUSH_24"))
	collectionType := e.collection.getGtype()
	if collectionType.getKind() == G_MAP {
		e.emitMapSetFromStack24()
		return
	}

	assert(collectionType.getKind() == G_ARRAY || collectionType.getKind() == G_SLICE, e.token(), S("unexpected kind"))
	e.collection.emit()
	emit(S("PUSH_8 # head address of collection"))
	e.index.emit()
	emit(S("PUSH_8 # index"))
	elmType := collectionType.elementType
	size := elmType.getSize()
	assert(size > 0, nil, S("size > 0"))
	emit(S("LOAD_NUMBER %d # elementSize"), size)
	emit(S("PUSH_8"))
	emit(S("IMUL_FROM_STACK # index * elementSize"))
	emit(S("PUSH_8 # index * elementSize"))
	emit(S("SUM_FROM_STACK # (index * size) + address"))
	emit(S("PUSH_8"))
	emit(S("STORE_24_INDIRECT_FROM_STACK"))
}

func (e *ExprIndex) emitOffsetSavePrimitive(offset int) {
	collectionType := e.collection.getGtype()
	switch {
	case collectionType.getKind() == G_ARRAY, collectionType.getKind() == G_SLICE:
		e.emitArrayOrSliceSavePrimitive(offset)
	case collectionType.getKind() == G_MAP:
		emit(S("PUSH_8")) // push RHS value
		e.emitMapSetFromStack8()
		return
	default:
		TBI(e.token(), S("unable to handle %s"), collectionType)
	}
}

func (e *ExprStructField) emitSavePrimitive() {
	fieldType := e.getGtype()
	if e.strct.getGtype().getKind() == G_POINTER {
		emit(S("PUSH_8 # rhs"))

		e.strct.emit()
		emit(S("ADD_NUMBER %d"), fieldType.offset)
		emit(S("PUSH_8"))

		emit(S("STORE_8_INDIRECT_FROM_STACK"))
	} else {
		emitOffsetSavePrimitive(e.strct, 8, fieldType.offset)
	}
}

// take slice values from stack
func emitSave24(lhs Expr, offset int) {
	assertInterface(lhs)
	lhs = unwrapRel(lhs)
	//emit(S("# emitSave24(%T, offset %d)"), lhs, offset)
	emit(S("# emitSave24(?, offset %d)"), offset)
	switch lhs.(type) {
	case *ExprVariable:
		variable := lhs.(*ExprVariable)
		variable.emitSave24(offset)
	case *ExprStructField:
		structfield := lhs.(*ExprStructField)
		fieldType := structfield.getGtype()
		fieldOffset := fieldType.offset
		emit(S("# fieldOffset=%d (%s)"), fieldOffset, bytes(fieldType.fieldname))
		emitSave24(structfield.strct, fieldOffset+offset)
	case *ExprIndex:
		indexExpr := lhs.(*ExprIndex)
		indexExpr.emitSave24()
	default:
		errorft(lhs.token(), S("unkonwn type %T"), lhs)
	}
}

func (variable *ExprVariable) emitSave24(offset int) {
	emit(S("PUSH_24"))
	emit(S("# *ExprVariable.emitSave24()"))
	emit(S("pop %%rax # 3rd"))
	variable.emitOffsetSavePrimitive(8, offset+16, false)
	emit(S("pop %%rax # 2nd"))
	variable.emitOffsetSavePrimitive(8, offset+8, false)
	emit(S("pop %%rax # 1st"))
	variable.emitOffsetSavePrimitive(8, offset+0, true)
}

func (e *ExprIndex) emitArrayOrSliceSavePrimitive(offset int) {
	collection := e.collection
	index := e.index
	collectionType := collection.getGtype()
	assert(collectionType.getKind() == G_ARRAY || collectionType.getKind() == G_SLICE, collection.token(), S("should be collection"))

	elmType := collectionType.elementType
	elmSize := elmType.getSize()
	assert(elmSize > 0, nil, S("elmSize > 0"))

	emit(S("PUSH_8 # rhs"))

	collection.emit()
	emit(S("PUSH_8 # addr"))

	index.emit()
	emit(S("IMUL_NUMBER %d # index * elmSize"), elmSize)
	emit(S("PUSH_8"))

	emit(S("SUM_FROM_STACK # (index * elmSize) + addr"))
	emit(S("ADD_NUMBER %d # offset"), offset)
	emit(S("PUSH_8"))

	if elmSize == 1 {
		emit(S("STORE_1_INDIRECT_FROM_STACK"))
	} else {
		emit(S("STORE_8_INDIRECT_FROM_STACK"))
	}
	emitNewline()
}
