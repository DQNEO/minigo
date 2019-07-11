package main

type EType int

const undefinedSize = -1

const (
	G_UNKOWNE   EType = iota
	G_DEPENDENT       // depends on other expression
	G_NAMED
	// below are primitives which are declared in the universe block
	G_INT
	G_BOOL
	G_BYTE
	// end of primitives
	G_STRUCT
	G_STRUCT_FIELD
	G_ARRAY
	G_SLICE
	G_CLIKE_STRING
	G_MAP
	G_POINTER
	G_FUNC
	G_INTERFACE
)

type signature struct {
	fname      goidentifier
	paramTypes []*Gtype
	rettypes   []*Gtype
}

type Gtype struct {
	kind           EType                       // Elementary type
	receiverTypeId int                         // for receiverTypeId. 0:unkonwn
	dependendson   Expr                        // for G_DEPENDENT
	relation       *Relation                   // for G_NAMED
	size           int                         // for scalar type like int, bool, byte, for struct
	origType       *Gtype                      // for pointer
	fields         []*Gtype                    // for struct
	fieldname      goidentifier                  // for struct field
	offset         int                         // for struct field
	padding        int                         // for struct field
	length         int                         // for array, string (len without the terminating \0)
	elementType    *Gtype                      // for array, slice
	imethods       map[identifier]*signature   // for interface
	methods        map[identifier]*ExprFuncRef // for G_NAMED
	mapKey         *Gtype                      // for map
	mapValue       *Gtype                      // for map
}

func imethodGet(imethods map[identifier]*signature, name goidentifier) (*signature, bool) {
	ref, ok := imethods[toKey(name)]
	return ref, ok
}

func imethodSet(imethods map[identifier]*signature, name goidentifier, sig *signature) {
	imethods[toKey(name)] = sig
}

func methodGet(methods map[identifier]*ExprFuncRef, name goidentifier) (*ExprFuncRef, bool) {
	ref, ok := methods[toKey(name)]
	return ref, ok
}

func methodSet(methods map[identifier]*ExprFuncRef, name goidentifier, ref *ExprFuncRef) {
	methods[toKey(name)] = ref
}


func (gtype *Gtype) isNil() bool {
	if gtype == nil {
		return true
	}
	if gtype.kind == G_NAMED {
		return gtype.relation.gtype == nil

	}
	return false
}

func (gtype *Gtype) Underlying() *Gtype {
	if gtype.kind == G_NAMED {
		if gtype.relation.gtype == nil {
			// nil type
			return gtype
		}
		return gtype.relation.gtype.Underlying()

	}
	return gtype
}

func (gtype *Gtype) getKind() EType {
	if gtype == nil {
		return G_UNKOWNE
	}
	return gtype.Underlying().kind
}

// is array or slice
func (gtype *Gtype) isArrayLike() bool {
	kind := gtype.getKind()
	return kind == G_ARRAY || kind == G_SLICE
}

func (gtype *Gtype) is24WidthType() bool {
	switch gtype.getKind() {
	case G_INTERFACE, G_SLICE:
		return true
	default:
		return false
	}
}

func (gtype *Gtype) isBytesSlice() bool {
	underLying := gtype.Underlying()
	if underLying.kind == G_SLICE && underLying.elementType.getKind() == G_BYTE {
		return true
	}
	return false
}

func (gtype *Gtype) isClikeString() bool {
	return gtype.getKind() == G_CLIKE_STRING
}

func (gtype *Gtype) getSize() int {
	assert(gtype != nil, nil, S("gtype should not be nil"))
	assert(gtype.kind != G_DEPENDENT, nil, S("type should be inferred"))
	if gtype.kind == G_NAMED {
		if gtype.relation.gtype == nil {
			errorf(S("relation not resolved: %s"), gtype)
		}
		return gtype.relation.gtype.getSize()
	} else {
		switch gtype.kind {
		case G_ARRAY:
			assertNotNil(gtype.elementType != nil, nil)
			return gtype.length * gtype.elementType.getSize()
		case G_STRUCT:
			if gtype.size == undefinedSize {
				gtype.calcStructOffset()
			}
			return gtype.size
		case G_POINTER:
			return ptrSize
		case G_INTERFACE:
			//     data    ,  receiverTypeId, dtype
			return ptrSize + ptrSize + ptrSize
		case G_SLICE:
			return ptrSize + IntSize + IntSize
		case G_MAP:
			return ptrSize
		case G_CLIKE_STRING:
			return ptrSize
		default:
			return gtype.size
		}
	}
}

