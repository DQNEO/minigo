package main


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
	sc.set(name, funcref)
}

func (sc *scope) setConst(name identifier, cnst *ExprConstVariable) {
	sc.set(name, cnst)
}

func (sc *scope) setVar(name identifier, variable *ExprVariable) {
	sc.set(name, variable)
}

func (sc *scope) setGtype(name identifier, gtype *Gtype) {
	sc.set(name, gtype)
}

func (sc *scope) set(name identifier, v interface{}) {
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
