package main

// Inferer infers types
type Inferrer interface {
	infer()
}

func inferTypes(globals []*ExprVariable, locals []Inferrer) {
	//debugf(S("infering globals"))
	for _, variable := range globals {
		variable.infer()
	}
	//debugf(S("infering locals"))
	for _, ast := range locals {
		ast.infer()
	}
}

//  infer recursively all the types of global variables
func (variable *ExprVariable) infer() {
	//debugf(S("infering ExprVariable"))
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
		errorft(e.token(), S("unexpected type %T"), e)
	}
	vr, ok := rel.expr.(*ExprVariable)
	vr.infer() // recursive call
	variable.gtype = e.getGtype()
	//debugf(S("infered type=%s"), variable.gtype)
}

// local decl infer
func (decl *DeclVar) infer() {
	//debugf(S("infering DeclVar"))
	gtype := decl.initval.getGtype()
	assertNotNil(gtype != nil, decl.initval.token())
	decl.variable.gtype = gtype
}

func (clause *ForRangeClause) infer() {
	//debugf(S("infering ForRangeClause"))
	collectionType := clause.rangeexpr.getGtype()
	//debugf(S("collectionType = %s"), collectionType)
	indexvarRel, ok := clause.indexvar.(*Relation)
	assert(ok, nil, S("ok"))
	indexvar, ok := indexvarRel.expr.(*ExprVariable)
	assert(ok, nil, S("ok"))

	var indexType *Gtype
	switch collectionType.getKind() {
	case G_ARRAY, G_SLICE:
		indexType = gInt
	case G_MAP:
		indexType = collectionType.Underlying().mapKey
	default:
		// @TODO consider map etc.
		TBI(clause.tok, S("unable to handle %d "), collectionType.getKind())
	}
	indexvar.gtype = indexType

	if clause.valuevar != nil {
		valuevarRel, ok := clause.valuevar.(*Relation)
		assert(ok, nil, S("ok"))
		valuevar, ok := valuevarRel.expr.(*ExprVariable)
		assert(ok, nil, S("ok"))

		var elementType *Gtype
		if collectionType.getKind() == G_ARRAY {
			elementType = collectionType.Underlying().elementType
		} else if collectionType.getKind() == G_SLICE {
			elementType = collectionType.Underlying().elementType
		} else if collectionType.getKind() == G_MAP {
			elementType = collectionType.Underlying().mapValue
		} else {
			errorft(clause.token(), S("internal error"))
		}
		//debugf(S("for i, v %s := rannge %v"), elementType, collectionType)
		valuevar.gtype = elementType
	}
}

func (ast *StmtShortVarDecl) infer() {
	//debugf(S("infering StmtShortVarDecl"))
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
					errorft(fcall.token(), S("funcdef of %s is not found"), fcall.fname)
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
			//debugf(S("rightExpr.gtype=%s"), gtype)
			secondGtype := rightExpr.(*ExprIndex).getSecondGtype()
			if secondGtype != nil {
				rightTypes = append(rightTypes, secondGtype)
			}
		default:
			if rightExpr == nil {
				errorft(ast.token(), S("rightExpr is nil"))
			}
			gtype := rightExpr.getGtype()
			if gtype == nil {
				errorft(ast.token(), S("rightExpr %T gtype is nil"), rightExpr)
			}
			//debugf(S("infered type %s"), gtype)
			rightTypes = append(rightTypes, gtype)
		}
	}

	if len(ast.lefts) > len(rightTypes) {
		// @TODO this check is too loose.
		errorft(ast.tok, S("number of lhs and rhs does not match (%d <=> %d)"), len(ast.lefts), len(rightTypes))
	}
	for i, e := range ast.lefts {
		rel := e.(*Relation) // a brand new rel
		variable := rel.expr.(*ExprVariable)
		rightType := rightTypes[i]
		variable.gtype = rightType
	}

}
