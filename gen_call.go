package main

import "fmt"

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
	if origType.kind == G_INTERFACE {
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

func (methodCall *ExprMethodcall) emit() {
	origType := methodCall.getOrigType()
	if origType.kind == G_INTERFACE {
		methodCall.emitInterfaceMethodCall()
		return
	}

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
	var staticCall *IrStaticCall = &IrStaticCall{
		symbol:       getFuncSymbol(pkgname, name),
		callee:       funcref.funcdef,
		isMethodCall: true,
	}
	staticCall.emit(args)
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

func (funcall *ExprFuncallOrConversion) emit() {
	if funcall.rel.expr == nil && funcall.rel.gtype != nil {
		// Conversion
		conversion := &ExprConversion{
			tok:   funcall.token(),
			gtype: funcall.rel.gtype,
			expr:  funcall.args[0],
		}
		conversion.emit()
		return
	}

	assert(funcall.rel.expr != nil && funcall.rel.gtype == nil, funcall.token(), "this is conversion")
	assert(funcall.getFuncDef() != nil, funcall.token(), "funcdef is nil")
	decl := funcall.getFuncDef()

	// check if it's a builtin function
	switch decl {
	case builtinLen:
		assert(len(funcall.args) == 1, funcall.token(), "invalid arguments for len()")
		arg := funcall.args[0]
		exprLen := &ExprLen{
			tok: arg.token(),
			arg: arg,
		}
		exprLen.emit()
	case builtinCap:
		arg := funcall.args[0]
		e := &ExprCap{
			tok: arg.token(),
			arg: arg,
		}
		e.emit()
	case builtinAppend:
		assert(len(funcall.args) == 2, funcall.token(), "append() should take 2 argments")
		slice := funcall.args[0]
		valueToAppend := funcall.args[1]
		emit("# append(%s, %s)", slice.getGtype().String(), valueToAppend.getGtype().String())
		var staticCall *IrStaticCall = &IrStaticCall{
			callee: decl,
		}
		switch slice.getGtype().elementType.getSize() {
		case 1:
			staticCall.symbol = getFuncSymbol("iruntime", "append1")
			staticCall.emit(funcall.args)
		case 8:
			staticCall.symbol = getFuncSymbol("iruntime", "append8")
			staticCall.emit(funcall.args)
		case 24:
			if slice.getGtype().elementType.getKind() == G_INTERFACE && valueToAppend.getGtype().getKind() != G_INTERFACE {
				eConvertion := &ExprConversionToInterface{
					tok:  valueToAppend.token(),
					expr: valueToAppend,
				}
				funcall.args[1] = eConvertion
			}
			staticCall.symbol = getFuncSymbol("iruntime", "append24")
			staticCall.emit(funcall.args)
		default:
			TBI(slice.token(), "")
		}
	case builtinMakeSlice:
		assert(len(funcall.args) == 3, funcall.token(), "append() should take 3 argments")
		var staticCall *IrStaticCall = &IrStaticCall{
			callee: decl,
		}
		staticCall.symbol = getFuncSymbol("iruntime", "makeSlice")
		staticCall.emit(funcall.args)
	case builtinDumpSlice:
		arg := funcall.args[0]

		emit("lea .%s, %%rax", builtinStringKey2)
		emit("PUSH_8")

		arg.emit()
		emit("PUSH_SLICE")

		numRegs := 4
		for i := numRegs - 1; i >= 0; i-- {
			emit("POP_TO_ARG_%d", i)
		}

		emit("FUNCALL %s", "printf")
		emitNewline()
	case builtinDumpInterface:
		arg := funcall.args[0]

		emit("lea .%s, %%rax", builtinStringKey1)
		emit("PUSH_8")

		arg.emit()
		emit("PUSH_INTERFACE")

		numRegs := 4
		for i := numRegs - 1; i >= 0; i-- {
			emit("POP_TO_ARG_%d", i)
		}

		emit("FUNCALL %s", "printf")
		emitNewline()
	case builtinAssertInterface:
		emit("# builtinAssertInterface")
		labelEnd := makeLabel()
		arg := funcall.args[0]
		arg.emit() // rax=ptr, rbx=receverTypeId, rcx=dynamicTypeId

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

	case builtinAsComment:
		arg := funcall.args[0]
		if stringLiteral, ok := arg.(*ExprStringLiteral); ok {
			emitWithoutIndent("# %s", stringLiteral.val)
		}
	default:
		var staticCall *IrStaticCall = &IrStaticCall{
			symbol: getFuncSymbol(decl.pkg, funcall.fname),
			callee: decl,
		}
		staticCall.emit(funcall.args)
	}
}

type IrStaticCall struct {
	// https://sourceware.org/binutils/docs-2.30/as/Symbol-Intro.html#Symbol-Intro
	// A symbol is one or more characters chosen from the set of all letters (both upper and lower case), digits and the three characters ‘_.$’.
	symbol       string
	callee       *DeclFunc
	isMethodCall bool
}

