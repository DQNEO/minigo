package main

// built-in types
const sizeOfInterface = 8 * 3

var sInterface = Gtype{kind: G_INTERFACE, size: sizeOfInterface}
var gInterface = &sInterface
var sInt = Gtype{kind: G_INT, size: 8}
var gInt = &sInt
var sUintptr = Gtype{kind: G_UINT_PTR, size: 8}
var gUintptr = &sUintptr
var sByte = Gtype{kind: G_BYTE, size: 1}
var gByte = &sByte
var sUint16 = Gtype{kind: G_UINT_16, size: 2}
var gUint16 = &sUint16
var gBool = &Gtype{kind: G_BOOL, size: 8} // we treat bool as quad length data for now
var gString = &Gtype{
	kind:        G_STRING,
	elementType: &sByte,
}

var sTrue = IrExprBoolVal{bol: true}
var sFalse = IrExprBoolVal{bol: false}

var eTrue = &sTrue
var eFalse = &sFalse

var builtinTypesAsString []string = []string{
	"bool",
	"byte",
	"int",
	"uint16",
	"unintptr",
	"string",
	"func",
}

var eIota = &ExprConstVariable{
	name: identifier("iota"),
}

var builtinPanic = &DeclFunc{
	pkgPath:  "/builtin",
	rettypes: []*Gtype{},
}

var builtinLen = &DeclFunc{
	pkgPath:  "/builtin",
	rettypes: []*Gtype{&sInt},
}

var builtinCap = &DeclFunc{
	pkgPath:  "/builtin",
	rettypes: []*Gtype{&sInt},
}

var builtinAppend = &DeclFunc{
	pkgPath:  "/builtin",
	rettypes: []*Gtype{&sInt},
}

var builtinMake = &DeclFunc{
	pkgPath:  "/builtin",
	rettypes: []*Gtype{},
}

var builtinSyscall = &DeclFunc{
	pkgPath:  "/builtin",
	rettypes: []*Gtype{&sInt},
}

var builtinDumpSlice = &DeclFunc{
	pkgPath:  "/builtin",
	rettypes: []*Gtype{},
}

var builtinDumpInterface = &DeclFunc{
	pkgPath:  "/builtin",
	rettypes: []*Gtype{},
}

var builtinAssertInterface = &DeclFunc{
	pkgPath:  "/builtin",
	rettypes: []*Gtype{},
}

var builtinAsComment = &DeclFunc{
	pkgPath:  "/builtin",
	rettypes: []*Gtype{},
}

var sSliceType Gtype = Gtype{
	kind: G_SLICE,
	size: IntSize * 3,
	elementType: &Gtype{
		kind: G_STRING,
	},
}

func newUniverse() *Scope {
	universe := newScope(nil, identifier("universe"))
	setPredeclaredIdentifiers(universe)
	return universe
}

// https://golang.org/ref/spec#Predeclared_identifiers
func setPredeclaredIdentifiers(universe *Scope) {
	predeclareNil(universe)
	predeclareTypes(universe)
	predeclareConsts(universe)
	predeclareFuncs(universe)
}

func predeclareFuncs(universe *Scope) {
	var builtinFuncs map[identifier]*DeclFunc = map[identifier]*DeclFunc{
		// Inject genuine builtin funcs
		"panic": builtinPanic,
		"len": builtinLen,
		"cap": builtinCap,
		"append": builtinAppend,
		"make": builtinMake,

		// Inject my builtin funcs
		"Syscall": builtinSyscall,
		"dumpSlice": builtinDumpSlice,
		"dumpInterface": builtinDumpInterface,
		"assertInterface": builtinAssertInterface,
		"asComment": builtinAsComment,
	}

	for name, def := range builtinFuncs {
		universe.setFunc(name,  &ExprFuncRef{
			funcdef: def,
		})
	}
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
	universe.setGtype(identifier("uint16"), gUint16)
	universe.setGtype(identifier("uintptr"), gUintptr)
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
