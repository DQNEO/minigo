package main

type IrRoot struct {
	packages      []*AstPackage
	methodTable   map[int][]string
	uniquedDTypes []string
	importOS      bool
}

func makeIR(universe *AstPackage, iruntime *AstPackage, csl *compiledStdlib, mainPkg *AstPackage) *IrRoot {
	var packages []*AstPackage

	importedPackages := csl.getPackages()
	for _, pkg := range importedPackages {
		packages = append(packages, pkg)
	}

	packages = append(packages, universe)
	packages = append(packages, iruntime)
	packages = append(packages, mainPkg)

	var dynamicTypes []*Gtype
	var funcs []*DeclFunc

	for _, pkg := range packages {
		collectDecls(pkg)
		if pkg == universe {
			setStringLables(pkg, "universe")
		} else {
			setStringLables(pkg, string(pkg.name))
		}
		for _, dt := range pkg.dynamicTypes {
			dynamicTypes = append(dynamicTypes, dt)
		}
		for _, fn := range pkg.funcs {
			funcs = append(funcs, fn)
		}
		setTypeIds(pkg.namedTypes)
	}

	root := &IrRoot{}
	root.packages = packages
	root.setDynamicTypes(dynamicTypes)
	root.importOS = in_array("os", csl.uniqImportedPackageNames)
	root.methodTable = composeMethodTable(funcs)
	return root
}

func (root *IrRoot) setDynamicTypes(dynamicTypes []*Gtype) {
	var uniquedDTypes []string = builtinTypesAsString
	for _, gtype := range dynamicTypes {
		gs := gtype.String()
		if !in_array(gs, uniquedDTypes) {
			uniquedDTypes = append(uniquedDTypes, gs)
		}
	}

	root.uniquedDTypes = uniquedDTypes
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
