// builder builds packages
package main

// analyze imports of given go files
func parseImports(sourceFiles []gostring) []gostring {

	// "fmt" depends on "os. So inject it in advance.
	// Actually, dependency graph should be analyzed.
	var imported []gostring = []gostring{gostring("os")}
	for _, sourceFile := range sourceFiles {
		p := &parser{}
		astFile := p.ParseFile(string(sourceFile), nil, true)
		for _, importDecl := range astFile.importDecls {
			for _, spec := range importDecl.specs {
				baseName := getBaseNameFromImport(spec.path)
				if !inArray2(baseName, imported) {
					imported = append(imported, baseName)
				}
			}
		}
	}

	return imported
}

// inject builtin functions into the universe scope
func compileUniverse(universe *Scope) *AstPackage {
	p := &parser{
		packageName: goidentifier(""),
	}
	f := p.ParseString("internal_universe.go", internalUniverseCode, universe, false)
	attachMethodsToTypes(f.methods, p.packageBlockScope)
	inferTypes(f.uninferredGlobals, f.uninferredLocals)
	calcStructSize(f.dynamicTypes)
	return &AstPackage{
		name:           "",
		files:          []*AstFile{f},
		stringLiterals: f.stringLiterals,
		dynamicTypes:   f.dynamicTypes,
	}
}

// inject runtime things into the universe scope
func compileRuntime(universe *Scope) *AstPackage {
	p := &parser{
		packageName: goidentifier("iruntime"),
	}
	f := p.ParseString("internal_runtime.go", internalRuntimeCode, universe, false)
	attachMethodsToTypes(f.methods, p.packageBlockScope)
	inferTypes(f.uninferredGlobals, f.uninferredLocals)
	calcStructSize(f.dynamicTypes)
	return &AstPackage{
		name:           "",
		files:          []*AstFile{f},
		stringLiterals: f.stringLiterals,
		dynamicTypes:   f.dynamicTypes,
	}
}

func makePkg(pkg *AstPackage, universe *Scope) *AstPackage {
	resolveIdents(pkg, universe)
	attachMethodsToTypes(pkg.methods, pkg.scope)
	inferTypes(pkg.uninferredGlobals, pkg.uninferredLocals)
	calcStructSize(pkg.dynamicTypes)
	return pkg
}

// compileFiles parses files into *AstPackage
func compileFiles(universe *Scope, sourceFiles []gostring) *AstPackage {
	// compile the main package
	var pkgName packageName = "main"
	mainPkg := ParseFiles(pkgName, sourceFiles, false)
	if parseOnly {
		if debugAst {
			mainPkg.dump()
		}
		return nil
	}
	mainPkg = makePkg(mainPkg, universe)
	if debugAst {
		mainPkg.dump()
	}

	if resolveOnly {
		return nil
	}
	return mainPkg
}

// parse standard libraries
func compileStdLibs(universe *Scope, imported []gostring) *compiledStdlib {
	var libs *compiledStdlib = &compiledStdlib{
		compiledPackages:         map[packageName]*AstPackage{},
		uniqImportedPackageNames: nil,
	}
	stdPkgs := makeStdLib()

	for _, spkgName := range imported {
		pkgName := packageName(spkgName)
		pkgCode, ok := stdPkgs[pkgName]
		if !ok {
			errorf("package '" + string(spkgName) + "' is not a standard library.")
		}
		var codes []gostring = []gostring{gostring(pkgCode)}
		pkg := ParseFiles(pkgName, codes, true)
		pkg = makePkg(pkg, universe)
		libs.AddPackage(pkg)
		symbolTable.allScopes[pkgName] = pkg.scope
	}

	return libs
}

type compiledStdlib struct {
	compiledPackages         map[packageName]*AstPackage
	uniqImportedPackageNames []string
}

func (csl *compiledStdlib) getPackages() []*AstPackage {
	var importedPackages []*AstPackage

	for _, pkgName := range csl.uniqImportedPackageNames {
		compiledPkg := csl.compiledPackages[packageName(pkgName)]
		importedPackages = append(importedPackages, compiledPkg)
	}
	return importedPackages
}

func (csl *compiledStdlib) AddPackage(pkg *AstPackage) {
	csl.compiledPackages[pkg.name] = pkg
	if !in_array(string(pkg.name), csl.uniqImportedPackageNames) {
		csl.uniqImportedPackageNames = append(csl.uniqImportedPackageNames, string(pkg.name))
	}
}

type Program struct {
	packages    []*AstPackage
	methodTable map[int][]string
	importOS    bool
}

func build(universe *AstPackage, iruntime *AstPackage, csl *compiledStdlib, mainPkg *AstPackage) *Program {
	var packages []*AstPackage

	packages = append(packages, universe)

	importedPackages := csl.getPackages()
	for _, pkg := range importedPackages {
		packages = append(packages, pkg)
	}

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

	//  Do restructuring of local nodes
	for _, pkg := range packages {
		for _, fnc := range pkg.funcs {
			fnc = walkFunc(fnc)
		}
	}

	symbolTable.uniquedDTypes = uniqueDynamicTypes(dynamicTypes)

	program := &Program{}
	program.packages = packages
	program.importOS = in_array("os", csl.uniqImportedPackageNames)
	program.methodTable = composeMethodTable(funcs)
	return program
}
