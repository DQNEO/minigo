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
	namedTypes          map[identifier]methods
	importedNames       map[identifier]bool
	requireBlock        bool // workaround for parsing "{" as a block starter
}

type methods map[identifier]*ExprFuncRef

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

// skip one token
func (p *parser) skip() {
	tok := p.readToken()
	debugf("skipped %s", tok)
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

	// either of expr(var, const, funcref) or gtype
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

type AstStructFieldLiteral struct {
	key   identifier
	value Expr
}

type ExprStructLiteral struct {
	strctname *Relation
	fields    []*AstStructFieldLiteral
	invisiblevar *ExprVariable // to have offfset for &T{}
}

func (e *ExprStructLiteral) emit() {
	errorf("This cannot be emitted alone")
}

func (e *ExprStructLiteral) dump() {
	errorf("TBD")
}

type AstStructFieldAccess struct {
	strct     Expr
	fieldname identifier
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
		debugf("Reference to outer entity %s.%s", pkg, firstIdent)
	}

	rel := &Relation{
		name: firstIdent,
	}
	p.tryResolve(rel)

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
		p.skip()
		// assure operand is array, slice, or map
		tok := p.readToken()
		if tok.isPunct(":") {
			lowIndex := &ExprNumberLiteral{
				val: 0,
			}
			highIndex := p.parseExpr()
			p.expect("]")
			e = &ExprSliced{
				ref:  nil, // TBI
				low:  lowIndex,
				high: highIndex,
			}
		} else {
			p.unreadToken()
			index := p.parseExpr()
			tok := p.readToken()
			if tok.isPunct("]") {
				e = &ExprArrayIndex{
					array:   rel,
					index: index,
				}
			} else if tok.isPunct(":") {
				highIndex := p.parseExpr()
				p.expect("]")
				e = &ExprSliced{
					ref:  nil, // TBI
					low:  index,
					high: highIndex,
				}

			} else {
				tok.errorf("invalid token in index access")
			}
		}
	} else {
		// solo ident
		e = rel
	}

	return p.succeedingExpr(e)
}

func (p *parser) succeedingExpr(e Expr) Expr {
	defer p.traceOut(p.traceIn())

	var r Expr
	next := p.peekToken()
	if next.isPunct(".") {
		// https://golang.org/ref/spec#Selectors
		p.skip()
		ident := p.readIdent()
		debugf("read ident: %s", ident)
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
		} else {
			// (expr).field
			//   strct.field
			//   (strct.field).field
			//   fncall().field
			//   obj.method().field
			//   array[0].field
			r =  &AstStructFieldAccess{
				strct:     e,
				fieldname: ident,
			}
		}
	} else if next.isPunct("["){
		// https://golang.org/ref/spec#Index_expressions
		// (expr)[i]
		errorf("TBI")
	} else {
		// https://golang.org/ref/spec#OperandName
		r = e
	}

	return r

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
		return p.parseIdentExpr(tok.getIdent())
	}

	tok.errorf("unable to handle")
	return nil
}

func (p *parser) parseArrayLiteral() Expr {
	assert(p.lastToken().isPunct("["), "[ is read")
	defer p.traceOut(p.traceIn())
	var tlen *Token
	if !p.peekToken().isPunct("]") {
		tlen = p.readToken()
	}
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

	var length int
	if tlen == nil {
		length = len(values)
	} else {
		if len(values) != tlen.getIntval() {
			debugPrintV(values)
			errorf("array length does not match (%d != %d)",
				len(values), tlen.getIntval())
		}
		length = tlen.getIntval()
	}

	gtype := &Gtype{
		typ:    G_ARRAY,
		length: length,
		ptr:    typ,
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
				relation: rel,
			}
			return gtype
		} else if tok.isPunct("*") {
			// pointer
			gtype = &Gtype{
				typ: G_POINTER,
				ptr: p.parseType(),
			}
			debugf("ptr=%v", gtype.ptr)
			return gtype
		} else if tok.isKeyword("struct") {
			return p.parseStructDef()
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
	r := &AstVarDecl{
		variable: variable,
		initval:  initval,
	}
	p.currentScope.setVar(newName, variable)
	return r
}

func (p *parser) parseConstDeclSingle() *ExprConstVariable {
	newName := p.readIdent()

	// Type or "="
	var val Expr
	if !p.peekToken().isPunct("=") {
		// expect Type
		_ = p.parseType()
	}

	p.expect("=")
	val = p.parseExpr()

	p.expect(";")
	variable := &ExprConstVariable{
		name: newName,
		val:  val,
	}

	p.currentScope.setConst(newName, variable)
	return variable
}

