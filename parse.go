package main

import "strconv"
import (
	"fmt"
)

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

type ExprVariable struct {
	varname string
	gtype string
	offset int // for local variable
	isGlobal bool
}

type ExprFuncall struct {
	// funcall
	fname string
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

type AstDeclVar struct {
	variable *ExprVariable // lvar or gvar
	initval  Expr
}

type AstAssignment struct {
	left  *ExprVariable // lvalue
	right Expr
}

type AstStmt struct {
	declvar    *AstDeclVar
	assignment *AstAssignment
	expr       Expr
}

type AstPkgClause struct {
	name identifier
}

type AstImport struct {
	paths []string
}

type AstCompountStmt struct {
	// compound
	stmts []*AstStmt
}

type AstFuncDef struct {
	// funcdef
	fname string
	localvars []*ExprVariable
	body *AstCompountStmt
}

type AstFile struct {
	pkg *AstPkgClause
	imports []*AstImport
	funcdefs []*AstFuncDef
	decls []*AstDeclVar
}

var localvars []*ExprVariable
var localenv  = make(map[string]*ExprVariable)
var globalvars []*ExprVariable
var globalenv = make(map[string]*ExprVariable)

var ts *TokenStream

func (stream *TokenStream) getToken(i int) interface{} {
	return ts.tokens[i]
}

func readToken() *Token {
	for {
		tok := ts.readToken()
		if !tok.isTypeSpace() {
			return tok
		}
	}
}

func unreadToken() {
	ts.unreadToken()
}

func expectPunct(punct string) {
	tok := readToken()
	if !tok.isPunct(punct) {
		errorf("punct '%s' expected but got '%s'", punct, tok.sval)
	}
}

func readFuncallArgs() []Expr {
	var r []Expr
	for {
		tok := readToken()
		if tok.isPunct(")") {
			return r
		}
		unreadToken()
		arg := parseExpr(nil)
		r = append(r, arg)
		tok = readToken()
		if tok.isPunct(")") {
			return r
		} else if tok.isPunct(",") {
			continue
		} else {
			errorf("invalid token in funcall arguments: %s", tok)
		}
	}
}

func parseIdentOrFuncall(name string) Expr {

	tok := readToken()
	if tok.isPunct("(") {
		// try funcall
		args := readFuncallArgs()

		// workaround: replace "println" -> "puts"
		if name == "println" {
			name = "puts"
		}
		return &ExprFuncall{
			fname: name,
			args:  args,
		}
	}
	unreadToken()

	if lvar,ok := localenv[name]; ok {
		return lvar
	}

	if gvar,ok := globalenv[name]; ok {
		return gvar
	}

	errorf("Undefined variable %s", name)
	return nil
}

var stringIndex = 0
var strings []*ExprStringLiteral

func newAstString(sval string) *ExprStringLiteral {
	ast := &ExprStringLiteral{
		val:   sval,
		slabel: fmt.Sprintf("L%d", stringIndex),
	}
	stringIndex++
	strings = append(strings, ast)
	return ast
}

func parseUnaryExpr() Expr {
	tok := readToken()
	switch tok.typ {
	case "ident":
		return parseIdentOrFuncall(tok.sval)
	case "string":
		return newAstString(tok.sval)
	case "number":
		ival, _ := strconv.Atoi(tok.sval)
		return &ExprNumberLiteral{
			val: ival,
		}
	default:
		errorf("unable to handle token %v", tok)
	}
	return nil
}

func priority(op string) int {
	switch op {
	case "-","+":
		return 10
	case "*":
		return 20
	default :
		errorf("unkown operator")
	}
	return 0;
}

func parseExpr(ast Expr) Expr {
	return parseExprInt(-1, ast)
}

func parseExprInt(prior int, ast Expr) Expr {
	if ast == nil {
		ast = parseUnaryExpr()
	}
	for {
		tok := readToken()
		if tok.isSemicolon() {
			return ast
		}
		if !tok.isTypePunct() {
			return ast
		}
		if tok.sval == "+" || tok.sval == "*" || tok.sval == "-" {
			prior2 := priority(tok.sval)
			if prior < prior2 {
				right := parseExprInt(prior2,nil)
				ast = &ExprBinop{
					op:    tok.sval,
					left:  ast,
					right: right,
				}
				continue
			} else {
				unreadToken()
				return ast
			}
		} else if tok.sval == "," || tok.sval == ")" { // end of funcall argument
			unreadToken()
			return ast
		} else {
			dumpToken(tok)
			errorf("unable to handle token=\"%s\"\n", tok.sval)
		}
	}

	return ast
}

func parseDeclVar(isGlobal bool) *AstStmt {
	// read varname
	tok := readToken()
	if !tok.isTypeIdent() {
		errorf("var expects ident, but got %s", tok)
	}
	varname := tok.sval

	// read type
	tok2 := readToken()
	if !tok2.isTypeIdent() {
		errorf("Type expected, but got %s", tok2)
	}
	gtype := tok2.sval

	var variable *ExprVariable
	if isGlobal {
		variable = &ExprVariable{
			varname: varname,
			gtype: gtype,
			isGlobal: true,
		}
		globalvars = append(globalvars, variable)
		globalenv[varname] = variable
	} else {
		variable = &ExprVariable{
			varname: varname,
			gtype: gtype,
		}
		localvars = append(localvars, variable)
		localenv[varname] = variable
	}

	// read "="
	tok3 := readToken()
	var initval Expr
	if tok3.isPunct("=") {
		initval = parseUnaryExpr()
	} else {
		unreadToken()
	}
	expectPunct(";")
	return &AstStmt{declvar:&AstDeclVar{
		variable: variable,
		initval:  initval,
	}}
}

func parseAssignment(left *ExprVariable) *AstAssignment {
	rexpr := parseExpr(nil)
	return &AstAssignment{
		left:  left,
		right: rexpr,
	}
}

func parseStmt() *AstStmt {
	tok := readToken()
	if tok.isKeyword("var") {
		return parseDeclVar(false)
	}
	unreadToken()
	ast := parseUnaryExpr()
	tok2 := readToken()
	if tok2.isPunct("=") {
		//assure_lvalue(ast)
		return &AstStmt{assignment:	parseAssignment(ast.(*ExprVariable))}
	}
	unreadToken()
	return &AstStmt{expr:parseExpr(ast)}
}

func parseCompoundStmt() *AstCompountStmt {
	r := &AstCompountStmt{}
	for {
		tok := readToken()
		if tok.isPunct("}") {
			return r
		}
		if tok.isSemicolon() {
			continue
		}
		unreadToken()
		stmt := parseStmt()
		r.stmts = append(r.stmts, stmt)
	}
	return nil
}

func parseFuncDef() *AstFuncDef {
	localvars = make([]*ExprVariable, 0)
	localenv = make(map[string]*ExprVariable)
	fname := readToken()
	if !fname.isTypeIdent() {
		errorf("identifer expected, but got %v", fname)
	}
	expectPunct("(")
	expectPunct(")")
	// expect Type
	expectPunct("{")
	body := parseCompoundStmt()
	_localvars := localvars
	localvars = nil
	return &AstFuncDef{
		fname: fname.sval,
		localvars: _localvars,
		body: body,
	}
}

func expectType(typ string) {
	tok := readToken()
	if tok.typ != typ {
		errorf("token type %s expected, but got %s", typ, tok)
	}
}

func expectNewline() {
	expectType("newline")
}

func parseImport() *AstImport {
	tok := readToken()
	var paths []string
	if tok.isPunct("(") {
		for {
			tok := readToken()
			if tok.isTypeString() {
				paths = append(paths, tok.sval)
				expectPunct(";")
			} else if tok.isPunct(")") {
				break
			} else {
				errorf("invalid import path %s", tok)
			}
		}
	} else {
		if !tok.isTypeString() {
			errorf("import expects package name")
		}
		paths = []string{tok.sval}
	}

	expectPunct(";")

	return &AstImport{
		paths: paths,
	}
}

func parseTopLevels() *AstFile {
	var r = &AstFile{}
	for {
		tok := readToken()
		if tok.isEOF() {
			return r
		}
		if tok.isSemicolon() {
			continue
		}
		if tok.isKeyword("package") {
			tok = readToken()
			assert(tok.isTypeIdent(), "expect ident")
			r.pkg = &AstPkgClause{name :tok.getIdent()}
			readToken()
			continue
		} else if tok.isKeyword("import") {
			ast := parseImport()
			r.imports = append(r.imports, ast)
			continue
		} else if tok.isKeyword("var") {
			decl := parseDeclVar(true)
			r.decls = append(r.decls, decl.declvar)
		} else if tok.isKeyword("func") {
			ast := parseFuncDef()
			r.funcdefs = append(r.funcdefs, ast)
			continue
		} else {
			errorf("unable to handle token %v", tok)
		}
	}
	return r
}

func parse(t *TokenStream) *AstFile {
	ts = t
	return parseTopLevels()
}
