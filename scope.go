package main

type Scope struct {
	idents map[identifier]*IdentBody
	name   bytes
	outer  *Scope
}

type IdentBody struct {
	typ   int // 1:*Gtype, 2:Expr
	gtype *Gtype
	expr  Expr
}

func (sc *Scope) get(name identifier) *IdentBody {
	for s := sc; s != nil; s = s.outer {
		v, ok := s.idents[name]
		if ok {
			return v
		}
	}
	return nil
}

func (sc *Scope) setFunc(name identifier, funcref *ExprFuncRef) {
	sc.set(name, &IdentBody{
		expr: funcref,
	})
}

func (sc *Scope) setConst(name identifier, cnst *ExprConstVariable) {
	sc.set(name, &IdentBody{
		expr: cnst,
	})
}

func (sc *Scope) setVar(name identifier, variable *ExprVariable) {
	sc.set(name, &IdentBody{
		expr: variable,
	})
}

func (sc *Scope) setGtype(name identifier, gtype *Gtype) {
	sc.set(name, &IdentBody{
		gtype: gtype,
	})
}

func (sc *Scope) set(name identifier, elm *IdentBody) {
	if elm == nil {
		panic(S("nil cannot be set"))
	}
	sc.idents[identifier(name)] = elm
}

func (sc *Scope) getGtype(name identifier) *Gtype {
	if sc == nil {
		errorf("sc is nil")
	}
	idents := sc.idents
	elm, ok := idents[identifier(name)]
	if !ok {
		return nil
	}
	if elm.gtype == nil {
		errorf("type %s is not defined", name)
	}
	return elm.gtype
}

func newScope(outer *Scope, name bytes) *Scope {
	return &Scope{
		outer:  outer,
		name:   name,
		idents: map[identifier]*IdentBody{},
	}
}
