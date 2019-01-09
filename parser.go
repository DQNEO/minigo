package main

import (
	"fmt"
	"runtime"
	"strings"

	"os"
)

type parser struct {
	tokenStream         *TokenStream
	unresolvedRelations []*Relation
	packageBlockScope   *scope
	currentScope        *scope
	scopes              map[identifier]*scope
	stringLiterals []*ExprStringLiteral
	globalvars          []*ExprVariable
	localvars           []*ExprVariable
	namedTypes          map[identifier]methods
	shortassignments    []*StmtShortVarDecl
	importedNames       map[identifier]bool
	requireBlock        bool // workaround for parsing "{" as a block starter
	inCase              int // > 0  while in reading case compound stmts
	constSpecIndex      int
	currentPackageName  identifier
}

type methods map[identifier]*ExprFuncRef

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

// skip one token
func (p *parser) skip() {
	p.readToken()
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
		errorf("Keyword %s expected but got %s", name, tok)
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
	debugf("func %s end after %s", getCallerName(2), p.lastToken())
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
		if p.peekToken().isPunct("...") {
			p.expect("...")
			arg = &ExprVaArg{expr:arg}
			r = append(r, arg)
			p.expect(")")
			return r
		}
		r = append(r, arg)
		tok = p.readToken()
		if tok.isPunct(")") {
			return r
		} else if tok.isPunct(",") {
			continue
		} else {
			tok.errorf("invalid token in funcall arguments")
		}
	}
}

//var outerPackages map[identifier](map[identifier]interface{})

func (p *parser) addStringLiteral(ast *ExprStringLiteral) {
	p.stringLiterals = append(p.stringLiterals, ast)
}

// expr which begins with an ident.
// e.g. ident, ident() or ident.*, etc
func (p *parser) parseIdentExpr(firstIdent identifier) Expr {
	defer p.traceOut(p.traceIn())

	// https://golang.org/ref/spec#QualifiedIdent
	// read QualifiedIdent
	var pkg identifier // ignored for now
	if _, ok := p.importedNames[firstIdent]; ok {
		pkg = firstIdent
		p.expect(".")
		// shift firstident
		firstIdent = p.readIdent()
	}

	rel := &Relation{
		name: firstIdent,
	}
	p.tryResolve(pkg, rel)

	next := p.peekToken()

	var e Expr
	if next.isPunct("{") {
		if p.requireBlock {
			return rel
		}
		// struct literal
		p.skip()
		e = p.parseStructLiteral(rel)
	} else if next.isPunct("(") {
		// funcall or method call
		p.skip()
		args := p.readFuncallArgs()
		fname := string(rel.name)
		e = &ExprFuncall{
			rel:rel,
			fname: fname,
			args:  args,
		}
	} else if next.isPunct("[") {
		// index access
		e = p.parseArrayIndex(rel)
	} else {
		// solo ident
		e = rel
	}

	return p.succeedingExpr(e)
}

func (p *parser) parseArrayIndex(e Expr) Expr {
	defer p.traceOut(p.traceIn())
	p.expect("[")

	var r Expr
	// assure operand is array, slice, or map
	tok := p.readToken()
	if tok.isPunct(":") {
		lowIndex := &ExprNumberLiteral{
			val: 0,
		}
		highIndex := p.parseExpr()
		p.expect("]")
		r = &ExprSliced{
			ref:  nil, // TBI
			low:  lowIndex,
			high: highIndex,
		}
	} else {
		p.unreadToken()
		index := p.parseExpr()
		tok := p.readToken()
		if tok.isPunct("]") {
			r = &ExprArrayIndex{
				array:   e,
				index: index,
			}
		} else if tok.isPunct(":") {
			highIndex := p.parseExpr()
			p.expect("]")
			r = &ExprSliced{
				ref:  nil, // TBI
				low:  index,
				high: highIndex,
			}

		} else {
			tok.errorf("invalid token in index access")
		}
	}
	return r
}

// https://golang.org/ref/spec#Type_assertions
func (p *parser) parseTypeAssertion(e Expr) Expr {
	defer p.traceOut(p.traceIn())
	p.expect("(")
	if p.peekToken().isKeyword("type") {
		p.skip()
		p.expect(")")
		return &ExprTypeSwitchGuard{
			expr:e,
		}
	} else {
		gtype := p.parseType()
		p.expect(")")
		e = &ExprTypeAssertion{
			expr: e,
			gtype: gtype,
		}
		return p.succeedingExpr(e)
	}
	errorf("internal error")
	return nil
}

