package main

import (
	"fmt"
)

var gBool = &Gtype{typ:"bool",}
var gInt = &Gtype{typ:"int",}
var gString = &Gtype{typ:"string",}

var predeclaredTypes = []*Gtype{
	gInt,
	gBool,
	gString,
}

var predeclaredConsts = []*AstConstDecl{
	&AstConstDecl{
		variable:&ExprConstVariable{
			name:"true",
			gtype:gBool,
			val: &ExprNumberLiteral{1},
		},
	},
	&AstConstDecl{
		variable:&ExprConstVariable{
			name:"false",
			gtype:gBool,
			val: &ExprNumberLiteral{0},
		},
	},
}


var universeblockscope *scope
var packageblockscope *scope
var currentscope *scope

type scope struct {
	idents map[identifier]interface{}
	outer *scope
}

func (sc *scope) get(name identifier) interface{} {
	for s := sc; s != nil; s = s.outer {
		v := s.idents[name]
		if v != nil {
			return v
		}
	}
	return nil
}

func (sc *scope) set(name identifier, v interface{}) {
	if v == nil {
		panic("nil cannot be set")
	}
	sc.idents[name] = v
}

func (sc *scope) getGtype(name identifier) *Gtype {
	v := sc.get(name)
	if v == nil {
		errorf("type %s is not defined", name)
	}
	gtype, ok := v.(*Gtype)
	if !ok {
		errorf("type %s is not defined", name)
	}
	return gtype
}

var localvars []*ExprVariable
var globalvars []*ExprVariable

var ts *TokenStream

func (stream *TokenStream) getToken(i int) interface{} {
	return ts.tokens[i]
}

func peekToken() *Token {
	tok := ts.peekToken()
	return tok
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


func expect(punct string) {
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

// https://golang.org/ref/spec#Operands
type AstOperandName struct {
	pkg identifier
	ident identifier
}

func parseIdentOrFuncall(firstIdent identifier) Expr {

	// https://golang.org/ref/spec#QualifiedIdent
	// read QualifiedIdent
	tok := readToken()
	var pkg identifier
	var ident identifier
	if tok.isPunct(".") {
		pkg = firstIdent
		ident = readIdent()
	} else {
		unreadToken()
		pkg = ""
		ident = firstIdent
	}

	operand := AstOperandName{
		pkg:   pkg,
		ident: ident,
	}

	tok = readToken()
	if tok.isPunct("(") {
		// try funcall
		args := readFuncallArgs()

		if operand.pkg == "" && operand.ident == "println" {
			// dirty hack: "println" -> "puts"
			operand.ident = "puts"
		} else if operand.pkg == "fmt" && operand.ident == "Printf" {
			// dirty hack: "fmt" -> "Printf"
			operand.ident = "printf"
		}
		return &ExprFuncall{
			fname: operand.ident,
			args:  args,
		}
	}
	unreadToken()


	v := currentscope.get(firstIdent)
	if v == nil{
		errorf("Undefined variable %s", firstIdent)
		return nil
	}
	variable, ok := v.(*ExprVariable)
	if ok {
		return variable
	}
	variable2, ok := v.(*ExprConstVariable)
	if ok {
		return variable2
	}

	return nil
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
	switch  {
	case tok.isTypeIdent():
		return parseIdentOrFuncall(tok.getIdent())
	case tok.isTypeString():
		return newAstString(tok.sval)
	case tok.isTypeInt():
		ival := tok.getIntval()
		return &ExprNumberLiteral{
			val: ival,
		}
	case tok.isPunct("["):
		return parseArrayLiteral()
	default:
		errorf("unable to handle %s", tok)
	}
	return nil
}

type ExprArrayLiteral struct {
	gtype *Gtype
	values []Expr
}

func (e ExprArrayLiteral) emit() {

}

func (e ExprArrayLiteral) dump() {

}

func parseArrayLiteral() Expr {
	debugf("func %s start with %s" , "parseArrayLiteral", peekToken().String())
	debugNest++
	defer func() {
		debugNest--
		debugf("func %s end", "parseArrayLiteral")
	}()
	expect("]")
	typ := parseType()
	expect("{")
	var values []Expr
	for {
		tok := readToken()
		if tok.isPunct("}") {
			break
		}
		var v Expr
		if tok.isTypeString() {
			v = &ExprStringLiteral{val: tok.sval}
		} else if tok.isTypeInt() {
			v = &ExprNumberLiteral{val: tok.getIntval()}
		}
		tok = readToken()
		if tok.isPunct(",") {
			continue
		} else if tok.isPunct("}") {
			break
		} else {
			errorf("unpexpected %s", tok)
		}
		values = append(values, v)
	}

	gtype := &Gtype{
		typ:   "array",
		count: len(values),
		ptr:   typ,
	}

	r := &ExprArrayLiteral{
		gtype: gtype,
		values: values,
	}
	debugAstConstructed(r)
	return r
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
	debugf("func %s start with %s" , "parseExpr", peekToken().String())
	debugNest++
	defer func() {
		debugNest--
		debugf("func %s end", "parseExpr")
	}()
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
		} else if tok.sval == "," || tok.sval == ")" || tok.sval == "{" || tok.sval == "}" || tok.isPunct(";"){ // end of funcall argument
			unreadToken()
			return ast
		} else {
			tok.errorf("Unexpected")
		}
	}

	return ast
}

