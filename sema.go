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

func (n *DeclFunc) walk() *DeclFunc {
	return n
}

func (n *DeclVar) walk() *DeclVar {
	return n
}

var symbolTable *SymbolTable

type SymbolTable struct {
	allScopes map[identifier]*Scope
	uniquedDTypes []string
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
