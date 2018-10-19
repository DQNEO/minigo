package main

import (
	"fmt"
)

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

type ExprSliced struct {
	ref *AstOperandName
	low Expr
	high Expr
}

func (e *ExprSliced) dump() {

}
func (e *ExprSliced) emit() {

}

type ExprIndexAccess struct {
	ref *AstOperandName
	index Expr
}

func (e *ExprIndexAccess) dump() {

}

func (e *ExprIndexAccess) emit() {

}

func parseIdentOrFuncall(firstIdent identifier) Expr {
	debugf("func %s start with %s" , "parseIdentOrFuncall", peekToken().String())
	debugNest++
	defer func() {
		debugNest--
		debugf("func %s end", "parseIdentOrFuncall")
	}()

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

	operand := &AstOperandName{
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
	} else if tok.isPunct("[") {
		// index access
		// assure operand is array, slice, or map
		tok := readToken()
		if tok.isPunct(":") {
			lowIndex := &ExprNumberLiteral{
				val:0,
			}
			highIndex := parseExpr()
			expect("]")
			return &ExprSliced{
				ref : operand,
				low: lowIndex,
				high: highIndex,
			}
		} else {
			unreadToken()
			index := parseExpr()
			tok := readToken()
			if tok.isPunct("]") {
				return &ExprIndexAccess{
					ref: operand,
					index: index,
				}
			} else if tok.isPunct(":") {
				highIndex := parseExpr()
				expect("]")
				return &ExprSliced{
					ref : operand,
					low: index,
					high: highIndex,
				}

			} else {
				tok.errorf("invalid token in index access")
			}
		}
	} else {
		unreadToken()
	}


	v := currentscope.get(firstIdent)
	if v == nil{
		errorf("Undefined variable: %s", firstIdent)
		return nil
	}
	vardecl, _ := v.(*AstVarDecl)
	if vardecl != nil {
		return vardecl.variable
	}
	constdecl, _ := v.(*AstConstDecl)
	if constdecl != nil {
		return constdecl.variable
	}
	errorf("variable not found %v",firstIdent )
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
		length: len(values),
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
	case "==","!=", "<",">", ">=", "<=":
		return 10
	case "-","+":
		return 10
	case "*":
		return 20
	default :
		errorf("unkown operator %s", op)
	}
	return 0;
}

func parseExpr() Expr {
	return parseExprInt(-1)
}

var  binops  = []string{
	"+","*","-", "==","!=","<",">","<=","=>",
}

