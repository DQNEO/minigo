package main

import "fmt"

type GTYPE_TYPE int

const undefinedSize = -1

const (
	G_UNKOWNE   GTYPE_TYPE = iota
	G_DEPENDENT            // depends on other expression
	G_REL
	// below are primitives which are declared in the universe block
	G_INT
	G_BOOL
	G_BYTE
	// end of primitives
	G_STRUCT
	G_STRUCT_FIELD
	G_ARRAY
	G_SLICE
	G_STRING
	G_MAP
	G_POINTER
	G_FUNC
	G_INTERFACE
)

type signature struct {
	fname      identifier
	paramTypes []*Gtype
	isVariadic bool
	rettypes   []*Gtype
}

type Gtype struct {
	typ            GTYPE_TYPE
	receiverTypeId int                     // for receiverTypeId. 0:unkonwn
	dependendson   Expr                    // for G_DEPENDENT
	relation       *Relation               // for G_REL
	size           int                     // for scalar type like int, bool, byte, for struct
	origType       *Gtype                  // for pointer
	fields         []*Gtype                // for struct
	fieldname      identifier              // for struct field
	offset         int                     // for struct field
	padding        int                     // for struct field
	length       int                       // for array, string(len without the terminating \0)
	elementType  *Gtype                    // for array, slice
	imethods     map[identifier]*signature // for interface
	methods      map[identifier]*ExprFuncRef
	// for fixed array
	mapKey   *Gtype // for map
	mapValue *Gtype // for map
}

func (gtype *Gtype) isNil() bool {
	if gtype == nil {
		return true
	}
	if gtype.typ == G_REL {
		return gtype.relation.gtype == nil

	}
	return false
}

func (gtype *Gtype) getSource() *Gtype {
	if gtype.typ == G_REL {
		return gtype.relation.gtype.getSource()

	}
	return gtype
}

func (gtype *Gtype) getPrimType() GTYPE_TYPE {
	if gtype == nil {
		return G_UNKOWNE
	}
	return gtype.getSource().typ
}

func (gtype *Gtype) is24Width() bool {
	switch gtype.getPrimType() {
	case G_INTERFACE, G_MAP, G_SLICE:
		return true
	default:
		return false
	}
}

func (gtype *Gtype) isString() bool {
	if gtype.getPrimType() == G_STRING {
		return true
	}
	return false
}

func (gtype *Gtype) getSize() int {
	assertNotNil(gtype != nil, nil)
	assert(gtype.typ != G_DEPENDENT, nil, "type should be inferred")
	if gtype.typ == G_REL {
		if gtype.relation.gtype == nil {
			errorf("relation not resolved: %s", gtype)
		}
		return gtype.relation.gtype.getSize()
	} else {
		if gtype.typ == G_ARRAY {
			assertNotNil(gtype.elementType != nil, nil)
			return gtype.length * gtype.elementType.getSize()
		} else if gtype.typ == G_STRUCT {
			// @TODO consider the case of real zero e.g. struct{}
			if gtype.size == undefinedSize {
				gtype.calcStructOffset()
			}
			return gtype.size
		} else if gtype.typ == G_POINTER || gtype.typ == G_STRING {
			return ptrSize
		} else if gtype.typ == G_INTERFACE {
			//     data    ,  receiverTypeId, dtype
			return ptrSize + ptrSize + ptrSize
		} else if gtype.typ == G_SLICE {
			return ptrSize + IntSize + IntSize
		} else if gtype.typ == G_MAP {
			return ptrSize + IntSize + IntSize
		} else {
			return gtype.size
		}
	}
}

