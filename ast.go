package main

type AstPackage struct {
	name              goidentifier
	scope             *Scope
	files             []*AstFile
	namedTypes        []*DeclType
	dynamicTypes      []*Gtype
	uninferredGlobals []*ExprVariable
	uninferredLocals  []Inferrer // VarDecl, StmtShortVarDecl or RangeClause
	stringLiterals    []*ExprStringLiteral
	vars              []*DeclVar
	funcs             []*DeclFunc
	methods           map[identifier]methods
}

type AstFile struct {
	tok               *Token
	name              gostring
	packageClause     *PackageClause
	importDecls       []*ImportDecl
	DeclList          []*TopLevelDecl
	unresolved        []*Relation
	uninferredGlobals []*ExprVariable
	uninferredLocals  []Inferrer // VarDecl, StmtShortVarDecl or RangeClause
	stringLiterals    []*ExprStringLiteral
	dynamicTypes      []*Gtype
	namedTypes        []*DeclType
	methods           map[identifier]methods
}

type Expr interface {
	token() *Token
	emit()
	dump()
	getGtype() *Gtype
}

type Stmt interface {
	token() *Token
	emit()
	dump()
}

type Node interface {
	token() *Token
}

type Relation struct {
	pkg  goidentifier
	name goidentifier
	tok  *Token

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
	tok    *Token
	val    gostring
	slabel gostring
}

// local or global variable
type ExprVariable struct {
	tok        *Token
	varname    goidentifier
	gtype      *Gtype
	offset     int // for local variable
	isGlobal   bool
	isVariadic bool
}

type ExprConstVariable struct {
	tok       *Token
	name      goidentifier
	gtype     *Gtype
	val       Expr // like ExprConstExpr ?
	iotaIndex int  // for iota
}

// ident( ___ )
type ExprFuncallOrConversion struct {
	tok   *Token
	rel   *Relation
	fname goidentifier
	args  []Expr
}

type ExprMethodcall struct {
	tok      *Token
	receiver Expr
	fname    goidentifier
	args     []Expr
}

type ExprBinop struct {
	tok   *Token
	op    gostring
	left  Expr
	right Expr
}

type IrExprStringComparison struct {
	tok *Token
	op gostring
	cstringLeft  Expr
	cstringRight Expr
}

type ExprUop struct {
	tok     *Token
	op      gostring
	operand Expr
}

// call of builtin len()
type ExprLen struct {
	tok *Token
	arg Expr
}

// call of builtin cap()
type ExprCap struct {
	tok *Token
	arg Expr
}

// local or global
type DeclVar struct {
	tok      *Token
	pkg      goidentifier
	varname  *Relation
	variable *ExprVariable
	initval  Expr
}

type DeclConst struct {
	tok    *Token
	consts []*ExprConstVariable
}

type StmtAssignment struct {
	tok    *Token
	lefts  []Expr
	rights []Expr
}

type StmtShortVarDecl struct {
	tok    *Token
	lefts  []Expr
	rights []Expr
}

type ForRangeClause struct {
	tok                 *Token
	invisibleMapCounter *ExprVariable
	indexvar            Expr
	valuevar            Expr
	rangeexpr           Expr
}

type ForForClause struct {
	tok  *Token
	init Stmt
	cond Stmt
	post Stmt
}

const (
	FOR_KIND_RANGE_MAP  int = 1
	FOR_KIND_RANGE_LIST int = 2
	FOR_KIND_CLIKE      int = 3
)

type LoopLabels struct {
	labelBegin    string
	labelEndBlock string
	labelEndLoop  string
}

type StmtFor struct {
	tok  *Token
	kind int // 1:range map, 2:range list, 3:c-like
	// either of rng or cls is set
	rng    *ForRangeClause
	cls    *ForForClause
	block  *StmtSatementList
	labels *LoopLabels
	outer  *StmtFor // to manage lables in nested for-statements
}

type StmtIf struct {
	tok        *Token
	simplestmt Stmt
	cond       Expr
	then       Stmt
	els        Stmt
}

type StmtReturn struct {
	tok               *Token
	exprs             []Expr
	rettypes          []*Gtype
	labelDeferHandler gostring
}

