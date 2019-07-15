package main

import "fmt"

const globalPrefix string = "_gbl_"

func (vr *ExprVariable) globalSymbol() string {
	assert(vr.isGlobal, vr.token(), " Not global var")
	return fmt.Sprintf("%s.%s", globalPrefix, string(vr.varname))
}

// gloabal var which should be initialized with zeros
// https://en.wikipedia.org/wiki/.bss
func (decl *DeclVar) emitBss() {
	emit(".data")
	// https://sourceware.org/binutils/docs-2.30/as/Lcomm.html#Lcomm
	emit(".lcomm %s, %d", decl.variable.globalSymbol(), decl.variable.getGtype().getSize())
}

func (decl *DeclVar) emitData() {
	ptok := decl.token()
	gtype := decl.variable.gtype
	right := decl.initval

	emitWithoutIndent("%s: # gtype=%s", decl.variable.globalSymbol(), gtype.String())
	emitWithoutIndent("# right.gtype = %s", right.getGtype().String())
	emitWithoutIndent(".data 0")
	doEmitData(ptok, right.getGtype(), right, "", 0)
}

func (e *ExprStructLiteral) lookup(fieldname identifier) Expr {
	for _, field := range e.fields {
		if string(field.key)  ==  string(fieldname) {
			return field.value
		}
	}

	return nil
}

func doEmitData(ptok *Token /* left type */, gtype *Gtype, value /* nullable */ Expr, containerName string, depth int) {
	value = unwrapRel(value)
	emit("# doEmitData: containerName=%s, depth=%d", containerName, depth)
	primType := gtype.getKind()
	if primType == G_ARRAY {
		arrayliteral, ok := value.(*ExprArrayLiteral)
		var values []Expr
		if ok {
			values = arrayliteral.values
		}
		assert(ok || arrayliteral == nil, ptok, "*ExprArrayLiteral expected, but got ")
		elmType := gtype.elementType
		assertNotNil(elmType != nil, nil)
		for i := 0; i < gtype.length; i++ {
			var selector string
			if i >= len(values) {
				// zero value
				doEmitData(ptok, elmType, nil, "", depth)
			} else {
				value := arrayliteral.values[i]
				assertNotNil(value != nil, nil)
				size := elmType.getSize()
				if size == 8 {
					switch value.(type) {
					case *ExprUop:
						uop := value.(*ExprUop)
						operand := unwrapRel(uop.operand)
						vr, ok := operand.(*ExprVariable)
						assert(ok, uop.token(), "only variable is allowed")
						emit(".quad %s # %s %s", vr.globalSymbol(), value.getGtype().String(), selector)
					case *ExprVariable:
						assert(false, value.token(), "variable here is not allowed")
					default:
						emit(".quad %d # %s %s", evalIntExpr(value), value.getGtype().String(), selector)
					}
				} else if size == 1 {
					emit(".byte %d", evalIntExpr(value))
				} else {
					doEmitData(ptok, gtype.elementType, value, selector, depth)
				}
			}
		}
		emit(".quad 0 # nil terminator")

	} else if primType == G_SLICE || primType == G_STRING {
		switch value.(type) {
		case nil:
			emit(".quad 0")
			emit(".quad 0")
			emit(".quad 0")
		case *ExprNilLiteral:
			emit(".quad 0")
			emit(".quad 0")
			emit(".quad 0")
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
		case *ExprStringLiteral:
			stringLiteral := value.(*ExprStringLiteral)
			emit(".quad .%s", stringLiteral.slabel)
			var length int = len(stringLiteral.val)
			emit(".quad %d", length)
			emit(".quad %d", length)
		case *ExprFuncallOrConversion:
			call := value.(*ExprFuncallOrConversion)
			assert(call.rel.gtype != nil, value.token(), "should be Conversion")
			stringLiteral,ok := call.args[0].(*ExprStringLiteral)
			assert(ok, call.token(), "arg0 should be stringliteral")
			emit(".quad .%s", stringLiteral.slabel)
			var length int = len(stringLiteral.val)
			emit(".quad %d", length)
			emit(".quad %d", length)
		default:
			TBI(ptok, "unable to handle gtype %s", gtype.String())
		}
	} else if primType == G_INTERFACE {
		emit(".quad 0")
		emit(".quad 0")
		emit(".quad 0")
	} else if primType == G_BOOL {
		if value == nil {
			// zero value
			emit(".quad 0 # %s %s",  gtype.String(), containerName)
			return
		}
		var val int = evalIntExpr(value)
		emit(".quad %d # %s %s", val, gtype.String(), containerName)
	} else if primType == G_STRUCT {
		s := string(containerName) + "."  + string(gtype.relation.name)
		containerName = s
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
				s2 :=  string(containerName) + "." + string(field.fieldname)
				doEmitData(ptok, field, nil,s2, depth)
				continue
			}
			structLiteral, ok := value.(*ExprStructLiteral)
			assert(ok, nil, "ok")
			value := structLiteral.lookup(field.fieldname)
			if value == nil {
				// zero value
				//continue
			}
			gtype := field
			s3 :=  string(containerName) + "." + string(field.fieldname)
			doEmitData(ptok, gtype, value, s3, depth)
		}
	} else {
		var val int
		var gtypeString string = gtype.String()
		switch value.(type) {
		case nil:
			emit(".quad %d # %s %s zero value", val, gtypeString, containerName)
		case *ExprNumberLiteral:
			val = value.(*ExprNumberLiteral).val
			emit(".quad %d # %s %s", val, gtypeString, containerName)
		case *ExprConstVariable:
			cnst := value.(*ExprConstVariable)
			val = evalIntExpr(cnst)
			emit(".quad %d # %s ", val, gtypeString)
		case *ExprVariable:
			vr := value.(*ExprVariable)
			val = evalIntExpr(vr)
			emit(".quad %d # %s ", val, gtypeString)
		case *ExprBinop:
			val = evalIntExpr(value)
			emit(".quad %d # %s ", val, gtypeString)
		case *ExprUop:
			uop := value.(*ExprUop)
			assert(uop.op ==  "&", ptok, "only uop & is allowed")
			operand := unwrapRel(uop.operand)
			vr, ok := operand.(*ExprVariable)
			if ok {
				assert(vr.isGlobal, value.token(), "operand should be a global variable")
				emit(".quad %s", vr.globalSymbol())
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
	emitWithoutIndent("# emitGlobal for %s", decl.variable.globalSymbol())
	assertNotNil(decl.variable.gtype != nil, nil)

	if decl.initval == nil {
		decl.emitBss()
	} else {
		decl.emitData()
	}
}