func (gtype *Gtype) String() string {
	if gtype == nil {
		return "NO_TYPE"
	}
	switch gtype.typ {
	case G_DEPENDENT:
		return "dependent"
	case G_REL:
		if gtype.relation.pkg == "" {
			//errorf("pkg is empty: %s", gtype.relation.name)
		}
		return fmt.Sprintf("G_REL(%s.%s)",
			gtype.relation.pkg, gtype.relation.name)
	case G_INT:
		return "int"
	case G_BOOL:
		return "bool"
	case G_BYTE:
		return "byte"
	case G_ARRAY:
		elm := gtype.elementType
		return fmt.Sprintf("[%d]%s", gtype.length, elm.String())
	case G_STRUCT:
		var r = "struct{"
		for _, field := range gtype.fields {
			r += field.String() + ","
		}
		r += "}"
		return r
	case G_STRUCT_FIELD:
		return "structfield"
	case G_POINTER:
		origType := gtype.origType
		return fmt.Sprintf("*%s", origType.String())
	case G_SLICE:
		return fmt.Sprintf("[]%s", gtype.elementType.String())
	case G_STRING:
		return "string"
	case G_FUNC:
		return "func"
	case G_INTERFACE:
		if len(gtype.imethods) == 0 {
			return "interface{}"
		} else {
			return fmt.Sprintf("interface {...}")
		}
	case G_MAP:
		return "map"
	default:
		errorf("gtype.String() error: invalid gtype.type=%d", gtype.typ)
	}
	return ""
}

func (strct *Gtype) getField(name identifier) *Gtype {
	assertNotNil(strct != nil, nil)
	assert(strct.typ == G_STRUCT, nil, "assume G_STRUCT type")
	for _, field := range strct.fields {
		if field.fieldname == name {
			return field
		}
	}
	errorf("field %s not found in the struct", name)
	return nil
}

func (strct *Gtype) calcStructOffset() {
	assert(strct.typ == G_STRUCT, nil, "assume G_STRUCT type, but got "+strct.String())
	var offset int
	for _, fieldtype := range strct.fields {
		var align int
		if fieldtype.getSize() < MaxAlign {
			align = fieldtype.getSize()
			assert(align > 0, nil, "field size should be > 0: filed="+fieldtype.String())
		} else {
			align = MaxAlign
		}
		if offset%align != 0 {
			padding := align - offset%align
			fieldtype.padding = padding
			offset += padding
		}
		fieldtype.offset = offset
		offset += fieldtype.getSize()
	}

	strct.size = offset
}

func (rel *Relation) getGtype() *Gtype {
	if rel.expr == nil {
		//errorft(rel.token(), "rel.expr is nil for %s", rel)
		return nil
	}
	return rel.expr.getGtype()
}

func (e *ExprStructLiteral) getGtype() *Gtype {
	return &Gtype{
		typ:      G_REL,
		relation: e.strctname,
	}
}

func (e *ExprFuncallOrConversion) getGtype() *Gtype {
	assert(e.rel.expr != nil || e.rel.gtype != nil, e.token(), "")
	if e.rel.expr != nil {
		funcref, ok := e.rel.expr.(*ExprFuncRef)
		assert(ok, e.token(), "it should be a ExprFuncRef")
		firstRetType := funcref.funcdef.rettypes[0]
		return firstRetType
	} else if e.rel.gtype != nil {
		return e.rel.gtype
	}
	errorf("should not reach here")
	return nil
}

func (e *ExprMethodcall) getGtype() *Gtype {
	gtype := e.receiver.getGtype()
	if gtype.typ == G_POINTER {
		gtype = gtype.origType
	}

	// refetch gtype from the package block scope
	// I forgot the reason to do this :(
	_, ok := allScopes[gtype.relation.pkg]
	if !ok {
		errorft(e.token(), "ExprMethodcall.getGtype(): socope \"%s\" does not exist in allScopes ", gtype.relation.pkg)
	}
	pgtype := allScopes[gtype.relation.pkg].getGtype(gtype.relation.name)
	if pgtype == nil {
		errorft(e.token(), "%s is not found in the scope", gtype)
	}
	if pgtype.typ == G_INTERFACE {
		methodsig, ok := pgtype.imethods[e.fname]
		if !ok {
			errorft(e.token(), "method %s not found in %s %s", e.fname, gtype, e.tok)
		}
		assertNotNil(methodsig != nil, e.tok)
		return methodsig.rettypes[0]
	} else {
		method, ok := pgtype.methods[e.fname]
		if !ok {
			errorft(e.token(), "method %s not found in %s %s", e.fname, gtype, e.tok)
		}
		assertNotNil(method != nil, e.tok)
		return method.funcdef.rettypes[0]
	}
}

