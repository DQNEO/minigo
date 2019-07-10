package main

type Scope struct {
	idents map[identifier]*IdentBody
	name   gostring
	outer  *Scope
}

type IdentBody struct {
	typ   int // 1:*Gtype, 2:Expr
	gtype *Gtype
	expr  Expr
}

func (sc *Scope) get(name goidentifier) *IdentBody {
	for s := sc; s != nil; s = s.outer {
		v, ok := s.idents[toKey(name)]
		if ok {
			return v
		}
	}
	return nil
}

func (sc *Scope) setFunc(name goidentifier, funcref *ExprFuncRef) {
	sc.set(name, &IdentBody{
		expr: funcref,
	})
}

func (sc *Scope) setConst(name goidentifier, cnst *ExprConstVariable) {
	sc.set(name, &IdentBody{
		expr: cnst,
	})
}

func (sc *Scope) setVar(name goidentifier, variable *ExprVariable) {
	sc.set(name, &IdentBody{
		expr: variable,
	})
}

func (sc *Scope) setGtype(name goidentifier, gtype *Gtype) {
	sc.set(name, &IdentBody{
		gtype: gtype,
	})
}

func (sc *Scope) set(name goidentifier, elm *IdentBody) {
	if elm == nil {
		panic(S("nil cannot be set"))
	}
	sc.idents[toKey(name)] = elm
}

func (sc *Scope) getGtype(name goidentifier) *Gtype {
	if sc == nil {
		errorf(S("sc is nil"))
	}
	idents := sc.idents
	elm, ok := idents[identifier(name)]
	if !ok {
		return nil
	}
	if elm.gtype == nil {
		errorf(S("type %s is not defined"), name)
	}
	return elm.gtype
}

func newScope(outer *Scope, name gostring) *Scope {
	return &Scope{
		outer:  outer,
		name:   name,
		idents: map[identifier]*IdentBody{},
	}
}
