package main

var gInt = &Gtype{typ: "scalar", size: 8}
var gBool = &Gtype{typ: "scalar", size: 8}

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

func (sc *scope) setVarDecl(name identifier, decl *AstVarDecl) {
	sc._set(name, decl)
}

func (sc *scope) setTypeDecl(name identifier, decl *AstTypeDecl) {
	sc._set(name, decl)
}

func (sc *scope) setConstDecl(name identifier, decl *AstConstDecl) {
	sc._set(name, decl)
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
		errorf("type %s is not defined", name)
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
	r.setTypeDecl("int", &AstTypeDecl{gtype: gInt})
	r.setTypeDecl("bool", &AstTypeDecl{gtype: gBool})

	r.setConstDecl(identifier("iota"), &AstConstDecl{})
	r.setConstDecl("true", &AstConstDecl{
		variable: &ExprConstVariable{
			name:  "true",
			gtype: gBool,
			val:   &ExprNumberLiteral{1},
		},
	})
	r.setConstDecl("false", &AstConstDecl{
		variable: &ExprConstVariable{
			name:  "false",
			gtype: gBool,
			val:   &ExprNumberLiteral{0},
		},
	})


	r.setFuncDecl("len", &AstFuncDecl{fname: "len",})
	return r
}
