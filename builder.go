// builder builds packages
package main

func extractImports(astFile *AstFile) map[string]bool {
	var imports map[string]bool = map[string]bool{}
	for _, importDecl := range astFile.importDecls {
		for _, spec := range importDecl.specs {
			baseName := getBaseNameFromImport(spec.path)
			imports[baseName] = true
		}
	}
	return imports
}

// analyze imports of given go files
func parseImports(sourceFiles []string) map[string]bool {

	var imported map[string]bool = map[string]bool{}
	for _, sourceFile := range sourceFiles {
		var importedInFile map[string]bool = map[string]bool{}
		p := &parser{}
		astFile := p.ParseFile(sourceFile, nil, true)
		importedInFile = extractImports(astFile)
		for name, _ := range importedInFile {
			imported[name] = true
		}
	}

	return imported
}

// inject builtin functions into the universe scope
func compileUniverse(universe *Scope) *AstPackage {
	p := &parser{
		packageName: identifier(""),
	}
	f := p.ParseString("internal_universe.go", internalUniverseCode, universe, false)
	attachMethodsToTypes(f.methods, p.packageBlockScope)
	inferTypes(f.uninferredGlobals, f.uninferredLocals)
	calcStructSize(f.dynamicTypes)
	return &AstPackage{
		name:           identifier(""),
		files:          []*AstFile{f},
		stringLiterals: f.stringLiterals,
		dynamicTypes:   f.dynamicTypes,
	}
}

// inject runtime things into the universe scope
func compileRuntime(universe *Scope) *AstPackage {
	p := &parser{
		packageName: identifier("iruntime"),
	}
	f := p.ParseString("internal_runtime.go", internalRuntimeCode, universe, false)
	attachMethodsToTypes(f.methods, p.packageBlockScope)
	inferTypes(f.uninferredGlobals, f.uninferredLocals)
	calcStructSize(f.dynamicTypes)
	return &AstPackage{
		name:           identifier("iruntime"),
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
func compileFiles(universe *Scope, sourceFiles []string) *AstPackage {
	// compile the main package
	pkgName := identifier("main")
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
func compileStdLibs(universe *Scope, directDependencies map[string]bool) map[identifier]*AstPackage {

	// "fmt" depends on "os. So inject it in advance.
	// Actually, dependency graph should be analyzed.
	primPackages := []string{"syscall", "io", "bytes", "os", "strconv"}
	var sortedUniqueImports []string
	sortedUniqueImports = primPackages
	for pkg, _ := range directDependencies {
		if !inArray(pkg, sortedUniqueImports) {
			sortedUniqueImports = append(sortedUniqueImports, pkg)
		}
	}

	var compiledStdPkgs map[identifier]*AstPackage = map[identifier]*AstPackage{}

	for _, spkgName := range sortedUniqueImports {
		pkgName := identifier(spkgName)
		pkgCode, ok := stdSources[pkgName]
		if !ok {
			errorf("package '%s' is not a standard library.", spkgName)
		}
		codes := []string{pkgCode}
		pkg := ParseFiles(pkgName, codes, true)
		pkg = makePkg(pkg, universe)
		compiledStdPkgs[pkgName] = pkg
		symbolTable.allScopes[pkgName] = pkg.scope
	}

	return compiledStdPkgs
}

type Program struct {
	packages    []*AstPackage
	methodTable map[int][]string
	importOS    bool
}

func build(universe *AstPackage, iruntime *AstPackage, stdPkgs map[identifier]*AstPackage, mainPkg *AstPackage) *Program {
	var packages []*AstPackage

	packages = append(packages, universe)
	packages = append(packages, iruntime)

	for _, pkg := range stdPkgs {
		packages = append(packages, pkg)
	}

	packages = append(packages, mainPkg)

	var dynamicTypes []*Gtype
	var funcs []*DeclFunc

	for _, pkg := range packages {
		collectDecls(pkg)
		if pkg == universe {
			setStringLables(pkg, "universe")
		} else {
			setStringLables(pkg, pkg.name)
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
	_, importOS := stdPkgs["os"]
	program.importOS = importOS
	program.methodTable = composeMethodTable(funcs)
	return program
}