func (gtype *Gtype) String() gostring {
	var gs gostring
	if gtype == nil {
		return S("NO_TYPE")
	}
	switch gtype.kind {
	case G_DEPENDENT:
		return S("dependent")
	case G_NAMED:
		if len(gtype.relation.pkg) == 0 {
			//errorf(S("pkg is empty: %s"), gtype.relation.name)
		}
		child := gtype.relation.gtype
		if child != nil  {
			switch child.kind {
			case G_INT:
				return S("int")
			case G_BOOL:
				return S("bool")
			case G_BYTE:
				return S("byte")
			case G_CLIKE_STRING:
				return S("clikestring")
			case G_FUNC:
				return S("func")
			}
		}
		 gs = Sprintf(S("G_NAMED(%s.%s)"),
			gostring(gtype.relation.pkg), gostring(gtype.relation.name))
		return gs
	case G_INT:
		return S("int")
	case G_BOOL:
		return S("bool")
	case G_BYTE:
		return S("byte")
	case G_ARRAY:
		elm := gtype.elementType
		gs = Sprintf(S("[%d]%s"), gtype.length, elm.String())
		return gs
	case G_STRUCT:
		var r gostring = S("struct{")
		for _, field := range gtype.fields {
			var tmp gostring = concat(field.String(), S(","))
			r = concat(r, tmp)
		}
		r = concat(r, S("}"))
		return r
	case G_STRUCT_FIELD:
		return S("structfield")
	case G_POINTER:
		origType := gtype.origType
		gs = Sprintf(S("*%s"), origType.String())
		return gs
	case G_SLICE:
		gs = Sprintf(S("[]%s"), gtype.elementType.String())
		return gs
	case G_CLIKE_STRING:
		return S("clikestring")
	case G_FUNC:
		return S("func")
	case G_INTERFACE:
		if len(gtype.imethods) == 0 {
			return S("interface{}")
		} else {
			return S("interface {...}")
		}
	case G_MAP:
		return S("map")
	default:
		errorf(S("gtype.String() error: invalid gtype.type=%d"), gtype.kind)
	}
	return S("")
}

func (strct *Gtype) getField(name goidentifier) *Gtype {
	assertNotNil(strct != nil, nil)
	assert(strct.kind == G_STRUCT, nil, S("assume G_STRUCT type"))
	for _, field := range strct.fields {
		if eq(gostring(field.fieldname), gostring(name)) {
			return field
		}
	}
	errorf(S("field %s not found in the struct"), name)
	return nil
}

