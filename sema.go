// Semantic Analyzer to produce IR struct
package main

var symbolTable *SymbolTable

type SymbolTable struct {
	allScopes     map[packageName]*Scope
	uniquedDTypes []string
}

func makeDynamicTypeLabel(id int) gostring {
	return Sprintf(S("DynamicTypeId%d"), id)
}

func (symbolTable *SymbolTable) getTypeLabel(gtype *Gtype) gostring {
	dynamicTypeId := get_index(gtype.String(), symbolTable.uniquedDTypes)
	if dynamicTypeId == -1 {
		errorft(nil, "type %s not found in uniquedDTypes", gtype.String())
	}
	return makeDynamicTypeLabel(dynamicTypeId)
}

var typeId int = 1 // start with 1 because we want to zero as error

func setTypeIds(namedTypes []*DeclType) {
	for _, namedType := range namedTypes {
		namedType.gtype.receiverTypeId = typeId
		typeId++
	}
}

func resolve(sc *Scope, rel *Relation) *IdentBody {
	relbody := sc.get(rel.name)
	if relbody != nil {
		if relbody.gtype != nil {
			rel.gtype = relbody.gtype
		} else if relbody.expr != nil {
			rel.expr = relbody.expr
		} else {
			errorft(rel.token(), "Bad type relbody %v", relbody)
		}
	}
	return relbody
}

func resolveIdents(pkg *AstPackage, universe *Scope) {
	packageScope := pkg.scope
	packageScope.outer = universe
	for _, file := range pkg.files {
		for _, rel := range file.unresolved {
			relbody := resolve(packageScope, rel)
			if relbody == nil {
				errorft(rel.token(), "unresolved identifier %s", rel.name)
			}
		}
	}
}

// copy methods from p.nameTypes to gtype.methods of each type
func attachMethodsToTypes(pmethods map[identifier]methods, packageScope *Scope) {
	for typeName, methods := range pmethods {
		gtype := packageScope.getGtype(typeName)
		if gtype == nil {
			debugf("%#v", packageScope.idents)
			errorf("typaneme %s is not found in the package scope %s", typeName, packageScope.name)
		}
		gtype.methods = methods
	}
}

func collectDecls(pkg *AstPackage) {
	for _, f := range pkg.files {
		for _, decl := range f.DeclList {
			if decl.vardecl != nil {
				pkg.vars = append(pkg.vars, decl.vardecl)
			} else if decl.funcdecl != nil {
				pkg.funcs = append(pkg.funcs, decl.funcdecl)
			}
		}
	}
}

func setStringLables(pkg *AstPackage, prefix string) {
	for id, sl := range pkg.stringLiterals {
		var no int = id + 1
		sl.slabel = Sprintf(S("%s.S%d"), gostring(prefix), no)
	}
}

func calcStructSize(gtypes []*Gtype) {
	for _, gtype := range gtypes {
		if gtype.getKind() == G_STRUCT {
			gtype.calcStructOffset()
		} else if gtype.getKind() == G_POINTER && gtype.origType.getKind() == G_STRUCT {
			gtype.origType.calcStructOffset()
		}
	}
}

func uniqueDynamicTypes(dynamicTypes []*Gtype) []string {
	var r []string = builtinTypesAsString
	for _, gtype := range dynamicTypes {
		gs := gtype.String()
		if !in_array(gs, r) {
			r = append(r, gs)
		}
	}
	return r
}

func composeMethodTable(funcs []*DeclFunc) map[int][]string {
	var methodTable map[int][]string = map[int][]string{} // receiverTypeId : []methodTable

	for _, funcdecl := range funcs {
		if funcdecl.receiver == nil {
			continue
		}

		gtype := funcdecl.receiver.getGtype()
		if gtype.kind == G_POINTER {
			gtype = gtype.origType
		}
		if gtype.relation == nil {
			errorf("no relation for %#v", funcdecl.receiver.getGtype())
		}
		typeId := gtype.relation.gtype.receiverTypeId
		symbol := funcdecl.getSymbol()
		methods := methodTable[typeId]
		methods = append(methods, string(symbol))
		methodTable[typeId] = methods
	}
	debugf("set methodTable")
	return methodTable
}

