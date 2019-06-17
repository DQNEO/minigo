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

func walkFunc(fnc *DeclFunc) *DeclFunc {
	fnc.body = walkStmtList(fnc.body)
	return fnc
}

func walkStmtList(stmtList *StmtSatementList) *StmtSatementList {
	for _, stmt := range stmtList.stmts {
		stmt = walkStmt(stmt)
	}
	return stmtList
}

func walkStmt(stmt Stmt) Stmt {
	switch stmt.(type) {
	case nil:
		return nil
	case *StmtFor:
		s := stmt.(*StmtFor)
		return s.convert()
	case *ForRangeListEmitter:
		s := stmt.(*ForRangeListEmitter)
		s.block = walkStmtList(s.block)
	case *RangeMapEmitter:
		s := stmt.(*RangeMapEmitter)
		s.block = walkStmtList(s.block)
	case *PlainForEmitter:
		s := stmt.(*PlainForEmitter)
		s.block = walkStmtList(s.block)
	case *StmtIf:
		s := stmt.(*StmtIf)
		return s
	case *StmtReturn:
		s := stmt.(*StmtReturn)
		return s
	case *StmtInc:
		s := stmt.(*StmtInc)
		return s
	case *StmtDec:
		s := stmt.(*StmtDec)
		return s
	case *StmtSatementList:
		s := stmt.(*StmtSatementList)
		return s
	case *StmtAssignment:
		s := stmt.(*StmtAssignment)
		return s
	case *StmtShortVarDecl:
		s := stmt.(*StmtShortVarDecl)
		return s
	case *StmtContinue:
		s := stmt.(*StmtContinue)
		return s
	case *StmtBreak:
		s := stmt.(*StmtBreak)
		return s
	case *StmtExpr:
		s := stmt.(*StmtExpr)
		return s
	case *StmtDefer:
		s := stmt.(*StmtDefer)
		return s
	case *StmtSwitch:
		s := stmt.(*StmtSwitch)
		return s
	}

	return stmt
}

