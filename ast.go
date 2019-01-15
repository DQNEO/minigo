package main

type Expr interface {
	emit()
	dump()
	getGtype() *Gtype
}

type Stmt interface {
	emit()
	dump()
}

type Inferer interface {
	infer()
}

type Relation struct {
	tok *Token
	name identifier

	// either of expr(var, const, funcref) or gtype
	expr  Expr
	gtype *Gtype
}

type ExprNilLiteral struct {
	tok *Token
}

type ExprNumberLiteral struct {
	tok *Token
	val int
}

type ExprStringLiteral struct {
	tok *Token
	val    string
	slabel string
}

// local or global variable
type ExprVariable struct {
	tok *Token
	varname  identifier
	gtype    *Gtype
	offset   int // for local variable
	isGlobal bool
}

type ExprConstVariable struct {
	tok *Token
	name      identifier
	gtype     *Gtype
	val       Expr // like ExprConstExpr ?
	iotaIndex int  // for iota
}

// ident(...)
type ExprFuncallOrConversion struct {
	tok   *Token
	rel   *Relation
	fname string
	args  []Expr
}

type ExprMethodcall struct {
	tok      *Token
	receiver Expr
	fname    identifier
	args     []Expr
}

type ExprBinop struct {
	tok *Token
	op    string
	left  Expr
	right Expr
}

type ExprUop struct {
	tok *Token
	op      string
	operand Expr
}

// local or global
type DeclVar struct {
	tok *Token
	pkg      identifier
	variable *ExprVariable
	initval  Expr
}

type DeclConst struct {
	tok *Token
	consts []*ExprConstVariable
}

type StmtAssignment struct {
	tok *Token
	lefts  []Expr
	rights []Expr
}

type StmtShortVarDecl struct {
	tok    *Token
	lefts  []Expr
	rights []Expr
}

type ForRangeClause struct {
	tok *Token
	indexvar  *Relation
	valuevar  *Relation
	rangeexpr Expr
}

type ForForClause struct {
	tok *Token
	init Stmt
	cond Stmt
	post Stmt
}

type StmtFor struct {
	tok *Token
	// either of rng or cls is set
	rng   *ForRangeClause
	cls   *ForForClause
	block *StmtSatementList
}

type StmtIf struct {
	tok *Token
	simplestmt Stmt
	cond       Expr
	then       *StmtSatementList
	els        Stmt
}

type StmtReturn struct {
	tok *Token
	exprs []Expr
}

type StmtInc struct {
	tok *Token
	operand Expr
}

type StmtDec struct {
	tok *Token
	operand Expr
}

type PackageClause struct {
	tok *Token
	name identifier
}

type ImportSpec struct {
	tok *Token
	path string
}

type ImportDecl struct {
	tok *Token
	specs []*ImportSpec
}

type StmtSatementList struct {
	tok *Token
	stmts []Stmt
}

type ExprFuncRef struct {
	tok *Token
	funcdef *DeclFunc
}

type DeclFunc struct {
	tok *Token
	pkg        identifier
	receiver   *ExprVariable
	fname      identifier
	rettypes   []*Gtype
	params     []*ExprVariable
	isVariadic bool
	localvars  []*ExprVariable
	body       *StmtSatementList
}

type TopLevelDecl struct {
	tok *Token
	// either of followings
	funcdecl  *DeclFunc // includes method declaration
	vardecl   *DeclVar
	constdecl *DeclConst
	typedecl  *DeclType
}

type SourceFile struct {
	packageClause *PackageClause
	importDecls   []*ImportDecl
	topLevelDecls []*TopLevelDecl
}

type DeclType struct {
	tok *Token
	name  identifier
	gtype *Gtype
}

// https://golang.org/ref/spec#Slice_expressions
type ExprSlice struct {
	tok *Token
	collection Expr
	low  Expr
	high Expr
}

// Expr e.g. array[2]
type ExprIndex struct {
	tok        *Token
	collection Expr
	index      Expr
}

type ExprArrayLiteral struct {
	tok *Token
	gtype  *Gtype
	values []Expr
}

type ExprTypeAssertion struct {
	tok *Token
	expr  Expr
	gtype *Gtype
}

type StmtContinue struct {
	tok *Token
}

type StmtBreak struct {
	tok *Token
}

type StmtExpr struct {
	tok *Token
	expr Expr
}

type StmtDefer struct {
	tok *Token
	expr Expr
}

type ExprVaArg struct {
	tok *Token
	expr Expr
}

type ExprConversion struct {
	tok *Token
	gtype *Gtype
	expr  Expr
}

type ExprCaseClause struct {
	tok *Token
	exprs    []Expr
	compound *StmtSatementList
}

type StmtSwitch struct {
	tok *Token
	cond  Expr
	cases []*ExprCaseClause
	dflt  *StmtSatementList
}

type KeyedElement struct {
	tok *Token
	key   identifier
	value Expr
}

type ExprStructLiteral struct {
	tok *Token
	strctname    *Relation
	fields       []*KeyedElement
	invisiblevar *ExprVariable // to have offfset for &T{}
}

type ExprStructField struct {
	tok       *Token
	strct     Expr
	fieldname identifier
}

type ExprTypeSwitchGuard struct {
	tok *Token
	expr Expr
}
