package main

import "fmt"

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
		assert(typeToBeloing != nil, methodCall.token(), "shoudl not be nil:"+gtype.String())
	} else {
		typeToBeloing = gtype
	}
	assert(typeToBeloing.kind == G_NAMED, methodCall.tok, "method must belong to a named type")
	origType := typeToBeloing.relation.gtype
	assert(typeToBeloing.relation.gtype != nil, methodCall.token(), fmt.Sprintf("origType should not be nil:%#v", typeToBeloing.relation))
	return origType
}


func (methodCall *ExprMethodcall) getRettypes() []*Gtype {
	origType := methodCall.getOrigType()
	if origType == nil {
		errorft(methodCall.token(), "origType should not be nil")
	}
	if origType.getKind() == G_INTERFACE {
		return origType.imethods[methodCall.fname].rettypes
	} else {
		funcref, ok := origType.methods[methodCall.fname]
		if !ok {
			errorft(methodCall.token(), "method %s is not found in type %s", methodCall.fname, methodCall.receiver.getGtype().String())
		}
		return funcref.funcdef.rettypes
	}
}

type IrInterfaceMethodCall struct {
	receiver   Expr
	methodName identifier
}

func (methodCall *ExprMethodcall) emitInterfaceMethodCall() {
	args := []Expr{methodCall.receiver}
	for _, arg := range methodCall.args {
		args = append(args, arg)
	}
	call := &IrInterfaceMethodCall{
		receiver:   methodCall.receiver,
		methodName: methodCall.fname,
	}
	call.emit(args)
}

func (methodCall *ExprMethodcall) emitDynamicTypeMethodCall() {
	origType := methodCall.getOrigType()
	args := []Expr{methodCall.receiver}
	for _, arg := range methodCall.args {
		args = append(args, arg)
	}

	funcref, ok := origType.methods[methodCall.fname]
	if !ok {
		errorft(methodCall.token(), "method %s is not found in type %s", methodCall.fname, methodCall.receiver.getGtype().String())
	}
	pkgname := funcref.funcdef.pkg
	name := methodCall.getUniqueName()
	var staticCall Expr = &IrStaticCall{
		tok: methodCall.token(),
		symbol:       getFuncSymbol(pkgname, name),
		callee:       funcref.funcdef,
		isMethodCall: true,
		args:args,
		origExpr:methodCall,
	}
	staticCall.emit()
}

func (methodCall *ExprMethodcall) emit() {
	origType := methodCall.getOrigType()
	if origType.getKind() == G_INTERFACE {
		methodCall.emitInterfaceMethodCall()
		return
	}
	methodCall.emitDynamicTypeMethodCall()
}

func (funcall *ExprFuncallOrConversion) getFuncDef() *DeclFunc {
	relexpr := funcall.rel.expr
	assert(relexpr != nil, funcall.token(), fmt.Sprintf("relexpr should NOT be nil for %s", funcall.fname))
	funcref, ok := relexpr.(*ExprFuncRef)
	if !ok {
		errorft(funcall.token(), "Compiler error: funcref is not *ExprFuncRef (%s)", funcall.fname)
	}
	assertNotNil(funcref.funcdef != nil, nil)
	return funcref.funcdef
}

