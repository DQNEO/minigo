package main

import (
	"fmt"
	"runtime"

	"os"
)

type parser struct {
	tokenStream         *TokenStream
	unresolvedRelations []*Relation
	packageBlockScope   *scope
	currentScope        *scope
	globalvars          []*ExprVariable
	localvars           []*ExprVariable
	importedNames       map[identifier]bool
}

type TokenStream struct {
	tokens []*Token
	index  int
}

func (p *parser) peekToken() *Token {
	ts := p.tokenStream
	if ts.index > len(ts.tokens)-1 {
		return makeToken("EOF", "")
	}
	r := ts.tokens[ts.index]
	return r
}

func (p *parser) lastToken() *Token {
	ts := p.tokenStream
	return ts.tokens[ts.index-1]
}

func (p *parser) readToken() *Token {
	ts := p.tokenStream
	if ts.index > len(ts.tokens)-1 {
		return makeToken("EOF", "")
	}
	r := ts.tokens[ts.index]
	ts.index++
	return r
}

func (p *parser) unreadToken() {
	p.tokenStream.index--
}

func (p *parser) readIdent() identifier {
	tok := p.readToken()
	if !tok.isTypeIdent() {
		errorf("Identifier expected, but got %s", tok)
	}
	return tok.getIdent()
}

func (p *parser) expectKeyword(name string) {
	tok := p.readToken()
	if !tok.isKeyword(name) {
		errorf("Keyword %s expected but got %s", name)
	}
}

func (p *parser) expect(punct string) {
	tok := p.readToken()
	if !tok.isPunct(punct) {
		errorf("punct '%s' expected but got '%s'", punct, tok)
	}
}

func getCallerName(n int) string {
	pc, _, _, ok := runtime.Caller(n)
	if !ok {
		errorf("Unable to get caller")
	}
	details := runtime.FuncForPC(pc)
	//r := (strings.Split(details.Name(), "."))[2]
	return details.Name()
}

func (p *parser) traceIn() int {
	if !debugParser {
		return 0
	}
	debugf("func %s is gonna read %s", getCallerName(2), p.peekToken())
	debugNest++
	return 0
}

func (p *parser) traceOut(_ int) {
	if !debugParser {
		return
	}
	if r := recover(); r != nil {
		fmt.Println(r)
		os.Exit(1)
	}
	debugNest--
	debugf("func %s end", getCallerName(2))
}

func (p *parser) readFuncallArgs() []Expr {
	defer p.traceOut(p.traceIn())
	var r []Expr
	for {
		tok := p.readToken()
		if tok.isPunct(")") {
			return r
		}
		p.unreadToken()
		arg := p.parseExpr()
		r = append(r, arg)
		tok = p.readToken()
		if tok.isPunct(")") {
			return r
		} else if tok.isPunct(",") {
			continue
		} else {
			errorf("invalid token in funcall arguments: %s", tok)
		}
	}
}

//var outerPackages map[identifier](map[identifier]interface{})

type Relation struct {
	name identifier

	// either of expr or gtype
	expr  Expr
	gtype *Gtype
}

var labelSeq = 0
var stringLiterals []*ExprStringLiteral

func (p *parser) newAstString(sval string) *ExprStringLiteral {
	ast := &ExprStringLiteral{
		val:    sval,
		slabel: fmt.Sprintf("L%d", labelSeq),
	}
	labelSeq++
	stringLiterals = append(stringLiterals, ast)
	return ast
}

