package main

import "fmt"

type IrRoot struct {
	vars           []*DeclVar
	funcs          []*DeclFunc
	packages       []*AstPackage
	methodTable    map[int][]string
	uniquedDTypes  []string
	importOS       bool
}

var groot *IrRoot

var vars []*DeclVar
var funcs []*DeclFunc

func collectDecls(pkg *AstPackage) {
	for _, f := range pkg.files {
		for _, decl := range f.topLevelDecls {
			if decl.vardecl != nil {
				vars = append(vars, decl.vardecl)
			} else if decl.funcdecl != nil {
				funcs = append(funcs, decl.funcdecl)
			}
		}
	}
}


func setStringLables(pkg *AstPackage, prefix string) {
	for id, sl := range pkg.stringLiterals {
		sl.slabel = fmt.Sprintf("%s.S%d", prefix, id+1)
	}
}

func makeIR(internalUniverse *AstPackage, internalRuntime *AstPackage, csl *compiledStdlib, mainPkg *AstPackage) *IrRoot {
	var packages []*AstPackage
	var dynamicTypes []*Gtype

	packages = append(packages, internalUniverse)
	packages = append(packages, internalRuntime)

	setStringLables(internalUniverse, "")
	setStringLables(internalRuntime, "iruntime")
	setStringLables(mainPkg, string(mainPkg.name))

	for _, dt := range internalUniverse.dynamicTypes {
		dynamicTypes = append(dynamicTypes, dt)
	}
	for _, dt := range internalRuntime.dynamicTypes {
		dynamicTypes = append(dynamicTypes, dt)
	}

	var importedPackages []*AstPackage

	for _, pkgName := range csl.uniqImportedPackageNames {
		compiledPkg := csl.compiledPackages[identifier(pkgName)]
		importedPackages = append(importedPackages, compiledPkg)
	}

	for _, pkg := range importedPackages {
		collectDecls(pkg)
		packages = append(packages, pkg)

		setStringLables(pkg, string(pkg.name))

		for _, dt := range pkg.dynamicTypes {
			dynamicTypes = append(dynamicTypes, dt)
		}
	}

	collectDecls(internalUniverse)
	collectDecls(internalRuntime)
	collectDecls(mainPkg)

	for _, dt := range mainPkg.dynamicTypes {
		dynamicTypes = append(dynamicTypes, dt)
	}


	packages = append(packages, mainPkg)

	root := &IrRoot{}
	root.packages = packages
	root.vars = vars
	root.funcs = funcs
	root.setDynamicTypes(dynamicTypes)
	root.importOS = in_array("os", csl.uniqImportedPackageNames)
	root.methodTable = composeMethodTable(funcs)
	return root
}

func (ir *IrRoot) setDynamicTypes(dynamicTypes []*Gtype) {
	var uniquedDTypes []string = builtinTypesAsString
	for _, gtype := range dynamicTypes {
		gs := gtype.String()
		if !in_array(gs, uniquedDTypes) {
			uniquedDTypes = append(uniquedDTypes, gs)
		}
	}

	ir.uniquedDTypes = uniquedDTypes
}

func composeMethodTable(funcs []*DeclFunc)  map[int][]string  {
	var methodTable map[int][]string = map[int][]string{} // receiverTypeId : []methodTable

	for _, funcdecl := range funcs {
		if funcdecl.receiver == nil {
			continue
		}
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
	debugf("set methodTable")
	return methodTable
}

