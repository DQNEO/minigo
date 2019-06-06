package main

import "fmt"

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
					switch left.getGtype().getKind() {
					case G_SLICE:
						emit("POP_24")
						emitSave24(left, 0)
					case G_INTERFACE:
						emit("POP_24")
						emitSave24(left, 0)
					default:
						emit("pop %%rax")
						emitSavePrimitive(left)
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
			okRegister := mapOkRegister(right.getGtype().is24WidthType())
			emit("mov %%%s, %%rax # emit okValue", okRegister)
			emitSavePrimitive(okVariable)
		}
		return
	}

}

func emitAssignPrimitive(left Expr, right Expr) {
	assert(left.getGtype().getSize() <= 8, left.token(), fmt.Sprintf("invalid type for lhs: %s", left.getGtype()))
	assert(right != nil || right.getGtype().getSize() <= 8, right.token(), fmt.Sprintf("invalid type for rhs: %s", right.getGtype()))
	right.emit()            //   expr => %rax
	emitSavePrimitive(left) //   %rax => memory
}

func assignToStruct(lhs Expr, rhs Expr) {
	emit("# assignToStruct start")

	if rel, ok := lhs.(*Relation); ok {
		lhs = rel.expr
	}
	assert(rhs == nil || (rhs.getGtype().getKind() == G_STRUCT),
		lhs.token(), "rhs should be struct type")
	// initializes with zero values
	emit("# initialize struct with zero values: start")
	for _, fieldtype := range lhs.getGtype().relation.gtype.fields {
		switch fieldtype.getKind() {
		case G_ARRAY:
			arrayType := fieldtype
			elementType := arrayType.elementType
			elmSize := arrayType.elementType.getSize()
			switch {
			case elementType.getKind() == G_STRUCT:
				left := &ExprStructField{
					strct:     lhs,
					fieldname: fieldtype.fieldname,
				}
				assignToArray(left, nil)
			default:
				assert(0 <= elmSize && elmSize <= 8, lhs.token(), "invalid size")
				for i := 0; i < arrayType.length; i++ {
					emit("mov $0, %%rax")
					emitOffsetSavePrimitive(lhs, elmSize, fieldtype.offset+i*elmSize)
				}
			}

		case G_SLICE:
			emit("LOAD_EMPTY_SLICE")
			emitSave24(lhs, fieldtype.offset)
		case G_MAP:
			emit("LOAD_EMPTY_MAP")
			emitSave24(lhs, fieldtype.offset)
		case G_STRUCT:
			left := &ExprStructField{
				strct:     lhs,
				fieldname: fieldtype.fieldname,
			}
			assignToStruct(left, nil)
		case G_INTERFACE:
			emit("LOAD_EMPTY_INTERFACE")
			emitSave24(lhs, fieldtype.offset)
		default:
			emit("mov $0, %%rax")
			regSize := fieldtype.getSize()
			assert(0 < regSize && regSize <= 8, lhs.token(), fieldtype.String())
			emitOffsetSavePrimitive(lhs, regSize, fieldtype.offset)
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
		emit("PUSH_8")
		emitAddress(rhs)
		emit("PUSH_8")
		emitCopyStructFromStack(lhs.getGtype().getSize())
	case *ExprUop:
		re := rhs.(*ExprUop)
		if re.op == "*" {
			// copy struct
			emitAddress(lhs)
			emit("PUSH_8")
			re.operand.emit()
			emit("PUSH_8")
			emitCopyStructFromStack(lhs.getGtype().getSize())
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

			switch fieldtype.getKind() {
			case G_ARRAY:
				initvalues, ok := field.value.(*ExprArrayLiteral)
				assert(ok, nil, "ok")
				arrayType := strcttyp.getField(field.key)
				elementType := arrayType.elementType
				elmSize := elementType.getSize()
				switch {
				case elementType.getKind() == G_STRUCT:
					left := &ExprStructField{
						strct:     lhs,
						fieldname: fieldtype.fieldname,
					}
					assignToArray(left, field.value)
				default:
					for i, val := range initvalues.values {
						val.emit()
						emitOffsetSavePrimitive(variable, elmSize, arrayType.offset+i*elmSize)
					}
				}
			case G_SLICE:
				left := &ExprStructField{
					tok:       variable.token(),
					strct:     lhs,
					fieldname: field.key,
				}
				assignToSlice(left, field.value)
			case G_MAP:
				left := &ExprStructField{
					tok:       variable.token(),
					strct:     lhs,
					fieldname: field.key,
				}
				assignToMap(left, field.value)
			case G_INTERFACE:
				left := &ExprStructField{
					tok:       lhs.token(),
					strct:     lhs,
					fieldname: field.key,
				}
				assignToInterface(left, field.value)
			case G_STRUCT:
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
				emitOffsetSavePrimitive(variable, regSize, fieldtype.offset)
			}
		}
	default:
		TBI(rhs.token(), "")
	}

	emit("# assignToStruct end")
}

func assignToMap(lhs Expr, rhs Expr) {
	emit("# assignToMap")
	if rhs == nil {
		emit("# initialize map with a zero value")
		emit("LOAD_EMPTY_MAP")
		emitSave24(lhs, 0)
		return
	}
	switch rhs.(type) {
	case *ExprMapLiteral:
		emit("# map literal")

		lit := rhs.(*ExprMapLiteral)
		lit.emit()
	case *Relation, *ExprVariable, *ExprIndex, *ExprStructField, *ExprFuncallOrConversion, *ExprMethodcall:
		rhs.emit()
	default:
		TBI(rhs.token(), "unable to handle %T", rhs)
	}
	emitSave24(lhs, 0)
}


func assignToInterface(lhs Expr, rhs Expr) {
	emit("# assignToInterface")
	if rhs == nil || isNil(rhs) {
		emit("LOAD_EMPTY_INTERFACE")
		emitSave24(lhs, 0)
		return
	}

	assert(rhs.getGtype() != nil, rhs.token(), fmt.Sprintf("rhs gtype is nil:%T", rhs))
	if rhs.getGtype().getKind() == G_INTERFACE {
		rhs.emit()
		emitSave24(lhs, 0)
		return
	}

	emitConversionToInterface(rhs)
	emitSave24(lhs, 0)
}

func assignToSlice(lhs Expr, rhs Expr) {
	emit("# assignToSlice")
	assertInterface(lhs)
	//assert(rhs == nil || rhs.getGtype().kind == G_SLICE, nil, "should be a slice literal or nil")
	if rhs == nil {
		emit("LOAD_EMPTY_SLICE")
		emitSave24(lhs, 0)
		return
	}

	//	assert(rhs.getGtype().getKind() == G_SLICE, rhs.token(), "rsh should be slice type")

	switch rhs.(type) {
	case *Relation:
		rel := rhs.(*Relation)
		if _, ok := rel.expr.(*ExprNilLiteral); ok {
			emit("LOAD_EMPTY_SLICE")
			emitSave24(lhs, 0)
			return
		}
		rvariable, ok := rel.expr.(*ExprVariable)
		assert(ok, nil, "ok")
		rvariable.emit()
	case *ExprSliceLiteral:
		lit := rhs.(*ExprSliceLiteral)
		lit.emit()
	case *ExprSlice:
		e := rhs.(*ExprSlice)
		e.emit()
	case *ExprConversion:
		// https://golang.org/ref/spec#Conversions
		// Converting a value of a string type to a slice of bytes type
		// yields a slice whose successive elements are the bytes of the string.
		//
		// see also https://blog.golang.org/strings
		conversion := rhs.(*ExprConversion)
		assert(conversion.gtype.getKind() == G_SLICE, rhs.token(), "must be a slice of bytes")
		assert(conversion.expr.getGtype().getKind() == G_STRING, rhs.token(), "must be a string type, but got "+conversion.expr.getGtype().String())
		stringVarname, ok := conversion.expr.(*Relation)
		assert(ok, rhs.token(), "ok")
		stringVariable := stringVarname.expr.(*ExprVariable)
		stringVariable.emit()
		emit("PUSH_8 # ptr")
		strlen := &ExprLen{
			arg: stringVariable,
		}
		strlen.emit()
		emit("PUSH_8 # len")
		emit("PUSH_8 # cap")
		emit("POP_SLICE")
	default:
		//emit("# emit rhs of type %T %s", rhs, rhs.getGtype().String())
		rhs.emit() // it should put values to rax,rbx,rcx
	}

	emitSave24(lhs, 0)
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
	assert(rhs == nil || rhs.getGtype().getKind() == G_ARRAY, nil, "rhs should be array")
	switch elementType.getKind() {
	case G_STRUCT:
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
						emitSave24(lhs, offsetByIndex)
						continue
					} else if arrayLiteral.values[i].getGtype().getKind() != G_INTERFACE {
						// conversion of dynamic type => interface type
						dynamicValue := arrayLiteral.values[i]
						emitConversionToInterface(dynamicValue)
						emit("LOAD_EMPTY_INTERFACE")
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

			emitOffsetSavePrimitive(lhs, elmSize, offsetByIndex)
		}
	}
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