func (p *parser) parsePrim() Expr {
	defer p.traceOut(p.traceIn())
	tok := p.readToken()
	switch {
	case tok.isSemicolon():
		return nil
	case tok.isTypeString(): // string literal
		return p.newAstString(tok.sval)
	case tok.isTypeInt(): // int literal
		ival := tok.getIntval()
		return &ExprNumberLiteral{
			val: ival,
		}
	case tok.isTypeChar(): // char literal
		sval := tok.sval
		c := sval[0]
		return &ExprNumberLiteral{
			val: int(c),
		}
	case tok.isPunct("["): // array literal
		return p.parseArrayLiteral()
	case tok.isTypeIdent():
		firstIdent := tok.getIdent()
		// https://golang.org/ref/spec#QualifiedIdent
		// read QualifiedIdent
		tok := p.readToken()
		var pkg identifier
		var ident identifier
		if tok.isPunct(":=") {
			p.unreadToken()
			return &Relation{
				name: firstIdent,
			}
		} else if tok.isPunct(".") {
			// Assume firstIdent is a package name
			pkg = firstIdent
			_, ok := p.importedNames[pkg]
			if ok {
				ident = p.readIdent()
				debugf("Reference to outer entity %s.%s", pkg, ident)
			} else {
				//return nil
			}
		} else {
			p.unreadToken()
			pkg = ""
			ident = firstIdent
		}

		operand := &AstOperandName{
			pkg:   pkg,
			ident: ident,
		}

		tok = p.readToken()
		if tok.isPunct("(") {
			// try funcall
			args := p.readFuncallArgs()

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
			tok := p.readToken()
			if tok.isPunct(":") {
				lowIndex := &ExprNumberLiteral{
					val: 0,
				}
				highIndex := p.parseExpr()
				p.expect("]")
				return &ExprSliced{
					ref:  operand,
					low:  lowIndex,
					high: highIndex,
				}
			} else {
				p.unreadToken()
				index := p.parseExpr()
				tok := p.readToken()
				if tok.isPunct("]") {
					rel := &Relation{
						name: firstIdent,
					}
					p.tryResolve(rel)
					return &ExprArrayIndex{
						rel:   rel,
						index: index,
					}
				} else if tok.isPunct(":") {
					highIndex := p.parseExpr()
					p.expect("]")
					return &ExprSliced{
						ref:  operand,
						low:  index,
						high: highIndex,
					}

				} else {
					tok.errorf("invalid token in index access")
				}
			}
		} else {
			p.unreadToken()
		}

		rel := &Relation{
			name: firstIdent,
		}
		p.tryResolve(rel)
		return rel
	default:
		errorf("unable to handle %s", tok)
	}
	errorf("unable to handle %s", tok)
	return nil
}

func (p *parser) parseArrayLiteral() Expr {
	assert(p.lastToken().isPunct("["), "[ is read")
	defer p.traceOut(p.traceIn())
	tlen := p.readToken()
	p.expect("]")
	typ := p.parseType()
	p.expect("{")
	var values []Expr
	for {
		tok := p.readToken()
		if tok.isPunct("}") {
			break
		}
		var v Expr
		if tok.isTypeString() {
			v = &ExprStringLiteral{val: tok.sval}
		} else if tok.isTypeInt() {
			v = &ExprNumberLiteral{val: tok.getIntval()}
		} else if tok.isTypeChar() {
			v = &ExprNumberLiteral{
				val: int(tok.sval[0]),
			}
		} else {
			tok.errorf("TBD")
		}
		assert(v != nil, "v is not nil")
		values = append(values, v)
		tok = p.readToken()
		if tok.isPunct(",") {
			continue
		} else if tok.isPunct("}") {
			break
		} else {
			errorf("unpexpected %s", tok)
		}
	}
	if len(values) != tlen.getIntval() {
		debugPrintV(values)
		errorf("array length does not match (%d != %d)",
			len(values), tlen.getIntval())
	}

	gtype := &Gtype{
		typ:    G_ARRAY,
		length: len(values),
		ptr:    typ,
	}

	r := &ExprArrayLiteral{
		gtype:  gtype,
		values: values,
	}

	return r
}

func (p *parser) parseUnaryExpr() Expr {
	defer p.traceOut(p.traceIn())
	return p.parsePrim()
}

