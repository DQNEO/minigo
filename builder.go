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

func compileMainFiles(universe *Scope, sourceFiles []string) *AstPackage {
	// compile the main package
	pm := &parser{}
	mainPkg := ParseSources(pm, identifier("main"), sourceFiles, false)
	if parseOnly {
		if debugAst {
			mainPkg.dump()
		}
		return nil
	}
	pm.resolve(universe)
	allScopes[mainPkg.name] = mainPkg.scope
	inferTypes(pm.packageUninferredGlobals, pm.packageUninferredLocals)
	setTypeIds(mainPkg.namedTypes)
	if debugAst {
		mainPkg.dump()
	}

	if resolveOnly {
		return nil
	}

	return mainPkg
}

func compileStdLibs(universe *Scope, imported []string) *compiledStdlib {

	// add std packages
	// parse std packages
	var libs *compiledStdlib = &compiledStdlib{
		compiledPackages:         map[identifier]*AstPackage{},
		uniqImportedPackageNames: nil,
	}
	stdPkgs := makeStdLib()

	for _, spkgName := range imported {
		pkgName := identifier(spkgName)
		pkgCode, ok := stdPkgs[pkgName]
		if !ok {
			errorf("package '" + spkgName + "' is not a standard library.")
		}
		var codes []string = []string{pkgCode}
		p := &parser{}
		pkg := ParseSources(p, pkgName, codes, true)
		p.resolve(universe)
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
