package main

type Emitter interface {
	emit()
}

func (funcall *ExprFuncallOrConversion) getRettypes() []*Gtype {
	if funcall.rel.gtype != nil {
		// Conversion
		return []*Gtype{funcall.rel.gtype}
	}

	return funcall.getFuncDef().rettypes
}

func (ast *ExprMethodcall) getUniqueName() string {
	gtype := ast.receiver.getGtype()
	return getMethodUniqueName(gtype, ast.fname)
}

func (methodCall *ExprMethodcall) getOrigType() *Gtype {
	gtype := methodCall.receiver.getGtype()
	assertNotNil(methodCall.receiver != nil, methodCall.token())
	assertNotNil(gtype != nil, methodCall.tok)
	assert(gtype.kind == G_NAMED || gtype.kind == G_POINTER || gtype.kind == G_INTERFACE, methodCall.tok, "method must be an interface or belong to a named type")
	var typeToBeloing *Gtype
	if gtype.kind == G_POINTER {
		typeToBeloing = gtype.origType
		assert(typeToBeloing != nil, methodCall.token(), "shoudl not be nil:%s", gtype.String())
	} else {
		typeToBeloing = gtype
	}
	assert(typeToBeloing.kind == G_NAMED, methodCall.tok, "method must belong to a named type")
	origType := typeToBeloing.relation.gtype
	assert(typeToBeloing.relation.gtype != nil, methodCall.token(), "origType should not be nil")
	return origType
}

func (methodCall *ExprMethodcall) getRettypes() []*Gtype {
	origType := methodCall.getOrigType()
	if origType == nil {
		errorft(methodCall.token(), "origType should not be nil")
	}

	if origType.getKind() == G_INTERFACE {
		imethod, _ := imethodGet(origType.imethods, methodCall.fname)
		return imethod.rettypes
	} else {
		funcref, ok := methodGet(origType.methods, methodCall.fname)
		if !ok {
			errorft(methodCall.token(), "method %s is not found in type %s", methodCall.fname, methodCall.receiver.getGtype().String())
		}
		return funcref.funcdef.rettypes
	}
}

type IrInterfaceMethodCall struct {
	receiver   Expr
	methodName identifier
	args       []Expr
	callee     *signature
}

func (methodCall *ExprMethodcall) interfaceMethodCall() Emitter {
	iType := methodCall.getOrigType()
	sig, ok := iType.imethods[methodCall.fname]
	if !ok {
		panic("signature not found in the interface")
	}
	call := &IrInterfaceMethodCall{
		receiver:   methodCall.receiver,
		methodName: methodCall.fname,
		args:       methodCall.args,
		callee:sig,
	}
	return call
}

func (methodCall *ExprMethodcall) dynamicTypeMethodCall() Emitter {
	origType := methodCall.getOrigType()
	funcref, ok := methodGet(origType.methods, methodCall.fname)
	if !ok {
		errorft(methodCall.token(), "method %s is not found in type %s", methodCall.fname, methodCall.receiver.getGtype().String())
	}

	name := methodCall.getUniqueName()
	var staticCall Expr = &IrCall{
		tok:          methodCall.token(),
		symbol:       getFuncSymbol(funcref.funcdef.pkgPath, name),
		callee:       funcref.funcdef,
		receiver:     methodCall.receiver,
		args:         methodCall.args,
		origExpr:     methodCall,
	}
	return staticCall
}

func (methodCall *ExprMethodcall) emit() {
	origType := methodCall.getOrigType()
	var e Emitter
	if origType.getKind() == G_INTERFACE {
		e = methodCall.interfaceMethodCall()
	} else {
		e = methodCall.dynamicTypeMethodCall()
	}

	e.emit()
}

func (funcall *ExprFuncallOrConversion) getFuncDef() *DeclFunc {
	relexpr := funcall.rel.expr
	assert(relexpr != nil, funcall.token(), "relexpr should NOT be nil")
	funcref, ok := relexpr.(*ExprFuncRef)
	if !ok {
		errorft(funcall.token(), "Compiler error: funcref is not *ExprFuncRef (%s)", funcall.fname)
	}
	assertNotNil(funcref.funcdef != nil, nil)
	return funcref.funcdef
}