func priority(op string) int {
	switch op {
	case "&&", "||":
		return 5
	case "==", "!=", "<", ">", ">=", "<=":
		return 10
	case "-", "+":
		return 10
	case "/", "%":
		return 15
	case "*":
		return 20
	default:
		errorf("unkown operator %s", op)
	}
	return 0
}

func (p *parser) parseExpr() Expr {
	return p.parseExprInt(-1)
}

var binops = []string{
	"+", "*", "-", "==", "!=", "<", ">", "<=", ">=", "&&", "||", "/", "%",
}

func (p *parser) parseExprInt(prior int) Expr {
	defer p.traceOut(p.traceIn())

	ast := p.parseUnaryExpr()

	if ast == nil {
		return nil
	}
	for {
		tok := p.readToken()
		if tok.isSemicolon() {
			p.unreadToken()
			return ast
		}

		// if bion
		if in_array(tok.sval, binops) {
			prior2 := priority(tok.sval)
			if prior < prior2 {
				right := p.parseExprInt(prior2)
				if ast == nil {
					tok.errorf("bad lefts unary expr:%v", ast)
				}
				ast = &ExprBinop{
					op:    tok.sval,
					left:  ast,
					right: right,
				}

				continue
			} else {
				p.unreadToken()
				return ast
			}
			/*
				} else if tok.sval == "," || tok.sval == ")" ||
					tok.sval == "{" || tok.sval == "}" ||
					tok.isPunct(";") || tok.isPunct(":") { // end of funcall argument
					p.unreadToken()
					return ast
			*/
		} else {
			p.unreadToken()
			return ast
			tok.errorf("Unexpected")
		}
	}

	return ast
}

func (p *parser) newVariable(varname identifier, gtype *Gtype, isGlobal bool) *ExprVariable {
	var variable *ExprVariable
	if isGlobal {
		variable = &ExprVariable{
			varname:  varname,
			gtype:    gtype,
			isGlobal: true,
		}
		p.globalvars = append(p.globalvars, variable)
	} else {
		variable = &ExprVariable{
			varname: varname,
			gtype:   gtype,
		}
		p.localvars = append(p.localvars, variable)
	}
	return variable
}

// https://golang.org/ref/spec#Type
func (p *parser) parseType() *Gtype {
	defer p.traceOut(p.traceIn())
	var gtype *Gtype

	for {
		tok := p.readToken()
		if tok.isTypeIdent() {
			ident := tok.getIdent()
			// unresolved
			rel := &Relation{
				name: ident,
			}
			p.tryResolve(rel)
			gtype = &Gtype{
				typ:      G_REL,
				relname:  ident,
				relation: rel,
			}
			return gtype
		} else if tok.isPunct("*") {
			// pointer
		} else if tok.isKeyword("struct") {
			_ = p.parseStructDef()
		} else if tok.isKeyword("interface") {
			_ = p.parseInterfaceDef()
		} else if tok.isPunct("[") {
			// array or slice
			tok := p.readToken()
			if tok.isPunct("]") {
				// slice
				typ := p.parseType()
				return &Gtype{
					typ:      G_SLICE,
					ptr:      typ, // element type
					length:   0,
					capacity: 0,
				}
			} else {
				// array
				p.expect("]")
				typ := p.parseType()
				return &Gtype{
					typ:    G_ARRAY,
					length: tok.getIntval(),
					ptr:    typ,
				}
			}
		} else if tok.isPunct("]") {

		} else {
			p.unreadToken()
			break
		}

	}
	errorf("Unkown type")
	return nil
}

func (p *parser) parseVarDecl(isGlobal bool) *AstVarDecl {
	assert(p.lastToken().isKeyword("var"), "last token is \"var\"")
	defer p.traceOut(p.traceIn())
	// read newName
	newName := p.readIdent()
	var typ *Gtype
	var initval Expr
	// "=" or Type
	tok := p.readToken()
	if tok.isPunct("=") {
		// Infer mode
		initval = p.parseExpr()
		if typ == nil {
			typ = gInt // FIXME: infer type
		}
	} else {
		p.unreadToken()
		typ = p.parseType()
		assert(typ != nil, "has typ")
		tok := p.readToken()
		if tok.isPunct("=") {
			initval = p.parseExpr()
		}
	}
	//p.expect(";")

	variable := p.newVariable(newName, typ, isGlobal)
	r := &AstVarDecl{
		variable: variable,
		initval:  initval,
	}
	p.currentScope.setVar(newName, variable)
	return r
}

