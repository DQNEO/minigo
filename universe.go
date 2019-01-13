package main

// built-in types
var gInterface = &Gtype{typ: G_ANY, size: 8}
var gInt = &Gtype{typ: G_INT, size: 8}
var gByte = &Gtype{typ: G_BYTE, size: 1}
var gBool = &Gtype{typ: G_BOOL, size: 8}
var gString = &Gtype{typ: G_STRING}
var eIota = &ExprConstVariable{
	name:"iota",
}

// https://golang.org/ref/spec#Predeclared_identifiers
func setPredeclaredIdentifiers(universe *scope) {
	predeclareTypes(universe)
	predeclareConsts(universe)
	predeclareLibcFuncs(universe)
}

func predeclareTypes(universe *scope) {
	universe.setGtype("int", gInt)
	universe.setGtype("byte", gByte)
	universe.setGtype("bool", gBool)
	universe.setGtype("string", gString)
}

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
	universe.setFunc("puts", &ExprFuncRef{
		funcdef: &DeclFunc{
			pkg: "libc",
			// No implementation thanks to the libc function.
		},
	})
	universe.setFunc("printf", &ExprFuncRef{
		funcdef: &DeclFunc{
			pkg: "libc",
			// No implementation thanks to the libc function.
		},
	})
}
