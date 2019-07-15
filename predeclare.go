package main

// built-in types
const sizeOfInterface = 8 * 3

var sInterface = Gtype{kind: G_INTERFACE, size: sizeOfInterface}
var gInterface = &sInterface
var sInt = Gtype{kind: G_INT, size: 8}
var gInt = &sInt
var sByte = Gtype{kind: G_BYTE, size: 1}
var gByte = &sByte
var gBool = &Gtype{kind: G_BOOL, size: 8} // we treat bool as quad length data for now
var gString = &Gtype{
	kind: G_STRING,
	elementType:&sByte,
}
var sTrue = ExprNumberLiteral{val: 1}
var sFalse = ExprNumberLiteral{val: 0}
var eTrue = &sTrue
var eFalse = &sFalse

var builtinTypesAsString []string = []string{
	"bool",
	"byte",
	"int",
	"string",
	"func",
}

var eIota = &ExprConstVariable{
	name: identifier("iota"),
}

var builtinPanic = &DeclFunc{
	rettypes: []*Gtype{},
}

var builtinLen = &DeclFunc{
	rettypes: []*Gtype{&sInt},
}

var builtinCap = &DeclFunc{
	rettypes: []*Gtype{&sInt},
}

var builtinAppend = &DeclFunc{
	rettypes: []*Gtype{&sInt},
}

var builtinMakeSlice = &DeclFunc{
	rettypes: []*Gtype{&sBuiltinRunTimeArgsRettypes1},
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

var sBuiltinRunTimeArgsRettypes1 Gtype = Gtype{
	kind: G_SLICE,
	size: IntSize * 3,
	elementType: &Gtype{
		kind: G_STRING,
	},
}

var builtinPrintstring = &DeclFunc{
	rettypes: []*Gtype{},
}

func newUniverse() *Scope {
	universe := newScope(nil, "universe")
	setPredeclaredIdentifiers(universe)
	return universe
}

// https://golang.org/ref/spec#Predeclared_identifiers
func setPredeclaredIdentifiers(universe *Scope) {
	predeclareNil(universe)
	predeclareTypes(universe)
	predeclareConsts(universe)
	predeclareLibcFuncs(universe)

	universe.setFunc(identifier("panic"), &ExprFuncRef{
		funcdef: builtinPanic,
	})
	universe.setFunc(identifier("len"), &ExprFuncRef{
		funcdef: builtinLen,
	})
	universe.setFunc(identifier("cap"), &ExprFuncRef{
		funcdef: builtinCap,
	})
	universe.setFunc(identifier("append"), &ExprFuncRef{
		funcdef: builtinAppend,
	})
	universe.setFunc(identifier("makeSlice"), &ExprFuncRef{
		funcdef: builtinMakeSlice,
	})

	universe.setFunc(identifier("dumpSlice"), &ExprFuncRef{
		funcdef: builtinDumpSlice,
	})

	universe.setFunc(identifier("dumpInterface"), &ExprFuncRef{
		funcdef: builtinDumpInterface,
	})

	universe.setFunc(identifier("assertInterface"), &ExprFuncRef{
		funcdef: builtinAssertInterface,
	})

	universe.setFunc(identifier("asComment"), &ExprFuncRef{
		funcdef: builtinAsComment,
	})

	universe.setFunc(identifier("printstring"), &ExprFuncRef{
		funcdef: builtinPrintstring,
	})

}

// Zero value:
// nil
func predeclareNil(universe *Scope) {
	universe.set(identifier("nil"), &IdentBody{
		expr: &ExprNilLiteral{},
	})
}

// Types:
// bool byte complex64 complex128 error float32 float64
// int int8 int16 int32 int64 rune string
// uint uint8 uint16 uint32 uint64 uintptr
func predeclareTypes(universe *Scope) {
	universe.setGtype(identifier("bool"), gBool)
	universe.setGtype(identifier("byte"), gByte)
	universe.setGtype(identifier("int"), gInt)
	universe.setGtype(identifier("string"), gString)
	universe.setGtype(identifier("uint8"), gByte)
}

// Constants:
// true false iota
func predeclareConsts(universe *Scope) {
	universe.setConst(identifier("true"), &ExprConstVariable{
		name:  identifier("true"),
		gtype: gBool,
		val:   eTrue,
	})
	universe.setConst(identifier("false"), &ExprConstVariable{
		name:  identifier("false"),
		gtype: gBool,
		val:   eFalse,
	})

	universe.setConst(identifier("iota"), eIota)
}

func predeclareLibcFuncs(universe *Scope) {
	libc := identifier("libc")

	universe.setFunc(identifier("exit"), &ExprFuncRef{
		funcdef: &DeclFunc{
			pkg: libc,
		},
	})
	universe.setFunc(identifier("open"), &ExprFuncRef{
		funcdef: &DeclFunc{
			pkg: libc,
			rettypes: []*Gtype{gInt},
		},
	})

	universe.setFunc(identifier("read"), &ExprFuncRef{
		funcdef: &DeclFunc{
			pkg: libc,
			rettypes: []*Gtype{gInt},
		},
	})

	universe.setFunc(identifier("write"), &ExprFuncRef{
		funcdef: &DeclFunc{
			pkg: libc,
			rettypes: []*Gtype{gInt},
		},
	})
}