func (p *parser) succeedingExpr(e Expr) Expr {
	defer p.traceOut(p.traceIn())

	var r Expr
	next := p.peekToken()
	if next.isPunct(".") {
		p.skip()
		if p.peekToken().isPunct("(") {
			// type assertion
			return p.parseTypeAssertion(e)
		}

		// https://golang.org/ref/spec#Selectors
		ident := p.readIdent()
		next = p.peekToken()
		if next.isPunct("(") {
			// (expr).method()
			p.expect("(")
			args := p.readFuncallArgs()
			r = &ExprMethodcall{
				receiver: e,
				fname:    ident,
				args:     args,
			}
			return p.succeedingExpr(r)
		} else {
			// (expr).field
			//   strct.field
			//   (strct.field).field
			//   fncall().field
			//   obj.method().field
			//   array[0].field
			r =  &ExprStructField{
				strct:     e,
				fieldname: ident,
			}
			return p.succeedingExpr(r)
		}
	} else if next.isPunct("["){
		// https://golang.org/ref/spec#Index_expressions
		// (expr)[i]
		e = p.parseArrayIndex(e)
		return p.succeedingExpr(e)
	} else {
		// https://golang.org/ref/spec#OperandName
		r = e
	}

	return r

}

func (p *parser) parseMakeExpr() Expr {
	defer p.traceOut(p.traceIn())
	tok := p.readToken()
	assert(tok.isIdent("make"), "read make")
	p.expect("(")
	mapType := p.parseMapType()
	_ = mapType
	p.expect(")")
	return &ExprNilLiteral{}
}

func (p *parser) parseMapType() *Gtype {
	defer p.traceOut(p.traceIn())
	p.expectKeyword("map")
	p.expect("[")
	mapKey := p.parseType()
	p.expect("]")
	mapValue := p.parseType()
	return &Gtype{
		typ: G_MAP,
		mapKey:mapKey,
		mapValue:mapValue,
	}
}

// https://golang.org/ref/spec#Conversions
func (p *parser) parseTypeConversion(gtype *Gtype) Expr {
	defer p.traceOut(p.traceIn())
	p.expect("(")
	e := p.parseExpr()
	p.expect(")")
	return &ExprConversion{
		expr : e,
		gtype: gtype,
	}
}

func (p *parser) parsePrim() Expr {
	defer p.traceOut(p.traceIn())
	tok := p.readToken()

	switch {
	case tok.isSemicolon():
		return nil
	case tok.isTypeString(): // string literal
		ast := &ExprStringLiteral{
			val:    tok.sval,
		}
		p.addStringLiteral(ast)
		return ast
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
	case tok.isPunct("["): // array literal or type casting
		p.unreadToken()
		gtype := p.parseType()
		if p.peekToken().isPunct("(") {
			// Conversion
			return p.parseTypeConversion(gtype)
		}
		return p.parseArrayLiteral(gtype)
	case tok.isIdent("make"):
		p.unreadToken()
		return p.parseMakeExpr()
	case tok.isTypeIdent():
		return p.parseIdentExpr(tok.getIdent())
	}

	tok.errorf("unable to handle")
	return nil
}