func funcall2emitter(funcall *ExprFuncallOrConversion) Emitter {

	assert(funcall.rel.expr != nil && funcall.rel.gtype == nil, funcall.token(), "this is conversion")
	assert(funcall.getFuncDef() != nil, funcall.token(), "funcdef is nil")
	declFunc := funcall.getFuncDef()

	// check if it's a builtin function
	switch declFunc {
	case builtinDumpSlice:
		arg := funcall.args[0]
		return &builtinDumpSliceEmitter{
			arg: arg,
		}
	case builtinDumpInterface:
		arg := funcall.args[0]
		return &builtinDumpInterfaceEmitter{
			arg: arg,
		}
	case builtinAssertInterface:
		arg := funcall.args[0]
		return &builtinAssertInterfaceEmitter{
			arg: arg,
		}
	case builtinAsComment:
		arg := funcall.args[0]
		return &builtinAsCommentEmitter{
			arg: arg,
		}
	default:
		return &IrCall{
			tok:      funcall.token(),
			symbol:   declFunc.getSymbol(),
			callee:   declFunc,
			args:     funcall.args,
			origExpr: funcall,
		}
	}

}

func (funcall *ExprFuncallOrConversion) emit() {
	e := funcall2emitter(funcall)
	e.emit()
}

type IrCall struct {
	// https://sourceware.org/binutils/docs-2.30/as/Symbol-Intro.html#Symbol-Intro
	// A symbol is one or more characters chosen from the set of all letters (both upper and lower case), digits and the three characters ‘_.$’.
	tok          *Token
	isInterfaceMethodCall bool
	symbol       string
	icallee      *signature
	callee       *DeclFunc
	receiver     Expr
	args         []Expr
	origExpr     Expr
}

func (ircall *IrCall) token() *Token {
	return ircall.tok
}

func (ircall *IrCall) dump() {
	ircall.origExpr.dump()
}

func (ircall *IrCall) getGtype() *Gtype {
	return ircall.origExpr.getGtype()
}

type builtinDumpSliceEmitter struct {
	arg Expr
}

func (em *builtinDumpSliceEmitter) emit() {
	// @TODO Implement me
	emit("# builtinDumpSliceEmitter is not implemented")
}

type builtinDumpInterfaceEmitter struct {
	arg Expr
}

func (em *builtinDumpInterfaceEmitter) emit() {
	// @TODO Implement me
	emit("# builtinDumpInterfaceEmitter is not implemented")
}

type builtinAssertInterfaceEmitter struct {
	arg Expr
}

func (em *builtinAssertInterfaceEmitter) emit() {
	emit("# builtinAssertInterface")
	labelEnd := makeLabel()
	em.arg.emit() // rax=ptr, rbx=receverTypeId, rcx=dynamicTypeId

	// (ptr != nil && rcx == nil) => Error

	emit("CMP_NE_ZERO")
	emit("TEST_IT")
	emit("je %s", labelEnd)

	emit("mov %%rcx, %%rax")

	emit("CMP_EQ_ZERO")
	emit("TEST_IT")
	emit("je %s", labelEnd)

	slabel := makeLabel()
	emit(".data 0")
	emitWithoutIndent("%s:", slabel)
	emit(".string \"%s\"", "assertInterface failed")
	emit(".text")
	emit("lea %s, %%rax", slabel)
	emit("PUSH_8")
	emit("POP_TO_ARG_0")
	emit("FUNCALL %s", getFuncSymbol(IRuntimePath, "panic"))

	emitWithoutIndent("%s:", labelEnd)
	emitNewline()
}

type builtinAsCommentEmitter struct {
	arg Expr
}

func (em *builtinAsCommentEmitter) emit() {
	if stringLiteral, ok := em.arg.(*ExprStringLiteral); ok {
		emitWithoutIndent("# %s", stringLiteral.val)
	}
}
