package main

type Scope struct {
	idents map[identifier]*IdentBody
	name   string
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

func (sc *Scope) setConst(name goidentifier, cnst *ExprConstVariable) {
	sc.set(identifier(name), &IdentBody{
		expr: cnst,
	})
}

func (sc *Scope) setVar(name goidentifier, variable *ExprVariable) {
	sc.set(identifier(name), &IdentBody{
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
		panic("nil cannot be set")
	}
	sc.idents[name] = elm
}

func (sc *Scope) getGtype(name identifier) *Gtype {
	if sc == nil {
		errorf("sc is nil")
	}
	idents := sc.idents
	elm, ok := idents[name]
	if !ok {
		return nil
	}
	if elm.gtype == nil {
		errorf("type %s is not defined", name)
	}
	return elm.gtype
}

func newScope(outer *Scope, name string) *Scope {
	return &Scope{
		outer:  outer,
		name:   name,
		idents: map[identifier]*IdentBody{},
	}
}
