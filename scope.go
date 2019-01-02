package main

// built-in types
var gInt = &Gtype{typ: G_INT, size: 8}
var gByte = &Gtype{typ: G_BYTE, size: 1}
var gBool = &Gtype{typ: G_BOOL, size: 8}
var gString = &Gtype{typ: G_SLICE,}

const ptrSize int = 8

type scope struct {
	idents map[identifier]interface{}
	outer  *scope
}

func (sc *scope) get(name identifier) interface{} {
	for s := sc; s != nil; s = s.outer {
		v := s.idents[name]
		if v != nil {
			return v
		}
	}
	return nil
}

func (sc *scope) setFunc(name identifier, funcref *ExprFuncRef) {
	sc._set(name, funcref)
}

func (sc *scope) setConst(name identifier, cnst *ExprConstVariable) {
	sc._set(name, cnst)
}

func (sc *scope) setVar(name identifier, variable *ExprVariable) {
	sc._set(name, variable)
}

func (sc *scope) setGtype(name identifier, gtype *Gtype) {
	sc._set(name, gtype)
}

func (sc *scope) _set(name identifier, v interface{}) {
	if v == nil {
		panic("nil cannot be set")
	}
	sc.idents[name] = v
}

func (sc *scope) getGtype(name identifier) *Gtype {
	v := sc.get(name)
	if v == nil {
		return nil
	}
	gtype, ok := v.(*Gtype)
	if !ok {
		errorf("type %s is not defined", name)
	}
	return gtype
}

func newScope(outer *scope) *scope {
	return &scope{
		outer:  outer,
		idents: make(map[identifier]interface{}),
	}
}

func newUniverseBlockScope() *scope {
	r := newScope(nil)

	r.setGtype("int", gInt)
	r.setGtype("byte", gByte)
	r.setGtype("bool", gBool)
	r.setConst(identifier("iota"), &ExprConstVariable{})
	r.setConst("true", &ExprConstVariable{
		name:  "true",
		gtype: gBool,
		val:   &ExprNumberLiteral{1},
	})
	r.setConst("false", &ExprConstVariable{
		name:  "false",
		gtype: gBool,
		val:   &ExprNumberLiteral{0},
	})

	r.setFunc("len", &ExprFuncRef{
		// @FIXME
	})
	r.setFunc("println", &ExprFuncRef{
		// @FIXME
	})
	r.setFunc("Printf", &ExprFuncRef{
		// @FIXME : should be fmt.Printf
	})
	return r
}
