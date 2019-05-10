package main

// built-in types
const sizeOfInterface = 8 * 3

var gInterface *Gtype = &Gtype{typ: G_INTERFACE, size: sizeOfInterface}
var gInt  *Gtype  = &Gtype{typ: G_INT, size: 8}
var gByte  *Gtype = &Gtype{typ: G_BYTE, size: 1}
var gBool  *Gtype = &Gtype{typ: G_BOOL, size: 8} // we treat bool as quad length data for now
var gString  *Gtype = &Gtype{
	typ:G_STRING,
}

var builtinTypesAsString []string = []string{"bool", "byte", "int", "string"}

var eIota = &ExprConstVariable{
	name: "iota",
}

var builtinLen = &DeclFunc{
	rettypes: []*Gtype{gInt},
}

var builtinCap = &DeclFunc{
	rettypes: []*Gtype{gInterface},
}

var builtinAppend = &DeclFunc{
	rettypes: []*Gtype{gInt},
}

var builtinMakeSlice = &DeclFunc{
	rettypes: []*Gtype{builtinRunTimeArgsRettypes1},
}

var builtinDumpSlice = &DeclFunc{
	rettypes: []*Gtype{},
}

var builtinDumpInterface = &DeclFunc{
	rettypes: []*Gtype{},
}

var builtinAssertInterface = &DeclFunc{
	rettypes: []*Gtype{},
}

var builtinAsComment = &DeclFunc{
	rettypes: []*Gtype{},
}

var builtinRunTimeArgsRettypes1 *Gtype =&Gtype{
	typ: G_SLICE,
	size: IntSize * 3,
	elementType: &Gtype{
		typ:G_STRING,
	},
}

var builtinRunTimeArgs = &DeclFunc{
	rettypes: []*Gtype{
		builtinRunTimeArgsRettypes1,
	},
	}


func newUniverse() *scope {
	universe := newScope(nil, "universe")
	setPredeclaredIdentifiers(universe)
	return universe
}

// https://golang.org/ref/spec#Predeclared_identifiers
func setPredeclaredIdentifiers(universe *scope) {
	predeclareNil(universe)
	predeclareTypes(universe)
	predeclareConsts(universe)
	predeclareLibcFuncs(universe)

	universe.setFunc("len", &ExprFuncRef{
		funcdef: builtinLen,
	})
	universe.setFunc("cap", &ExprFuncRef{
		funcdef: builtinCap,
	})
	universe.setFunc("append", &ExprFuncRef{
		funcdef: builtinAppend,
	})
	universe.setFunc("makeSlice", &ExprFuncRef{
		funcdef: builtinMakeSlice,
	})

	universe.setFunc("dumpSlice", &ExprFuncRef{
		funcdef: builtinDumpSlice,
	})

	universe.setFunc("dumpInterface", &ExprFuncRef{
		funcdef: builtinDumpInterface,
	})

	universe.setFunc("assertInterface", &ExprFuncRef{
		funcdef: builtinAssertInterface,
	})

	universe.setFunc("asComment", &ExprFuncRef{
		funcdef: builtinAsComment,
	})

	universe.setFunc("runtime_args", &ExprFuncRef{
		funcdef: builtinRunTimeArgs,
	})
}

// Zero value:
// nil
func predeclareNil(universe *scope) {
	universe.set("nil", &IdentBody{
		expr: &ExprNilLiteral{},
	})
}

// Types:
// bool byte complex64 complex128 error float32 float64
// int int8 int16 int32 int64 rune string
// uint uint8 uint16 uint32 uint64 uintptr
func predeclareTypes(universe *scope) {
	universe.setGtype("bool", gBool)
	universe.setGtype("byte", gByte)
	universe.setGtype("int", gInt)
	universe.setGtype("string", gString)
	universe.setGtype("uint8", gByte)
}

// Constants:
// true false iota
func predeclareConsts(universe *scope) {
	universe.setConst("true", &ExprConstVariable{
		name:  "true",
		gtype: gBool,
		val:   &ExprNumberLiteral{val: 1},
	})
	universe.setConst("false", &ExprConstVariable{
		name:  "false",
		gtype: gBool,
		val:   &ExprNumberLiteral{val: 0},
	})

	universe.setConst("iota", eIota)
}

func predeclareLibcFuncs(universe *scope) {
	universe.setFunc("printf", &ExprFuncRef{
		funcdef: &DeclFunc{
			pkg: "libc",
		},
	})
	universe.setFunc("sprintf", &ExprFuncRef{
		funcdef: &DeclFunc{
			pkg: "libc",
			rettypes: []*Gtype{gInt},
		},
	})
	universe.setFunc("exit", &ExprFuncRef{
		funcdef: &DeclFunc{
			pkg: "libc",
		},
	})
	universe.setFunc("open", &ExprFuncRef{
		funcdef: &DeclFunc{
			pkg:      "libc",
			rettypes: []*Gtype{gInt},
		},
	})

	universe.setFunc("read", &ExprFuncRef{
		funcdef: &DeclFunc{
			pkg:      "libc",
			rettypes: []*Gtype{gInt},
		},
	})

	universe.setFunc("write", &ExprFuncRef{
		funcdef: &DeclFunc{
			pkg:      "libc",
			rettypes: []*Gtype{gInt},
		},
	})

	universe.setFunc("atoi", &ExprFuncRef{
		funcdef: &DeclFunc{
			pkg:      "libc",
			rettypes: []*Gtype{gInt},
		},
	})
}