func walkFunc(f *DeclFunc) *DeclFunc {
	f.prologue = f.prepare()
	f.body = walkStmtList(f.body)
	return f
}

func walkStmtList(stmtList *StmtSatementList) *StmtSatementList {
	if stmtList == nil {
		return nil
	}
	for i, stmt := range stmtList.stmts {
		stmt2 := walkStmt(stmt)
		stmtList.stmts[i] = stmt2
	}
	return stmtList
}

func walkExpr(expr Expr) Expr {
	var r Expr
	switch expr.(type) {
	case nil:
		return r
	case *Relation:
		e := expr.(*Relation)
		return e.expr
	case *ExprNilLiteral:
	case *ExprNumberLiteral:
	case *ExprStringLiteral:
	case *ExprVariable:
	case *ExprConstVariable:
	case *ExprFuncallOrConversion:
		funcall := expr.(*ExprFuncallOrConversion)
		for i := 0; i < len(funcall.args); i++ {
			arg := funcall.args[i]
			arg = walkExpr(arg)
			funcall.args[i] = arg
		}
		if funcall.rel.expr == nil && funcall.rel.gtype != nil {
			// Conversion
			r = &IrExprConversion{
				tok:     funcall.token(),
				toGtype: funcall.rel.gtype,
				arg:     funcall.args[0],
			}
			return r
		}
		decl := funcall.getFuncDef()
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
		case builtinMakeSlice:
			assert(len(funcall.args) == 3, funcall.token(), "append() should take 3 argments")
			var staticCall *IrStaticCall = &IrStaticCall{
				tok:      funcall.token(),
				origExpr: funcall,
				callee:   decl,
			}
			staticCall.symbol = getFuncSymbol(S("iruntime"), S("makeSlice"))
			staticCall.args = funcall.args
			return staticCall
		case builtinAppend:
			assert(len(funcall.args) == 2, funcall.token(), "append() should take 2 argments")
			slice := funcall.args[0]
			valueToAppend := funcall.args[1]
			emit("# append(%s, %s)", gostring(slice.getGtype().String()), gostring(valueToAppend.getGtype().String()))
			var staticCall *IrStaticCall = &IrStaticCall{
				tok:      funcall.token(),
				origExpr: funcall,
				callee:   decl,
			}
			switch slice.getGtype().elementType.getSize() {
			case 1:
				staticCall.symbol = getFuncSymbol(S("iruntime"), S("append1"))
			case 8:
				staticCall.symbol = getFuncSymbol(S("iruntime"), S("append8"))
			case 24:
				if slice.getGtype().elementType.getKind() == G_INTERFACE && valueToAppend.getGtype().getKind() != G_INTERFACE {
					eConvertion := &IrExprConversionToInterface{
						tok: valueToAppend.token(),
						arg: valueToAppend,
					}
					funcall.args[1] = eConvertion
				}
				staticCall.symbol = getFuncSymbol(S("iruntime"), S("append24"))
			default:
				TBI(slice.token(), "")
			}
			staticCall.args = funcall.args
			return staticCall
		}
		return funcall
	case *ExprMethodcall:
		methodCall := expr.(*ExprMethodcall)
		for i := 0; i < len(methodCall.args); i++ {
			arg := methodCall.args[i]
			arg = walkExpr(arg)
			methodCall.args[i] = arg
		}
		methodCall.receiver = walkExpr(methodCall.receiver)
		expr = methodCall
		return expr
	case *ExprBinop:
		e := expr.(*ExprBinop)
		e.left = walkExpr(e.left)
		e.right = walkExpr(e.right)
		return e
	case *ExprUop:
		e := expr.(*ExprUop)
		e.operand = walkExpr(e.operand)
		return e
	case *ExprFuncRef:
	case *ExprSlice:
		e := expr.(*ExprSlice)
		e.collection = walkExpr(e.collection)
		e.low = walkExpr(e.low)
		e.high = walkExpr(e.high)
		e.max = walkExpr(e.max)
		return e
	case *ExprIndex:
		e := expr.(*ExprIndex)
		e.index = walkExpr(e.index)
		e.collection = walkExpr(e.collection)
		return e
	case *ExprArrayLiteral:
	case *ExprSliceLiteral:
		e := expr.(*ExprSliceLiteral)
		for i, v := range e.values {
			v2 := walkExpr(v)
			e.values[i] = v2
		}
		return e
	case *ExprTypeAssertion:
	case *ExprVaArg:
		e := expr.(*ExprVaArg)
		e.expr = walkExpr(e.expr)
		return e
	case *ExprStructLiteral:
		e := expr.(*ExprStructLiteral)
		for _, field := range e.fields {
			field.value = walkExpr(field.value)
		}
		return e
	case *ExprStructField:
		e := expr.(*ExprStructField)
		e.strct = walkExpr(e.strct)
		return e
	case *ExprTypeSwitchGuard:
	case *ExprMapLiteral:
	case *ExprLen:

	case *ExprCap:
		//case *ExprConversionToInterface:
	}
	return expr
}

