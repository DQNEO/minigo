// Semantic Analyzer to produce IR struct
package main

import "fmt"

var symbolTable *SymbolTable

type SymbolTable struct {
	allScopes map[identifier]*Scope
	uniquedDTypes []string
}


func makeDynamicTypeLabel(id int) string {
	return fmt.Sprintf("DynamicTypeId%d", id)
}

func (symbolTable *SymbolTable) getTypeLabel(gtype *Gtype) string {
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
		sl.slabel = fmt.Sprintf("%s.S%d", prefix, id+1)
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
		methods = append(methods, symbol)
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
	if expr == nil {
		return nil
	}

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
		for i:=0;i<len(funcall.args);i++ {
			arg := funcall.args[i]
			arg = walkExpr(arg)
			funcall.args[i] = arg
		}
		if funcall.rel.expr == nil && funcall.rel.gtype != nil {
			// Conversion
			r = &IrExprConversion{
				tok:   funcall.token(),
				gtype: funcall.rel.gtype,
				expr:  funcall.args[0],
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
				tok: funcall.token(),
				origExpr:funcall,
				callee: decl,
			}
			staticCall.symbol = getFuncSymbol("iruntime", "makeSlice")
			staticCall.args = funcall.args
			return staticCall
		case builtinAppend:
			assert(len(funcall.args) == 2, funcall.token(), "append() should take 2 argments")
			slice := funcall.args[0]
			valueToAppend := funcall.args[1]
			emit("# append(%s, %s)", slice.getGtype().String(), valueToAppend.getGtype().String())
			var staticCall *IrStaticCall = &IrStaticCall{
				tok: funcall.token(),
				origExpr:funcall,
				callee: decl,
			}
			switch slice.getGtype().elementType.getSize() {
			case 1:
				staticCall.symbol = getFuncSymbol("iruntime", "append1")
			case 8:
				staticCall.symbol = getFuncSymbol("iruntime", "append8")
			case 24:
				if slice.getGtype().elementType.getKind() == G_INTERFACE && valueToAppend.getGtype().getKind() != G_INTERFACE {
					eConvertion := &IrExprConversionToInterface{
						tok:  valueToAppend.token(),
						expr: valueToAppend,
					}
					funcall.args[1] = eConvertion
				}
				staticCall.symbol = getFuncSymbol("iruntime", "append24")
			default:
				TBI(slice.token(), "")
			}
			staticCall.args = funcall.args
			return staticCall
		}
		return funcall
	case *ExprMethodcall:
		return expr
		/*
		methodCall,_ := expr.(*ExprMethodcall)
		methodCall.receiver = walkExpr(methodCall.receiver)
		for i:=0 ;i<len(methodCall.args); i++ {
			arg := methodCall.args[i]
			arg = walkExpr(arg)
			methodCall.args[i] = arg
		}

		if methodCall.getOrigType().getKind() == G_INTERFACE {
			return methodCall.toInterfaceMethodCall()
		}
		return methodCall

		 */
	case *ExprBinop:
		e,_ := expr.(*ExprBinop)
		e.left = walkExpr(e.left)
		e.right = walkExpr(e.right)
		return e
	case *ExprUop:
		e,_ := expr.(*ExprUop)
		e.operand = walkExpr(e.operand)
		return e
	case *ExprFuncRef:
	case *ExprSlice:
		e,_ := expr.(*ExprSlice)
		e.collection = walkExpr(e.collection)
		e.low = walkExpr(e.low)
		e.high = walkExpr(e.high)
		e.max = walkExpr(e.max)
		return e
	case *ExprIndex:
		e,_ := expr.(*ExprIndex)
		e.index = walkExpr(e.index)
		e.collection = walkExpr(e.collection)
		return e
	case *ExprArrayLiteral:
	case *ExprSliceLiteral:
	case *ExprTypeAssertion:
	case *ExprVaArg:
		e,_ := expr.(*ExprVaArg)
		e.expr = walkExpr(e.expr)
		return e
		/*
	case *ExprConversion:
		e,_ := expr.(*ExprConversion)
		e.expr = walkExpr(e.expr)
		return e
		 */
	case *ExprStructLiteral:
		e,_ := expr.(*ExprStructLiteral)
		for _, field := range e.fields {
			field.value = walkExpr(field.value)
		}
		return e
	case *ExprStructField:
		e,_ := expr.(*ExprStructField)
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
	case *StmtFor:
		s := stmt.(*StmtFor)
		s2 = s.convert()
		s2 = walkStmt(s2)
		return s2
	case *RangeMapEmitter:
		s := stmt.(*RangeMapEmitter)
		s.rangeexpr = walkExpr(s.rangeexpr)
		s.block  = walkStmtList(s.block)
		s2 = s
		return s2
	case *ForRangeListEmitter:
		s := stmt.(*ForRangeListEmitter)
		s.cond = walkExpr(s.cond)
		s.block  = walkStmtList(s.block)
		s2 = s
		return s2
	case *PlainForEmitter:
		s := stmt.(*PlainForEmitter)
		cls := s.cls
		cls.init = walkStmt(cls.init)
		cls.cond = walkStmt(cls.cond)
		cls.post = walkStmt(cls.post)
		s.block  = walkStmtList(s.block)
		s2 = s
		return s2
	case *StmtIf:
		s := stmt.(*StmtIf)
		s.simplestmt = walkStmt(s.simplestmt)
		s.cond = walkExpr(s.cond)
		s.then = walkStmt(s.then)
		s.els = walkStmt(s.els)
		s2 = s
		return s2
	case *StmtReturn:
		s := stmt.(*StmtReturn)
		for i, expr := range s.exprs {
			e := walkExpr(expr)
			s.exprs[i] = e
		}
		s2 = s
		return s2
	case *StmtInc:
		s := stmt.(*StmtInc)
		s.operand = walkExpr(s.operand)
		s2 = s
		return s2
	case *StmtDec:
		s := stmt.(*StmtDec)
		s.operand = walkExpr(s.operand)
		s2 = s
		return s2
	case *StmtSatementList:
		s := stmt.(*StmtSatementList)
		s = walkStmtList(s)
		s2 = s
		return s2
	case *StmtAssignment:
		s := stmt.(*StmtAssignment)
		for i, right := range s.rights {
			right = walkExpr(right)
			s.rights[i] = right
		}
		/*
		for i, left := range s.lefts {
			left = walkExpr(left)
			s.lefts[i] = left
		}
			 */
		s2 = s
		return s2
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
		s2 = s
		return s2
	case *StmtBreak:
		s := stmt.(*StmtBreak)
		s2 = s
		return s2
	case *StmtExpr:
		s := stmt.(*StmtExpr)
		s.expr = walkExpr(s.expr)
		s2 = s
		return s2
	case *StmtDefer:
		s := stmt.(*StmtDefer)
		s.expr = walkExpr(s.expr)
		s2 = s
		return s2
	case *StmtSwitch:
		s := stmt.(*StmtSwitch)
		for _, cse := range s.cases {
			cse.compound = walkStmtList(cse.compound)
		}
		s.dflt = walkStmtList(s.dflt)
		s2 = s
		return s2
	}

	return stmt
}