func (p *parser) parseConstDecl() *AstConstDecl {
	defer p.traceOut(p.traceIn())
	// ident or "("
	var cnsts []*ExprConstVariable
	if p.peekToken().isPunct("(") {
		p.readToken()
		for {
			// multi definitions
			cnst := p.parseConstDeclSingle()
			cnsts = append(cnsts, cnst)
			if p.peekToken().isPunct(")") {
				p.readToken()
				break
			}
		}
	} else {
		// single definition
		cnsts = []*ExprConstVariable{p.parseConstDeclSingle()}
	}

	r := &AstConstDecl{
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
func (p *parser) parseForStmt() *AstForStmt {
	defer p.traceOut(p.traceIn())
	var r = &AstForStmt{}
	p.enterNewScope()
	defer p.exitScope()

	p.requireBlock = true
	expr := p.parseExpr()
	p.requireBlock = false
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
	p.requireBlock = true
	r.rng.rangeexpr = p.parseExpr()
	p.requireBlock = false
	p.expect("{")
	r.block = p.parseCompoundStmt()
	return r
}

func (p *parser) parseIfStmt() *AstIfStmt {
	defer p.traceOut(p.traceIn())
	var r = &AstIfStmt{}
	p.enterNewScope()
	defer p.exitScope()
	p.requireBlock = true
	r.cond = p.parseExpr()
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

func (p *parser) parseReturnStmt() *AstReturnStmt {
	defer p.traceOut(p.traceIn())
	var r *AstReturnStmt
	exprs := p.parseExpressionList(nil)
	// workaround for {nil}
	if len(exprs) == 1 && exprs[0] == nil {
		exprs = nil
	}
	r = &AstReturnStmt{exprs: exprs}
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
	assert(rights[0] != nil , "rights[0] is an expr")
	return &AstAssignment{
		lefts:  lefts,
		rights: rights,
	}
}

func (p *parser) parseAssignmentOperation(left Expr, assignop string) *AstAssignment {
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
	return &AstAssignment{
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
		}
	default:
		return gInt
	}
	return nil
}

func (p *parser) parseShortAssignment(lefts []Expr) *AstAssignment {
	rights := p.parseExpressionList(nil)
	for i, e := range lefts {
		rel := e.(*Relation) // a brand new rel
		right := rights[i] // @FIXME this is not correct any more
		gtype := inferType(right)
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
	} else if tok2.isPunct("+=") || tok2.isPunct("-=") || tok2.isPunct("*=") {
		return p.parseAssignmentOperation(expr1, tok2.sval)
	} else if tok2.isPunct("++") {
		return &AstIncrStmt{
			operand:expr1,
		}
	} else if tok2.isPunct("--")  {
		return &AstDecrStmt{
			operand:expr1,
		}
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

func (p *parser) parseFuncSignature() (identifier, []*ExprVariable, []*Gtype) {
	tok := p.readToken()
	fname := tok.getIdent()
	p.expect("(")

	tok = p.readToken()

	var params []*ExprVariable

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

	next := p.peekToken()
	if next.isPunct("{") || next.isSemicolon() {
		return fname, params, nil
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
	return fname, params, rettypes
}

func (p *parser) parseFuncDef() *AstFuncDecl {
	defer p.traceOut(p.traceIn())
	p.localvars = make([]*ExprVariable, 0)
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

	fname, params, rettypes := p.parseFuncSignature()

	p.expect("{")
	debugf("scope:%s", p.currentScope)

	r := &AstFuncDecl{
		receiver:receiver,
		fname:     fname,
		rettypes:   rettypes,
		params:    params,
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

func (p *parser) parseInterfaceDef(newName identifier) *AstTypeDecl {
	defer p.traceOut(p.traceIn())
	p.expectKeyword("interface")
	p.expect("{")
	var methods []*signature
	for {
		if p.peekToken().isPunct("}") {
			break
		}

		fname, params, rettypes := p.parseFuncSignature()
		p.expect(";")

		var paramTypes []*Gtype
		for _, param := range params {
			paramTypes = append(paramTypes, param.gtype)
		}
		method := &signature{
			fname:fname,
			paramTypes:paramTypes,
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
	r := &AstTypeDecl{
		name:  newName,
		gtype: gtype,
	}
	return r
}

func (p *parser) tryResolve(rel *Relation) {
	if rel.gtype != nil || rel.expr != nil {
		return
	}
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
	if p.peekToken().isKeyword("interface") {
		return p.parseInterfaceDef(newName)
	}

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

	p.resolveMethods()
}

// copy methods from p.nameTypes to gtype.methods of each type
func (p *parser) resolveMethods() {
	for typeName, methods := range p.namedTypes {
		gtype := p.packageBlockScope.getGtype(typeName)
		gtype.methods = methods
	}
}