func (p *parser) parseConstDecl() *AstConstDecl {
	defer p.traceOut(p.traceIn())
	// read newName
	newName := p.readIdent()

	// Type or "="
	tok := p.readToken()
	var val Expr
	if tok.isPunct("=") {
		// infer mode: const x = EXPR
		val = p.parseExpr()
		p.expect(";")
	} else {
		p.unreadToken()
		// expect Type
		_ = p.parseType()
		// const x T = EXPR
		p.expect("=")
		val = p.parseExpr()
		p.expect(";")
	}

	variable := &ExprConstVariable{
		name: newName,
		val:  val,
	}

	r := &AstConstDecl{
		variable: variable,
	}

	p.currentScope.setConst(newName, variable)
	return r
}

func (p *parser) parseIdentList() []identifier {
	defer p.traceOut(p.traceIn())
	var r []identifier
	for {
		tok := p.readToken()
		if tok.isTypeIdent() {
			r = append(r, tok.getIdent())
		} else if len(r) == 0 {
			// at least one ident is needed
			tok.errorf("Ident expected")
		}

		tok = p.readToken()
		if tok.isPunct(",") {
			continue
		} else {
			p.unreadToken()
			return r
		}
	}
	return r
}

func (p *parser) enterNewScope() {
	p.currentScope = newScope(p.currentScope)
}

func (p *parser) exitScope() {
	p.currentScope = p.currentScope.outer
}

// https://golang.org/ref/spec#For_statements
func (p *parser) parseForStmt() *AstForStmt {
	defer p.traceOut(p.traceIn())
	var r = &AstForStmt{}
	p.enterNewScope()
	defer p.exitScope()

	expr := p.parseExpr()
	if p.peekToken().isPunct("{") {
		// single cond
		r.cls = &ForForClause{
			cond: expr,
		}
	} else {
		// for clause or range clause
		var initstmt Stmt
		lefts := p.parseExpressionList(expr)
		tok2 := p.readToken()
		if tok2.isPunct("=") {
			if p.peekToken().isKeyword("range") {
				return p.parseForRange(lefts)
			} else {
				initstmt = p.parseAssignment(lefts)
			}
		} else if tok2.isPunct(":=") {
			if p.peekToken().isKeyword("range") {
				for _, e := range lefts {
					rel := e.(*Relation) // a brand new rel
					gtype := gInt        // FIXME infer type
					variable := p.newVariable(rel.name, gtype, false)
					rel.expr = variable
					p.currentScope.setVar(rel.name, variable)
				}
				return p.parseForRange(lefts)
			} else {
				initstmt = p.parseShortAssignment(lefts)
			}
		} else {
			p.unreadToken()
		}

		cls := &ForForClause{}
		// regular for cond
		cls.init = initstmt
		p.expect(";")
		cls.cond = p.parseStmt()
		p.expect(";")
		cls.post = p.parseStmt()
		r.cls = cls
	}

	p.expect("{")
	r.block = p.parseCompoundStmt()
	return r
}

func (p *parser) parseForRange(exprs []Expr) *AstForStmt {
	p.expectKeyword("range")

	if len(exprs) > 2 {
		errorf("range values should be 1 or 2")
	}
	indexvar, ok := exprs[0].(*Relation)
	if !ok {
		errorf(" rng.lefts[0]. is not relation")
	}
	var valuevar *Relation
	if len(exprs) == 2 {
		valuevar = exprs[1].(*Relation)
	}

	var r = &AstForStmt{
		rng: &ForRangeClause{
			indexvar: indexvar,
			valuevar: valuevar,
		},
	}

	r.rng.rangeexpr = p.parseExpr()
	p.expect("{")
	r.block = p.parseCompoundStmt()
	return r
}

