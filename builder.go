// builder builds packages
package main

// analyze imports of given go files
func parseImports(sourceFiles []gostring) []gostring {

	// "fmt" depends on "os. So inject it in advance.
	// Actually, dependency graph should be analyzed.
	var imported []gostring = []gostring{gostring("os"),gostring("strconv")}
	for _, sourceFile := range sourceFiles {
		p := &parser{}
		astFile := p.ParseFile(sourceFile, nil, true)
		for _, importDecl := range astFile.importDecls {
			for _, spec := range importDecl.specs {
				baseName := getBaseNameFromImport(spec.path)
				if !inArray(baseName, imported) {
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
	f := p.ParseString(S("internal_universe.go"), gostring(internalUniverseCode), universe, false)
	attachMethodsToTypes(f.methods, p.packageBlockScope)
	inferTypes(f.uninferredGlobals, f.uninferredLocals)
	calcStructSize(f.dynamicTypes)
	return &AstPackage{
		name:           goidentifier(""),
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
	f := p.ParseString(S("internal_runtime.go"), gostring(internalRuntimeCode), universe, false)
	attachMethodsToTypes(f.methods, p.packageBlockScope)
	inferTypes(f.uninferredGlobals, f.uninferredLocals)
	calcStructSize(f.dynamicTypes)
	return &AstPackage{
		name:           goidentifier(""),
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
	var pkgName goidentifier = goidentifier("main")
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

func toKey(gid goidentifier) identifier {
	return identifier(gid)
}

// parse standard libraries
func compileStdLibs(universe *Scope, imported []gostring) *compiledStdlib {
	var libs *compiledStdlib = &compiledStdlib{
		compiledPackages:         map[identifier]*AstPackage{},
		uniqImportedPackageNames: nil,
	}
	stdPkgs := makeStdLib()

	for _, spkgName := range imported {
		pkgName := goidentifier(spkgName)
		pkgCode, ok := stdPkgs[toKey(pkgName)]
		if !ok {
			errorf(S("package '%s' is not a standard library."), spkgName)
		}
		var codes []gostring = []gostring{gostring(pkgCode)}
		pkg := ParseFiles(pkgName, codes, true)
		pkg = makePkg(pkg, universe)
		libs.AddPackage(pkg)
		symbolTable.allScopes[toKey(pkgName)] = pkg.scope
	}

	return libs
}

type compiledStdlib struct {
	compiledPackages         map[identifier]*AstPackage
	uniqImportedPackageNames []gostring
}

func (csl *compiledStdlib) getPackages() []*AstPackage {
	var importedPackages []*AstPackage

	for _, pkgName := range csl.uniqImportedPackageNames {
		compiledPkg := csl.compiledPackages[toKey(goidentifier(pkgName))]
		importedPackages = append(importedPackages, compiledPkg)
	}
	return importedPackages
}

func (csl *compiledStdlib) AddPackage(pkg *AstPackage) {
	csl.compiledPackages[toKey(pkg.name)] = pkg
	if !inArray(gostring(pkg.name), csl.uniqImportedPackageNames) {
		csl.uniqImportedPackageNames = append(csl.uniqImportedPackageNames, gostring(pkg.name))
	}
}

type Program struct {
	packages    []*AstPackage
	methodTable map[int][]gostring
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
			setStringLables(pkg, S("universe"))
		} else {
			setStringLables(pkg, gostring(pkg.name))
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
	program.importOS = inArray(S("os"), csl.uniqImportedPackageNames)
	program.methodTable = composeMethodTable(funcs)
	return program
}
