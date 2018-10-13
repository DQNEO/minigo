package main

import "strconv"
import (
	"fmt"
)

type Ast struct {
	//
	typ string
	// int
	ival int
	// unary
	operand *Ast
	// binop
	op    string
	left  *Ast
	right *Ast
	// string
	sval   string
	slabel string
	// compound
	stmts []*Ast
	// funcall
	fname string
	args  []*Ast
	// package
	pkgname string
	// imports
	packages []string
	// decl
	declvar *Ast
	// lvar (local var)
	varname string
	gtype string
	offset int
}

type AstPkgDecl struct {
	name string
}

type AstImport struct {
	paths []string
}

type AstFuncDef struct {
	// funcdef
	fname string
	localvars []*Ast
	body *Ast
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

func readFuncallArgs() []*Ast {
	var r []*Ast
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

func parseIdentOrFuncall(name string) *Ast {

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
		return &Ast{
			typ:   "funcall",
			fname: name,
			args:  args,
		}
	}

	errorf("TBD")
	return nil
}

var stringIndex = 0
var strings []*Ast

func newAstString(sval string) *Ast {
	ast := &Ast{
		typ:    "string",
		sval:   sval,
		slabel: fmt.Sprintf("L%d", stringIndex),
	}
	stringIndex++
	strings = append(strings, ast)
	return ast
}

func parseUnaryExpr() *Ast {
	tok := readToken()
	switch tok.typ {
	case "string":
		return newAstString(tok.sval)
	case "ident","keyword":
		return parseIdentOrFuncall(tok.sval)
	case "number":
		ival, _ := strconv.Atoi(tok.sval)
		return &Ast{
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

func parseExpr() *Ast {
	return parseExprInt(-1)
}

func parseExprInt(prior int) *Ast {
	ast := parseUnaryExpr()
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
				right := parseExprInt(prior2)
				ast = &Ast{
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
		} else if tok.sval == "=" {
			//assure_lvalue(ast)
			assert(ast.typ == "lvar", "assure lvaue")
			rexpr := parseExpr()
			return &Ast{
				typ:"assign",
				left: ast,
				right:rexpr,
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

var localvars []*Ast
var localenv map[string]*Ast

func readDeclVar() *Ast {
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

	lvar := &Ast{
		typ: "lvar",
		varname: varname,
		gtype: gtype,
	}
	localvars = append(localvars, lvar)
	localenv[varname] = lvar
	return &Ast{
		typ:     "decl",
		declvar: lvar,
	}
}

func parseStmt() *Ast {
	tok := readToken()
	if tok.isKeyword("var") {
		return readDeclVar()
	}
	unreadToken()
	return parseExpr()
}

func parseCompoundStmt() []*Ast {
	var r []*Ast
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
		r = append(r, stmt)
	}
	return nil
}

func parseFuncDef() *AstFuncDef {
	localvars = make([]*Ast, 0)
	localenv = make(map[string]*Ast)
	fname := readToken()
	if !fname.isTypeIdent() {
		errorf("identifer expected, but got %v", fname)
	}
	expectPunct("(")
	expectPunct(")")
	//skipSpaceToken()
	// expect Type
	expectPunct("{")
	stmts := parseCompoundStmt()
	_localvars := localvars
	localvars = nil
	return &AstFuncDef{
		fname: fname.sval,
		localvars: _localvars,
		body: &Ast{
			typ:   "compound",
			stmts: stmts,
		},
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
	//skipSpaceToken()
	tok := readToken()
	if !tok.isTypeString() {
		errorf("import expects package name")
	}
	packageName := tok.sval
	expectNewline()
	return &AstImport{
		paths: []string{packageName},
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
			readToken() // expect newline
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