func (p *parser) parseIfStmt() *AstIfStmt {
	defer p.traceOut(p.traceIn())
	var r = &AstIfStmt{}
	p.enterNewScope()
	defer p.exitScope()
	r.cond = p.parseExpr()
	p.expect("{")
	r.then = p.parseCompoundStmt()
	tok := p.readToken()
	if tok.isKeyword("else") {
		tok2 := p.readToken()
		if tok2.isKeyword("if") {
			// we regard "else if" as a kind of a nested if statement
			// else if => else { if .. { } else {} }
			r.els = p.parseIfStmt()
		} else if tok2.isPunct("{") {
			r.els = p.parseCompoundStmt()
		} else {
			tok2.errorf("Unexpected token")
		}
	} else {
		p.unreadToken()
	}
	return r
}

func (p *parser) parseReturnStmt() *AstReturnStmt {
	defer p.traceOut(p.traceIn())
	var r *AstReturnStmt
	r = &AstReturnStmt{expr: p.parseExpr()}
	return r
}

func (p *parser) parseExpressionList(first Expr) []Expr {
	defer p.traceOut(p.traceIn())
	var r []Expr
	if first == nil {
		first = p.parseExpr()
		// should skip "," if exists
	}
	r = append(r, first)
	for {
		tok := p.readToken()
		if tok.isSemicolon() {
			p.unreadToken()
			return r
		} else if tok.isPunct("=") || tok.isPunct(":=") {
			p.unreadToken()
			return r
		} else if tok.isPunct(",") {
			expr := p.parseExpr()
			r = append(r, expr)
			continue
		} else {
			p.unreadToken()
			return r
		}
	}
	return r
}

func (p *parser) parseAssignment(lefts []Expr) *AstAssignment {
	rights := p.parseExpressionList(nil)
	return &AstAssignment{
		lefts:  lefts,
		rights: rights,
	}
}

func (p *parser) parseShortAssignment(lefts []Expr) *AstAssignment {
	rights := p.parseExpressionList(nil)
	for _, e := range lefts {
		rel := e.(*Relation) // a brand new rel
		gtype := gInt        // FIXME infer type
		variable := p.newVariable(rel.name, gtype, false)
		rel.expr = variable
		p.currentScope.setVar(rel.name, variable)
	}
	return &AstAssignment{
		lefts:  lefts,
		rights: rights,
	}
}

// this is in function scope
func (p *parser) parseStmt() Stmt {
	defer p.traceOut(p.traceIn())
	tok := p.readToken()
	if tok.isKeyword("var") {
		return p.parseVarDecl(false)
	} else if tok.isKeyword("const") {
		return p.parseConstDecl()
	} else if tok.isKeyword("type") {
		return p.parseTypeDecl()
	} else if tok.isKeyword("for") {
		return p.parseForStmt()
	} else if tok.isKeyword("if") {
		return p.parseIfStmt()
	} else if tok.isKeyword("return") {
		return p.parseReturnStmt()
	}
	p.unreadToken()
	expr1 := p.parseExpr()
	tok2 := p.readToken()
	if tok2.isPunct(",") {
		p.unreadToken()
		// Multi value assignment
		lefts := p.parseExpressionList(expr1)
		tok3 := p.readToken()
		if tok3.isPunct("=") {
			return p.parseAssignment(lefts)
		} else if tok3.isPunct(":=") {
			return p.parseShortAssignment(lefts)
		} else {
			tok3.errorf("TBD")
		}
	} else if tok2.isPunct("=") {
		return p.parseAssignment([]Expr{expr1})
	} else if tok2.isPunct(":=") {
		// Single value ShortVarDecl
		return p.parseShortAssignment([]Expr{expr1})
	} else {
		p.unreadToken()
		return expr1
	}
	return nil
}

