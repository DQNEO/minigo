package main

var gInt = &Gtype{typ: G_INT, size: 8}
var gByte = &Gtype{typ: G_BYTE, size: 1}
var gBool = &Gtype{typ: G_BOOL, size: 8}

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

func (sc *scope) setFuncDecl(name identifier, decl *AstFuncDecl) {
	sc._set(name, decl)
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
	decl, ok := v.(*AstTypeDecl)
	if !ok {
		errorf("type %s is not defined", name)
	}
	return decl.gtype
}

func newScope(outer *scope) *scope {
	return &scope{
		outer:  outer,
		idents: make(map[identifier]interface{}),
	}
}

func newUniverseBlockScope() *scope {
	r := newScope(nil)
	r._set("int", gInt)
	r._set("byte", gByte)
	r._set("bool", gBool)
	r._set(identifier("iota"), &ExprConstVariable{})
	r._set("true",  &ExprConstVariable{
			name:  "true",
			gtype: gBool,
			val:   &ExprNumberLiteral{1},
	})
	r._set("false", &ExprConstVariable{
			name:  "false",
			gtype: gBool,
			val:   &ExprNumberLiteral{0},
	})

	r.setFuncDecl("len", &AstFuncDecl{fname: "len",})
	return r
}
