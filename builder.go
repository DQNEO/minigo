// builder builds packages
package main

import "fmt"

func extractImports(astFile *AstFile) importMap {
	var imports importMap = map[string]bool{}
	for _, importDecl := range astFile.importDecls {
		for _, spec := range importDecl.specs {
			baseName := getBaseNameFromImport(spec.path)
			imports[baseName] = true
		}
	}
	return imports
}

// analyze imports of a given go source
func parseImportsFromString(source string) importMap {

	var imports importMap = map[string]bool{}

	p := &parser{}
	astFile := p.ParseString("", source, nil, true)
	imports = extractImports(astFile)

	return imports
}

// analyze imports of given go files
func parseImports(sourceFiles []string) importMap {

	var imported importMap = map[string]bool{}
	for _, sourceFile := range sourceFiles {
		var importedInFile importMap = map[string]bool{}
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

type importMap map[string]bool

func parseImportRecursive(dep map[string]importMap , directDependencies importMap, stdSources map[identifier]string) {
	for spkgName , _ := range directDependencies {
		pkgName := identifier(spkgName)
		pkgCode, ok := stdSources[pkgName]
		if !ok {
			errorf("package '%s' is not a standard library.", spkgName)
		}
		var imports = parseImportsFromString(pkgCode)
		dep[spkgName] = imports
		parseImportRecursive(dep, imports, stdSources)
	}
}

func removeResolvedPkg(dep map[string]importMap, pkgToRemove string) map[string]importMap {
	var dep2 map[string]importMap = map[string]importMap{}

	for pkg1, imports := range dep {
		if pkg1 == pkgToRemove {
			continue
		}
		var newimports importMap = map[string]bool{}
		for pkg2, _ := range imports {
			if pkg2 == pkgToRemove {
				continue
			}
			newimports[pkg2] = true
		}
		dep2[pkg1] = newimports
	}

	return dep2
}

func removeResolvedPackages(dep map[string]importMap, sortedUniqueImports []string) map[string]importMap {
	for _, resolved := range sortedUniqueImports {
		dep = removeResolvedPkg(dep, resolved)
	}
	return dep
}

func dumpDep(dep map[string]importMap) {
	debugf("#------------- dep -----------------")
	for spkgName, imports := range dep {
		debugf("#  %s has %d imports:", spkgName, len(imports))
		for sspkgName, _ := range imports {
			debugf("#    %s", sspkgName)
		}
	}
}

func get0dependentPackages(dep map[string]importMap) []string {
	var moved []string
	if len(dep) == 0 {
		return nil
	}
	for spkgName, imports := range dep {
		var numDepends int
		for _, flg  := range imports {
			if flg {
				numDepends++
			}
		}
		if numDepends == 0 {
			moved = append(moved, spkgName)
		}
	}
	return moved
}


func resolveDependency(directDependencies importMap, stdSources map[identifier]string) []string {
	var sortedUniqueImports []string
	var dep map[string]importMap = map[string]importMap{}
	parseImportRecursive(dep, directDependencies, stdSources)

	for  {
		//dumpDep(dep)
		moved := get0dependentPackages(dep)
		if len(moved) == 0 {
			break
		}
		dep = removeResolvedPackages(dep, moved)
		for _, pkg := range moved {
			sortedUniqueImports = append(sortedUniqueImports, pkg)
		}

	}
	return sortedUniqueImports
}

func getStdFileName(pkgName string) string {
	return fmt.Sprintf("stdlib/%s/%s.go", pkgName, pkgName)
}

// Compile standard libraries
func compileStdLibs(universe *Scope, directDependencies importMap, stdSources map[identifier]string) map[identifier]*AstPackage {

	sortedUniqueImports := resolveDependency(directDependencies, stdSources)

	var compiledStdPkgs map[identifier]*AstPackage = map[identifier]*AstPackage{}

	for _, spkgName := range sortedUniqueImports {
		pkgName := identifier(spkgName)
		file := getStdFileName(spkgName)
		files := []string{file}
		pkg := ParseFiles(pkgName, files, false)
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