func newVariable(varname identifier, gtype *Gtype, isGlobal bool) *ExprVariable {
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
	return variable
}

// https://golang.org/ref/spec#Type
func parseType() *Gtype {
	for {
		tok := readToken()
		if tok.isPunct("*") {
			// pointer
		} else if tok.isKeyword("struct") {

		} else if tok.isTypeIdent() {
			typename := tok.getIdent()
			gtype := currentscope.getGtype(typename)
			return gtype
		} else if tok.isPunct("[") {
		} else if tok.isPunct("]") {

		} else {
			unreadToken()
			break
		}

	}
	return gInt // FIXME
}

func parseDeclVar(isGlobal bool) *AstVarDecl {
	// read name
	name := readIdent()

	// "=" or Type
	tok := readToken()
	var gtype *Gtype
	var initval Expr
	if tok.isPunct("=") {
		//var x = EXPR
		initval = parseExpr()
		gtype = gInt  // FIXME: infer type
		expect(";")
	} else {
		unreadToken()
		// expect Type
		gtype = parseType()
		tok3 := readToken()
		if tok3.isPunct("=") {
			initval = parseExpr()
			expect(";")
		} else if tok3.isPunct(";") {
			// k
		} else {
			errorf("Invalid token %s", tok3)
		}
	}

	variable := newVariable(name, gtype, isGlobal)
	currentscope.set(name, variable)

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
	var gtype *Gtype
	var val Expr
	if tok.isPunct("=") {
		// infer mode: const x = EXPR
		val = parseExpr()
		gtype = gInt// FIXME: infer type
		expect(";")
	} else {
		unreadToken()
		// expect Type
		gtype = parseType()
		// const x T = EXPR
		expect("=")
		val = parseExpr()
		expect(";")
	}

	variable := &ExprConstVariable{
		name: name,
		gtype: gtype,
		val: val,
	}
	currentscope.set(name, variable)

	return &AstConstDecl{
		variable: variable,
	}
}

func parseAssignment() *AstAssignment {
	tleft := readToken()
	v := currentscope.get(tleft.getIdent())
	variable, ok := v.(*ExprVariable)
	if !ok {
		errorf("%s is not a variable", tleft)
	}
	expect("=")
	rexpr := parseExpr()
	expect(";")
	return &AstAssignment{
		left:  variable,
		right: rexpr,
	}
}

func parseIdentList() []identifier {
	var r []identifier
	for {
		tok := readToken()
		if tok.isTypeIdent() {
			r = append(r, tok.getIdent())
		} else if len(r) == 0 {
			// at least one ident is needed
			tok.errorf("Ident expected")
		}

		tok = readToken()
		if tok.isPunct(",") {
			continue
		} else {
			unreadToken()
			return r
		}
	}
	return r
}

func parseForStmt() *AstForStmt {
	debugf("func %s start with %s" , "parseForStmt", peekToken())
	debugNest++
	defer func() {
		debugNest--
		debugf("func %s end", "parseForStmt")
	}()
	var r = &AstForStmt{}
	// Assume "range" style
	idents := parseIdentList()
	r.idents = idents
	expect(":=")
	expectKeyword("range")
	r.list = parseExpr()
	expect("{")
	r.block = parseCompoundStmt()
	return r
}

