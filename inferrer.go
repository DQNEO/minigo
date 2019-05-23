package main

// Inferer infers types
type Inferrer interface {
	infer()
}

func inferTypes(globals []*ExprVariable, locals []Inferrer) {
	//debugf("infering globals")
	for _, variable := range globals {
		variable.infer()
	}
	//debugf("infering locals")
	for _, ast := range locals {
		ast.infer()
	}
}

//  infer recursively all the types of global variables
func (variable *ExprVariable) infer() {
	//debugf("infering ExprVariable")
	if variable.gtype.kind != G_DEPENDENT {
		// done
		return
	}
	e := variable.gtype.dependendson
	dependType := e.getGtype()
	if dependType.kind != G_DEPENDENT {
		variable.gtype = dependType
		return
	}

	rel, ok := e.(*Relation)
	if !ok {
		errorft(e.token(), "unexpected type %T", e)
	}
	vr, ok := rel.expr.(*ExprVariable)
	vr.infer() // recursive call
	variable.gtype = e.getGtype()
	//debugf("infered type=%s", variable.gtype)
}

// local decl infer
func (decl *DeclVar) infer() {
	//debugf("infering DeclVar")
	gtype := decl.initval.getGtype()
	assertNotNil(gtype != nil, decl.initval.token())
	decl.variable.gtype = gtype
}

func (clause *ForRangeClause) infer() {
	//debugf("infering ForRangeClause")
	collectionType := clause.rangeexpr.getGtype()
	//debugf("collectionType = %s", collectionType)
	indexvar, ok := clause.indexvar.expr.(*ExprVariable)
	assert(ok, nil, "ok")

	var indexType *Gtype
	switch collectionType.kind {
	case G_ARRAY, G_SLICE:
		indexType = gInt
	case G_MAP:
		indexType = collectionType.mapKey
	default:
		// @TODO consider map etc.
		TBI(clause.tok, "unable to handle %s", collectionType)
	}
	indexvar.gtype = indexType

	if clause.valuevar != nil {
		valuevar, ok := clause.valuevar.expr.(*ExprVariable)
		assert(ok, nil, "ok")

		var elementType *Gtype
		if collectionType.kind == G_ARRAY {
			elementType = collectionType.elementType
		} else if collectionType.kind == G_SLICE {
			elementType = collectionType.elementType
		} else if collectionType.kind == G_MAP {
			elementType = collectionType.mapValue
		} else {
			errorft(clause.token(), "internal error")
		}
		//debugf("for i, v %s := rannge %v", elementType, collectionType)
		valuevar.gtype = elementType
	}
}

func (ast *StmtShortVarDecl) infer() {
	//debugf("infering StmtShortVarDecl")
	var rightTypes []*Gtype
	for _, rightExpr := range ast.rights {
		switch rightExpr.(type) {
		case *ExprFuncallOrConversion:
			fcallOrConversion := rightExpr.(*ExprFuncallOrConversion)
			if fcallOrConversion.rel.gtype != nil {
				// Conversion
				rightTypes = append(rightTypes, fcallOrConversion.rel.gtype)
			} else {
				fcall := fcallOrConversion
				funcdef := fcall.getFuncDef()
				if funcdef == nil {
					errorft(fcall.token(), "funcdef of %s is not found", fcall.fname)
				}
				if funcdef == builtinLen {
					rightTypes = append(rightTypes, gInt)
				} else {
					for _, gtype := range fcall.getFuncDef().rettypes {
						rightTypes = append(rightTypes, gtype)
					}
				}
			}
		case *ExprMethodcall:
			fcall := rightExpr.(*ExprMethodcall)
			rettypes := fcall.getRettypes()
			for _, gtype := range rettypes {
				rightTypes = append(rightTypes, gtype)
			}
		case *ExprTypeAssertion:
			assertion := rightExpr.(*ExprTypeAssertion)
			rightTypes = append(rightTypes, assertion.gtype)
			rightTypes = append(rightTypes, gBool)
		case *ExprIndex:
			e := rightExpr.(*ExprIndex)
			gtype := e.getGtype()
			assertNotNil(gtype != nil, e.tok)
			rightTypes = append(rightTypes, gtype)
			//debugf("rightExpr.gtype=%s", gtype)
			secondGtype := rightExpr.(*ExprIndex).getSecondGtype()
			if secondGtype != nil {
				rightTypes = append(rightTypes, secondGtype)
			}
		default:
			if rightExpr == nil {
				errorft(ast.token(), "rightExpr is nil")
			}
			gtype := rightExpr.getGtype()
			if gtype == nil {
				errorft(ast.token(), "rightExpr %T gtype is nil", rightExpr)
			}
			//debugf("infered type %s", gtype)
			rightTypes = append(rightTypes, gtype)
		}
	}

	if len(ast.lefts) > len(rightTypes) {
		// @TODO this check is too loose.
		errorft(ast.tok, "number of lhs and rhs does not match (%d <=> %d)", len(ast.lefts), len(rightTypes))
	}
	for i, e := range ast.lefts {
		rel := e.(*Relation) // a brand new rel
		variable := rel.expr.(*ExprVariable)
		rightType := rightTypes[i]
		variable.gtype = rightType
	}

}
