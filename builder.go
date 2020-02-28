// builder builds packages
package main

import "./stdlib/fmt"

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
func parseImportsFromFile(sourceFile string) importMap {
	p := &parser{}
	astFile := p.ParseFile(sourceFile, nil, true)
	imports := extractImports(astFile)
	return imports
}

// analyze imports of given go files
func parseImports(sourceFiles []string) importMap {

	var imported importMap = map[string]bool{}
	for _, sourceFile := range sourceFiles {
		importsInFile := parseImportsFromFile(sourceFile)
		for name, _ := range importsInFile {
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
	f := p.ParseFile("internal/universe/universe.go", universe, false)
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

// inject unsafe package
func compileUnsafe(universe *Scope) *AstPackage {
	pkgName := identifier("unsafe")
	pkgScope := newScope(nil, pkgName)
	symbolTable.allScopes[pkgName] = pkgScope

	p := &parser{
		packageName: pkgName,
	}
	f := p.ParseFile("stdlib/unsafe/unsafe.go", pkgScope, false)
	attachMethodsToTypes(f.methods, p.packageBlockScope)
	inferTypes(f.uninferredGlobals, f.uninferredLocals)
	calcStructSize(f.dynamicTypes)
	return &AstPackage{
		name:           pkgName,
		files:          []*AstFile{f},
		stringLiterals: f.stringLiterals,
		dynamicTypes:   f.dynamicTypes,
	}
}

// inject runtime things into the universe scope
func compileRuntime(universe *Scope) *AstPackage {
	pkgName := identifier("iruntime")
	pkgScope := newScope(nil, pkgName)
	symbolTable.allScopes[pkgName] = pkgScope
	p := &parser{
		packageName: pkgName,
	}
	f := p.ParseFile("internal/runtime/runtime.go", universe, false)
	attachMethodsToTypes(f.methods, p.packageBlockScope)
	inferTypes(f.uninferredGlobals, f.uninferredLocals)
	calcStructSize(f.dynamicTypes)
	return &AstPackage{
		name:           pkgName,
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
	pkgScope := newScope(nil, pkgName)
	mainPkg := ParseFiles(pkgScope, sourceFiles)
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

func parseImportRecursive(dep map[string]importMap , directDependencies importMap) {
	for spkgName , _ := range directDependencies {
		file := getStdFileName(spkgName)
		imports := parseImportsFromFile(file)
		dep[spkgName] = imports
		parseImportRecursive(dep, imports)
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


func resolveDependency(directDependencies importMap) []string {
	var sortedUniqueImports []string
	var dep map[string]importMap = map[string]importMap{}
	parseImportRecursive(dep, directDependencies)

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
func compileStdLibs(universe *Scope, directDependencies importMap) map[identifier]*AstPackage {

	sortedUniqueImports := resolveDependency(directDependencies)

	var compiledStdPkgs map[identifier]*AstPackage = map[identifier]*AstPackage{}

	for _, spkgName := range sortedUniqueImports {
		pkgName := identifier(spkgName)
		file := getStdFileName(spkgName)
		files := []string{file}
		pkgScope := newScope(nil, pkgName)
		symbolTable.allScopes[pkgName] = pkgScope
		pkg := ParseFiles(pkgScope, files)
		pkg = makePkg(pkg, universe)
		compiledStdPkgs[pkgName] = pkg
	}

	return compiledStdPkgs
}


type Program struct {
	packages    []*AstPackage
	methodTable map[int][]string
	importOS    bool
}

func build(pkgUniverse *AstPackage, pkgUnsafe *AstPackage, pkgIRuntime *AstPackage, stdPkgs map[identifier]*AstPackage, pkgMain *AstPackage) *Program {
	var packages []*AstPackage

	packages = append(packages, pkgUniverse)
	packages = append(packages, pkgUnsafe)
	packages = append(packages, pkgIRuntime)

	for _, pkg := range stdPkgs {
		packages = append(packages, pkg)
	}

	packages = append(packages, pkgMain)

	var dynamicTypes []*Gtype
	var funcs []*DeclFunc

	for _, pkg := range packages {
		collectDecls(pkg)
		if pkg == pkgUniverse {
			setStringLables(pkg, "pkgUniverse")
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