type StmtInc struct {
	tok     *Token
	operand Expr
}

type StmtDec struct {
	tok     *Token
	operand Expr
}

type PackageClause struct {
	tok  *Token
	name goidentifier
}

type ImportSpec struct {
	tok  *Token
	path gostring
}

type ImportDecl struct {
	tok   *Token
	specs []*ImportSpec
}

type StmtSatementList struct {
	tok   *Token
	stmts []Stmt
}

type ExprFuncRef struct {
	tok     *Token
	funcdef *DeclFunc
}

type DeclFunc struct {
	tok       *Token
	pkg       goidentifier
	receiver  *ExprVariable
	fname     goidentifier
	rettypes  []*Gtype
	params    []*ExprVariable
	localvars []*ExprVariable
	body      *StmtSatementList
	stmtDefer *StmtDefer
	// every function has a defer handler
	labelDeferHandler gostring
	prologue          Emitter
}

type TopLevelDecl struct {
	tok *Token
	// either of followings
	funcdecl  *DeclFunc // includes method declaration
	vardecl   *DeclVar
	constdecl *DeclConst
	typedecl  *DeclType
}

type DeclType struct {
	tok   *Token
	name  goidentifier
	gtype *Gtype
}

// https://golang.org/ref/spec#Slice_expressions
type ExprSlice struct {
	tok        *Token
	collection Expr
	low        Expr
	high       Expr
	max        Expr
}

// Expr e.g. array[2], myap["foo"]
type ExprIndex struct {
	tok        *Token
	collection Expr
	index      Expr
}

type ExprArrayLiteral struct {
	tok    *Token
	gtype  *Gtype
	values []Expr
}

// https://golang.org/ref/spec#Composite_literals
// A slice literal describes the entire underlying array literal.
// A slice literal has the form []T{x1, x2, … xn}
type ExprSliceLiteral struct {
	tok          *Token
	gtype        *Gtype
	values       []Expr
	invisiblevar *ExprVariable // the underlying array
}

func (e *ExprSliceLiteral) getGtype() *Gtype {
	return e.gtype
}

// https://golang.org/ref/spec#Type_assertions
// x.(T)
type ExprTypeAssertion struct {
	tok   *Token
	expr  Expr   // x
	gtype *Gtype // T
}

type StmtContinue struct {
	tok    *Token
	labels *LoopLabels
}

type StmtBreak struct {
	tok    *Token
	labels *LoopLabels
}

type StmtExpr struct {
	tok  *Token
	expr Expr
}

type StmtDefer struct {
	tok   *Token
	expr  Expr
	label gostring // start of defer
}

// f( ,...slice)
type ExprVaArg struct {
	tok  *Token
	expr Expr // slice
}

type IrExprConversion struct {
	tok     *Token
	toGtype *Gtype // to
	arg     Expr   // from
}

type IrExprConversionToInterface struct {
	tok *Token
	arg Expr // dynamic type
}

// ExprCaseClause or TypeCaseClause
type ExprCaseClause struct {
	tok      *Token
	exprs    []Expr
	gtypes   []*Gtype
	compound *StmtSatementList
}

// https://golang.org/ref/spec#Switch_statements
type StmtSwitch struct {
	tok          *Token
	cond         Expr
	cases        []*ExprCaseClause
	dflt         *StmtSatementList
}

type KeyedElement struct {
	tok   *Token
	key   goidentifier // should be Expr ?
	value Expr
}

type ExprStructLiteral struct {
	tok          *Token
	strctname    *Relation
	fields       []*KeyedElement
	invisiblevar *ExprVariable // to have offfset for &T{}
}

type ExprStructField struct {
	tok       *Token
	strct     Expr
	fieldname goidentifier
}

type MapElement struct {
	tok   *Token
	key   Expr
	value Expr
}

type ExprMapLiteral struct {
	tok      *Token
	gtype    *Gtype
	elements []*MapElement
}

type ExprTypeSwitchGuard struct {
	tok  *Token
	expr Expr
}

func (node *Relation) token() *Token { return node.tok }