func (p *parser) parseCompoundStmt() *AstCompountStmt {
	defer p.traceOut(p.traceIn())

	r := &AstCompountStmt{}
	for {
		tok := p.readToken()
		if tok.isPunct("}") {
			return r
		}
		if tok.isSemicolon() {
			continue
		}
		p.unreadToken()
		stmt := p.parseStmt()
		r.stmts = append(r.stmts, stmt)
	}
	return nil
}

func (p *parser) parseFuncDef() *AstFuncDecl {
	defer p.traceOut(p.traceIn())
	p.localvars = make([]*ExprVariable, 0)
	p.enterNewScope()
	defer p.exitScope()
	fname := p.readToken().getIdent()
	p.expect("(")
	var params []*ExprVariable

	tok := p.readToken()
	if !tok.isPunct(")") {
		p.unreadToken()
		for {
			tok := p.readToken()
			pname := tok.getIdent()
			ptype := p.parseType()
			// assureType(tok.sval)
			variable := &ExprVariable{
				varname: pname,
				gtype:   ptype,
			}
			params = append(params, variable)
			p.currentScope.setVar(pname, variable)
			tok = p.readToken()
			if tok.isPunct(")") {
				break
			}
			if !tok.isPunct(",") {
				errorf("Invalid token %s", tok)
			}
		}
	}

	// read func rettype
	tok = p.readToken()
	var rettype string
	if tok.isTypeIdent() {
		// rettype
		rettype = tok.sval
		p.expect("{")
	} else {
		assert(tok.isPunct("{"), "begin of func body")
	}
	debugf("scope:%s", p.currentScope)
	body := p.parseCompoundStmt()
	r := &AstFuncDecl{
		fname:     fname,
		rettype:   rettype,
		params:    params,
		localvars: p.localvars,
		body:      body,
	}
	p.localvars = nil
	return r
}

