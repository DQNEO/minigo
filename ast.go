package main

import "fmt"

type Expr interface {
	emit()
	dump()
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
	varname         identifier
	typeConstructor interface{}
	gtype           *Gtype
	offset          int // for local variable
	isGlobal        bool
}

type ExprConstVariable struct {
	name            identifier
	typeConstructor interface{}
	gtype           *Gtype
	val             Expr // like ExprConstExpr ?
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
	op      string
	operand Expr
}

// local or global
type AstVarDecl struct {
	variable *ExprVariable
	initval  Expr
}

type AstConstDecl struct {
	variable *ExprConstVariable
}

type AstAssignment struct {
	left  Expr
	right Expr
}

type AstForStmt struct {
	idents []identifier
	list   Expr
	// or
	left   *AstStmt
	middle *AstStmt
	right  *AstStmt

	block *AstCompountStmt
}

type AstIfStmt struct {
	cond Expr
	then *AstCompountStmt
	els  *AstStmt
}

type AstStmt struct {
	compound   *AstCompountStmt
	declvar    *AstVarDecl
	constdecl  *AstConstDecl
	typedecl   *AstTypeDecl
	assignment *AstAssignment
	forstmt    *AstForStmt
	ifstmt     *AstIfStmt
	expr       Expr
}

type AstPackageClause struct {
	name identifier
}

type AstImportSpec struct {
	packageName identifier
	path string
}
type AstImportDecl struct {
	specs[] *AstImportSpec
}

type AstCompountStmt struct {
	// compound
	stmts []*AstStmt
}

type AstFuncDecl struct {
	// funcdef
	fname     identifier
	rettype   string
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
	pkg          *AstPackageClause
	imports      []*AstImportDecl
	decls        []*AstTopLevelDecl
}

type AstPackageRef struct {
	name identifier
	path string
}

type AstTypeDef struct {
	name            identifier  // we need this ?
	typeConstructor interface{} // (identifier | QualifiedIdent) | TypeLiteral
}

type AstTypeDecl struct {
	typedef *AstTypeDef
	gtype   *Gtype // resolved later
}

type GTYPE_TYPE int
const (
	G_UNKOWN GTYPE_TYPE = iota
	G_INT
	G_BOOL
	G_BYTE
	G_STRUCT
	G_ARRAY
	G_SLICE
	G_STRING
	G_MAP
	G_POINTER
)

type Gtype struct {
	typ       GTYPE_TYPE
	size      int // for scalar type like int, bool, byte
	ptr       *Gtype // for array, pointer, etc
	structdef *AstStructDef // for struct type
	length    int // for fixed array
}

func (gtype *Gtype) String() string {
	switch gtype.typ {
	case G_INT:
		return "int"
	case G_ARRAY:
		elm := gtype.ptr
		return fmt.Sprintf("[]%s", elm)
	default:
		errorf("unkown type: %d", gtype.typ)
	}
	return ""
}

type AstInterfaceDef struct {
	methods []identifier // for interface
}

type AstStructDef struct {
	fields []*StructField // for struct
}
type StructField struct {
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

type ExprIndexAccess struct {
	variable Expr // identexpr or variableexpr
	index Expr
}

func (e *ExprIndexAccess) dump() {
	errorf("TBD")

}

type ExprArrayLiteral struct {
	gtype  *Gtype
	values []Expr
}

func (e ExprArrayLiteral) emit() {
	errorf("DO NOT EMIT")
}
