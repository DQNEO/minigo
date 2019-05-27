// builder builds packages
package main

// analyze imports of given go files
func parseImports(sourceFiles []string) []string {

	pForImport := &parser{}
	// "fmt" depends on "os. So inject it in advance.
	// Actually, dependency graph should be analyzed.
	var imported []string = []string{"os"}
	for _, sourceFile := range sourceFiles {
		astFile := pForImport.parseFile(sourceFile, nil, true)
		for _, importDecl := range astFile.importDecls {
			for _, spec := range importDecl.specs {
				baseName := getBaseNameFromImport(spec.path)
				if !in_array(baseName, imported) {
					imported = append(imported, baseName)
				}
			}
		}
	}

	return imported
}

// inject builtin functions into the universe scope
func compileUniverse(universe *Scope) *AstPackage {
	p := &parser{}
	p.initPackage("")
	internalUniverse := p.parseString("internal_universe.go", internalUniverseCode, universe, false)
	p.resolveMethods()
	inferTypes(p.packageUninferredGlobals, p.packageUninferredLocals)
	return  &AstPackage{
		name:"",
		files:[]*AstFile{internalUniverse},
		stringLiterals:p.packageStringLiterals,
		dynamicTypes:p.packageDynamicTypes,
	}
}

// inject runtime things into the universe scope
func compileRuntime(universe *Scope) *AstPackage {
	p := &parser{}
	p.initPackage("")
	internalRuntime := p.parseString("internal_runtime.go", internalRuntimeCode, universe, false)
	p.resolveMethods()
	inferTypes(p.packageUninferredGlobals, p.packageUninferredLocals)
	return &AstPackage{
		name:"",
		files:[]*AstFile{internalRuntime},
		stringLiterals:p.packageStringLiterals,
		dynamicTypes:p.packageDynamicTypes,
	}
}

func compileMainPackage(universe *Scope, sourceFiles []string) *AstPackage {
	// compile the main package
	p := &parser{}
	mainPkg := ParseSources(p, identifier("main"), sourceFiles, false)
	if parseOnly {
		if debugAst {
			mainPkg.dump()
		}
		return nil
	}
	resolveInPackage(mainPkg, universe)
	p.resolveMethods()
	allScopes[mainPkg.name] = mainPkg.scope
	inferTypes(p.packageUninferredGlobals, p.packageUninferredLocals)
	setTypeIds(mainPkg.namedTypes)
	if debugAst {
		mainPkg.dump()
	}

	if resolveOnly {
		return nil
	}
	return mainPkg
}

// parse standard libraries
func compileStdLibs(universe *Scope, imported []string) *compiledStdlib {
	var libs *compiledStdlib = &compiledStdlib{
		compiledPackages:         map[identifier]*AstPackage{},
		uniqImportedPackageNames: nil,
	}
	stdPkgs := makeStdLib()

	for _, spkgName := range imported {
		p := &parser{}
		pkgName := identifier(spkgName)
		pkgCode, ok := stdPkgs[pkgName]
		if !ok {
			errorf("package '" + spkgName + "' is not a standard library.")
		}
		var codes []string = []string{pkgCode}
		pkg := ParseSources(p, pkgName, codes, true)
		resolveInPackage(pkg, universe)
		p.resolveMethods()
		allScopes[pkgName] = pkg.scope
		inferTypes(p.packageUninferredGlobals, p.packageUninferredLocals)
		setTypeIds(pkg.namedTypes)
		libs.AddPackage(pkg)
	}

	return libs
}

type compiledStdlib struct {
	compiledPackages         map[identifier]*AstPackage
	uniqImportedPackageNames []string
}

func (csl *compiledStdlib) AddPackage(pkg *AstPackage) {
	csl.compiledPackages[pkg.name] = pkg
	if !in_array(string(pkg.name), csl.uniqImportedPackageNames) {
		csl.uniqImportedPackageNames = append(csl.uniqImportedPackageNames, string(pkg.name))
	}
}
