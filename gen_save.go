package main

func emitPop(gtype *Gtype) {
	if gtype.is24WidthType() {
		emit("POP_24")
	} else {
		emit("POP_8")
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
		errorft(lhs.token(), "Unknown case %T", lhs)
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
		errorft(lhs.token(), "unkonwn type %T", lhs)
	default:
		errorft(lhs.token(), "unkonwn type %T", lhs)
	}
}

// e.g. *x = 1, or *x++
func (uop *ExprUop) emitSavePrimitive() {
	emit("# *ExprUop.emitSavePrimitive()")
	assert(eqGostrings(uop.op , gostring("*")), uop.tok, "uop op should be *")
	emit("PUSH_8 # what")
	uop.operand.emit()
	emit("PUSH_8 # where")
	emit("STORE_8_INDIRECT_FROM_STACK")
}

// x = 1
func (variable *ExprVariable) emitOffsetSavePrimitive(size int, offset int, forceIndirection bool) {
	assert(0 <= size && size <= 8, variable.token(), "invalid size")
	if variable.getGtype().getKind() == G_POINTER && (offset > 0 || forceIndirection) {
		assert(variable.getGtype().getKind() == G_POINTER, variable.token(), "")
		emit("PUSH_8 # what")
		variable.emit()
		emit("ADD_NUMBER %d", offset)
		emit("PUSH_8 # where")
		emit("STORE_8_INDIRECT_FROM_STACK # %s", gostring(variable.varname))
		return
	}
	if variable.isGlobal {
		emit("STORE_%d_TO_GLOBAL %s %d # %s ", size, gostring(variable.varname), offset, gostring(variable.varname))
	} else {
		emit("STORE_%d_TO_LOCAL %d+%d # %s", size, variable.offset, offset, gostring(variable.varname))
	}
}

// save data from stack
func (e *ExprIndex) emitSave24() {
	// load head address of the array
	// load index
	// multi index * size
	// calc address = head address + offset
	// copy value to the address
	emit("PUSH_24")
	collectionType := e.collection.getGtype()
	if collectionType.getKind() == G_MAP {
		e.emitMapSet(true)
		return
	}

	assert(collectionType.getKind() == G_ARRAY || collectionType.getKind() == G_SLICE || collectionType.getKind() == G_STRING, e.token(), "unexpected kind")
	e.collection.emit()
	emit("PUSH_8 # head address of collection")
	e.index.emit()
	emit("PUSH_8 # index")
	var elmType *Gtype
	if collectionType.isString() {
		elmType = gByte
	} else {
		elmType = collectionType.elementType
	}
	size := elmType.getSize()
	assert(size > 0, nil, "size > 0")
	emit("LOAD_NUMBER %d # elementSize", size)
	emit("PUSH_8")
	emit("IMUL_FROM_STACK # index * elementSize")
	emit("PUSH_8 # index * elementSize")
	emit("SUM_FROM_STACK # (index * size) + address")
	emit("PUSH_8")
	emit("STORE_24_INDIRECT_FROM_STACK")
}

func (e *ExprIndex) emitOffsetSavePrimitive(offset int) {
	collectionType := e.collection.getGtype()
	switch {
	case collectionType.getKind() == G_ARRAY, collectionType.getKind() == G_SLICE, collectionType.getKind() == G_STRING:
		e.emitArrayOrSliceSavePrimitive(offset)
	case collectionType.getKind() == G_MAP:
		emit("PUSH_8") // push RHS value
		e.emitMapSet(false)
		return
	default:
		TBI(e.token(), "unable to handle %s", collectionType)
	}
}

func (e *ExprStructField) emitSavePrimitive() {
	fieldType := e.getGtype()
	if e.strct.getGtype().getKind() == G_POINTER {
		emit("PUSH_8 # rhs")

		e.strct.emit()
		emit("ADD_NUMBER %d", fieldType.offset)
		emit("PUSH_8")

		emit("STORE_8_INDIRECT_FROM_STACK")
	} else {
		emitOffsetSavePrimitive(e.strct, 8, fieldType.offset)
	}
}

// take slice values from stack
func emitSave24(lhs Expr, offset int) {
	assertInterface(lhs)
	lhs = unwrapRel(lhs)
	//emit("# emitSave24(%T, offset %d)", lhs, offset)
	emit("# emitSave24(?, offset %d)", offset)
	switch lhs.(type) {
	case *ExprVariable:
		variable := lhs.(*ExprVariable)
		variable.emitSave24(offset)
	case *ExprStructField:
		structfield := lhs.(*ExprStructField)
		fieldType := structfield.getGtype()
		fieldOffset := fieldType.offset
		emit("# fieldOffset=%d (%s)", fieldOffset, gostring(fieldType.fieldname))
		emitSave24(structfield.strct, fieldOffset+offset)
	case *ExprIndex:
		indexExpr := lhs.(*ExprIndex)
		indexExpr.emitSave24()
	default:
		errorft(lhs.token(), "unkonwn type %T", lhs)
	}
}

func (variable *ExprVariable) emitSave24(offset int) {
	emit("PUSH_24")
	emit("# *ExprVariable.emitSave24()")
	emit("pop %%rax # 3rd")
	variable.emitOffsetSavePrimitive(8, offset+16, false)
	emit("pop %%rax # 2nd")
	variable.emitOffsetSavePrimitive(8, offset+8, false)
	emit("pop %%rax # 1st")
	variable.emitOffsetSavePrimitive(8, offset+0, true)
}

func (e *ExprIndex) emitArrayOrSliceSavePrimitive(offset int) {
	collection := e.collection
	index := e.index
	collectionType := collection.getGtype()
	assert(collectionType.getKind() == G_ARRAY || collectionType.getKind() == G_SLICE || collectionType.getKind() == G_STRING, collection.token(), "should be collection")

	var elmType *Gtype
	if collectionType.isString() {
		elmType = gByte
	} else {
		elmType = collectionType.elementType
	}
	elmSize := elmType.getSize()
	assert(elmSize > 0, nil, "elmSize > 0")

	emit("PUSH_8 # rhs")

	collection.emit()
	emit("PUSH_8 # addr")

	index.emit()
	emit("IMUL_NUMBER %d # index * elmSize", elmSize)
	emit("PUSH_8")

	emit("SUM_FROM_STACK # (index * elmSize) + addr")
	emit("ADD_NUMBER %d # offset", offset)
	emit("PUSH_8")

	if elmSize == 1 {
		emit("STORE_1_INDIRECT_FROM_STACK")
	} else {
		emit("STORE_8_INDIRECT_FROM_STACK")
	}
	emitNewline()
}