func walkStmt(stmt Stmt) Stmt {
	var s2 Stmt
	switch stmt.(type) {
	case nil:
		return s2
	case *DeclVar:
		s := stmt.(*DeclVar)
		s.initval = walkExpr(s.initval)
		return s
	case *StmtFor:
		s := stmt.(*StmtFor)
		s2 = s.convert()
		s2 = walkStmt(s2)
		return s2
	case *IrStmtRangeMap:
		s := stmt.(*IrStmtRangeMap)
		s.rangeexpr = walkExpr(s.rangeexpr)
		s.block = walkStmtList(s.block)
		return s
	case *IrStmtForRangeList:
		s := stmt.(*IrStmtForRangeList)
		s.init = walkStmt(s.init)
		s.cond = walkExpr(s.cond)
		s.block = walkStmtList(s.block)
		return s
	case *IrStmtClikeFor:
		s := stmt.(*IrStmtClikeFor)
		cls := s.cls
		cls.init = walkStmt(cls.init)
		cls.cond = walkStmt(cls.cond)
		cls.post = walkStmt(cls.post)
		s.block = walkStmtList(s.block)
		return s
	case *StmtIf:
		s := stmt.(*StmtIf)
		s.simplestmt = walkStmt(s.simplestmt)
		s.cond = walkExpr(s.cond)
		s.then = walkStmt(s.then)
		s.els = walkStmt(s.els)
		return s
	case *StmtReturn:
		s := stmt.(*StmtReturn)
		for i, expr := range s.exprs {
			e := walkExpr(expr)
			s.exprs[i] = e
		}
		return s
	case *StmtInc:
		s := stmt.(*StmtInc)
		s.operand = walkExpr(s.operand)
		return s
	case *StmtDec:
		s := stmt.(*StmtDec)
		s.operand = walkExpr(s.operand)
		return s
	case *StmtSatementList:
		s := stmt.(*StmtSatementList)
		s = walkStmtList(s)
		return s
	case *StmtAssignment:
		s := stmt.(*StmtAssignment)
		for i, right := range s.rights {
			right = walkExpr(right)
			s.rights[i] = right
		}

		for i, left := range s.lefts {
			left = walkExpr(left)
			s.lefts[i] = left
		}

		return s
	case *StmtShortVarDecl:
		s := stmt.(*StmtShortVarDecl)
		var s2 Stmt = &StmtAssignment{
			tok:    s.tok,
			lefts:  s.lefts,
			rights: s.rights,
		}
		s2 = walkStmt(s2)
		return s2
	case *StmtContinue:
		s := stmt.(*StmtContinue)
		return s
	case *StmtBreak:
		s := stmt.(*StmtBreak)
		return s
	case *StmtExpr:
		s := stmt.(*StmtExpr)
		s.expr = walkExpr(s.expr)
		return s
	case *StmtDefer:
		s := stmt.(*StmtDefer)
		s.expr = walkExpr(s.expr)
		return s
	case *StmtSwitch:
		s := stmt.(*StmtSwitch)
		s.cond = walkExpr(s.cond)
		for _, cse := range s.cases {
			cse.compound = walkStmtList(cse.compound)
		}
		s.dflt = walkStmtList(s.dflt)
		return s
	}

	return stmt
}
