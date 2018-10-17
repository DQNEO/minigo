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

func readIdent() identifier {
	tok := readToken()
	if !tok.isTypeIdent() {
		errorf("Identifier expected, but got %s", tok)
	}
	return tok.getIdent()
}

func expectKeyword(name string) {
	tok := readToken()
	if !tok.isKeyword(name) {
		errorf("Keyword %s expected but got %s", tok)
	}
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
			errorf("unable to handle token=\"%s\"\n", tok.sval)
		}
	}

	return ast
}

func registerVariable(varname identifier, gtype string, isGlobal bool) *ExprVariable {
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
	return variable
}

func parseDeclVar(isGlobal bool) *AstVarDecl {
	// read varname
	varname := readIdent()

	// Type or "="
	tok := readToken()
	var gtype string
	var initval Expr
	if tok.isPunct("=") {
		//var x = EXPR
		initval = parseExpr()
		gtype = "int" // should infer type
		expectPunct(";")
	} else if tok.isTypeIdent() {
		// var x T (= EXPR)
		gtype = tok.sval
		tok3 := readToken()
		if tok3.isPunct("=") {
			initval = parseExpr()
			expectPunct(";")
		} else if tok3.isPunct(";") {
			// k
		} else {
			errorf("Invalid token %s", tok3)
		}
	} else {
		errorf("Type or = expected, but got %s", tok)
	}

	variable := registerVariable(varname, gtype, isGlobal)

	return &AstVarDecl{
		variable: variable,
		initval:  initval,
	}
}

func parseConstDecl(isGlobal bool) *AstConstDecl {
	// read name
	name := readIdent()

	// Type or "="
	tok := readToken()
	var gtype string
	var val Expr
	if tok.isPunct("=") {
		// infer mode: const x = EXPR
		val = parseExpr()
		gtype = "int" // TODO: infer type
		expectPunct(";")
	} else if tok.isTypeIdent() {
		// const x T = EXPR
		gtype = tok.sval
		expectPunct("=")
		val = parseExpr()
		expectPunct(";")
	} else {
		errorf("Type or = expected, but got %s", tok)
	}

	variable := registerVariable(name, gtype, isGlobal)

	return &AstConstDecl{
		variable: variable,
		initval:  val,
	}
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
		return  &AstStmt{declvar:parseDeclVar(false)}
	} else if tok.isKeyword("const") {
		return  &AstStmt{constdecl:parseConstDecl(false)}
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

func parseFuncDef() *AstFuncDecl {
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
	return &AstFuncDecl{
		fname: fname,
		rettype:rettype,
		params: params,
		localvars: _localvars,
		body: body,
	}
}

func parseImport() *AstImportDecl {
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

	return &AstImportDecl{
		paths: paths,
	}
}

func shouldHavePackageClause() *AstPackageClause {
	expectKeyword("package")
	r := &AstPackageClause{name :readIdent()}
	expectPunct(";")
	return r
}

func mayHaveImportDecls() []*AstImportDecl {
	var r []*AstImportDecl
	for {
		tok := readToken()
		if !tok.isKeyword("import") {
			unreadToken()
			return r
		}
		r = append(r, parseImport())
	}
}

func mayHaveTopLevelDecls() []*AstTopLevelDecl {
	var r []*AstTopLevelDecl

	for {
		tok := readToken()
		if tok.isEOF() {
			return r
		}
		if tok.isKeyword("var") {
			vardecl := parseDeclVar(true)
			r = append(r, &AstTopLevelDecl{vardecl: vardecl})
		} else if tok.isKeyword("const") {
			constdecl := parseConstDecl(true)
			r = append(r, &AstTopLevelDecl{constdecl:constdecl})
		} else if tok.isKeyword("func") {
			funcdecl := parseFuncDef()
			r  = append(r, &AstTopLevelDecl{funcdecl:funcdecl})
		} else if tok.isPunct(";") {
			continue
		} else {
			errorf("unable to handle token %v", tok)
		}
	}
	return r
}

// https://golang.org/ref/spec#Source_file_organization
// Each source file consists of
// a package clause defining the package to which it belongs,
// followed by a possibly empty set of import declarations that declare packages whose contents it wishes to use,
// followed by a possibly empty set of declarations of functions, types, variables, and constants.
func parseSourceFile() *AstSourceFile {
	return &AstSourceFile{
		pkg :    shouldHavePackageClause(),
		imports: mayHaveImportDecls(),
		decls: mayHaveTopLevelDecls(),
	}
}

func newScope(outer *scope) *scope {
	return &scope{
		outer:outer,
		idents:make(map[identifier]*ExprVariable),
	}
}

func parse(t *TokenStream) *AstSourceFile {
	ts = t
	globalscope = newScope(nil)
	currentscope = globalscope
	return parseSourceFile()
}