func (strct *Gtype) calcStructOffset() {
	assert(strct.getKind() == G_STRUCT, nil, S("assume G_STRUCT type, but got %s"), strct.String())
	var offset int
	for _, fieldtype := range strct.fields {
		var align int
		if fieldtype.getSize() < MaxAlign {
			align = fieldtype.getSize()
			assert(align > 0, nil, S("field size should be > 0: filed=%s"), fieldtype.String())
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
		//errorft(rel.token(), S("rel.expr is nil for %s"), rel)
		return nil
	}
	return rel.expr.getGtype()
}

func (e *ExprStructLiteral) getGtype() *Gtype {
	return &Gtype{
		kind:     G_NAMED,
		relation: e.strctname,
	}
}

func (e *ExprFuncallOrConversion) getGtype() *Gtype {
	assert(e.rel.expr != nil || e.rel.gtype != nil, e.token(), S(""))
	if e.rel.expr != nil {
		funcref, ok := e.rel.expr.(*ExprFuncRef)
		assert(ok, e.token(), S("it should be a ExprFuncRef"))
		firstRetType := funcref.funcdef.rettypes[0]
		return firstRetType
	} else if e.rel.gtype != nil {
		return e.rel.gtype
	}
	assertNotReached(e.token())
	return nil
}

func (e *ExprMethodcall) getGtype() *Gtype {
	gtype := e.receiver.getGtype()
	if gtype.kind == G_POINTER {
		gtype = gtype.origType
	}

	underlyingType := gtype.relation.gtype // I forgot the reason to do this :(
	if underlyingType.kind == G_INTERFACE {
		methodsig, ok := imethodGet(underlyingType.imethods, e.fname)
		if !ok {
			errorft(e.token(), S("method %s not found in %s %s"), e.fname, gtype, e.tok)
		}
		assertNotNil(methodsig != nil, e.tok)
		return methodsig.rettypes[0]
	} else {
		method, ok := methodGet(underlyingType.methods, e.fname)
		if !ok {
			errorft(e.token(), S("method %s not found in %s %s"), e.fname, gtype, e.tok)
		}
		assertNotNil(method != nil, e.tok)
		return method.funcdef.rettypes[0]
	}
}

func (e *ExprUop) getGtype() *Gtype {
	switch switchexpr(e.op) {
	case "&":
		return &Gtype{
			kind:     G_POINTER,
			origType: e.operand.getGtype(),
		}
	case "*":
		return e.operand.getGtype().origType
	case "!":
		return gBool
	case "-":
		return gInt
	}
	errorf(S("internal error"))
	return nil
}

func (f *ExprFuncRef) getGtype() *Gtype {
	return &Gtype{
		kind: G_FUNC,
	}
}

func (e *ExprSlice) getGtype() *Gtype {
	return &Gtype{
		kind:        G_SLICE,
		elementType: e.collection.getGtype().elementType,
	}
}

func (e *ExprIndex) getGtype() *Gtype {
	assert(e.collection.getGtype() != nil, e.token(), S("collection type should not be nil"))
	gtype := e.collection.getGtype()
	if gtype.kind == G_NAMED {
		gtype = gtype.relation.gtype
	}

	if gtype.kind == G_MAP {
		// map value
		return gtype.mapValue
	} else if gtype.kind == G_SLICE {
		return gtype.elementType
	} else {
		// array element
		return gtype.elementType
	}
}

// type of fist,second := m[a]
func (e *ExprIndex) getSecondGtype() *Gtype {
	assertNotNil(e.collection.getGtype() != nil, nil)
	if e.collection.getGtype().getKind() == G_MAP {
		// map
		return gBool
	}

	return nil
}

func (e *ExprStructField) getGtype() *Gtype {
	gstruct := e.strct.getGtype()

	assert(gstruct != gInt, e.tok, S("struct should not be gInt"))

	var strctType *Gtype
	if gstruct.kind == G_POINTER {
		strctType = gstruct.origType
	} else {
		strctType = gstruct
	}

	fields := strctType.relation.gtype.fields
	//debugf(S("fields=%v"), fields)
	for _, field := range fields {
		if eq(gostring(e.fieldname),  gostring(field.fieldname)) {
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
		kind:   G_CLIKE_STRING,
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
	switch switchexpr(e.op) {
	case "<", ">", "<=", ">=", "!=", "==", "&&", "||":
		return gBool
	case "+":
		return e.left.getGtype()
	case "-", "*", "%", "/":
		return gInt
	}
	errorf(S("internal error"))
	return nil
}

func (e *ExprNilLiteral) getGtype() *Gtype {
	return nil
}

func (e *IrExprConversion) getGtype() *Gtype {
	return e.toGtype
}

func (e *ExprTypeSwitchGuard) getGtype() *Gtype {
	return e.expr.getGtype()
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

func (e *IrExprConversionToInterface) getGtype() *Gtype {
	return gInterface
}
