package main

import "fmt"

type GTYPE_TYPE int

const (
	G_UNKOWNE GTYPE_TYPE = iota
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
	fname identifier
	paramTypes []*Gtype
	rettypes []*Gtype
}

type Gtype struct {
	typ             GTYPE_TYPE
	relation        *Relation  // for G_REL
	size            int        // for scalar type like int, bool, byte, for struct
	ptr             *Gtype     // for array, pointer
	fields          []*Gtype   // for struct
	fieldname       identifier // for struct field
	offset          int        // for struct field
	length          int        // for slice, array
	capacity        int        // for slice
	underlyingarray interface{}
	imethods        []*signature // for interface
	methods map[identifier]*ExprFuncRef
	// for fixed array
}

func (gtype *Gtype) getSize() int {
	if gtype.typ == G_REL {
		if gtype.relation.gtype == nil {
			errorf("relation not resolved: %s", gtype.relation)
		}
		return gtype.relation.gtype.getSize()
	} else {
		if gtype.typ == G_ARRAY {
			return gtype.length * gtype.ptr.getSize()
		} else if gtype.typ == G_STRUCT {
			// @TODO consider the case of real zero e.g. struct{}
			if gtype.size == 0 {
				gtype.calcStructOffset()
			}
			return gtype.size
		} else if gtype.typ == G_POINTER || gtype.typ == G_STRING ||gtype.typ == G_INTERFACE {
			return ptrSize
		} else {
			return gtype.size
		}
	}
}

func (gtype *Gtype) String() string {
	switch gtype.typ {
	case G_REL:
		return fmt.Sprintf("rel(%s)", gtype.relation.name)
	case G_INT:
		return "int"
	case G_BYTE:
		return "byte"
	case G_ARRAY:
		elm := gtype.ptr
		return fmt.Sprintf("[]%s", elm)
	case G_STRUCT:
		return "struct"
	case G_STRUCT_FIELD:
		return "structfield"
	case G_POINTER:
		elm := gtype.ptr
		return fmt.Sprintf("*%s", elm)
	case G_SLICE:
		return "slice"
	case G_STRING:
		return "string"
	default:
		errorf("gtype.String() error: type=%d", gtype.typ)
	}
	return ""
}

func (strct *Gtype) getField(name identifier) *Gtype {
	assert(strct != nil, "assume G_STRUCT type")
	assert(strct.typ == G_STRUCT, "assume G_STRUCT type")
	for _, field := range strct.fields {
		if field.fieldname == name {
			return field
		}
	}
	errorf("field %s not found in the struct", name)
	return nil
}

func (strct *Gtype) calcStructOffset() {
	assert(strct.typ == G_STRUCT, "assume G_STRUCT type")
	var offset int
	for _, fieldtype := range strct.fields {
		var align int
		if fieldtype.getSize() < MaxAlign {
			align = fieldtype.getSize()
		} else {
			align = MaxAlign
		}
		if offset%align != 0 {
			offset += align - offset%align
		}
		fieldtype.offset = offset
		offset += fieldtype.getSize()
	}

	strct.size = offset
}

func (rel *Relation) getGtype() *Gtype {
	return rel.expr.getGtype()
}

func (e *ExprStructLiteral) getGtype() *Gtype {
	return e.strctname.gtype // this can be nil while the relation is unresolved.
}

func (e *ExprFuncall) getGtype() *Gtype {
	return e.rel.expr.(*ExprFuncRef).funcdef.rettypes[0]
}

func (e *ExprMethodcall) getGtype() *Gtype {
	gtype := e.receiver.getGtype()
	method := gtype.methods[e.fname]
	return method.funcdef.rettypes[0]
}

func (e *ExprUop) getGtype() *Gtype {
	switch e.op {
	case "&":
		return &Gtype{
			typ: G_POINTER,
			ptr: e.operand.getGtype(),
		}
	case "*":
		return e.operand.getGtype().ptr
	case "!":
		return gBool
	}
	errorf("internal error")
	return nil
}

func (f *ExprFuncRef) getGtype() *Gtype {
	return &Gtype{
		typ:G_FUNC,
	}
}

func (e *ExprSliced) getGtype() *Gtype {
	errorf("TBI")
	return nil
}

func (e *ExprArrayIndex) getGtype() *Gtype {
	return e.array.getGtype().ptr
}

func (e *AstStructFieldAccess) getGtype() *Gtype {
	gstruct := e.strct.getGtype()
	for _, field := range gstruct.fields {
		if e.fieldname == field.fieldname {
			return field.ptr
		}
	}
	errorf("internal error")
	return nil
}

func (e ExprArrayLiteral) getGtype() *Gtype {
	return e.gtype
}

func (e *ExprNumberLiteral) getGtype() *Gtype {
	return gInt
}

func (e *ExprStringLiteral) getGtype() *Gtype {
	return gString
}

func (e *ExprVariable) getGtype() *Gtype {
	return e.gtype
}

func (e *ExprConstVariable) getGtype() *Gtype {
	return e.gtype
}

func (e *ExprBinop) getGtype() *Gtype {
	switch e.op {
	case "<",">","<=",">=","!=","==","&&", "||":
		return gBool
	case "+","-","*","%","/":
		return gInt
	}
	errorf("internal error")
	return nil
}
