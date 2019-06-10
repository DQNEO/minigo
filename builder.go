// builder builds packages
package main

// analyze imports of given go files
func parseImports(sourceFiles []string) []string {

	// "fmt" depends on "os. So inject it in advance.
	// Actually, dependency graph should be analyzed.
	var imported []string = []string{"os"}
	for _, sourceFile := range sourceFiles {
		p := &parser{}
		astFile := p.ParseFile(sourceFile, nil, true)
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
	p := &parser{
		packageName: "",
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
		packageName: "iruntime",
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
func compileFiles(universe *Scope, sourceFiles []string) *AstPackage {
	// compile the main package
	var pkgName identifier = "main"
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
func compileStdLibs(universe *Scope, imported []string) *compiledStdlib {
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
		pkg := ParseFiles(pkgName, codes, true)
		pkg = makePkg(pkg,universe)
		libs.AddPackage(pkg)
		symbolTable.allScopes[pkgName] = pkg.scope
	}

	return libs
}

type compiledStdlib struct {
	compiledPackages         map[identifier]*AstPackage
	uniqImportedPackageNames []string
}

func (csl *compiledStdlib) getPackages() []*AstPackage {
	var importedPackages []*AstPackage

	for _, pkgName := range csl.uniqImportedPackageNames {
		compiledPkg := csl.compiledPackages[identifier(pkgName)]
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
	packages      []*AstPackage
	methodTable   map[int][]string
	importOS      bool
}

func build(universe *AstPackage, iruntime *AstPackage, csl *compiledStdlib, mainPkg *AstPackage) *Program {
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
		for _, f := range pkg.files {
			f = f.walk()
		}
	}

	symbolTable.uniquedDTypes = uniqueDynamicTypes(dynamicTypes)

	program := &Program{}
	program.packages = packages
	program.importOS = in_array("os", csl.uniqImportedPackageNames)
	program.methodTable = composeMethodTable(funcs)
	return program
}