func (p *parser) parseArrayLiteral(gtype *Gtype) Expr {
	defer p.traceOut(p.traceIn())
	gtype.typ = G_ARRAY // convert []T from slice to 0 length array
	p.expect("{")
	var values []Expr
	for {
		tok := p.peekToken()
		if tok.isPunct("}") {
			p.skip()
			break
		}

		v := p.parseExpr()
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

	if gtype.length == 0 {
		gtype.length = len(values)
	} else {
		if len(values) != gtype.length {
			errorf("array length does not match (%d != %d)",
				len(values), gtype.length)
		}
	}

	r := &ExprArrayLiteral{
		gtype:  gtype,
		values: values,
	}

	return r
}

func (p *parser) parseStructLiteral(rel *Relation) *ExprStructLiteral {
	assert(p.lastToken().isPunct("{"), "{ is read")
	defer p.traceOut(p.traceIn())
	r := &ExprStructLiteral{
		strctname: rel,
	}

	for {
		tok := p.readToken()
		if tok.isPunct("}") {
			break
		}
		p.expect(":")
		assert(tok.isTypeIdent(), "field name is ident")
		value := p.parseExpr()
		f := &AstStructFieldLiteral{
			key:   tok.getIdent(),
			value: value,
		}
		r.fields = append(r.fields, f)
		if p.peekToken().isPunct("}") {
			p.expect("}")
			break
		}
		p.expect(",")
	}

	return r
}

func (p *parser) parseUnaryExpr() Expr {
	defer p.traceOut(p.traceIn())
	tok := p.readToken()
	switch {
	case tok.isPunct("("):
		e := p.parseExpr()
		p.expect(")")
		return e
	case tok.isPunct("&"):
		uop := &ExprUop{
			op:      tok.sval,
			operand: p.parsePrim(),
		}
		// when &T{}, allocate stack memory
		if strctliteral, ok := uop.operand.(*ExprStructLiteral); ok {
			// newVariable
			strctliteral.invisiblevar = p.newVariable("", &Gtype{
				typ:G_REL,
				relation:strctliteral.strctname,
			}, false)
		}
		return uop
	case tok.isPunct("*"):
		return &ExprUop{
			op:      tok.sval,
			operand: p.parsePrim(),
		}
	case tok.isPunct("!"):
		return &ExprUop{
			op: tok.sval,
			operand:p.parsePrim(),
		}
	case tok.isPunct("-"):
		return &ExprUop{
			op: tok.sval,
			operand:p.parsePrim(),
		}
	default:
		p.unreadToken()
	}
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
	defer p.traceOut(p.traceIn())
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
			p.tryResolve("", rel)
			gtype = &Gtype{
				typ:      G_REL,
				relation: rel,
			}
			return gtype
		} else if tok.isKeyword("interface") {
			p.expect("{")
			p.expect("}")
			return gInterface
		} else if tok.isPunct("*") {
			// pointer
			gtype = &Gtype{
				typ: G_POINTER,
				ptr: p.parseType(),
			}
			return gtype
		} else if tok.isKeyword("struct") {
			return p.parseStructDef()
		} else if tok.isKeyword("map") {
			p.unreadToken()
			return p.parseMapType()
		} else if tok.isPunct("[") {
			// array or slice
			tok := p.readToken()
			if tok.isPunct("]") {
				// slice
				typ := p.parseType()
				return &Gtype{
					typ:      G_SLICE,
					length:   0,
					ptr:      typ, // element type
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

		} else if tok.isPunct("...") {
			// vaargs
			tok.errorf("TBI: VAARGS(...)")
		} else {
			tok.errorf("Unkonwn token")
			p.unreadToken()
			break
		}

	}
	errorf("Unkown type")
	return nil
}

func (p *parser) parseVarDecl(isGlobal bool) *DeclVar {
	assert(p.lastToken().isKeyword("var"), "last token is \"var\"")
	defer p.traceOut(p.traceIn())
	// read newName
	newName := p.readIdent()
	var typ *Gtype
	var initval Expr
	// "=" or Type
	tok := p.readToken()
	if tok.isPunct("=") {
		initval = p.parseExpr()
		if typ == nil {
			typ = inferType(initval)
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
	r := &DeclVar{
		pkg: p.currentPackageName,
		variable: variable,
		initval:  initval,
	}
	p.currentScope.setVar(newName, variable)
	return r
}

func (p *parser) parseConstDeclSingle(lastExpr Expr, iotaIndex int) *ExprConstVariable {
	defer p.traceOut(p.traceIn())
	newName := p.readIdent()

	// Type or "=" or ";"
	var val Expr
	if !p.peekToken().isPunct("=") && !p.peekToken().isPunct(";") {
		// expect Type
		_ = p.parseType()
	}

	if p.peekToken().isPunct(";") && lastExpr != nil {
		val = lastExpr
	} else {
		p.expect("=")
		val = p.parseExpr()
	}
	p.expect(";")

	variable := &ExprConstVariable{
		name:      newName,
		val:       val,
		iotaIndex: iotaIndex,

	}

	p.currentScope.setConst(newName, variable)
	return variable
}

func (p *parser) parseConstDecl() *DeclConst {
	defer p.traceOut(p.traceIn())
	// ident or "("
	var cnsts []*ExprConstVariable
	var iotaIndex int
	var lastExpr Expr

	if p.peekToken().isPunct("(") {
		p.readToken()
		for {
			// multi definitions
			cnst := p.parseConstDeclSingle(lastExpr, iotaIndex)
			lastExpr = cnst.val
			iotaIndex++
			cnsts = append(cnsts, cnst)
			if p.peekToken().isPunct(")") {
				p.readToken()
				break
			}
		}
	} else {
		// single definition
		cnsts = []*ExprConstVariable{p.parseConstDeclSingle(nil, 0)}
	}

	r := &DeclConst{
		consts: cnsts,
	}

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
func (p *parser) parseForStmt() *StmtFor {
	defer p.traceOut(p.traceIn())
	var r = &StmtFor{}
	p.enterNewScope()
	defer p.exitScope()

	var cond Expr
	if p.peekToken().isPunct("{") {
		// inifinit loop : for { ... }
	} else {
		p.requireBlock = true
		cond = p.parseExpr()
		p.requireBlock = false
	}
	if p.peekToken().isPunct("{") {
		// single cond
		r.cls = &ForForClause{
			cond: cond,
		}
	} else {
		// for clause or range clause
		var initstmt Stmt
		lefts := p.parseExpressionList(cond)
		tok2 := p.readToken()
		if tok2.isPunct("=") {
			if p.peekToken().isKeyword("range") {
				return p.parseForRange(lefts)
			} else {
				initstmt = p.parseAssignment(lefts)
			}
		} else if tok2.isPunct(":=") {
			if p.peekToken().isKeyword("range") {
				assert(len(lefts) == 1 || len(lefts) == 2 , "lefts is not empty")
				e := lefts[0]
				rel := e.(*Relation) // a brand new rel
				gtype := gInt // index is int
				variable := p.newVariable(rel.name, gtype, false)
				rel.expr = variable
				p.currentScope.setVar(rel.name, variable)

				if len(lefts) == 2 {
					e := lefts[1]
					rel := e.(*Relation) // a brand new rel
					gtype := inferType(rel.expr)
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

func (p *parser) parseForRange(exprs []Expr) *StmtFor {
	defer p.traceOut(p.traceIn())
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

	var r = &StmtFor{
		rng: &ForRangeClause{
			indexvar: indexvar,
			valuevar: valuevar,
		},
	}
	p.requireBlock = true
	r.rng.rangeexpr = p.parseExpr()
	p.requireBlock = false
	p.expect("{")
	r.block = p.parseCompoundStmt()
	return r
}

func (p *parser) parseIfStmt() *StmtIf {
	defer p.traceOut(p.traceIn())
	var r = &StmtIf{}
	p.enterNewScope()
	defer p.exitScope()
	p.requireBlock = true
	stmt := p.parseStmt()
	if p.peekToken().isPunct(";") {
		p.skip()
		r.simplestmt = stmt
		r.cond = p.parseExpr()
	} else {
		es, ok := stmt.(*StmtExpr)
		if !ok {
			errorf("internal error")
		}
		r.cond = es.expr
	}
	p.expect("{")
	p.requireBlock = false
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

func (p *parser) parseReturnStmt() *StmtReturn {
	defer p.traceOut(p.traceIn())
	var r *StmtReturn
	exprs := p.parseExpressionList(nil)
	// workaround for {nil}
	if len(exprs) == 1 && exprs[0] == nil {
		exprs = nil
	}
	r = &StmtReturn{exprs: exprs}
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

func (p *parser) parseAssignment(lefts []Expr) *StmtAssignment {
	defer p.traceOut(p.traceIn())
	rights := p.parseExpressionList(nil)
	assert(rights[0] != nil , "rights[0] is an expr")
	return &StmtAssignment{
		lefts:  lefts,
		rights: rights,
	}
}

func (p *parser) parseAssignmentOperation(left Expr, assignop string) *StmtAssignment {
	defer p.traceOut(p.traceIn())
	var op string
	switch assignop {
	case "+=":
		op = "+"
	case "-=":
		op = "-"
	case "*=":
		op = "*"
	default:
		errorf("internal error")
	}
	rights := p.parseExpressionList(nil)
	assert(len(rights) == 1, "num of rights is 1")
	binop := &ExprBinop{
		op:    op,
		left:  left,
		right: rights[0],
	}
	return &StmtAssignment{
		lefts: []Expr{left},
		rights: []Expr{binop},
	}
}

func inferType(e Expr) *Gtype {
	switch e.(type) {
	case *ExprArrayLiteral:
		return e.(*ExprArrayLiteral).gtype
	case *ExprStructLiteral:
		strct := e.(*ExprStructLiteral)
		return &Gtype{
			typ: G_REL,
			relation: strct.strctname,
		}
	case *ExprUop:
		uop := e.(*ExprUop)
		if uop.op == "&" {
			return &Gtype{
				typ: G_POINTER,
				ptr: inferType(uop.operand),
			}
		} else if uop.op == "-" {
			return gInt
		}
	default:
		return gInt
	}
	return nil
}

func (p *parser) parseShortAssignment(lefts []Expr) *StmtShortVarDecl {
	defer p.traceOut(p.traceIn())
	rights := p.parseExpressionList(nil)
	for _, e := range lefts {
		rel := e.(*Relation) // a brand new rel
		variable := p.newVariable(rel.name, nil, false)
		rel.expr = variable
		p.currentScope.setVar(rel.name, variable)
	}
	r := &StmtShortVarDecl{
		lefts:  lefts,
		rights: rights,
	}
	p.shortassignments = append(p.shortassignments, r)
	return r
}

func (p *parser) parseSwitchStmt() Stmt {
	defer p.traceOut(p.traceIn())
	p.expectKeyword("switch")
	var cond Expr
	if p.peekToken().isPunct("{") {

	} else {
		p.requireBlock = true
		cond = p.parseExpr()
		p.requireBlock = false
	}
	p.expect("{")
	r := &StmtSwitch{
		cond:cond,
	}

	for {
		tok := p.peekToken()
		if tok.isKeyword("case") {
			p.skip()
			var exprs []Expr
			expr := p.parseExpr()
			exprs = append(exprs, expr)
			for {
				tok := p.peekToken()
				if tok.isPunct(",") {
					p.skip()
					expr := p.parseExpr()
					exprs = append(exprs, expr)
				} else if tok.isPunct(":") {
					break
				}
			}
			p.expect(":")
			p.inCase++
			compound := p.parseCompoundStmt()
			casestmt := &ExprCaseClause{
				exprs: exprs,
				compound: compound,
			}
			p.inCase--
			r.cases = append(r.cases, casestmt)
			if p.lastToken().isPunct("}") {
				break
			}
		} else if tok.isKeyword("default") {
			p.skip()
			p.expect(":")
			compound := p.parseCompoundStmt()
			r.dflt = compound
			break
		} else {
			errorf("internal error")
		}
	}

	return r
}

func (p *parser) parseDeferStmt() *StmtDefer {
	p.expectKeyword("defer")
	callExpr := p.parsePrim()
	return &StmtDefer{
		expr :callExpr,
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
	} else if tok.isKeyword("switch") {
		p.unreadToken()
		return p.parseSwitchStmt()
	} else if tok.isKeyword("continue") {
		return &StmtContinue{}
	} else if tok.isKeyword("break") {
		return &StmtBreak{}
	} else if tok.isKeyword("defer") {
		p.unreadToken()
		return p.parseDeferStmt()
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
	} else if tok2.isPunct("+=") || tok2.isPunct("-=") || tok2.isPunct("*=") {
		return p.parseAssignmentOperation(expr1, tok2.sval)
	} else if tok2.isPunct("++") {
		return &StmtInc{
			operand:expr1,
		}
	} else if tok2.isPunct("--")  {
		return &StmtDec{
			operand:expr1,
		}
	} else {
		p.unreadToken()
		return &StmtExpr{
			expr: expr1,
		}
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
		if p.inCase > 0 && (tok.isKeyword("case") || tok.isKeyword("default")) {
			p.unreadToken()
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

func (p *parser) parseFuncSignature() (identifier, []*ExprVariable, bool, []*Gtype) {
	defer p.traceOut(p.traceIn())
	tok := p.readToken()
	fname := tok.getIdent()
	p.expect("(")

	tok = p.readToken()

	var params []*ExprVariable
	var isVariadic bool

	if !tok.isPunct(")") {
		p.unreadToken()
		for {
			tok := p.readToken()
			pname := tok.getIdent()
			if p.peekToken().isPunct("...") {
				p.expect("...")
				gtype := p.parseType()
				variable := &ExprVariable{
					varname: pname,
					gtype : gtype,
				}
				params = append(params, variable)
				p.expect(")")
				break
			}
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

	next := p.peekToken()
	if next.isPunct("{") || next.isSemicolon() {
		return fname, params, isVariadic, nil
	}

	var rettypes []*Gtype
	if next.isPunct("(") {
		p.skip()
		for {
			rettype := p.parseType()
			rettypes = append(rettypes, rettype)
			next := p.peekToken()
			if next.isPunct(")") {
				p.skip()
				break
			} else if next.isPunct(",") {
				p.skip()
			} else {
				next.errorf("invalid token")
			}
		}

	} else {
		rettypes = []*Gtype{p.parseType()}
	}

	return fname, params, isVariadic, rettypes
}

func (p *parser) parseFuncDef() *DeclFunc {
	defer p.traceOut(p.traceIn())
	p.localvars = nil
	var isMethod bool
	p.enterNewScope()
	defer p.exitScope()

	var receiver *ExprVariable

	if p.peekToken().isPunct("(") {
		isMethod = true
		p.expect("(")
		// method definition
		tok := p.readToken()
		pname := tok.getIdent()
		ptype := p.parseType()
		receiver = &ExprVariable{
			varname: pname,
			gtype:   ptype,
		}
		p.currentScope.setVar(pname, receiver)
		p.expect(")")
	}

	fname, params, isVariadic, rettypes := p.parseFuncSignature()

	p.expect("{")

	r := &DeclFunc{
		pkg : p.currentPackageName,
		receiver:receiver,
		fname:     fname,
		rettypes:   rettypes,
		params:    params,
		isVariadic:isVariadic,
	}
	ref := &ExprFuncRef{
		funcdef: r,
	}

	if isMethod {
		var typeToBelong *Gtype
		if receiver.gtype.typ == G_POINTER {
			typeToBelong = receiver.gtype.ptr
		} else {
			typeToBelong = receiver.gtype
		}

		assert(typeToBelong.typ == G_REL , "methods must belong to a named type")
		methods, ok := p.namedTypes[typeToBelong.relation.name]
		if !ok {
			methods = make(map[identifier]*ExprFuncRef)
			p.namedTypes[typeToBelong.relation.name] = methods
		}
		methods[fname] = ref
	} else {
		p.packageBlockScope.setFunc(fname,ref)
	}

	body := p.parseCompoundStmt()
	r.body = body
	r.localvars = p.localvars

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
				specs = append(specs, &AstImportSpec{
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
		specs = []*AstImportSpec{&AstImportSpec{
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

const MaxAlign = 16

func (p *parser) parseStructDef() *Gtype {
	assert(p.lastToken().isKeyword("struct"), `require "struct" is already read`)
	defer p.traceOut(p.traceIn())

	p.expect("{")
	var fields []*Gtype
	for {
		tok := p.readToken()
		if tok.isPunct("}") {
			break
		}
		fieldname := tok.getIdent()
		gtype := p.parseType()
		fieldtype := *gtype
		fieldtype.ptr = gtype
		fieldtype.fieldname = fieldname
		fieldtype.offset = 0 // will be calculated later
		fields = append(fields, &fieldtype)
		p.expect(";")
	}
	// calc offset
	p.expect(";")
	return &Gtype{
		typ:    G_STRUCT,
		size:   0, // will be calculated later
		fields: fields,
	}
}

func (p *parser) parseInterfaceDef(newName identifier) *DeclType {
	defer p.traceOut(p.traceIn())
	p.expectKeyword("interface")
	p.expect("{")
	var methods []*signature
	for {
		if p.peekToken().isPunct("}") {
			break
		}

		fname, params, isVariadic, rettypes := p.parseFuncSignature()
		p.expect(";")

		var paramTypes []*Gtype
		for _, param := range params {
			paramTypes = append(paramTypes, param.gtype)
		}
		method := &signature{
			fname:fname,
			paramTypes:paramTypes,
			isVariadic: isVariadic,
			rettypes:rettypes,
		}
		methods = append(methods, method)
	}
	p.expect("}")

	gtype := &Gtype{
		typ:      G_INTERFACE,
		imethods: methods,
	}

	p.currentScope.setGtype(newName, gtype)
	r := &DeclType{
		name:  newName,
		gtype: gtype,
	}
	return r
}

func (p *parser) tryResolve(pkg identifier, rel *Relation) {
	if rel.gtype != nil || rel.expr != nil {
		return
	}

	if pkg == "" {
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
	} else {
		// foreign package
		relfound := p.scopes[pkg].get(rel.name)
		if relfound == nil {
			errorf("name %s is not found in %s package", rel.name, pkg)
		}
		switch relfound.(type) {
		case *Gtype:
			rel.gtype = relfound.(*Gtype)
		case Expr:
			rel.expr = relfound.(Expr)
		default:
			errorf("Bad type relfound %v", relfound)
		}
	}
}

func (p *parser) parseTypeDecl() *DeclType {
	defer p.traceOut(p.traceIn())
	newName := p.readIdent()
	if p.peekToken().isKeyword("interface") {
		return p.parseInterfaceDef(newName)
	}

	gtype := p.parseType()
	r := &DeclType{
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
	case tok.isKeyword("func"):
		funcdecl := p.parseFuncDef()
		r = &AstTopLevelDecl{funcdecl: funcdecl}
	case tok.isKeyword("var"):
		vardecl := p.parseVarDecl(true)
		r = &AstTopLevelDecl{vardecl: vardecl}
	case tok.isKeyword("const"):
		constdecl := p.parseConstDecl()
		r = &AstTopLevelDecl{constdecl: constdecl}
	case tok.isKeyword("type"):
		typedecl := p.parseTypeDecl()
		r = &AstTopLevelDecl{typedecl: typedecl}
	default:
		errorf("TBD: unable to handle token %v", tok)
	}

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
func (p *parser) parseSourceFile(bs *ByteStream, packageBlockScope *scope) *AstFile {

	ts := NewTokenStream(bs)
	p.tokenStream = ts

	p.packageBlockScope = packageBlockScope
	p.currentScope = packageBlockScope

	r := &AstFile{}
	r.pkg = p.expectPackageClause()
	r.imports = p.parseImportDecls()

	p.importedNames = make(map[identifier]bool)
	for _, importdecl := range r.imports {
		for _, spec := range importdecl.specs {
			var pkgName identifier
			if strings.Contains(spec.path,"/") {
				words := strings.Split(spec.path, "/")
				pkgName = identifier(words[len(words) -1])
			} else {
				pkgName = identifier(spec.path)
			}

			p.importedNames[pkgName] = true
		}
	}

	r.decls = p.parseTopLevelDecls()
	return r
}



func (ast *StmtShortVarDecl) inferTypes() {
	var rightTypes []*Gtype
	for i, rightExpr := range ast.rights {
		switch rightExpr.(type) {
		case *ExprFuncall:
			fcall := rightExpr.(*ExprFuncall)
			if fcall.getFuncDef() == nil {
				errorf("funcdef of %s is not found", fcall.fname)
			}
			for _, gtype := range fcall.getFuncDef().rettypes {
				rightTypes = append(rightTypes, gtype)
			}
		case *ExprMethodcall:
			fcall := rightExpr.(*ExprMethodcall)
			for _, gtype := range fcall.getFuncDef().rettypes {
				rightTypes = append(rightTypes, gtype)
			}
		default:
			gtype := inferType(rightExpr)
			if gtype == nil {
				errorf("rights[%d] gtype is nil", i)
			}
			rightTypes = append(rightTypes, gtype)
		}
	}

	for i, e := range ast.lefts {
		rel := e.(*Relation) // a brand new rel
		variable := rel.expr.(*ExprVariable)
		rightType := rightTypes[i]
		variable.gtype = rightType
	}

}

func (p *parser) resolve(universe *scope) {
	p.packageBlockScope.outer = universe
	for _, rel := range p.unresolvedRelations {
		p.tryResolve("",rel)
	}

	p.resolveMethods()
	p.inferShortAssignmentTypes()
}

// copy methods from p.nameTypes to gtype.methods of each type
func (p *parser) resolveMethods() {
	for typeName, methods := range p.namedTypes {
		gtype := p.packageBlockScope.getGtype(typeName)
		if gtype == nil {
			errorf("typaneme %s is not found in the package scope", typeName)
		}
		gtype.methods = methods
	}
}

func (p *parser) inferShortAssignmentTypes() {
	for _, sa := range p.shortassignments {
		sa.inferTypes()
	}
}
