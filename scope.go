package main

var gInt = &Gtype{typ: "scalar", size: 8}
var gBool = &Gtype{typ: "scalar", size: 8}

var predeclaredConsts = []*AstConstDecl{
	&AstConstDecl{
		variable: &ExprConstVariable{
			name:  "true",
			gtype: gBool,
			val:   &ExprNumberLiteral{1},
		},
	},
	&AstConstDecl{
		variable: &ExprConstVariable{
			name:  "false",
			gtype: gBool,
			val:   &ExprNumberLiteral{0},
		},
	},
}

var predeclaredFunctions = []*AstFuncDecl{
	{
		fname: "len",
	},
}

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

	for _, c := range predeclaredConsts {
		r.setConstDecl(c.variable.name, c)
	}
	for _, f := range predeclaredFunctions {
		r.setFuncDecl(f.fname, &AstFuncDecl{})
	}
	return r
}
