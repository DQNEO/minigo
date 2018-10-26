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
	G_ARRAY
	G_SLICE
	G_STRING
	G_MAP
	G_POINTER
)

type Gtype struct {
	typ       GTYPE_TYPE
	relname   identifier    // for G_REL
	relation  *Relation     // for G_REL
	size      int           // for scalar type like int, bool, byte
	ptr       *Gtype        // for array, pointer
	structdef *AstStructDef // for struct type
	length    int           // for fixed array
}

type StructField struct {
	name  identifier
	gtype *Gtype
}

func (gtype *Gtype) String() string {
	switch gtype.typ {
	case G_REL:
		return "rel"
	case G_INT:
		return "int"
	case G_BYTE:
		return "byte"
	case G_ARRAY:
		elm := gtype.ptr
		return fmt.Sprintf("[]%s", elm)
	default:
		errorf("default: %d", gtype.typ)
	}
	return ""
}