func parseExprInt(prior int) Expr {
	debugf("func %s start with %s" , "parseExprInt", peekToken().String())
	debugNest++
	defer func() {
		debugNest--
		debugf("func %s end", "parseExprInt")
	}()
	ast := parseUnaryExpr()
	debugAstConstructed(ast)
	if ast == nil {
		return nil
	}
	for {
		tok := readToken()
		if tok.isSemicolon() {
			unreadToken()
			return ast
		}

		// if bion
		if in_array(tok.sval, binops) {
			prior2 := priority(tok.sval)
			if prior < prior2 {
				right := parseExprInt(prior2)
				if ast == nil {
					tok.errorf("bad left unary expr:%v", ast)
				}
				ast = &ExprBinop{
					op:    tok.sval,
					left:  ast,
					right: right,
				}
				debugAstConstructed(ast)
				continue
			} else {
				unreadToken()
				return ast
			}
		/*
		} else if tok.sval == "," || tok.sval == ")" ||
			tok.sval == "{" || tok.sval == "}" ||
			tok.isPunct(";") || tok.isPunct(":") { // end of funcall argument
			unreadToken()
			return ast
		*/
		} else {
			unreadToken()
			return ast
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
			//_ := tok.getIdent()
			//_ := currentscope.getGtype(typename)
			//return gtype
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
	//var gtype *Gtype
	var initval Expr
	if tok.isPunct("=") {
		//var x = EXPR
		initval = parseExpr()
		//gtype = gInt  // FIXME: infer type
		expect(";")
	} else {
		unreadToken()
		// expect Type
		_ = parseType()
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

	variable := newVariable(name, nil, isGlobal)

	r := &AstVarDecl{
		variable: variable,
		initval:  initval,
	}
	currentscope.setVarDecl(name, r)
	return r
}

func parseConstDecl(isGlobal bool) *AstConstDecl {
	// read name
	name := readIdent()

	// Type or "="
	tok := readToken()
	var val Expr
	if tok.isPunct("=") {
		// infer mode: const x = EXPR
		val = parseExpr()
		expect(";")
	} else {
		unreadToken()
		// expect Type
		_ = parseType()
		// const x T = EXPR
		expect("=")
		val = parseExpr()
		expect(";")
	}

	variable := &ExprConstVariable{
		name: name,
		val: val,
	}

	r := &AstConstDecl{
		variable: variable,
	}

	currentscope.setConstDecl(name, r)
	return r
}

func parseAssignment() *AstAssignment {
	tleft := readToken()
	item := currentscope.get(tleft.getIdent())
	if item == nil {
		errorf("variable %s is not found", tleft.getIdent())
	}
	vardecl, ok := item.(*AstVarDecl)
	if !ok {
		errorf("%s is not a variable", tleft)
	}
	expect("=")
	rexpr := parseExpr()
	expect(";")
	return &AstAssignment{
		left:  vardecl.variable,
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
	currentscope = newScope(currentscope)
	// Assume "range" style
	idents := parseIdentList()
	expect(":=")
	// TODO register each ient to the scope
	for _, ident := range idents {
		currentscope.setVarDecl(ident, &AstVarDecl{variable:newVariable(ident,nil,false)})
	}
	r.idents = idents
	expectKeyword("range")
	r.list = parseExpr()
	expect("{")
	r.block = parseCompoundStmt()
	currentscope = currentscope.outer
	return r
}

func parseIfStmt() *AstIfStmt {
	debugf("func %s start with %s" , "parseForStmt", peekToken())
	debugNest++
	defer func() {
		debugNest--
		debugf("func %s end", "parseForStmt")
	}()
	var r = &AstIfStmt{}
	currentscope = newScope(currentscope)
	r.cond = parseExpr()
	expect("{")
	r.then = parseCompoundStmt()
	tok := readToken()
	if (tok.isKeyword("else")) {
		tok := readToken()
		if tok.isPunct("{") {
			r.els = &AstStmt{compound:parseCompoundStmt()}
		} else if tok.isKeyword("if") {
			r.els = &AstStmt{ifstmt:parseIfStmt(),}
		} else {
			tok.errorf("Syntax error")
		}
	} else {
		unreadToken()
	}
	currentscope = currentscope.outer
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
	} else if tok.isKeyword("if") {
		return  &AstStmt{ifstmt:parseIfStmt()}
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
	debugf("func %s start with %s" , "parseCompoundStmt", peekToken().String())
	debugNest++
	defer func() {
		debugNest--
		debugf("func %s end", "parseCompoundStmt")
	}()


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
			currentscope.setVarDecl(pname, &AstVarDecl{
				variable:variable,
			})
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
	debugf("scope:%s", currentscope)
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

// read after "struct" token
func parseStructDef() *AstStructDef {
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
	return &AstStructDef{
		fields: fields,
	}
}

func parseInterfaceDef() *AstInterfaceDef {
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
	return &AstInterfaceDef{
		methods: methods,
	}
}

func parseTypeDecl() *AstTypeDecl  {
	name := readIdent()
	tok := readToken()
	var typeConstuctor interface{}
	if tok.isKeyword("struct" ) {
		typeConstuctor = parseStructDef()
	} else if tok.isKeyword("interface")  {
		typeConstuctor = parseInterfaceDef()
	} else if tok.isTypeIdent() {
		ident := tok.getIdent() // name of another type
		typeConstuctor = ident
	} else {
		tok.errorf("TBD")
	}
	r := &AstTypeDecl{
		typedef:&AstTypeDef{
			name:            name,
			typeConstructor: typeConstuctor,
		},
	}
	currentscope.setTypeDecl(name, r)
	return r
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
	r := &AstSourceFile{}
	r.pkg =   shouldHavePackageClause()
	r.imports =  mayHaveImportDecls()
	r.packageNames = make(map[identifier]string)
	for _, importdecl := range r.imports {
		for _, path := range importdecl.paths {
			ident := identifier(path)
			r.packageNames[ident] = path
		}
	}

	r.decls =   parseTopLevelDecls()
	return r
}

func resolveVar(decl *AstVarDecl) {
	if decl.variable.gtype != nil {
		return
	}

	constructor := decl.variable.typeConstructor
	switch constructor.(type) {
	case identifier:
		item := packageblockscope.get(constructor.(identifier))
		if item == nil {
			errorf("Undefined type %v", item)
		}
		typedecl, ok := item.(*AstTypeDecl)
		if !ok {
			errorf("%v is not a type", item)
		}
		if typedecl.gtype == nil {
			errorf("type is not resolved", item)
		}
		decl.variable.gtype = typedecl.gtype
	}
}


func resolveConst(decl *AstConstDecl) {
	if decl.variable.gtype != nil {
		return
	}

	constructor := decl.variable.typeConstructor
	switch constructor.(type) {
	case identifier:
		item := packageblockscope.get(constructor.(identifier))
		if item == nil {
			errorf("Undefined type %v", item)
		}
		typedecl, ok := item.(*AstTypeDecl)
		if !ok {
			errorf("%v is not a type", item)
		}
		if typedecl.gtype == nil {
			errorf("type is not resolved", item)
		}
		decl.variable.gtype = typedecl.gtype
	}
}

func resolveType(decl *AstTypeDecl) {
	if decl.gtype != nil {
		return
	}

	constructor := decl.typedef.typeConstructor
	switch constructor.(type) {
	case identifier:
		item := packageblockscope.get(constructor.(identifier))
		if item == nil {
			errorf("Undefined type %v", item)
		}
		typedecl, ok := item.(*AstTypeDecl)
		if !ok {
			errorf("%v is not a type", item)
		}
		if typedecl.gtype == nil {
			resolveType(typedecl)
		}
		decl.gtype = &Gtype{
			typ:"ref",
			size:typedecl.gtype.size,
			ptr:typedecl.gtype,
		}

	}
}

func resolve(file *AstSourceFile) {
	// resolve types
	for _, decl := range packageblockscope.idents {
		typedecl, ok := decl.(*AstTypeDecl)
		if ok {
			resolveType(typedecl)
			typedecl.dump()
		}
	}

	for _, decl := range packageblockscope.idents {
		constdecl, ok := decl.(*AstConstDecl)
		if ok {
			resolveConst(constdecl)
			constdecl.dump()
		}
	}

	for _, decl := range packageblockscope.idents {
		vardecl, ok := decl.(*AstVarDecl)
		if ok {
			debugf("resolve decl :")
			resolveVar(vardecl)
			vardecl.dump()
		}
	}
}

func parse(t *TokenStream) *AstSourceFile {
	ts = t
	universeblockscope = newUniverseBlockScope()
	packageblockscope = newScope(universeblockscope)
	currentscope = packageblockscope
	return parseSourceFile()
}
