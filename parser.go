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
	stringLiterals      []*ExprStringLiteral
	globalvars          []*ExprVariable
	localvars           []*ExprVariable
	namedTypes          map[identifier]methods
	globaluninferred    []*ExprVariable
	localuninferred     []Inferer // VarDecl, StmtShortVarDecl or RangeClause
	importedNames       map[identifier]bool
	requireBlock        bool // workaround for parsing "{" as a block starter
	inCase              int  // > 0  while in reading case compound stmts
	constSpecIndex      int
	currentPackageName  identifier
}

type methods map[identifier]*ExprFuncRef

func (p *parser) assert(cond bool, msg string) {
	assert(cond, p.lastToken(), msg)
}

func (p *parser) assertNotNil(x interface{}) {
	assertNotNil(x != nil, p.lastToken())
}

func (p *parser) peekToken() *Token {
	if p.tokenStream.isEnd() {
		return makeToken("EOF", "")
	}
	r := p.tokenStream.tokens[p.tokenStream.index]
	return r
}

func (p *parser) lastToken() *Token {
	return p.tokenStream.tokens[p.tokenStream.index-1]
}

// skip one token
func (p *parser) skip() {
	if p.tokenStream.isEnd() {
		return
	}
	p.tokenStream.index++
}

func (p *parser) readToken() *Token {
	tok := p.peekToken()
	p.skip()
	return tok
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

func (p *parser) expectKeyword(name string) *Token {
	tok := p.readToken()
	if !tok.isKeyword(name) {
		errorf("Keyword %s expected but got %s", name, tok)
	}
	return tok
}

func (p *parser) expect(punct string) *Token {
	tok := p.readToken()
	if !tok.isPunct(punct) {
		errorf("punct '%s' expected but got '%s'", punct, tok)
	}
	return tok
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
		tok := p.peekToken()
		if tok.isPunct(")") {
			p.skip()
			return r
		}
		arg := p.parseExpr()
		if p.peekToken().isPunct("...") {
			p.expect("...")
			arg = &ExprVaArg{expr: arg}
			r = append(r, arg)
			p.expect(")")
			return r
		}
		r = append(r, arg)
		tok = p.peekToken()
		if tok.isPunct(")") {
			p.skip()
			return r
		} else if tok.isPunct(",") {
			p.skip()
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
		e = p.parseStructLiteral(rel)
	} else if next.isPunct("(") {
		// funcall or method call
		p.skip()
		args := p.readFuncallArgs()
		fname := string(rel.name)
		e = &ExprFuncall{
			rel:   rel,
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
	tok := p.peekToken()
	if tok.isPunct(":") {
		p.skip()
		lowIndex := &ExprNumberLiteral{
			val: 0,
		}
		highIndex := p.parseExpr()
		p.expect("]")
		r = &ExprSliced{
			low:  lowIndex,
			high: highIndex,
		}
	} else {
		index := p.parseExpr()
		tok := p.peekToken()
		if tok.isPunct("]") {
			p.skip()
			r = &ExprArrayIndex{
				array: e,
				index: index,
			}
		} else if tok.isPunct(":") {
			p.skip()
			highIndex := p.parseExpr()
			p.expect("]")
			r = &ExprSliced{
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
			expr: e,
		}
	} else {
		gtype := p.parseType()
		p.expect(")")
		e = &ExprTypeAssertion{
			expr:  e,
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
		tok := p.readToken()
		next = p.peekToken()
		if next.isPunct("(") {
			// (expr).method()
			p.expect("(")
			args := p.readFuncallArgs()
			r = &ExprMethodcall{
				tok:      tok,
				receiver: e,
				fname:    tok.getIdent(),
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
			r = &ExprStructField{
				tok:       tok,
				strct:     e,
				fieldname: tok.getIdent(),
			}
			return p.succeedingExpr(r)
		}
	} else if next.isPunct("[") {
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
	p.assert(tok.isIdent("make"), "read make")

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
		typ:      G_MAP,
		mapKey:   mapKey,
		mapValue: mapValue,
	}
}

// https://golang.org/ref/spec#Conversions
func (p *parser) parseTypeConversion(gtype *Gtype) Expr {
	defer p.traceOut(p.traceIn())

	p.expect("(")
	e := p.parseExpr()
	p.expect(")")

	return &ExprConversion{
		expr:  e,
		gtype: gtype,
	}
}

func (p *parser) parsePrim() Expr {
	defer p.traceOut(p.traceIn())
	tok := p.peekToken()

	switch {
	case tok.isSemicolon():
		p.skip()
		return nil
	case tok.isTypeString(): // string literal
		p.skip()
		ast := &ExprStringLiteral{
			val: tok.sval,
		}
		p.addStringLiteral(ast)
		return ast
	case tok.isTypeInt(): // int literal
		p.skip()
		ival := tok.getIntval()
		return &ExprNumberLiteral{
			val: ival,
		}
	case tok.isTypeChar(): // char literal
		p.skip()
		sval := tok.sval
		c := sval[0]
		return &ExprNumberLiteral{
			val: int(c),
		}
	case tok.isPunct("["): // array literal or type casting
		gtype := p.parseType()
		if p.peekToken().isPunct("(") {
			// Conversion
			return p.parseTypeConversion(gtype)
		}
		return p.parseArrayLiteral(gtype)
	case tok.isIdent("make"):
		return p.parseMakeExpr()
	case tok.isTypeIdent():
		p.skip()
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
		p.assertNotNil(v)
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
	defer p.traceOut(p.traceIn())
	p.expect("{")

	r := &ExprStructLiteral{
		strctname: rel,
	}

	for {
		tok := p.readToken()
		if tok.isPunct("}") {
			break
		}
		p.expect(":")
		p.assert(tok.isTypeIdent(), "field name is ident")
		value := p.parseExpr()
		f := &KeyedElement{
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
				typ:      G_REL,
				relation: strctliteral.strctname,
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
			op:      tok.sval,
			operand: p.parsePrim(),
		}
	case tok.isPunct("-"):
		return &ExprUop{
			op:      tok.sval,
			operand: p.parsePrim(),
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
		tok := p.peekToken()
		if tok.isSemicolon() {
			return ast
		}

		// if bion
		if in_array(tok.sval, binops) {
			prior2 := priority(tok.sval)
			if prior < prior2 {
				p.skip()
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
				return ast
			}
		} else {
			return ast
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
			p.unreadToken()
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
		}

	}
	errorf("Unkown type")
	return nil
}

// local decl infer
func (decl *DeclVar) infer() {
	gtype := decl.initval.getGtype()
	assertNotNil(gtype != nil, nil)
	decl.variable.gtype = gtype
}

func (p *parser) parseVarDecl(isGlobal bool) *DeclVar {
	defer p.traceOut(p.traceIn())
	p.expectKeyword("var")

	// read newName
	newName := p.readIdent()
	var typ *Gtype
	var initval Expr
	// "=" or Type
	tok := p.peekToken()
	if tok.isPunct("=") {
		p.skip()
		initval = p.parseExpr()
	} else {
		typ = p.parseType()
		p.assertNotNil(typ)
		tok := p.readToken()
		if tok.isPunct("=") {
			initval = p.parseExpr()
		}
	}
	//p.expect(";")

	variable := p.newVariable(newName, typ, isGlobal)
	r := &DeclVar{
		pkg:      p.currentPackageName,
		variable: variable,
		initval:  initval,
	}
	if typ == nil {
		if isGlobal {
			variable.gtype = &Gtype{
				typ:          G_DEPENDENT,
				dependendson: initval,
			}
			p.globaluninferred = append(p.globaluninferred, variable)
		} else {
			p.localuninferred = append(p.localuninferred, r)
		}
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
	p.expectKeyword("const")

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

		tok = p.peekToken()
		if tok.isPunct(",") {
			p.skip()
			continue
		} else {
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

func (clause *ForRangeClause) infer() {
	collectionType := clause.rangeexpr.getGtype()
	debugf("collectionType = %s", collectionType)
	indexvar, ok := clause.indexvar.expr.(*ExprVariable)
	assert(ok, nil, "ok")

	var indexType *Gtype
	switch collectionType.typ {
	case G_ARRAY, G_SLICE:
		indexType = gInt
	default:
		// @TODO consider map etc.
		errorf("TBI")
	}
	indexvar.gtype = indexType

	if clause.valuevar != nil {
		valuevar, ok := clause.valuevar.expr.(*ExprVariable)
		assert(ok, nil, "ok")

		var elementType *Gtype
		if collectionType.typ == G_ARRAY {
			elementType = collectionType.ptr
		} else if collectionType.typ == G_SLICE {
			elementType = collectionType.ptr // @TODO is this right ?
		}
		debugf("for i, v %s := rannge %v", elementType, collectionType)
		valuevar.gtype = elementType
	}
}

// https://golang.org/ref/spec#For_statements
func (p *parser) parseForStmt() *StmtFor {
	defer p.traceOut(p.traceIn())
	p.expectKeyword("for")

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
		tok2 := p.peekToken()
		if tok2.isPunct("=") {
			p.skip()
			if p.peekToken().isKeyword("range") {
				return p.parseForRange(lefts)
			} else {
				initstmt = p.parseAssignment(lefts)
			}
		} else if tok2.isPunct(":=") {
			p.skip()
			if p.peekToken().isKeyword("range") {
				p.assert(len(lefts) == 1 || len(lefts) == 2, "lefts is not empty")
				p.shortVarDecl(lefts[0])

				if len(lefts) == 2 {
					p.shortVarDecl(lefts[1])
				}

				r := p.parseForRange(lefts)
				p.localuninferred = append(p.localuninferred, r.rng)
				return r
			} else {
				p.unreadToken()
				initstmt = p.parseShortAssignment(lefts)
			}
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

	p.requireBlock = true
	rangeExpr := p.parseExpr()
	p.requireBlock = false
	p.expect("{")

	var r = &StmtFor{
		rng: &ForRangeClause{
			indexvar:  indexvar,
			valuevar:  valuevar,
			rangeexpr: rangeExpr,
		},
	}
	r.block = p.parseCompoundStmt()
	return r
}

func (p *parser) parseIfStmt() *StmtIf {
	defer p.traceOut(p.traceIn())
	p.expectKeyword("if")

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
	tok := p.peekToken()
	if tok.isKeyword("else") {
		p.skip()
		tok2 := p.peekToken()
		if tok2.isKeyword("if") {
			// we regard "else if" as a kind of a nested if statement
			// else if => else { if .. { } else {} }
			r.els = p.parseIfStmt()
		} else if tok2.isPunct("{") {
			p.skip()
			r.els = p.parseCompoundStmt()
		} else {
			tok2.errorf("Unexpected token")
		}
	}
	return r
}

func (p *parser) parseReturnStmt() *StmtReturn {
	defer p.traceOut(p.traceIn())
	p.expectKeyword("return")

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
		tok := p.peekToken()
		if tok.isSemicolon() {
			return r
		} else if tok.isPunct("=") || tok.isPunct(":=") {
			return r
		} else if tok.isPunct(",") {
			p.skip()
			expr := p.parseExpr()
			r = append(r, expr)
			continue
		} else {
			return r
		}
	}
	return r
}

func (p *parser) parseAssignment(lefts []Expr) *StmtAssignment {
	defer p.traceOut(p.traceIn())

	rights := p.parseExpressionList(nil)
	p.assertNotNil(rights[0])
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
	p.assert(len(rights) == 1, "num of rights is 1")
	binop := &ExprBinop{
		op:    op,
		left:  left,
		right: rights[0],
	}
	return &StmtAssignment{
		lefts:  []Expr{left},
		rights: []Expr{binop},
	}
}

func (p *parser) shortVarDecl(e Expr) {
	rel := e.(*Relation) // a brand new rel
	variable := p.newVariable(rel.name, nil, false)
	p.currentScope.setVar(rel.name, variable)
	rel.expr = variable
}

func (p *parser) parseShortAssignment(lefts []Expr) *StmtShortVarDecl {
	defer p.traceOut(p.traceIn())
	separator := p.expect(":=")

	rights := p.parseExpressionList(nil)
	for _, e := range lefts {
		p.shortVarDecl(e)
	}
	r := &StmtShortVarDecl{
		tok:    separator,
		lefts:  lefts,
		rights: rights,
	}
	p.localuninferred = append(p.localuninferred, r)
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
		cond: cond,
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
				exprs:    exprs,
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
	defer p.traceOut(p.traceIn())
	p.expectKeyword("defer")

	callExpr := p.parsePrim()
	return &StmtDefer{
		expr: callExpr,
	}
}

// this is in function scope
func (p *parser) parseStmt() Stmt {
	defer p.traceOut(p.traceIn())

	tok := p.peekToken()
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
		return p.parseSwitchStmt()
	} else if tok.isKeyword("continue") {
		p.expectKeyword("continue")
		return &StmtContinue{}
	} else if tok.isKeyword("break") {
		p.expectKeyword("break")
		return &StmtBreak{}
	} else if tok.isKeyword("defer") {
		return p.parseDeferStmt()
	}

	expr1 := p.parseExpr()
	tok2 := p.peekToken()
	if tok2.isPunct(",") {
		// Multi value assignment
		lefts := p.parseExpressionList(expr1)
		tok3 := p.peekToken()
		if tok3.isPunct("=") {
			p.skip()
			return p.parseAssignment(lefts)
		} else if tok3.isPunct(":=") {
			return p.parseShortAssignment(lefts)
		} else {
			tok3.errorf("TBD")
		}
	} else if tok2.isPunct("=") {
		p.skip()
		return p.parseAssignment([]Expr{expr1})
	} else if tok2.isPunct(":=") {
		// Single value ShortVarDecl
		return p.parseShortAssignment([]Expr{expr1})
	} else if tok2.isPunct("+=") || tok2.isPunct("-=") || tok2.isPunct("*=") {
		p.skip()
		return p.parseAssignmentOperation(expr1, tok2.sval)
	} else if tok2.isPunct("++") {
		p.skip()
		return &StmtInc{
			operand: expr1,
		}
	} else if tok2.isPunct("--") {
		p.skip()
		return &StmtDec{
			operand: expr1,
		}
	} else {
		return &StmtExpr{
			expr: expr1,
		}
	}
	return nil
}

func (p *parser) parseCompoundStmt() *StmtSatementList {
	defer p.traceOut(p.traceIn())

	r := &StmtSatementList{}
	for {
		tok := p.peekToken()
		if tok.isPunct("}") {
			p.skip()
			return r
		}
		if p.inCase > 0 && (tok.isKeyword("case") || tok.isKeyword("default")) {
			return r
		}
		if tok.isSemicolon() {
			p.skip()
			continue
		}
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

	var params []*ExprVariable
	var isVariadic bool

	tok = p.peekToken()
	if tok.isPunct(")") {
		p.skip()
	} else {
		for {
			tok := p.readToken()
			pname := tok.getIdent()
			if p.peekToken().isPunct("...") {
				p.expect("...")
				gtype := p.parseType()
				variable := &ExprVariable{
					varname: pname,
					gtype:   gtype,
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
	p.expectKeyword("func")

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
		pkg:        p.currentPackageName,
		receiver:   receiver,
		fname:      fname,
		rettypes:   rettypes,
		params:     params,
		isVariadic: isVariadic,
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

		p.assert(typeToBelong.typ == G_REL, "methods must belong to a named type")
		methods, ok := p.namedTypes[typeToBelong.relation.name]
		if !ok {
			methods = make(map[identifier]*ExprFuncRef)
			p.namedTypes[typeToBelong.relation.name] = methods
		}
		methods[fname] = ref
	} else {
		p.packageBlockScope.setFunc(fname, ref)
	}

	body := p.parseCompoundStmt()
	r.body = body
	r.localvars = p.localvars

	p.localvars = nil
	return r
}

func (p *parser) parseImport() *ImportDecl {
	defer p.traceOut(p.traceIn())
	tokImport := p.expectKeyword("import")

	tok := p.readToken()
	var specs []*ImportSpec
	if tok.isPunct("(") {
		for {
			tok := p.readToken()
			if tok.isTypeString() {
				specs = append(specs, &ImportSpec{
					path: tok.sval,
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
		specs = []*ImportSpec{&ImportSpec{
			path: tok.sval,
		},
		}
	}
	p.expect(";")
	return &ImportDecl{
		tok: tokImport,
		specs: specs,
	}
}

func (p *parser) parsePackageClause() *PackageClause {
	defer p.traceOut(p.traceIn())
	tokPkg := p.expectKeyword("package")

	name := p.readIdent()
	p.expect(";")
	return &PackageClause{
		tok: tokPkg,
		name: name,
	}
}

func (p *parser) parseImportDecls() []*ImportDecl {
	defer p.traceOut(p.traceIn())
	var r []*ImportDecl
	for p.peekToken().isKeyword("import") {
		r = append(r, p.parseImport())
	}
	return r
}

const MaxAlign = 16

func (p *parser) parseStructDef() *Gtype {
	defer p.traceOut(p.traceIn())
	p.expectKeyword("struct")

	p.expect("{")
	var fields []*Gtype
	for {
		tok := p.peekToken()
		if tok.isPunct("}") {
			p.skip()
			break
		}
		fieldname := tok.getIdent()
		p.skip()
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
			fname:      fname,
			paramTypes: paramTypes,
			isVariadic: isVariadic,
			rettypes:   rettypes,
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
	p.expectKeyword("type")

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

// https://golang.org/ref/spec#TopLevelDecl
// TopLevelDecl  = Declaration | FunctionDecl | MethodDecl .
func (p *parser) parseTopLevelDecl(nextToken *Token) *TopLevelDecl {
	defer p.traceOut(p.traceIn())

	switch {
	case nextToken.isKeyword("func"):
		funcdecl := p.parseFuncDef()
		return &TopLevelDecl{funcdecl: funcdecl}
	case nextToken.isKeyword("var"):
		vardecl := p.parseVarDecl(true)
		return &TopLevelDecl{vardecl: vardecl}
	case nextToken.isKeyword("const"):
		constdecl := p.parseConstDecl()
		return &TopLevelDecl{constdecl: constdecl}
	case nextToken.isKeyword("type"):
		typedecl := p.parseTypeDecl()
		return &TopLevelDecl{typedecl: typedecl}
	}

	errorf("TBD: unable to handle token %v", nextToken)
	return nil
}

func (p *parser) parseTopLevelDecls() []*TopLevelDecl {
	defer p.traceOut(p.traceIn())

	var r []*TopLevelDecl
	for {
		tok := p.peekToken()
		if tok.isEOF() {
			return r
		}

		if tok.isPunct(";") {
			p.skip()
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
func (p *parser) parseSourceFile(bs *ByteStream, packageBlockScope *scope) *SourceFile {

	// initialize parser's status per file
	p.tokenStream = NewTokenStream(bs)
	p.packageBlockScope = packageBlockScope
	p.currentScope = packageBlockScope
	p.importedNames = make(map[identifier]bool)

	packageClause := p.parsePackageClause()
	importDecls := p.parseImportDecls()

	// regsiter imported names
	for _, importdecl := range importDecls {
		for _, spec := range importdecl.specs {
			var pkgName identifier
			if strings.Contains(spec.path, "/") {
				words := strings.Split(spec.path, "/")
				pkgName = identifier(words[len(words)-1])
			} else {
				pkgName = identifier(spec.path)
			}

			p.importedNames[pkgName] = true
		}
	}

	// @TODO import external decls here

	topLevelDecls := p.parseTopLevelDecls()

	return &SourceFile{
		packageClause: packageClause,
		importDecls:   importDecls,
		topLevelDecls: topLevelDecls,
	}
}

func (ast *StmtShortVarDecl) infer() {
	debugf("infering %s", ast.tok)
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
			debugf("receiver=%v", fcall.receiver)
			strctfield, ok := fcall.receiver.(*ExprStructField)
			debugf("strctfield.strct=%v, %v", strctfield.strct, ok)
			for _, gtype := range fcall.getFuncDef().rettypes {
				rightTypes = append(rightTypes, gtype)
			}
		default:
			gtype := rightExpr.getGtype()
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
		p.tryResolve("", rel)
	}

	p.resolveMethods()
	p.inferTypes()
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

//  infer recursively all the types of global variables
func (variable *ExprVariable) infer() {
	if variable.gtype.typ != G_DEPENDENT {
		// done
		return
	}
	e := variable.gtype.dependendson
	dependType := e.getGtype()
	if dependType.typ != G_DEPENDENT {
		variable.gtype = dependType
		return
	}

	rel, ok := e.(*Relation)
	if !ok {
		errorf("NG %#v", e)
	}
	vr, ok := rel.expr.(*ExprVariable)
	vr.infer() // recursive call
	variable.gtype = e.getGtype()
}

func (p *parser) inferTypes() {
	for _, variable := range p.globaluninferred {
		variable.infer()
	}

	for _, ast := range p.localuninferred {
		debugf("> inferring local statement %T", ast)
		ast.infer()
	}
}
