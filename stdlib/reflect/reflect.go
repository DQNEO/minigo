package reflect

type Type struct {
	name string
}

func TypeOf(arg interface{}) *Type {
	var name string
	switch arg.(type) {
	case int:
		name = "int"
	case string:
		name = "string"
	case byte:
		name = "uint8"
	case bool:
		name = "bool"
	case uintptr:
		name = "uintptr"
	}
	typ := &Type{
		name: name,
	}
	return typ
}

func (typ *Type) String() string {
	return typ.name
}
