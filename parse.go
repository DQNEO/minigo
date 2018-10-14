package main

import "strconv"
import (
	"fmt"
)

type AstExpr struct {
	//
	typ string
	// int
	ival int
	// unary
	operand *AstExpr
	// binop
	op    string
	left  *AstExpr
	right *AstExpr
	// string
	sval   string
	slabel string
	// funcall
	fname string
	args  []*AstExpr
	// lvar (local var)
	varname string
	gtype string
	offset int
}

type AstDeclLocalVar struct {
	localvar *AstExpr
	initval *AstExpr
}

type AstAssignment struct {
	left  *AstExpr
	right *AstExpr
}

type AstStmt struct {
	decl *AstDeclLocalVar
	assignment *AstAssignment
	expr *AstExpr
}

type AstPkgDecl struct {
	name string
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
	localvars []*AstExpr
	body *AstCompountStmt
}

type AstFile struct {
	pkg *AstPkgDecl
	imports []*AstImport
	funcdefs []*AstFuncDef
}

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
	if !tok.isTypePunct() {
		errorf("token type punct expected, but got %v", tok)
	}
	if !tok.isPunct(punct) {
		errorf("punct '%s' expected but got '%s'", punct, tok.sval)
	}
}

func readFuncallArgs() []*AstExpr {
	var r []*AstExpr
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

func parseIdentOrFuncall(name string) *AstExpr {

	if lvar,ok := localenv[name]; ok {
		return lvar
	}

	tok := readToken()
	if tok.isPunct("(") {
		// try funcall
		args := readFuncallArgs()

		// workaround: replace "println" -> "puts"
		if name == "println" {
			name = "puts"
		}
		return &AstExpr{
			typ:   "funcall",
			fname: name,
			args:  args,
		}
	}

	errorf("TBD")
	return nil
}

var stringIndex = 0
var strings []*AstExpr

func newAstString(sval string) *AstExpr {
	ast := &AstExpr{
		typ:    "string",
		sval:   sval,
		slabel: fmt.Sprintf("L%d", stringIndex),
	}
	stringIndex++
	strings = append(strings, ast)
	return ast
}

func parseUnaryExpr() *AstExpr {
	tok := readToken()
	switch tok.typ {
	case "string":
		return newAstString(tok.sval)
	case "ident","keyword":
		return parseIdentOrFuncall(tok.sval)
	case "number":
		ival, _ := strconv.Atoi(tok.sval)
		return &AstExpr{
			typ:  "int",
			ival: ival,
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

func parseExpr(ast *AstExpr) *AstExpr {
	return parseExprInt(-1, ast)
}

func parseExprInt(prior int, ast *AstExpr) *AstExpr {
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
				ast = &AstExpr{
					typ:   "binop",
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

var localvars []*AstExpr
var localenv map[string]*AstExpr

func readDeclLocalVar() *AstDeclLocalVar {
	tok := readToken()
	if !tok.isTypeIdent() {
		errorf("var expects ident, but got %s", tok)
	}
	varname := tok.sval

	tok2 := readToken()
	if !tok2.isTypeIdent() {
		errorf("Type expected, but got %s", tok2)
	}
	gtype := tok2.sval

	lvar := &AstExpr{
		typ: "lvar",
		varname: varname,
		gtype: gtype,
	}
	localvars = append(localvars, lvar)
	localenv[varname] = lvar
	return &AstDeclLocalVar{
		localvar: lvar,
	}
}

func parseAssignment(left *AstExpr) *AstAssignment {
	rexpr := parseExpr(nil)
	return &AstAssignment{
		left:  left,
		right: rexpr,
	}
}

func parseStmt() *AstStmt {
	tok := readToken()
	if tok.isKeyword("var") {
		return &AstStmt{decl:readDeclLocalVar()}
	}
	unreadToken()
	ast := parseUnaryExpr()
	tok2 := readToken()
	if tok2.isPunct("=") {
		//assure_lvalue(ast)
		assert(ast.typ == "lvar", "assure lvaue")
		return &AstStmt{assignment:	parseAssignment(ast)}
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
	localvars = make([]*AstExpr, 0)
	localenv = make(map[string]*AstExpr)
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
			r.pkg = &AstPkgDecl{name :tok.sval}
			readToken()
			continue
		} else if tok.isKeyword("import") {
			ast := parseImport()
			r.imports = append(r.imports, ast)
			continue
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