func funcall2emitter(funcall *ExprFuncallOrConversion) Emitter {
	if funcall.rel.expr == nil && funcall.rel.gtype != nil {
		// Conversion
		return &IrExprConversion{
			tok:   funcall.token(),
			gtype: funcall.rel.gtype,
			expr:  funcall.args[0],
		}
	}

	assert(funcall.rel.expr != nil && funcall.rel.gtype == nil, funcall.token(), "this is conversion")
	assert(funcall.getFuncDef() != nil, funcall.token(), "funcdef is nil")
	decl := funcall.getFuncDef()

	// check if it's a builtin function
	switch decl {
	case builtinLen:
		assert(len(funcall.args) == 1, funcall.token(), "invalid arguments for len()")
		arg := funcall.args[0]
		return &ExprLen{
			tok: arg.token(),
			arg: arg,
		}
	case builtinCap:
		arg := funcall.args[0]
		return &ExprCap{
			tok: arg.token(),
			arg: arg,
		}
	case builtinAppend:
		assert(len(funcall.args) == 2, funcall.token(), "append() should take 2 argments")
		slice := funcall.args[0]
		valueToAppend := funcall.args[1]
		emit("# append(%s, %s)", slice.getGtype().String(), valueToAppend.getGtype().String())
		var symbol string
		switch slice.getGtype().elementType.getSize() {
		case 1:
			symbol = getFuncSymbol("iruntime", "append1")
		case 8:
			symbol = getFuncSymbol("iruntime", "append8")
		case 24:
			if slice.getGtype().elementType.getKind() == G_INTERFACE && valueToAppend.getGtype().getKind() != G_INTERFACE {
				eConvertion := &IrExprConversionToInterface{
					tok:  valueToAppend.token(),
					expr: valueToAppend,
				}
				funcall.args[1] = eConvertion
			}
			symbol = getFuncSymbol("iruntime", "append24")
		default:
			TBI(slice.token(), "")
		}
		return &IrStaticCall{
			tok: funcall.token(),
			callee: decl,
			args: funcall.args,
			origExpr: funcall,
			symbol: symbol,
		}
	case builtinMakeSlice:
		assert(len(funcall.args) == 3, funcall.token(), "append() should take 3 argments")
		return &IrStaticCall{
			tok: funcall.token(),
			callee: decl,
			args: funcall.args,
			origExpr: funcall,
			symbol: getFuncSymbol("iruntime", "makeSlice"),
		}
	case builtinDumpSlice:
		arg := funcall.args[0]
		return &builtinDumpSliceEmitter{
			arg: arg,
		}
	case builtinDumpInterface:
		arg := funcall.args[0]
		return &builtinDumpInterfaceEmitter{
			arg:arg,
		}
	case builtinAssertInterface:
		arg := funcall.args[0]
		return &builtinAssertInterfaceEmitter{
			arg: arg,
		}
	case builtinAsComment:
		arg := funcall.args[0]
		return &builtinAsCommentEmitter{
			arg:arg,
		}
	default:
		return &IrStaticCall{
			tok: funcall.token(),
			symbol: getFuncSymbol(decl.pkg, funcall.fname),
			callee: decl,
			args: funcall.args,
			origExpr:funcall,
		}
	}

}

func (funcall *ExprFuncallOrConversion) emit() {
	e := funcall2emitter(funcall)
	e.emit()
}

type IrStaticCall struct {
	// https://sourceware.org/binutils/docs-2.30/as/Symbol-Intro.html#Symbol-Intro
	// A symbol is one or more characters chosen from the set of all letters (both upper and lower case), digits and the three characters ‘_.$’.
	tok          *Token
	symbol       string
	callee       *DeclFunc
	isMethodCall bool
	args         []Expr
	origExpr     Expr
}

func (ircall *IrStaticCall) token() *Token {
	return ircall.tok
}

func (ircall *IrStaticCall) dump() {
	ircall.origExpr.dump()
}

func (ircall *IrStaticCall) getGtype() *Gtype {
	return ircall.origExpr.getGtype()
}

type builtinDumpSliceEmitter struct {
	arg Expr
}

func (em *builtinDumpSliceEmitter) emit() {
	emit("lea .%s, %%rax", builtinStringKey2)
	emit("PUSH_8")

	em.arg.emit()
	emit("PUSH_SLICE")

	numRegs := 4
	for i := numRegs - 1; i >= 0; i-- {
		emit("POP_TO_ARG_%d", i)
	}

	emit("FUNCALL %s", "printf")
	emitNewline()
}

type builtinDumpInterfaceEmitter struct {
	arg Expr
}

func (em *builtinDumpInterfaceEmitter) emit() {
	emit("lea .%s, %%rax", builtinStringKey1)
	emit("PUSH_8")

	em.arg.emit()
	emit("PUSH_INTERFACE")

	numRegs := 4
	for i := numRegs - 1; i >= 0; i-- {
		emit("POP_TO_ARG_%d", i)
	}

	emit("FUNCALL %s", "printf")
	emitNewline()
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
	emit("FUNCALL %s", ".panic")

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
