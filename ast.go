package main

type Expr interface  {
	emit()
	dump()
}

type ExprNumberLiteral struct {
	val int
}

type ExprStringLiteral struct {
	val string
	slabel string
}

// local or global variable
type ExprVariable struct {
	varname identifier
	gtype *Gtype
	offset int // for local variable
	isGlobal bool
}

type ExprFuncall struct {
	fname identifier
	args  []Expr
}

type ExprBinop struct {
	op    string
	left  Expr
	right Expr
}

type ExprUop struct {
	op string
	operand Expr
}

// local or global
type AstVarDecl struct {
	variable *ExprVariable
	initval  Expr
}

type AstConstDecl struct {
	variable *ExprVariable
	initval  Expr
}

type AstAssignment struct {
	left  *ExprVariable // lvalue
	right Expr
}

type AstStmt struct {
	declvar    *AstVarDecl
	constdecl  *AstConstDecl
	assignment *AstAssignment
	expr       Expr
}

type AstPackageClause struct {
	name identifier
}

type AstImportDecl struct {
	paths []string
}

type AstCompountStmt struct {
	// compound
	stmts []*AstStmt
}

type AstFuncDecl struct {
	// funcdef
	fname identifier
	rettype string
	params []*ExprVariable
	localvars []*ExprVariable
	body *AstCompountStmt
}

type AstTopLevelDecl struct {
	funcdecl *AstFuncDecl
	vardecl *AstVarDecl
	constdecl *AstConstDecl
	typedecl *AstTypeDecl
}

type AstSourceFile struct {
	pkg *AstPackageClause
	imports []*AstImportDecl
	decls []*AstTopLevelDecl
}