func (node *AstFile) token() *Token       { return node.tok }
func (node *PackageClause) token() *Token { return node.tok }
func (node *ImportSpec) token() *Token    { return node.tok }
func (node *ImportDecl) token() *Token    { return node.tok }

func (node *TopLevelDecl) token() *Token { return node.tok }
func (node *DeclVar) token() *Token      { return node.tok }
func (node *DeclConst) token() *Token    { return node.tok }
func (node *DeclFunc) token() *Token     { return node.tok }
func (node *DeclType) token() *Token     { return node.tok }

func (node *StmtFor) token() *Token          { return node.tok }
func (node *StmtIf) token() *Token           { return node.tok }
func (node *StmtReturn) token() *Token       { return node.tok }
func (node *StmtInc) token() *Token          { return node.tok }
func (node *StmtDec) token() *Token          { return node.tok }
func (node *StmtSatementList) token() *Token { return node.tok }
func (node *StmtAssignment) token() *Token   { return node.tok }
func (node *StmtShortVarDecl) token() *Token { return node.tok }
func (node *StmtContinue) token() *Token     { return node.tok }
func (node *StmtBreak) token() *Token        { return node.tok }
func (node *StmtExpr) token() *Token         { return node.tok }
func (node *StmtDefer) token() *Token        { return node.tok }
func (node *StmtSwitch) token() *Token       { return node.tok }

func (node *ExprNilLiteral) token() *Token              { return node.tok }
func (node *ExprNumberLiteral) token() *Token           { return node.tok }
func (node *ExprStringLiteral) token() *Token           { return node.tok }
func (node *ExprVariable) token() *Token                { return node.tok }
func (node *ExprConstVariable) token() *Token           { return node.tok }
func (node *ExprFuncallOrConversion) token() *Token     { return node.tok }
func (node *ExprMethodcall) token() *Token              { return node.tok }
func (node *ExprBinop) token() *Token                   { return node.tok }
func (node *ExprUop) token() *Token                     { return node.tok }
func (node *ForRangeClause) token() *Token              { return node.tok }
func (node *ForForClause) token() *Token                { return node.tok }
func (node *ExprFuncRef) token() *Token                 { return node.tok }
func (node *ExprSlice) token() *Token                   { return node.tok }
func (node *ExprIndex) token() *Token                   { return node.tok }
func (node *ExprArrayLiteral) token() *Token            { return node.tok }
func (node *ExprSliceLiteral) token() *Token            { return node.tok }
func (node *ExprTypeAssertion) token() *Token           { return node.tok }
func (node *ExprVaArg) token() *Token                   { return node.tok }
func (node *IrExprConversion) token() *Token            { return node.tok }
func (node *ExprCaseClause) token() *Token              { return node.tok }
func (node *KeyedElement) token() *Token                { return node.tok }
func (node *ExprStructLiteral) token() *Token           { return node.tok }
func (node *ExprStructField) token() *Token             { return node.tok }
func (node *ExprTypeSwitchGuard) token() *Token         { return node.tok }
func (node *ExprMapLiteral) token() *Token              { return node.tok }
func (node *ExprLen) token() *Token                     { return node.tok }
func (node *ExprCap) token() *Token                     { return node.tok }
func (node *IrExprConversionToInterface) token() *Token { return node.tok }

// Internal node made by sema
type IrStmtForRangeList struct {
	tok       *Token
	init      Stmt
	cond      Expr
	assignVar *StmtAssignment
	cond2     Expr
	incr      Stmt
	block     *StmtSatementList
	labels    *LoopLabels
}

type IrStmtRangeMap struct {
	tok        *Token
	labels     *LoopLabels
	rangeexpr  Expr
	indexvar   Expr
	valuevar   Expr
	mapCounter *ExprVariable
	initstmt   Stmt
	condition  Expr
	indexIncr  Stmt
	block      *StmtSatementList
}

type IrStmtClikeFor struct {
	tok    *Token
	cls    *ForForClause
	block  *StmtSatementList
	labels *LoopLabels
}

func (node *IrStmtForRangeList) token() *Token { return node.tok }
func (node *IrStmtRangeMap) token() *Token     { return node.tok }
func (node *IrStmtClikeFor) token() *Token     { return node.tok }
