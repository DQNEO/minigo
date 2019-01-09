package main

type Expr interface {
	emit()
	dump()
	getGtype() *Gtype
}

type Stmt interface {
	emit()
}

type Relation struct {
	name identifier

	// either of expr(var, const, funcref) or gtype
	expr  Expr
	gtype *Gtype
}

type ExprNilLiteral struct {
}

type ExprNumberLiteral struct {
	val int
}

type ExprStringLiteral struct {
	val    string
	slabel string
}

// local or global variable
type ExprVariable struct {
	varname  identifier
	gtype    *Gtype
	offset   int // for local variable
	isGlobal bool
}

type ExprConstVariable struct {
	name  identifier
	gtype *Gtype
	val   Expr // like ExprConstExpr ?
	iotaIndex int  // for iota
}

type ExprFuncall struct {
	rel *Relation
	fname string
	args  []Expr
}

type ExprMethodcall struct {
	receiver Expr
	fname identifier
	args  []Expr
}

type ExprBinop struct {
	op    string
	left  Expr
	right Expr
}

type ExprUop struct {
	op      string
	operand Expr
}

// local or global
type DeclVar struct {
	pkg identifier
	variable *ExprVariable
	initval  Expr
}

type DeclConst struct {
	consts []*ExprConstVariable
}

type StmtAssignment struct {
	lefts  []Expr
	rights []Expr
}

type StmtShortVarDecl struct {
	lefts []Expr
	rights []Expr
}

type ForRangeClause struct {
	indexvar  *Relation
	valuevar  *Relation
	rangeexpr Expr
}

type ForForClause struct {
	init Stmt
	cond Stmt
	post Stmt
}

type StmtFor struct {
	// either of rng or cls is set
	rng   *ForRangeClause
	cls   *ForForClause
	block *AstCompountStmt
}

type StmtIf struct {
	simplestmt Stmt
	cond Expr
	then *AstCompountStmt
	els  Stmt
}

type StmtReturn struct {
	exprs []Expr
}

type StmtInc struct {
	operand Expr
}

type StmtDec struct {
	operand Expr
}

type AstPackageClause struct {
	name identifier
}

type AstImportSpec struct {
	path        string
}
type AstImportDecl struct {
	specs []*AstImportSpec
}

type AstCompountStmt struct {
	stmts []Stmt
}

type ExprFuncRef struct {
	funcdef *DeclFunc
}

type DeclFunc struct {
	pkg identifier
	receiver  *ExprVariable
	fname     identifier
	rettypes   []*Gtype
	params    []*ExprVariable
	isVariadic bool
	localvars []*ExprVariable
	body      *AstCompountStmt
}

type AstTopLevelDecl struct {
	// either of followings
	funcdecl  *DeclFunc // includes method declaration
	vardecl   *DeclVar
	constdecl *DeclConst
	typedecl  *DeclType
}

type AstFile struct {
	pkg     *AstPackageClause
	imports []*AstImportDecl
	decls   []*AstTopLevelDecl
}

type AstPackageRef struct {
	name identifier
	path string
}

type DeclType struct {
	name  identifier
	gtype *Gtype
}

// https://golang.org/ref/spec#Operands
type AstOperandName struct {
	pkg   identifier
	ident identifier
}

type ExprSliced struct {
	ref  *AstOperandName
	low  Expr
	high Expr
}

// Expr e.g. array[2]
type ExprArrayIndex struct {
	array   Expr
	index Expr
}

type ExprArrayLiteral struct {
	gtype  *Gtype
	values []Expr
}

type ExprTypeAssertion struct {
	expr Expr
	gtype *Gtype
}

type StmtContinue struct {

}

type StmtBreak struct {

}

type StmtExpr struct {
	expr Expr
}

type StmtDefer struct {
	expr Expr
}

type ExprVaArg struct {
	expr Expr
}

type ExprConversion struct {
	expr  Expr
	gtype *Gtype
}

type ExprCaseClause struct {
	exprs []Expr
	compound *AstCompountStmt
}

type AstSwitchStmt struct {
	cond Expr
	cases []*ExprCaseClause
	dflt *AstCompountStmt
}

type AstStructFieldLiteral struct {
	key   identifier
	value Expr
}

type ExprStructLiteral struct {
	strctname *Relation
	fields    []*AstStructFieldLiteral
	invisiblevar *ExprVariable // to have offfset for &T{}
}

type ExprStructField struct {
	strct     Expr
	fieldname identifier
}

type ExprTypeSwitchGuard struct {
	expr Expr
}
