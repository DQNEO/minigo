package main

import "strconv"
import (
	"fmt"
)

var globalscope *scope
var currentscope *scope

type scope struct {
	idents map[identifier]*ExprVariable
	outer *scope
}

func (sc *scope) get(name identifier) *ExprVariable {
	for s := sc; s != nil; s = s.outer {
		v := s.idents[name]
		if v != nil {
			return v
		}
	}
	return nil
}

func (sc *scope) set(name identifier, variable *ExprVariable) {
	sc.idents[name] = variable
}

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
	gtype string
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
type AstDeclVar struct {
	variable *ExprVariable
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
	fname identifier
	rettype string
	params []*ExprVariable
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
var globalvars []*ExprVariable

var ts *TokenStream

func (stream *TokenStream) getToken(i int) interface{} {
	return ts.tokens[i]
}

func readToken() *Token {
	tok := ts.readToken()
	return tok
}

func unreadToken() {
	ts.unreadToken()
}

func expectPunct(punct string) {
	tok := readToken()
	if !tok.isPunct(punct) {
		errorf("punct '%s' expected but got '%s'", punct, tok)
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
		arg := parseExpr()
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

func parseIdentOrFuncall(name identifier) Expr {

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


	variable := currentscope.get(name)
	if variable == nil {
		errorf("Undefined variable %s", name)
		return nil
	}

	return variable
}

var stringIndex = 0
var stringLiterals []*ExprStringLiteral

func newAstString(sval string) *ExprStringLiteral {
	ast := &ExprStringLiteral{
		val:   sval,
		slabel: fmt.Sprintf("L%d", stringIndex),
	}
	stringIndex++
	stringLiterals = append(stringLiterals, ast)
	return ast
}

func parsePrim() Expr {
	tok := readToken()
	switch tok.typ {
	case T_IDENT:
		return parseIdentOrFuncall(tok.getIdent())
	case T_STRING:
		return newAstString(tok.sval)
	case T_INT:
		ival, _ := strconv.Atoi(tok.sval)
		return &ExprNumberLiteral{
			val: ival,
		}
	default:
		errorf("unable to handle token %v", tok)
	}
	return nil
}

func parseUnaryExpr() Expr {
	return parsePrim()
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

func parseExpr() Expr {
	return parseExprInt(-1)
}

func parseExprInt(prior int) Expr {
	ast := parseUnaryExpr()
	for {
		tok := readToken()
		if tok.isSemicolon() {
			unreadToken()
			return ast
		}

		if tok.sval == "+" || tok.sval == "*" || tok.sval == "-" {
			prior2 := priority(tok.sval)
			if prior < prior2 {
				right := parseExprInt(prior2)
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
	varname := tok.getIdent()

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
	} else {
		variable = &ExprVariable{
			varname: varname,
			gtype: gtype,
		}
		localvars = append(localvars, variable)
	}
	currentscope.set(varname, variable)

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

func parseAssignment() *AstAssignment {
	tleft := readToken()
	variable := currentscope.get(tleft.getIdent())
	expectPunct("=")
	rexpr := parseExpr()
	expectPunct(";")
	return &AstAssignment{
		left:  variable,
		right: rexpr,
	}
}

func parseStmt() *AstStmt {
	tok := readToken()
	if tok.isKeyword("var") {
		return parseDeclVar(false)
	}
	tok2 := readToken()

	if tok2.isPunct("=") {
		unreadToken()
		unreadToken()
		//assure_lvalue(ast)
		return &AstStmt{assignment:	parseAssignment()}
	}
	unreadToken()
	unreadToken()
	return &AstStmt{expr:parseExpr()}
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
	currentscope = newScope(globalscope)
	fname := readToken().getIdent()
	expectPunct("(")
	var params []*ExprVariable

	tok := readToken()
	if !tok.isPunct(")") {
		unreadToken()
		for {
			tok := readToken()
			pname := tok.getIdent()
			ptype := readToken() //type
			// assureType(tok.sval)
			variable := &ExprVariable{
				varname:pname,
				gtype: ptype.sval,
			}
			params = append(params, variable)
			currentscope.set(pname, variable)
			tok = readToken()
			if tok.isPunct(")") {
				break
			}
			if !tok.isPunct(",") {
				errorf("Invalid token %s", tok)
			}
		}
	}

	// read func rettype
	tok = readToken()
	var rettype string
	if tok.isTypeIdent() {
		// rettype
		rettype = tok.sval
		expectPunct("{")
	} else {
		assert(tok.isPunct("{"), "begin of func body")
	}
	body := parseCompoundStmt()
	_localvars := localvars
	localvars = nil
	currentscope = globalscope
	return &AstFuncDef{
		fname: fname,
		rettype:rettype,
		params: params,
		localvars: _localvars,
		body: body,
	}
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

func newScope(outer *scope) *scope {
	return &scope{
		outer:outer,
		idents:make(map[identifier]*ExprVariable),
	}
}

func parse(t *TokenStream) *AstFile {
	ts = t
	globalscope = newScope(nil)
	currentscope = globalscope
	return parseTopLevels()
}
