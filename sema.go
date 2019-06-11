// Semantic Analyzer to produce IR struct
package main

import "fmt"

func (n *AstFile) walk() *AstFile {
	for _, decl := range n.DeclList {
		decl = decl.walk()
	}
	return n
}

func (n *TopLevelDecl) walk() *TopLevelDecl {
	if n.funcdecl != nil {
		n.funcdecl = n.funcdecl.walk()
	} else if n.vardecl != nil {
		n.vardecl = n.vardecl.walk()
	} else {
		// assert not reach here
		return nil
	}

	return n
}

func (f *DeclFunc) prepare() {
	var params []*ExprVariable

	// prepend receiver to params
	if f.receiver != nil {
		params = []*ExprVariable{f.receiver}
		for _, param := range f.params {
			params = append(params, param)
		}
	} else {
		params = f.params
	}

	f.passParams = params

	// offset for params and local variables
	var offset int

	var regIndex int
	for _, param := range f.passParams {
		var regConsume int
		switch param.getGtype().is24WidthType() {
		case true:
			regConsume = 3
		default:
			regConsume = 1
		}

		offset -= IntSize * regConsume
		param.offset = offset
		param.regIndex = regIndex
		regIndex += regConsume
	}

	var localarea int
	for _, lvar := range f.localvars {
		if lvar.gtype == nil {
			debugf("%s has nil gtype ", lvar)
		}
		size := lvar.gtype.getSize()
		assert(size != 0, lvar.token(), "size should  not be zero:"+lvar.gtype.String())
		loff := align(size, 8)
		localarea -= loff
		offset -= loff
		lvar.offset = offset
		//debugf("set offset %d to lvar %s, type=%s", lvar.offset, lvar.varname, lvar.gtype)
	}
	f.localarea = localarea
}

// transform nodes
func walkStmt(stmt Stmt) Stmt {
	if stmt == nil {
		return nil
	}
	switch stmt.(type) {
	case *StmtFor:
		f, _ := stmt.(*StmtFor)
		f.block = walkStmt(f.block)
		if f.rng != nil {
			if f.rng.rangeexpr.getGtype().getKind() == G_MAP {
				f.kind = FOR_KIND_RANGE_MAP
				mapCounter := f.rng.mapCounter

				// counter = 0
				f.rng.init = &StmtAssignment{
					lefts: []Expr{
						mapCounter,
					},
					rights: []Expr{
						&ExprNumberLiteral{
							val: 0,
						},
					},
				}

				// counter < len(list)
				f.rng.cond = &ExprBinop{
					op:   "<",
					left: mapCounter, // i
					// @TODO
					// The range expression x is evaluated once before beginning the loop
					right: &ExprLen{
						arg: f.rng.rangeexpr, // len(expr)
					},
				}

				// counter++
				f.rng.post = &StmtInc{
					operand: mapCounter,
				}
			} else {
				f.kind = FOR_KIND_RANGE_LIST
			}
		} else {
			//f.cls.init = walkStmt(f.cls.init) // This does not work
			f.cls.cond = walkStmt(f.cls.cond)
			f.cls.post = walkStmt(f.cls.post)
			f.kind = FOR_KIND_PLAIN
		}
		return f
	case *StmtIf:
		s, _ := stmt.(*StmtIf)
		s.simplestmt = walkStmt(s.simplestmt)
		s.then = walkStmt(s.then)
		s.els = walkStmt(s.els)
		return s
	case *StmtReturn:
		s, _ := stmt.(*StmtReturn)
		if len(s.exprs) > 7 {
			TBI(s.token(), "too many number of arguments")
		}
		for i := 0;i<len(s.exprs);i++ {
			e := s.exprs[i]
			e = walkExpr(e)
			s.exprs[i] = e
		}
		return s
	case *StmtInc:
	case *StmtDec:
	case *StmtSatementList:
		s, _ := stmt.(*StmtSatementList)
		for i:=0;i<len(s.stmts);i++ {
			stmt := s.stmts[i]
			s.stmts[i] = walkStmt(stmt)
		}
		return s
	case *StmtAssignment:
		s, _ := stmt.(*StmtAssignment)
		/*
		for i:=0; i<len(s.lefts); i++ {
			left := s.lefts[i]
			left = walkExpr(left)
			s.lefts[i] = left
		}
		for i:=0; i<len(s.rights); i++ {
			right := s.rights[i]
			right = walkExpr(right)
			s.rights[i] = right
		}
		*/
		return s
	case *StmtShortVarDecl:
		s, _ := stmt.(*StmtShortVarDecl)
		for i:=0; i<len(s.rights); i++ {
			right := s.rights[i]
			right = walkExpr(right)
			s.rights[i] = right
		}
		a := &StmtAssignment{
			tok:    s.tok,
			lefts:  s.lefts,
			rights: s.rights,
		}
		return walkStmt(a)
	case *StmtContinue:
	case *StmtBreak:
	case *StmtSwitch:
		s, _ := stmt.(*StmtSwitch)
		for _, xcase := range s.cases {
			xcase.compound = walkStmt(xcase.compound)
		}
		s.dflt = walkStmt(s.dflt)
	case *StmtDefer:
	case *StmtExpr:
		s, _ := stmt.(*StmtExpr)
		s.expr = walkExpr(s.expr)
		return s
	}
	return stmt
}

func walkExpr(expr Expr) Expr {
	if expr == nil {
		return nil
	}
	switch expr.(type) {
	case *Relation:
		e,_ := expr.(*Relation)
		return e.expr
	case *ExprNilLiteral:
	case *ExprNumberLiteral:
	case *ExprStringLiteral:
	case *ExprVariable:
	case *ExprConstVariable:
	case *ExprFuncallOrConversion:
		funcall,_ := expr.(*ExprFuncallOrConversion)
		for i:=0;i<len(funcall.args);i++ {
			arg := funcall.args[i]
			arg = walkExpr(arg)
			funcall.args[i] = arg
		}
		if funcall.rel.expr == nil && funcall.rel.gtype != nil {
			// Conversion
			return &ExprConversion{
				tok:   funcall.token(),
				gtype: funcall.rel.gtype,
				expr:  funcall.args[0],
			}
		}
		decl := funcall.getFuncDef()
		switch decl {
		case builtinLen:
			arg := funcall.args[0]
			return &ExprLen{
				tok: arg.token(),
				arg: arg,
			}
		}
		return funcall
	case *ExprMethodcall:
		methodCall,_ := expr.(*ExprMethodcall)
		for i:=0 ;i<len(methodCall.args); i++ {
			arg := methodCall.args[i]
			arg = walkExpr(arg)
			methodCall.args[i] = arg
		}
		return methodCall
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
	case *ExprIndex:
		e,_ := expr.(*ExprIndex)
		e.index = walkExpr(e.index)
		e.collection = walkExpr(e.collection)
		return e
	case *ExprArrayLiteral:
	case *ExprSliceLiteral:
	case *ExprTypeAssertion:
	case *ExprVaArg:
	case *ExprConversion:
	case *ExprStructLiteral:
	case *ExprStructField:
	case *ExprTypeSwitchGuard:
	case *ExprMapLiteral:
	case *ExprLen:

	case *ExprCap:
	case *ExprConversionToInterface:
	}
	return expr
}

func (f *DeclFunc) walk() *DeclFunc {
	f.prepare()
	f.body = walkStmt(f.body)
	return f
}

func (n *DeclVar) walk() *DeclVar {
	return n
}

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