func (p *parser) parseImport() *AstImportDecl {
	defer p.traceOut(p.traceIn())
	tok := p.readToken()
	var specs []*AstImportSpec
	if tok.isPunct("(") {
		for {
			tok := p.readToken()
			if tok.isTypeString() {
				name := identifier(tok.sval)
				specs = append(specs, &AstImportSpec{
					packageName: name,
					path:        tok.sval,
				})
				p.expect(";")
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
		name := identifier(tok.sval)
		specs = []*AstImportSpec{{
			packageName: name,
			path:        tok.sval,
		},
		}
	}
	p.expect(";")
	return &AstImportDecl{
		specs: specs,
	}
}

func (p *parser) expectPackageClause() *AstPackageClause {
	defer p.traceOut(p.traceIn())
	p.expectKeyword("package")
	r := &AstPackageClause{name: p.readIdent()}
	p.expect(";")
	return r
}

func (p *parser) parseImportDecls() []*AstImportDecl {
	defer p.traceOut(p.traceIn())
	var r []*AstImportDecl
	for {
		tok := p.readToken()
		if !tok.isKeyword("import") {
			p.unreadToken()
			return r
		}
		r = append(r, p.parseImport())
	}
}

// read after "struct" token
func (p *parser) parseStructDef() *AstStructDef {
	assert(p.lastToken().isKeyword("struct"), `require "struct" is already read`)
	defer p.traceOut(p.traceIn())
	p.expect("{")
	var fields []*StructField
	for {
		tok := p.readToken()
		if tok.isPunct("}") {
			break
		}
		fieldname := tok.getIdent()
		fieldtyep := p.parseType()
		fields = append(fields, &StructField{
			name:  fieldname,
			gtype: fieldtyep,
		})
		p.expect(";")
	}
	p.expect(";")
	return &AstStructDef{
		fields: fields,
	}
}

func (p *parser) parseInterfaceDef() *AstInterfaceDef {
	defer p.traceOut(p.traceIn())
	p.expect("{")
	var methods []identifier
	for {
		tok := p.readToken()
		if tok.isPunct("}") {
			break
		}
		fname := tok.getIdent()
		p.expect("(")
		p.expect(")")
		p.expect(";")
		methods = append(methods, fname)
	}
	p.expect(";")
	return &AstInterfaceDef{
		methods: methods,
	}
}

func (p *parser) tryResolve(rel *Relation) {
	relfound := p.currentScope.get(rel.name)
	if relfound != nil {
		switch relfound.(type) {
		case *Gtype:
			rel.gtype = relfound.(*Gtype)
		case Expr:
			rel.expr = relfound.(Expr)
		default:
			errorf("Bad type relfound %v", relfound)
		}
	} else {
		p.unresolvedRelations = append(p.unresolvedRelations, rel)
	}
}

func (p *parser) parseTypeDecl() *AstTypeDecl {
	defer p.traceOut(p.traceIn())
	newName := p.readIdent()
	gtype := p.parseType()
	r := &AstTypeDecl{
		name:  newName,
		gtype: gtype,
	}
	p.currentScope.setGtype(newName, gtype)
	return r
}

func (p *parser) parseTopLevelDecl(tok *Token) *AstTopLevelDecl {
	defer p.traceOut(p.traceIn())
	var r *AstTopLevelDecl
	switch {
	case tok.isKeyword("var"):
		vardecl := p.parseVarDecl(true)
		r = &AstTopLevelDecl{vardecl: vardecl}
	case tok.isKeyword("const"):
		constdecl := p.parseConstDecl()
		r = &AstTopLevelDecl{constdecl: constdecl}
	case tok.isKeyword("func"):
		funcdecl := p.parseFuncDef()
		r = &AstTopLevelDecl{funcdecl: funcdecl}
	case tok.isKeyword("type"):
		typedecl := p.parseTypeDecl()
		r = &AstTopLevelDecl{typedecl: typedecl}
	default:
		errorf("TBD: unable to handle token %v", tok)
	}

	//debugAstConstructed(r)
	return r
}

func (p *parser) parseTopLevelDecls() []*AstTopLevelDecl {
	defer p.traceOut(p.traceIn())
	var r []*AstTopLevelDecl
	for {
		tok := p.readToken()
		if tok.isEOF() {
			return r
		}

		if tok.isPunct(";") {
			continue
		}
		ast := p.parseTopLevelDecl(tok)
		r = append(r, ast)
	}
	return r
}

// https://golang.org/ref/spec#Source_file_organization
// Each source file consists of
// a package clause defining the package to which it belongs,
// followed by a possibly empty set of import declarations that declare packages whose contents it wishes to use,
// followed by a possibly empty set of declarations of functions, types, variables, and constants.
func (p *parser) parseSourceFile(sourceFile string, packageBlockScope *scope) *AstSourceFile {

	bs := NewByteStream(sourceFile)
	tokens := tokenize(bs)
	assert(len(tokens) > 0, "tokens should have length")

	/*
		if debugToken {
			renderTokens(tokens)
		}
	*/

	p.tokenStream = &TokenStream{
		tokens: tokens,
		index:  0,
	}

	p.packageBlockScope = packageBlockScope
	p.currentScope = packageBlockScope

	r := &AstSourceFile{}
	r.pkg = p.expectPackageClause()
	r.imports = p.parseImportDecls()

	p.importedNames = make(map[identifier]bool)
	for _, importdecl := range r.imports {
		for _, spec := range importdecl.specs {
			p.importedNames[spec.packageName] = true
		}
	}
	debugPrintV(p.importedNames)
	r.decls = p.parseTopLevelDecls()
	return r
}

func (p *parser) resolve() {
	// set the universe in the background
	universeblockscope := newUniverseBlockScope()
	p.packageBlockScope.outer = universeblockscope
	for _, rel := range p.unresolvedRelations {
		p.tryResolve(rel)
	}

}