func parseStmt() *AstStmt {
	tok := readToken()
	if tok.isKeyword("var") {
		return  &AstStmt{declvar:parseDeclVar(false)}
	} else if tok.isKeyword("const") {
		return  &AstStmt{constdecl:parseConstDecl(false)}
	} else if tok.isKeyword("type") {
		return  &AstStmt{typedecl:parseTypeDecl()}
	} else if tok.isKeyword("for") {
		return  &AstStmt{forstmt:parseForStmt()}
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
	debugf("func %s start with %s" , "parseFuncDef", peekToken())
	debugNest++
	defer func() {
		debugNest--
		debugf("func %s end", "parseFuncDef")
	}()
	localvars = make([]*ExprVariable, 0)
	currentscope = newScope(packageblockscope)
	fname := readToken().getIdent()
	expect("(")
	var params []*ExprVariable

	tok := readToken()
	if !tok.isPunct(")") {
		unreadToken()
		for {
			tok := readToken()
			pname := tok.getIdent()
			ptype := parseType()
			// assureType(tok.sval)
			variable := &ExprVariable{
				varname:pname,
				gtype: ptype,
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
		expect("{")
	} else {
		assert(tok.isPunct("{"), "begin of func body")
	}
	body := parseCompoundStmt()
	_localvars := localvars
	localvars = nil
	currentscope = packageblockscope
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
				expect(";")
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

	expect(";")

	return &AstImportDecl{
		paths: paths,
	}
}

func shouldHavePackageClause() *AstPackageClause {
	expectKeyword("package")
	r := &AstPackageClause{name :readIdent()}
	expect(";")
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

type Gtype struct {
	typ     identifier // "int", "string", "struct" , "interface",...
	ptr *Gtype
	count int // for fixed array
	methods []identifier // for interface
	fields []*StructField // for struct
}

type StructField struct {
	name identifier
	gtype *Gtype
}

type AstTypeDecl struct {
	name  identifier
	gtype *Gtype
}

// read after "struct" token
func parseStructDef() *Gtype {
	expect("{")
	var fields []*StructField
	for {
		tok := readToken()
		if tok.isPunct("}") {
			break
		}
		fieldname := tok.getIdent()
		fieldtyep := parseType()
		fields = append(fields, &StructField{
			name: fieldname,
			gtype: fieldtyep,
		})
		expect(";")
	}
	expect(";")
	return &Gtype{
		typ:"struct",
		fields: fields,
	}
}

func parseInterfaceDef() *Gtype {
	expect("{")
	var methods []identifier
	for {
		tok := readToken()
		if tok.isPunct("}") {
			break
		}
		fname := tok.getIdent()
		expect("(")
		expect(")")
		expect(";")
		methods = append(methods, fname)
	}
	expect(";")
	return &Gtype{
		typ:"interface",
		methods: methods,
	}
}

func parseTypeDecl() *AstTypeDecl  {
	name := readIdent()
	tok := readToken()
	var gtype *Gtype
	if tok.isKeyword("struct" ) {
		gtype = parseStructDef()
	} else if tok.isKeyword("interface")  {
		gtype = parseInterfaceDef()
	} else {
		ident := tok.getIdent() // "int", "bool", etc
		currentscope.getGtype(ident) // check existence
		gtype = currentscope.getGtype(ident)
	}
	currentscope.set(name, gtype)
	return &AstTypeDecl{
		name:  name,
		gtype: gtype,
	}
}

func parseTopLevelDecl(tok *Token) *AstTopLevelDecl {
	debugf("func %s start with %s" , "parseTopLevelDecl", peekToken())
	debugNest++
	defer func() {
		debugNest--
		debugf("func %s end", "parseTopLevelDecl")
	}()
	var r *AstTopLevelDecl
	switch  {
	case tok.isKeyword("var"):
		vardecl := parseDeclVar(true)
		r = &AstTopLevelDecl{vardecl: vardecl}
	case tok.isKeyword("const"):
		constdecl := parseConstDecl(true)
		r = &AstTopLevelDecl{constdecl:constdecl}
	case tok.isKeyword("func"):
		funcdecl := parseFuncDef()
		r =  &AstTopLevelDecl{funcdecl:funcdecl}
	case tok.isKeyword("type"):
		typedecl := parseTypeDecl()
		r = &AstTopLevelDecl{typedecl:typedecl}
	default:
		errorf("TBD: unable to handle token %v", tok)
	}

	debugAstConstructed(r)
	return r
}

func debugAstConstructed(ast interface{}) {
	debugPrintVar("Ast constructed", ast)
}

func parseTopLevelDecls() []*AstTopLevelDecl {
	debugf("func %s start with %s" , "parseTopLevelDecls", peekToken())
	debugNest++
	defer func() {
		debugNest--
		debugf("func %s end", "parseTopLevelDecls")
	}()
	var r []*AstTopLevelDecl
	for {
		tok := readToken()
		if tok.isEOF() {
			return r
		}

		if tok.isPunct(";") {
			continue
		}
		ast := parseTopLevelDecl(tok)
		r = append(r, ast)
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
		decls:   parseTopLevelDecls(),
	}
}

func newScope(outer *scope) *scope {
	return &scope{
		outer:outer,
		idents:make(map[identifier]interface{}),
	}
}

func newUniverseBlockScope() *scope {
	r := newScope(nil)
	for _, t := range predeclaredTypes {
		r.set(t.typ, t)
	}
	for _, c := range predeclaredConsts {
		r.set(c.variable.name, c.variable)
	}
	return r
}

func parse(t *TokenStream) *AstSourceFile {
	ts = t
	universeblockscope = newUniverseBlockScope()
	packageblockscope = newScope(universeblockscope)
	currentscope = packageblockscope
	return parseSourceFile()
}
