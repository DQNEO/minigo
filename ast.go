package main

type Expr interface {
	emit()
	dump()
	getGtype() *Gtype
}

type Stmt interface {
	emit()
}

type ExprNilLiteral struct {
}

func (e *ExprNilLiteral) emit() {
	emit("mov $0, %%rax")
}

func (e *ExprNilLiteral) dump() {
	debugf("nil")
}

func (e *ExprNilLiteral) getGtype() *Gtype {
	return nil
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
type AstVarDecl struct {
	variable *ExprVariable
	initval  Expr
}

type AstConstDecl struct {
	consts []*ExprConstVariable
}

type AstAssignment struct {
	lefts  []Expr
	rights []Expr
}

type AstShortAssignment struct {
	lefts []Expr
	rights []Expr
}

func (ast *AstShortAssignment) emit() {
	a := &AstAssignment{
		lefts:ast.lefts,
		rights:ast.rights,
	}
	a.emit()
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

type AstForStmt struct {
	// either of rng or cls is set
	rng   *ForRangeClause
	cls   *ForForClause
	block *AstCompountStmt
}

type AstIfStmt struct {
	cond Expr
	then *AstCompountStmt
	els  Stmt
}

type AstReturnStmt struct {
	exprs []Expr
}

type AstIncrStmt struct {
	operand Expr
}

type AstDecrStmt struct {
	operand Expr
}

type AstPackageClause struct {
	name identifier
}

type AstImportSpec struct {
	packageName identifier
	path        string
}
type AstImportDecl struct {
	specs []*AstImportSpec
}

type AstCompountStmt struct {
	stmts []Stmt
}

type ExprFuncRef struct {
	funcdef *AstFuncDecl
}

func (f *ExprFuncRef) emit() {
	emit("mov $1, %%rax") // emit 1 for now.  @FIXME
}

func (f *ExprFuncRef) dump() {
	f.funcdef.dump()
}

type AstFuncDecl struct {
	receiver  *ExprVariable
	fname     identifier
	rettypes   []*Gtype
	params    []*ExprVariable
	localvars []*ExprVariable
	body      *AstCompountStmt
}

type AstTopLevelDecl struct {
	// either of followings
	funcdecl  *AstFuncDecl
	vardecl   *AstVarDecl
	constdecl *AstConstDecl
	typedecl  *AstTypeDecl
}

type AstSourceFile struct {
	pkg     *AstPackageClause
	imports []*AstImportDecl
	decls   []*AstTopLevelDecl
}

type AstPackageRef struct {
	name identifier
	path string
}

type AstTypeDecl struct {
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

func (e *ExprSliced) dump() {
	errorf("TBD")
}
func (e *ExprSliced) emit() {
	errorf("TBD")
}

// Expr e.g. array[2]
type ExprArrayIndex struct {
	array   Expr
	index Expr
}

func (e *ExprArrayIndex) dump() {
	errorf("TBD")

}

type ExprArrayLiteral struct {
	gtype  *Gtype
	values []Expr
}

func (e ExprArrayLiteral) emit() {
	errorf("DO NOT EMIT")
}
