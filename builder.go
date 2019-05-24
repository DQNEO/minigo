// builder builds packages
package main

// analyze imports of given go files
func parseImports(sourceFiles []string) []string {

	pForImport := &parser{}
	// "fmt" depends on "os. So inject it in advance.
	// Actually, dependency graph should be analyzed.
	var imported []string = []string{"os"}
	for _, sourceFile := range sourceFiles {
		bs := NewByteStreamFromFile(sourceFile)
		astFile := pForImport.parseSourceFile(bs, nil, true)
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

func parseStdPkg(p *parser, universe *Scope, pkgname identifier, code string) *stdpkg {
	filename := string(pkgname) + ".memory"
	bs := NewByteStreamFromString(filename, code)

	// initialize a package
	p.initPackage(pkgname)
	p.scopes[pkgname] = newScope(nil, string(pkgname))

	asf := p.parseSourceFile(bs, p.scopes[pkgname], false)

	p.resolve(universe)
	if debugAst {
		asf.dump()
	}
	return &stdpkg{
		name:  pkgname,
		files: []*SourceFile{asf},
	}
}

func compileInputFiles(p *parser, pkgname identifier, sourceFiles []string) *Package {
	p.initPackage(pkgname)
	p.scopes[pkgname] = newScope(nil, string(pkgname))
	var astFiles []*SourceFile
	for _, sourceFile := range sourceFiles {
		bs := NewByteStreamFromFile(sourceFile)
		asf := p.parseSourceFile(bs, p.scopes[pkgname], false)
		astFiles = append(astFiles, asf)
	}

	mainPkg := &Package{
	}
	mainPkg.files = astFiles
	return mainPkg
}

type compiledStdlib struct {
	compiledPackages map[identifier]*stdpkg
	uniqImportedPackageNames []string
}

func (csl *compiledStdlib) AddPackage(pkg *stdpkg) {
	csl.compiledPackages[pkg.name] = pkg
	if !in_array(string(pkg.name), csl.uniqImportedPackageNames) {
		csl.uniqImportedPackageNames = append(csl.uniqImportedPackageNames, string(pkg.name))
	}
}

type stdpkg struct {
	name  identifier
	files []*SourceFile
}

type Package struct {
	name identifier
	files []*SourceFile
}