func (e *ExprUop) getGtype() *Gtype {
	switch e.op {
	case "&":
		return &Gtype{
			typ:      G_POINTER,
			origType: e.operand.getGtype(),
		}
	case "*":
		return e.operand.getGtype().origType
	case "!":
		return gBool
	case "-":
		return gInt
	}
	errorf("internal error")
	return nil
}

func (f *ExprFuncRef) getGtype() *Gtype {
	return &Gtype{
		typ: G_FUNC,
	}
}

func (e *ExprSlice) getGtype() *Gtype {
	return &Gtype{
		typ:         G_SLICE,
		elementType: e.collection.getGtype().elementType,
	}
}

func (e *ExprIndex) getGtype() *Gtype {
	assert(e.collection.getGtype() != nil, e.token(), "collection type should not be nil")
	gtype := e.collection.getGtype()
	if gtype.typ == G_REL {
		gtype = gtype.relation.gtype
	}

	if gtype.typ == G_MAP {
		// map value
		return gtype.mapValue
	} else if gtype.typ == G_STRING {
		// "hello"[i]
		return gByte
	} else if gtype.typ == G_SLICE {
		return gtype.elementType
	} else {
		// array element
		return gtype.elementType
	}
}

func (e *ExprIndex) getSecondGtype() *Gtype {
	assertNotNil(e.collection.getGtype() != nil, nil)
	if e.collection.getGtype().typ == G_MAP {
		// map
		return gBool
	}

	return nil
}

func (e *ExprStructField) getGtype() *Gtype {
	gstruct := e.strct.getGtype()

	assert(gstruct != gInt, e.tok, "struct should not be gInt")

	var strctType *Gtype
	if gstruct.typ == G_POINTER {
		strctType = gstruct.origType
	} else {
		strctType = gstruct
	}

	fields := strctType.relation.gtype.fields
	//debugf("fields=%v", fields)
	for _, field := range fields {
		if e.fieldname == field.fieldname {
			//return field.origType
			return field
		}
	}
	return nil
}

func (e *ExprArrayLiteral) getGtype() *Gtype {
	return e.gtype
}

func (e *ExprNumberLiteral) getGtype() *Gtype {
	return gInt
}

func (e *ExprStringLiteral) getGtype() *Gtype {
	return &Gtype{
		typ:    G_STRING,
		length: len(e.val),
	}
}

func (e *ExprLen) getGtype() *Gtype {
	return gInt
}

func (e *ExprCap) getGtype() *Gtype {
	return gInt
}

func (e *ExprVariable) getGtype() *Gtype {
	return e.gtype
}

func (e *ExprConstVariable) getGtype() *Gtype {
	return e.gtype
}

func (e *ExprBinop) getGtype() *Gtype {
	switch e.op {
	case "<", ">", "<=", ">=", "!=", "==", "&&", "||":
		return gBool
	case "+":
		return e.left.getGtype()
	case "-", "*", "%", "/":
		return gInt
	}
	errorf("internal error")
	return nil
}

func (e *ExprNilLiteral) getGtype() *Gtype {
	return nil
}

func (e *ExprConversion) getGtype() *Gtype {
	return e.gtype
}

func (e *ExprTypeSwitchGuard) getGtype() *Gtype {
	TBI(e.token(), "")
	return nil
}

func (e *ExprTypeAssertion) getGtype() *Gtype {
	return e.gtype
}

func (e *ExprVaArg) getGtype() *Gtype {
	return e.expr.getGtype()
}

func (e *ExprMapLiteral) getGtype() *Gtype {
	return e.gtype
}

func (e *ExprConversionToInterface) getGtype() *Gtype {
	return gInterface
}
