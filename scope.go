package main

type scope struct {
	idents map[identifier]*IdentBody
	name   string
	outer  *scope
}

type IdentBody struct {
	typ   int // 1:*Gtype, 2:Expr
	gtype *Gtype
	expr  Expr
}

func (sc *scope) get(name identifier) *IdentBody {
	for s := sc; s != nil; s = s.outer {
		v, ok := s.idents[name]
		if ok {
			return v
		}
	}
	return nil
}

func (sc *scope) setFunc(name identifier, funcref *ExprFuncRef) {
	sc.set(name, &IdentBody{
		expr: funcref,
	})
}

func (sc *scope) setConst(name identifier, cnst *ExprConstVariable) {
	sc.set(name, &IdentBody{
		expr: cnst,
	})
}

func (sc *scope) setVar(name identifier, variable *ExprVariable) {
	sc.set(name, &IdentBody{
		expr: variable,
	})
}

func (sc *scope) setGtype(name identifier, gtype *Gtype) {
	sc.set(name, &IdentBody{
		gtype: gtype,
	})
}

func (sc *scope) set(name identifier, elm *IdentBody) {
	if elm == nil {
		panic("nil cannot be set")
	}
	sc.idents[name] = elm
}

func (sc *scope) getGtype(name identifier) *Gtype {
	elm, ok := sc.idents[name]
	if !ok {
		return nil
	}
	if elm.gtype == nil {
		errorf("type %s is not defined", name)
	}
	return elm.gtype
}

func newScope(outer *scope, name string) *scope {
	return &scope{
		outer:  outer,
		name:   name,
		idents: map[identifier]*IdentBody{},
	}
}
