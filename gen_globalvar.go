package main

// gloabal var which should be initialized with zeros
// https://en.wikipedia.org/wiki/.bss
func (decl *DeclVar) emitBss() {
	emit(S(".data"))
	// https://sourceware.org/binutils/docs-2.30/as/Lcomm.html#Lcomm
	emit(S(".lcomm %s, %d"), gostring(decl.variable.varname), decl.variable.getGtype().getSize())
}

func (decl *DeclVar) emitData() {
	ptok := decl.token()
	gtype := decl.variable.gtype
	right := decl.initval

	emitWithoutIndent(S("%s: # gtype=%s"), gostring(decl.variable.varname), gtype.String())
	emitWithoutIndent(S("# right.gtype = %s"), right.getGtype().String())
	emitWithoutIndent(S(".data 0"))
	doEmitData(ptok, right.getGtype(), right, S(""), 0)
}

func (e *ExprStructLiteral) lookup(fieldname goidentifier) Expr {
	for _, field := range e.fields {
		if eq(gostring(field.key) , gostring(fieldname)) {
			return field.value
		}
	}

	return nil
}

func doEmitData(ptok *Token /* left type */, gtype *Gtype, value /* nullable */ Expr, containerName gostring, depth int) {
	value = unwrapRel(value)
	emit(S("# doEmitData: containerName=%s, depth=%d"), gostring(containerName), depth)
	primType := gtype.getKind()
	if primType == G_ARRAY {
		arrayliteral, ok := value.(*ExprArrayLiteral)
		var values []Expr
		if ok {
			values = arrayliteral.values
		}
		assert(ok || arrayliteral == nil, ptok, S("*ExprArrayLiteral expected, but got "))
		elmType := gtype.elementType
		assertNotNil(elmType != nil, nil)
		for i := 0; i < gtype.length; i++ {
			var selector gostring
			if i >= len(values) {
				// zero value
				doEmitData(ptok, elmType, nil, S(""), depth)
			} else {
				value := arrayliteral.values[i]
				assertNotNil(value != nil, nil)
				size := elmType.getSize()
				if size == 8 {
					if value.getGtype().getKind() == G_CLIKE_STRING {
						stringLiteral, ok := value.(*ExprStringLiteral)
						assert(ok, nil, S("ok"))
						emit(S(".quad .%s"), gostring(stringLiteral.slabel))
					} else {
						switch value.(type) {
						case *ExprUop:
							uop := value.(*ExprUop)
							operand := unwrapRel(uop.operand)
							vr, ok := operand.(*ExprVariable)
							assert(ok, uop.token(), S("only variable is allowed"))
							emit(S(".quad %s # %s %s"), gostring(vr.varname), value.getGtype().String(), gostring(selector))
						case *ExprVariable:
							assert(false, value.token(), S("variable here is not allowed"))
						default:
							emit(S(".quad %d # %s %s"), evalIntExpr(value), value.getGtype().String(), gostring(selector))
						}
					}
				} else if size == 1 {
					emit(S(".byte %d"), evalIntExpr(value))
				} else {
					doEmitData(ptok, gtype.elementType, value, selector, depth)
				}
			}
		}
		emit(S(".quad 0 # nil terminator"))

	} else if primType == G_SLICE {
		switch value.(type) {
		case nil:
			emit(S(".quad 0"))
			emit(S(".quad 0"))
			emit(S(".quad 0"))
		case *ExprNilLiteral:
			emit(S(".quad 0"))
			emit(S(".quad 0"))
			emit(S(".quad 0"))
		case *ExprSliceLiteral:
			// initialize a hidden array
			lit := value.(*ExprSliceLiteral)
			arrayLiteral := &ExprArrayLiteral{
				gtype:  lit.invisiblevar.gtype,
				values: lit.values,
			}

			emitDataAddr(arrayLiteral, depth)               // emit underlying array
			emit(S(".quad %d"), lit.invisiblevar.gtype.length) // len
			emit(S(".quad %d"), lit.invisiblevar.gtype.length) // cap
		case *ExprFuncallOrConversion:
			call := value.(*ExprFuncallOrConversion)
			assert(call.rel.gtype != nil, value.token(), S("should be Conversion"))
			toGtype := call.rel.gtype
			assert(toGtype.getKind() == G_SLICE && call.args[0].getGtype().isString(), call.token(), S("should be string to slice conversion"))
			stringLiteral,ok := call.args[0].(*ExprStringLiteral)
			assert(ok, call.token(), S("arg0 should be stringliteral"))
			emit(S(".quad .%s"), stringLiteral.slabel)
			var length int = len(stringLiteral.val)
			emit(S(".quad %d"), length)
			emit(S(".quad %d"), length)
		default:
			TBI(ptok, S("unable to handle gtype %s"), gtype.String())
		}
	} else if primType == G_INTERFACE {
		emit(S(".quad 0"))
		emit(S(".quad 0"))
		emit(S(".quad 0"))
	} else if primType == G_BOOL {
		if value == nil {
			// zero value
			emit(S(".quad 0 # %s %s"),  gtype.String(), gostring(containerName))
			return
		}
		var val int = evalIntExpr(value)
		emit(S(".quad %d # %s %s"), val, gtype.String(), gostring(containerName))
	} else if primType == G_STRUCT {
		containerName = concat3(containerName, S(".") , gostring(gtype.relation.name))
		for _, field := range gtype.relation.gtype.fields {
			emit(S("# padding=%d"), field.padding)
			switch field.padding {
			case 1:
				emit(S(".byte 0 # padding"))
			case 4:
				emit(S(".double 0 # padding"))
			case 8:
				emit(S(".quad 0 # padding"))
			default:
			}
			emit(S("# field:offesr=%d, fieldname=%s"), field.offset, gostring(field.fieldname))
			if value == nil {
				doEmitData(ptok, field, nil, concat3(containerName,S("."), gostring(field.fieldname)), depth)
				continue
			}
			structLiteral, ok := value.(*ExprStructLiteral)
			assert(ok, nil, S("ok"))
			value := structLiteral.lookup(field.fieldname)
			if value == nil {
				// zero value
				//continue
			}
			gtype := field
			doEmitData(ptok, gtype, value, concat3(containerName, S("."), gostring(field.fieldname)), depth)
		}
	} else {
		var val int
		var gtypeString gostring = gtype.String()
		switch value.(type) {
		case nil:
			emit(S(".quad %d # %s %s zero value"), val, gtypeString, gostring(containerName))
		case *ExprNumberLiteral:
			val = value.(*ExprNumberLiteral).val
			emit(S(".quad %d # %s %s"), val, gtypeString, gostring(containerName))
		case *ExprConstVariable:
			cnst := value.(*ExprConstVariable)
			val = evalIntExpr(cnst)
			emit(S(".quad %d # %s "), val, gtypeString)
		case *ExprVariable:
			vr := value.(*ExprVariable)
			val = evalIntExpr(vr)
			emit(S(".quad %d # %s "), val, gtypeString)
		case *ExprBinop:
			val = evalIntExpr(value)
			emit(S(".quad %d # %s "), val, gtypeString)
		case *ExprStringLiteral:
			stringLiteral := value.(*ExprStringLiteral)
			emit(S(".quad .%s"), stringLiteral.slabel)
		case *ExprUop:
			uop := value.(*ExprUop)
			assert(eq(uop.op, gostring("&")), ptok, S("only uop & is allowed"))
			operand := unwrapRel(uop.operand)
			vr, ok := operand.(*ExprVariable)
			if ok {
				assert(vr.isGlobal, value.token(), S("operand should be a global variable"))
				emit(S(".quad %s"), gostring(vr.varname))
			} else {
				// var gv = &Struct{_}
				emitDataAddr(operand, depth)
			}
		default:
			TBI(ptok, S("unable to handle %d"), primType)
		}
	}
}

// this logic is stolen from 8cc.
func emitDataAddr(operand Expr, depth int) {
	emit(S(".data %d"), depth+1)
	label := makeLabel()
	emit(S("%s:"), label)
	doEmitData(nil, operand.getGtype(), operand, S(""), depth+1)
	emit(S(".data %d"), depth)
	emit(S(".quad %s"), label)
}

func (decl *DeclVar) emitGlobal() {
	emitWithoutIndent(S("# emitGlobal for %s"), gostring(decl.variable.varname))
	assertNotNil(decl.variable.gtype != nil, nil)

	if decl.initval == nil {
		decl.emitBss()
	} else {
		decl.emitData()
	}
}
