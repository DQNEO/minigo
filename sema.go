// Semantic Analyzer to produce IR struct
package main

func setTypeIds(namedTypes []*DeclType) {
	var typeId = 1 // start with 1 because we want to zero as error
	for _, concreteNamedType := range namedTypes {
		concreteNamedType.gtype.receiverTypeId = typeId
		typeId++
	}
}

func makeIR(internalUniverse *AstFile, internalRuntime *AstFile, csl *compiledStdlib, mainPkg *AstPackage, stringLiterals []*ExprStringLiteral, allDynamicTypes []*Gtype) *IrRoot {
	var importedPackages []*AstPackage

	for _, pkgName := range csl.uniqImportedPackageNames {
		compiledPkg := csl.compiledPackages[identifier(pkgName)]
		importedPackages = append(importedPackages, compiledPkg)
	}

	var declvars []*DeclVar
	var funcs []*DeclFunc
	for _, pkg := range importedPackages {
		for _, f := range pkg.files {
			for _, decl := range f.topLevelDecls {
				if decl.vardecl != nil {
					declvars = append(declvars, decl.vardecl)
				} else if decl.funcdecl != nil {
					funcs = append(funcs, decl.funcdecl)
				}
			}
		}

		for _, sl := range pkg.stringLiterals {
			stringLiterals = append(stringLiterals, sl)
		}
	}

	var files []*AstFile
	files = append(files, internalUniverse)
	files = append(files, internalRuntime)
	for _, f := range mainPkg.files {
		files = append(files, f)
	}

	for _, f := range files {
		for _, decl := range f.topLevelDecls {
			if decl.vardecl != nil {
				declvars = append(declvars, decl.vardecl)
			} else if decl.funcdecl != nil {
				funcs = append(funcs, decl.funcdecl)
			}
		}
	}
	for _, sl := range mainPkg.stringLiterals {
		stringLiterals = append(stringLiterals, sl)
	}

	root := &IrRoot{}

	root.stringLiterals = stringLiterals
	root.vars = declvars
	root.funcs = funcs
	root.setDynamicTypes(allDynamicTypes)
	root.importOS = in_array("os", csl.uniqImportedPackageNames)
	root.composeMethodTable()
	return root
}

func (ir *IrRoot) setDynamicTypes(allDynamicTypes []*Gtype) {
	var uniquedDTypes []string = builtinTypesAsString
	for _, gtype := range allDynamicTypes {
		gs := gtype.String()
		if !in_array(gs, uniquedDTypes) {
			uniquedDTypes = append(uniquedDTypes, gs)
		}
	}

	ir.uniquedDTypes = uniquedDTypes
}

func (ir *IrRoot) composeMethodTable() {
	var methodTable map[int][]string = map[int][]string{} // receiverTypeId : []methodTable
	for _, funcdecl := range ir.funcs {
		if funcdecl.receiver != nil {
			//debugf("funcdecl:%v", funcdecl)
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
	}
	debugf("set methodTable")

	ir.methodTable = methodTable
}
