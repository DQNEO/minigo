// builder builds packages
package main

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

