// builder builds packages
package main

import "./stdlib/fmt"

// "fmt" => "/stdlib/fmt"
// "./stdlib/fmt" => "/stdlib/fmt"
func normalizeImportPath(path string) normalizedPackagePath {
	if len(path) > 9 && path[0] == '.' {
		// "./stdlib/fmt" => "/stdlib/fmt"
		bpath := []byte(path)
		bpath2 := bpath[1:]
		return normalizedPackagePath(string(bpath2))
	} else {
		// "fmt" => "/stdlib/fmt"
		return normalizedPackagePath("/stdlib/" + path)
	}
}

// analyze imports of a given go source
func parseImportsFromFile(sourceFile string) importMap {
	p := &parser{}
	astFile := p.ParseFile(sourceFile, nil, true)
	return astFile.imports
}

// analyze imports of given go files
func parseImports(sourceFiles []string) importMap {

	var imported importMap = map[normalizedPackagePath]bool{}
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
		packagePath: normalizeImportPath("builtin"), // anything goes
		packageName: identifier("builtin"),
	}
	f := p.ParseFile("stdlib/builtin/builtin.go", universe, false)
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
	pkgPath := normalizeImportPath("unsafe") // need to be normalized because it's imported by iruntime
	pkgScope := newScope(nil, pkgName)
	symbolTable.allScopes[pkgPath] = pkgScope

	p := &parser{
		packagePath:pkgPath,
		packageName: pkgName,
	}
	filepath := getStdFileName(pkgPath)
	f := p.ParseFile(filepath, pkgScope, false)
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

const IRuntimePath normalizedPackagePath = "iruntime"
const MainPath normalizedPackagePath = "main"

// inject runtime things into the universe scope
func compileRuntime(universe *Scope) *AstPackage {
	pkgName := identifier("iruntime")
	pkgPath := IRuntimePath
	pkgScope := newScope(nil, pkgName)
	symbolTable.allScopes[pkgPath] = pkgScope
	p := &parser{
		packagePath: pkgPath,
		packageName: pkgName,
	}
	f := p.ParseFile("internal/runtime/runtime.go", universe, false)
	attachMethodsToTypes(f.methods, p.packageBlockScope)
	inferTypes(f.uninferredGlobals, f.uninferredLocals)
	calcStructSize(f.dynamicTypes)
	return &AstPackage{
		normalizedPath: pkgPath,
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

// compileMainFiles parses files into *AstPackage
func compileMainFiles(universe *Scope, sourceFiles []string) *AstPackage {
	// compile the main package
	pkgName := identifier("main")
	pkgScope := newScope(nil, pkgName)
	mainPkg := ParseFiles("main", MainPath, pkgScope, sourceFiles)
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

type importMap map[normalizedPackagePath]bool

func parseImportRecursive(dep map[normalizedPackagePath]importMap, directDependencies importMap) {
	for normalizedPackagePath, _ := range directDependencies {
		file := getStdFileName(normalizedPackagePath)
		imports := parseImportsFromFile(file)
		dep[normalizedPackagePath] = imports
		parseImportRecursive(dep, imports)
	}
}

func removeResolvedPkg(dep map[normalizedPackagePath]importMap, pkgToRemove normalizedPackagePath) map[normalizedPackagePath]importMap {
	var dep2 map[normalizedPackagePath]importMap = map[normalizedPackagePath]importMap{}

	for pkg1, imports := range dep {
		if pkg1 == pkgToRemove {
			continue
		}
		var newimports importMap = map[normalizedPackagePath]bool{}
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

func removeResolvedPackages(dep map[normalizedPackagePath]importMap, sortedUniqueImports []normalizedPackagePath) map[normalizedPackagePath]importMap {
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

func get0dependentPackages(dep map[normalizedPackagePath]importMap) []normalizedPackagePath {
	var moved []normalizedPackagePath
	if len(dep) == 0 {
		return nil
	}
	for spkgName, imports := range dep {
		var numDepends int
		for _, flg := range imports {
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

func resolveDependency(directDependencies importMap) []normalizedPackagePath {
	var sortedUniqueImports []normalizedPackagePath
	var dep map[normalizedPackagePath]importMap = map[normalizedPackagePath]importMap{}
	parseImportRecursive(dep, directDependencies)

	for {
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

// "/stdlib/fmt" => "./stdlib/fmt/fmt.go"
func getStdFileName(path normalizedPackagePath) string {
	baseName := getBaseNameFromImport(string(path))
	filename := baseName + ".go"
	return fmt.Sprintf("./%s/%s", string(path), filename)
}

// Compile standard libraries
func compileStdLibs(universe *Scope, directDependencies importMap) map[normalizedPackagePath]*AstPackage {

	sortedUniqueImports := resolveDependency(directDependencies)

	var compiledStdPkgs map[normalizedPackagePath]*AstPackage = map[normalizedPackagePath]*AstPackage{}

	for _, path := range sortedUniqueImports {
		file := getStdFileName(path) // => "./stdlib/io/ioutil/ioutil.go"
		files := []string{file}
		pkgScope := newScope(nil, identifier(path))
		symbolTable.allScopes[path] = pkgScope
		pkgShortName := getBaseNameFromImport(string(path))
		pkg := ParseFiles(identifier(pkgShortName), path, pkgScope, files)
		pkg = makePkg(pkg, universe)
		compiledStdPkgs[path] = pkg
	}

	return compiledStdPkgs
}

type Program struct {
	packages    []*AstPackage
	methodTable map[int][]string
}

func build(pkgUniverse *AstPackage, pkgUnsafe *AstPackage, pkgIRuntime *AstPackage, stdPkgs map[normalizedPackagePath]*AstPackage, pkgMain *AstPackage) *Program {
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
	program.methodTable = composeMethodTable(funcs)
	return program
}
