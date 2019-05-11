package main

import (
	"fmt"
	"os"
	"runtime"
)

const __func__ string = "__func__"

type parser struct {
	// per function or block
	currentFunc    *DeclFunc
	localvars      []*ExprVariable
	requireBlock   bool // workaround for parsing "{" as a block starter
	inCase         int  // > 0  while in reading case compound stmts
	constSpecIndex int
	currentForStmt *StmtFor

	// per file
	tokenStream       *TokenStream
	packageBlockScope *scope
	currentScope      *scope
	importedNames     map[identifier]bool

	// per package
	currentPackageName  identifier
	methods             map[identifier]methods
	unresolvedRelations []*Relation
	globaluninferred    []*ExprVariable
	localuninferred     []Inferer // VarDecl, StmtShortVarDecl or RangeClause

	// global state
	scopes          map[identifier]*scope
	stringLiterals  []*ExprStringLiteral
	allNamedTypes   []*DeclType
	allDynamicTypes []*Gtype
}

func (p *parser) clearLocalState() {
	p.currentFunc = nil
	p.localvars = nil
	p.requireBlock = false
	p.inCase = 0
	p.constSpecIndex = 0
	p.currentForStmt = nil
}

type methods map[identifier]*ExprFuncRef

func (p *parser) initPackage(pkgname identifier) {
	p.currentPackageName = pkgname
	p.methods = map[identifier]methods{}
	p.unresolvedRelations = nil
	p.globaluninferred = nil
	p.localuninferred = nil
}

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

func (p *parser) expectIdent() identifier {
	tok := p.readToken()
	if !tok.isTypeIdent() {
		errorft(tok, "Identifier expected, but got %s", tok)
	}
	return tok.getIdent()
}

func (p *parser) expectKeyword(name string) *Token {
	tok := p.readToken()
	if !tok.isKeyword(name) {
		errorft(tok, "Keyword %s expected but got %s", name, tok)
	}
	return tok
}

func (p *parser) expect(punct string) *Token {
	tok := p.readToken()
	if !tok.isPunct(punct) {
		errorft(tok, "punct '%s' expected but got '%s'", punct, tok)
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

func (p *parser) traceIn(funcname string) int {
	if !debugParser {
		return 0
	}
	if GENERATION == 1 {
		funcname = getCallerName(2)
	}
	debugf("func %s is gonna read %s", funcname, p.peekToken().sval)
	debugNest++
	return 0
}

func (p *parser) traceOut(funcname string) {
	if !debugParser {
		return
	}
	if r := recover(); r != nil {
		fmt.Println(r)
		os.Exit(1)
	}
	if GENERATION == 1 {
		funcname = getCallerName(2)
	}
	debugNest--
	debugf("func %s end after %s", funcname, p.lastToken().sval)
}

func (p *parser) readFuncallArgs() []Expr {
	p.traceIn(__func__)
	defer p.traceOut(__func__)
	var r []Expr
	for {
		tok := p.peekToken()
		if tok.isPunct(")") {
			p.skip()
			return r
		}
		arg := p.parseExpr()
		if p.peekToken().isPunct("...") {
			ptok := p.expect("...")
			arg = &ExprVaArg{
				tok:  ptok,
				expr: arg,
			}
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
			errorft(tok, "invalid token in funcall arguments")
		}
	}
}

//var outerPackages map[identifier](map[identifier]interface{})

func (p *parser) addStringLiteral(ast *ExprStringLiteral) {
	p.stringLiterals = append(p.stringLiterals, ast)
}

// expr which begins with an ident.
// e.g. ident, ident() or ident.*, etc
func (p *parser) parseIdentExpr(firstIdentToken *Token) Expr {
	p.traceIn(__func__)
	defer p.traceOut(__func__)

	firstIdent := firstIdentToken.getIdent()
	// https://golang.org/ref/spec#QualifiedIdent
	// read QualifiedIdent
	var pkg identifier // ignored for now
	if _, ok := p.importedNames[firstIdent]; ok {
		pkg = firstIdent
		p.expect(".")
		// shift firstident
		firstIdent = p.expectIdent()
	}

	rel := &Relation{
		tok:  firstIdentToken,
		name: firstIdent,
		pkg:  p.currentPackageName, // @TODO is this right?
	}
	if rel.name == "__func__" {
		sliteral := &ExprStringLiteral{
			val: string(p.currentFunc.fname),
		}
		rel.expr = sliteral
		p.addStringLiteral(sliteral)
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
		e = &ExprFuncallOrConversion{
			tok:   next,
			rel:   rel,
			fname: fname,
			args:  args,
		}
	} else if next.isPunct("[") {
		// index access
		var collection Expr = rel // @TODO: it should do auto conversion on function call
		e = p.parseIndexOrSliceExpr(collection)
	} else {
		// solo ident
		e = rel
	}

	return p.succeedingExpr(e)
}

// https://golang.org/ref/spec#Index_expressions
func (p *parser) parseIndexOrSliceExpr(e Expr) Expr {
	p.traceIn(__func__)
	defer p.traceOut(__func__)
	p.expect("[")

	var r Expr
	// assure operand is array, slice, or map
	tok := p.peekToken()
	if tok.isPunct(":") {
		p.skip()
		// A missing low index defaults to zero
		lowIndex := &ExprNumberLiteral{
			tok: tok,
			val: 0,
		}
		var highIndex Expr
		tok := p.peekToken()
		if tok.isPunct("]") {
			p.skip()
			// a missing high index defaults to the length of the sliced operand:
			// this must be resolved after resolving types
			highIndex = nil
		} else {
			highIndex = p.parseExpr()
			p.expect("]")
		}
		r = &ExprSlice{
			tok:        tok,
			collection: e,
			low:        lowIndex,
			high:       highIndex,
		}
	} else {
		index := p.parseExpr()
		tok := p.peekToken()
		if tok.isPunct("]") {
			p.skip()
			r = &ExprIndex{
				tok:        tok,
				collection: e,
				index:      index,
			}
		} else if tok.isPunct(":") {
			p.skip()
			var highIndex Expr
			tok := p.peekToken()
			if tok.isPunct("]") {
				p.skip()
				// a missing high index defaults to the length of the sliced operand:
				r = &ExprSlice{
					tok:        tok,
					collection: e,
					low:        index,
					high:       nil,
				}
			} else {
				highIndex = p.parseExpr()
				tok := p.peekToken()
				if tok.isPunct("]") {
					p.skip()
					r = &ExprSlice{
						tok:        tok,
						collection: e,
						low:        index,
						high:       highIndex,
					}
				} else if tok.isPunct(":") {
					p.skip()
					maxIndex := p.parseExpr()
					r = &ExprSlice{
						tok:        tok,
						collection: e,
						low:        index,
						high:       highIndex,
						max:        maxIndex,
					}
					p.expect("]")
				} else {
					errorft(tok, "invalid token in index access")
				}
			}
		} else {
			errorft(tok, "invalid token in index access")
		}
	}
	if r == nil {
		errorft(tok, "should not be nil")
	}
	return r
}

// https://golang.org/ref/spec#Type_assertions
func (p *parser) parseTypeAssertionOrTypeSwitchGuad(e Expr) Expr {
	p.traceIn(__func__)
	defer p.traceOut(__func__)
	ptok := p.expect("(")

	if p.peekToken().isKeyword("type") {
		p.skip()
		ptok := p.expect(")")
		return &ExprTypeSwitchGuard{
			tok:  ptok,
			expr: e,
		}
	} else {
		gtype := p.parseType()
		ptok := p.expect(")")
		e = &ExprTypeAssertion{
			tok:   ptok,
			expr:  e,
			gtype: gtype,
		}
		return p.succeedingExpr(e)
	}
	errorft(ptok, "internal error")
	return nil
}

func (p *parser) succeedingExpr(e Expr) Expr {
	p.traceIn(__func__)
	defer p.traceOut(__func__)

	var r Expr
	next := p.peekToken()
	if next.isPunct(".") {
		p.skip()
		if p.peekToken().isPunct("(") {
			// type assertion
			return p.parseTypeAssertionOrTypeSwitchGuad(e)
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
		e = p.parseIndexOrSliceExpr(e)
		return p.succeedingExpr(e)
	} else {
		// https://golang.org/ref/spec#OperandName
		r = e
	}

	return r

}

func (p *parser) parseMakeExpr() Expr {
	p.traceIn(__func__)
	defer p.traceOut(__func__)
	tok := p.readToken()
	p.assert(tok.isIdent("make"), "read make")

	p.expect("(")
	mapType := p.parseMapType()
	_ = mapType
	p.expect(")")
	return &ExprNilLiteral{
		tok: tok,
	}
}

func (p *parser) parseMapType() *Gtype {
	p.traceIn(__func__)
	defer p.traceOut(__func__)
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
	p.traceIn(__func__)
	defer p.traceOut(__func__)

	ptok := p.expect("(")
	e := p.parseExpr()
	p.expect(")")

	return &ExprConversion{
		tok:   ptok,
		gtype: gtype,
		expr:  e,
	}
}

func (p *parser) parsePrim() Expr {
	p.traceIn(__func__)
	defer p.traceOut(__func__)
	tok := p.peekToken()

	switch {
	case tok.isSemicolon():
		p.skip()
		return nil
	case tok.isTypeString(): // string literal
		p.skip()
		ast := &ExprStringLiteral{
			tok: tok,
			val: tok.sval,
		}
		p.addStringLiteral(ast)
		return ast
	case tok.isTypeInt(): // int literal
		p.skip()
		ival := tok.getIntval()
		return &ExprNumberLiteral{
			tok: tok,
			val: ival,
		}
	case tok.isTypeChar(): // char literal
		p.skip()
		sval := tok.sval
		c := sval[0]
		return &ExprNumberLiteral{
			tok: tok,
			val: int(c),
		}
	case tok.isKeyword("map"): // map literal
		ptok := tok
		gtype := p.parseType()
		p.expect("{")
		var elements []*MapElement
		for {
			if p.peekToken().isPunct("}") {
				p.skip()
				break
			}
			key := p.parseExpr()
			p.expect(":")
			value := p.parseExpr()
			p.expect(",")
			element := &MapElement{
				tok:   key.token(),
				key:   key,
				value: value,
			}
			elements = append(elements, element)
		}
		return &ExprMapLiteral{
			tok:      ptok,
			gtype:    gtype,
			elements: elements,
		}
	case tok.isPunct("["): // array literal, slice literal or type casting
		gtype := p.parseType()
		tok = p.peekToken()
		if tok.isPunct("(") {
			// Conversion
			return p.parseTypeConversion(gtype)
		}

		values := p.parseArrayLiteral()
		switch gtype.typ {
		case G_ARRAY:
			if gtype.typ == G_ARRAY {
				if gtype.length == 0 {
					gtype.length = len(values)
				} else {
					if gtype.length < len(values) {
						errorft(tok, "array length does not match (%d != %d)",
							len(values), gtype.length)
					}
				}
			}

			return &ExprArrayLiteral{
				tok:    tok,
				gtype:  gtype,
				values: values,
			}
		case G_SLICE:
			return &ExprSliceLiteral{
				tok:    tok,
				gtype:  gtype,
				values: values,
				invisiblevar: p.newVariable("", &Gtype{
					typ:         G_ARRAY,
					elementType: gtype.elementType,
					length:      len(values),
				}),
			}
		default:
			errorft(tok, "internal error")
		}
	case tok.isIdent("make"):
		return p.parseMakeExpr()
	case tok.isTypeIdent():
		p.skip()
		return p.parseIdentExpr(tok)
	}

	errorft(tok, "unable to handle")
	return nil
}

// for now, this is suppose to be either of
// array literal or slice literal
func (p *parser) parseArrayLiteral() []Expr {
	p.traceIn(__func__)
	defer p.traceOut(__func__)
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
			errorft(tok, "unpexpected token")
		}
	}

	return values
}

func (p *parser) parseStructLiteral(rel *Relation) *ExprStructLiteral {
	p.traceIn(__func__)
	defer p.traceOut(__func__)
	ptok := p.expect("{")

	r := &ExprStructLiteral{
		tok:       ptok,
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
			tok:   tok,
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
	p.traceIn(__func__)
	defer p.traceOut(__func__)

	tok := p.readToken()
	switch {
	case tok.isPunct("("):
		e := p.parseExpr()
		p.expect(")")
		return e
	case tok.isPunct("&"):
		uop := &ExprUop{
			tok:     tok,
			op:      tok.sval,
			operand: p.parsePrim(),
		}
		// when &T{}, allocate stack memory
		if strctliteral, ok := uop.operand.(*ExprStructLiteral); ok {
			// newVariable
			strctliteral.invisiblevar = p.newVariable("", &Gtype{
				typ:      G_REL,
				relation: strctliteral.strctname,
			})
		}
		return uop
	case tok.isPunct("*"):
		return &ExprUop{
			tok:     tok,
			op:      tok.sval,
			operand: p.parsePrim(),
		}
	case tok.isPunct("!"):
		return &ExprUop{
			tok:     tok,
			op:      tok.sval,
			operand: p.parsePrim(),
		}
	case tok.isPunct("-"):
		return &ExprUop{
			tok:     tok,
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
		return 11
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
	p.traceIn(__func__)
	defer p.traceOut(__func__)
	return p.parseExprInt(-1)
}

var binops = []string{
	"+", "*", "-", "==", "!=", "<", ">", "<=", ">=", "&&", "||", "/", "%",
}

func (p *parser) parseExprInt(prior int) Expr {
	p.traceIn(__func__)
	defer p.traceOut(__func__)

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
					errorft(tok, "bad lefts unary expr:%v", ast)
				}
				ast = &ExprBinop{
					tok:   tok,
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

func (p *parser) newVariable(varname identifier, gtype *Gtype) *ExprVariable {
	var variable *ExprVariable
	if p.isGlobal() {
		variable = &ExprVariable{
			tok:      p.lastToken(),
			varname:  varname,
			gtype:    gtype,
			isGlobal: p.isGlobal(),
		}
	} else {
		variable = &ExprVariable{
			tok:      p.lastToken(),
			varname:  varname,
			gtype:    gtype,
			isGlobal: p.isGlobal(),
		}
		p.localvars = append(p.localvars, variable)
	}
	return variable
}

func (p *parser) registerDynamicType(gtype *Gtype) *Gtype {
	p.allDynamicTypes = append(p.allDynamicTypes, gtype)
	return gtype
}

// https://golang.org/ref/spec#Type
func (p *parser) parseType() *Gtype {
	p.traceIn(__func__)
	defer p.traceOut(__func__)
	var gtype *Gtype
	ptok := p.peekToken()

	for {
		tok := p.readToken()
		if tok.isTypeIdent() {
			ident := tok.getIdent()
			// unresolved
			rel := &Relation{
				tok:  tok,
				pkg:  p.currentPackageName,
				name: ident,
			}
			p.tryResolve("", rel)
			gtype = &Gtype{
				typ:      G_REL,
				relation: rel,
			}
			return p.registerDynamicType(gtype)
		} else if tok.isKeyword("interface") {
			p.expect("{")
			p.expect("}")
			return gInterface
		} else if tok.isPunct("*") {
			// pointer
			gtype = &Gtype{
				typ:      G_POINTER,
				origType: p.parseType(),
			}
			return p.registerDynamicType(gtype)
		} else if tok.isKeyword("struct") {
			p.unreadToken()
			gtype = p.parseStructDef()
			return p.registerDynamicType(gtype)
		} else if tok.isKeyword("map") {
			p.unreadToken()
			gtype = p.parseMapType()
			return p.registerDynamicType(gtype)
		} else if tok.isPunct("[") {
			// array or slice
			tok := p.readToken()
			// @TODO consider "..." case in a composite literal.
			// The notation ... specifies an array length
			// equal to the maximum element index plus one.
			if tok.isPunct("]") {
				// slice
				typ := p.parseType()
				gtype = &Gtype{
					typ:         G_SLICE,
					elementType: typ,
				}
				return p.registerDynamicType(gtype)
			} else {
				// array
				p.expect("]")
				typ := p.parseType()
				gtype = &Gtype{
					typ:         G_ARRAY,
					length:      tok.getIntval(),
					elementType: typ,
				}
				return p.registerDynamicType(gtype)
			}
		} else if tok.isPunct("]") {

		} else if tok.isPunct("...") {
			// vaargs
			TBI(tok, "VAARGS is not supported yet")
		} else {
			errorft(tok, "Unkonwn token")
		}

	}
	errorft(ptok, "Unkown type")
	return nil
}

// local decl infer
func (decl *DeclVar) infer() {
	debugf("infering DeclVar")
	gtype := decl.initval.getGtype()
	assertNotNil(gtype != nil, decl.initval.token())
	decl.variable.gtype = gtype
}

func (p *parser) parseVarDecl() *DeclVar {
	p.traceIn(__func__)
	defer p.traceOut(__func__)
	ptok := p.expectKeyword("var")

	// read newName
	newName := p.expectIdent()
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

	variable := p.newVariable(newName, typ)
	r := &DeclVar{
		tok: ptok,
		pkg: p.currentPackageName,
		varname: &Relation{
			expr: variable,
			pkg:  p.currentPackageName,
		},
		variable: variable,
		initval:  initval,
	}
	if typ == nil {
		if p.isGlobal() {
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
	p.traceIn(__func__)
	defer p.traceOut(__func__)
	newName := p.expectIdent()

	// Type or "=" or ";"
	var val Expr
	var gtype *Gtype
	if !p.peekToken().isPunct("=") && !p.peekToken().isPunct(";") {
		// expect Type
		gtype = p.parseType()
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
		gtype:     gtype,
	}

	p.currentScope.setConst(newName, variable)
	return variable
}

func (p *parser) parseConstDecl() *DeclConst {
	p.traceIn(__func__)
	defer p.traceOut(__func__)
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
	p.traceIn(__func__)
	defer p.traceOut(__func__)
	var r []identifier
	for {
		tok := p.readToken()
		if tok.isTypeIdent() {
			r = append(r, tok.getIdent())
		} else if len(r) == 0 {
			// at least one ident is needed
			errorft(tok, "Ident expected")
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

func (p *parser) enterNewScope(name string) {
	p.currentScope = newScope(p.currentScope, name)
}

func (p *parser) exitScope() {
	p.currentScope = p.currentScope.outer
}

func (p *parser) exitForBlock() {
	p.currentForStmt = p.currentForStmt.outer
}

func (clause *ForRangeClause) infer() {
	debugf("infering ForRangeClause")
	collectionType := clause.rangeexpr.getGtype()
	//debugf("collectionType = %s", collectionType)
	indexvar, ok := clause.indexvar.expr.(*ExprVariable)
	assert(ok, nil, "ok")

	var indexType *Gtype
	switch collectionType.typ {
	case G_ARRAY, G_SLICE:
		indexType = gInt
	case G_MAP:
		indexType = collectionType.mapKey
	default:
		// @TODO consider map etc.
		TBI(clause.tok, "unable to handle %s", collectionType)
	}
	indexvar.gtype = indexType

	if clause.valuevar != nil {
		valuevar, ok := clause.valuevar.expr.(*ExprVariable)
		assert(ok, nil, "ok")

		var elementType *Gtype
		if collectionType.typ == G_ARRAY {
			elementType = collectionType.elementType
		} else if collectionType.typ == G_SLICE {
			elementType = collectionType.elementType
		} else if collectionType.typ == G_MAP {
			elementType = collectionType.mapValue
		} else {
			errorft(clause.token(), "internal error")
		}
		//debugf("for i, v %s := rannge %v", elementType, collectionType)
		valuevar.gtype = elementType
	}
}

// https://golang.org/ref/spec#For_statements
func (p *parser) parseForStmt() *StmtFor {
	p.traceIn(__func__)
	defer p.traceOut(__func__)
	ptok := p.expectKeyword("for")

	var r = &StmtFor{
		tok:   ptok,
		outer: p.currentForStmt,
	}
	p.currentForStmt = r
	p.enterNewScope("for")
	var cond Expr
	if p.peekToken().isPunct("{") {
		// inifinit loop : for { ___ }
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
				return p.parseForRange(lefts, false)
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

				r := p.parseForRange(lefts, true)
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
	p.exitScope()
	p.exitForBlock()
	return r
}

func (p *parser) parseForRange(exprs []Expr, infer bool) *StmtFor {
	p.traceIn(__func__)
	defer p.traceOut(__func__)
	tokRange := p.expectKeyword("range")

	if len(exprs) > 2 {
		errorft(tokRange, "range values should be 1 or 2")
	}
	indexvar, ok := exprs[0].(*Relation)
	if !ok {
		errorft(tokRange, " rng.lefts[0]. is not relation")
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
		tok:   tokRange,
		outer: p.currentForStmt,
		rng: &ForRangeClause{
			tok:                 tokRange,
			invisibleMapCounter: p.newVariable("", gInt),
			indexvar:            indexvar,
			valuevar:            valuevar,
			rangeexpr:           rangeExpr,
		},
	}
	p.currentForStmt = r
	if infer {
		p.localuninferred = append(p.localuninferred, r.rng)
	}
	r.block = p.parseCompoundStmt()
	p.exitScope()
	p.exitForBlock()
	return r
}

func (p *parser) parseIfStmt() *StmtIf {
	p.traceIn(__func__)
	defer p.traceOut(__func__)
	ptok := p.expectKeyword("if")

	var r = &StmtIf{
		tok: ptok,
	}
	p.enterNewScope("if")
	p.requireBlock = true
	stmt := p.parseStmt()
	if p.peekToken().isPunct(";") {
		p.skip()
		r.simplestmt = stmt
		r.cond = p.parseExpr()
	} else {
		es, ok := stmt.(*StmtExpr)
		if !ok {
			errorft(stmt.token(), "internal error")
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
			errorft(tok2, "Unexpected token")
		}
	}
	p.exitScope()
	return r
}

func (p *parser) parseReturnStmt() *StmtReturn {
	p.traceIn(__func__)
	defer p.traceOut(__func__)
	ptok := p.expectKeyword("return")

	exprs := p.parseExpressionList(nil)
	// workaround for {nil}
	if len(exprs) == 1 && exprs[0] == nil {
		exprs = nil
	}
	return &StmtReturn{
		tok:               ptok,
		exprs:             exprs,
		rettypes:          p.currentFunc.rettypes,
		labelDeferHandler: p.currentFunc.labelDeferHandler,
	}
}

func (p *parser) parseExpressionList(first Expr) []Expr {
	p.traceIn(__func__)
	defer p.traceOut(__func__)
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
	p.traceIn(__func__)
	defer p.traceOut(__func__)
	ptok := p.lastToken()

	rights := p.parseExpressionList(nil)
	p.assertNotNil(rights[0])
	return &StmtAssignment{
		tok:    ptok,
		lefts:  lefts,
		rights: rights,
	}
}

func (p *parser) parseAssignmentOperation(left Expr, assignop string) *StmtAssignment {
	p.traceIn(__func__)
	defer p.traceOut(__func__)
	ptok := p.lastToken()

	var op string
	switch assignop {
	case "+=":
		op = "+"
	case "-=":
		op = "-"
	case "*=":
		op = "*"
	default:
		errorft(ptok, "internal error")
	}
	rights := p.parseExpressionList(nil)
	p.assert(len(rights) == 1, "num of rights is 1")
	binop := &ExprBinop{
		tok:   ptok,
		op:    op,
		left:  left,
		right: rights[0],
	}
	var right Expr = binop // FIXME: this is a workaround
	s := &StmtAssignment{
		tok:    ptok,
		lefts:  []Expr{left},
		rights: []Expr{right},
	}
	// dumpInterface(s.rights[0])
	return s
}

func (p *parser) shortVarDecl(e Expr) {
	rel := e.(*Relation) // a brand new rel
	assert(p.isGlobal() == false, e.token(), "should not be in global scope")
	variable := p.newVariable(rel.name, nil)
	p.currentScope.setVar(rel.name, variable)
	rel.expr = variable
}

func (p *parser) parseShortAssignment(lefts []Expr) *StmtShortVarDecl {
	p.traceIn(__func__)
	defer p.traceOut(__func__)
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

// https://golang.org/ref/spec#Switch_statements
func (p *parser) parseSwitchStmt() Stmt {
	p.traceIn(__func__)
	defer p.traceOut(__func__)
	ptok := p.expectKeyword("switch")

	var cond Expr
	if p.peekToken().isPunct("{") {

	} else {
		p.requireBlock = true
		cond = p.parseExpr()
		p.requireBlock = false
	}

	_, isTypeSwitch := cond.(*ExprTypeSwitchGuard)

	p.expect("{")
	r := &StmtSwitch{
		isTypeSwitch: isTypeSwitch,
		tok:          ptok,
		cond:         cond,
	}

	for {
		tok := p.peekToken()
		if tok.isKeyword("case") {
			p.skip()
			var exprs []Expr
			var gtypes []*Gtype
			if r.isTypeSwitch {
				gtype := p.parseType()
				gtypes = append(gtypes, gtype)
				for {
					tok := p.peekToken()
					if tok.isPunct(",") {
						p.skip()
						gtype := p.parseType()
						gtypes = append(gtypes, gtype)
					} else if tok.isPunct(":") {
						break
					}
				}
			} else {
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
			}
			ptok := p.expect(":")
			p.inCase++
			compound := p.parseCompoundStmt()
			casestmt := &ExprCaseClause{
				tok:      ptok,
				exprs:    exprs,
				gtypes:   gtypes,
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
			errorft(tok, "internal error")
		}
	}

	return r
}

func (p *parser) parseDeferStmt() *StmtDefer {
	p.traceIn(__func__)
	defer p.traceOut(__func__)
	ptok := p.expectKeyword("defer")

	callExpr := p.parsePrim()
	stmtDefer := &StmtDefer{
		tok:  ptok,
		expr: callExpr,
	}
	p.currentFunc.stmtDefer = stmtDefer
	return stmtDefer
}

// this is in function scope
func (p *parser) parseStmt() Stmt {
	p.traceIn(__func__)
	defer p.traceOut(__func__)

	tok := p.peekToken()
	if tok.isKeyword("var") {
		return p.parseVarDecl()
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
		ptok := p.expectKeyword("continue")
		return &StmtContinue{
			tok:     ptok,
			stmtFor: p.currentForStmt,
		}
	} else if tok.isKeyword("break") {
		ptok := p.expectKeyword("break")
		return &StmtBreak{
			tok:     ptok,
			stmtFor: p.currentForStmt,
		}
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
			TBI(tok3, "")
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
			tok:     tok2,
			operand: expr1,
		}
	} else if tok2.isPunct("--") {
		p.skip()
		return &StmtDec{
			tok:     tok2,
			operand: expr1,
		}
	} else {
		return &StmtExpr{
			tok:  tok2,
			expr: expr1,
		}
	}
	return nil
}

func (p *parser) parseCompoundStmt() *StmtSatementList {
	p.traceIn(__func__)
	defer p.traceOut(__func__)

	r := &StmtSatementList{
		tok: p.lastToken(),
	}
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
	p.traceIn(__func__)
	defer p.traceOut(__func__)

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
				sliceType := &Gtype{
					typ:         G_SLICE,
					elementType: gtype,
				}
				variable := &ExprVariable{
					tok:        tok,
					varname:    pname,
					gtype:      sliceType,
					isVariadic: true,
				}
				params = append(params, variable)
				p.currentScope.setVar(pname, variable)
				p.expect(")")
				break
			}
			ptype := p.parseType()
			// assureType(tok.sval)
			variable := &ExprVariable{
				tok:     tok,
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
				errorft(tok, "Invalid token")
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
				errorft(next, "invalid token")
			}
		}

	} else {
		rettypes = []*Gtype{p.parseType()}
	}

	return fname, params, isVariadic, rettypes
}

func (p *parser) parseFuncDef() *DeclFunc {
	p.traceIn(__func__)
	defer p.traceOut(__func__)
	ptok := p.expectKeyword("func")

	p.localvars = nil
	assert(len(p.localvars) == 0, ptok, "localvars should be zero")
	var isMethod bool
	p.enterNewScope("func")

	var receiver *ExprVariable

	if p.peekToken().isPunct("(") {
		isMethod = true
		p.expect("(")
		// method definition
		tok := p.readToken()
		pname := tok.getIdent()
		ptype := p.parseType()
		receiver = &ExprVariable{
			tok:     tok,
			varname: pname,
			gtype:   ptype,
		}
		p.currentScope.setVar(pname, receiver)
		p.expect(")")
	}

	fname, params, isVariadic, rettypes := p.parseFuncSignature()

	ptok2 := p.expect("{")

	r := &DeclFunc{
		tok:        ptok,
		pkg:        p.currentPackageName,
		receiver:   receiver,
		fname:      fname,
		rettypes:   rettypes,
		params:     params,
		isVariadic: isVariadic,
	}

	ref := &ExprFuncRef{
		tok:     ptok2,
		funcdef: r,
	}

	if isMethod {
		var typeToBelong *Gtype
		if receiver.gtype.typ == G_POINTER {
			typeToBelong = receiver.gtype.origType
		} else {
			typeToBelong = receiver.gtype
		}

		p.assert(typeToBelong.typ == G_REL, "methods must belong to a named type")
		var methods methods
		var ok bool
		methods, ok = p.methods[typeToBelong.relation.name]
		if !ok {
			methods = map[identifier]*ExprFuncRef{}
			p.methods[typeToBelong.relation.name] = methods
		}
		methods[fname] = ref
	} else {
		p.packageBlockScope.setFunc(fname, ref)
	}

	// every function has a defer_handler
	r.labelDeferHandler = makeLabel() + "_defer_handler"
	p.currentFunc = r
	body := p.parseCompoundStmt()
	r.body = body
	r.localvars = p.localvars

	p.localvars = nil
	p.exitScope()
	return r
}

func (p *parser) parseImport() *ImportDecl {
	p.traceIn(__func__)
	defer p.traceOut(__func__)
	tokImport := p.expectKeyword("import")

	tok := p.readToken()
	var specs []*ImportSpec
	if tok.isPunct("(") {
		for {
			tok := p.readToken()
			if tok.isTypeString() {
				specs = append(specs, &ImportSpec{
					tok:  tok,
					path: tok.sval,
				})
				p.expect(";")
			} else if tok.isPunct(")") {
				break
			} else {
				errorft(tok, "invalid import path")
			}
		}
	} else {
		if !tok.isTypeString() {
			errorft(tok, "import expects package name")
		}
		specs = []*ImportSpec{&ImportSpec{
			tok:  tok,
			path: tok.sval,
		},
		}
	}
	p.expect(";")
	return &ImportDecl{
		tok:   tokImport,
		specs: specs,
	}
}

func (p *parser) parsePackageClause() *PackageClause {
	p.traceIn(__func__)
	defer p.traceOut(__func__)
	tokPkg := p.expectKeyword("package")

	name := p.expectIdent()
	p.expect(";")
	return &PackageClause{
		tok:  tokPkg,
		name: name,
	}
}

func (p *parser) parseImportDecls() []*ImportDecl {
	p.traceIn(__func__)
	defer p.traceOut(__func__)
	var r []*ImportDecl
	for p.peekToken().isKeyword("import") {
		r = append(r, p.parseImport())
	}
	return r
}

const MaxAlign = 16

func (p *parser) parseStructDef() *Gtype {
	p.traceIn(__func__)
	defer p.traceOut(__func__)
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
		fieldtype := gtype
		//fieldtype.origType = gtype
		fieldtype.fieldname = fieldname
		fieldtype.offset = undefinedSize // will be calculated later
		fields = append(fields, fieldtype)
		p.expect(";")
	}
	// calc offset
	p.expect(";")
	return &Gtype{
		typ:    G_STRUCT,
		size:   undefinedSize, // will be calculated later
		fields: fields,
	}
}

func (p *parser) parseInterfaceDef(newName identifier) *DeclType {
	p.traceIn(__func__)
	defer p.traceOut(__func__)
	p.expectKeyword("interface")

	p.expect("{")
	var methods map[identifier]*signature = map[identifier]*signature{}

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
		methods[fname] = method
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
		relbody := p.currentScope.get(rel.name)
		if relbody != nil {
			if relbody.gtype != nil {
				rel.gtype = relbody.gtype
			} else if relbody.expr != nil {
				rel.expr = relbody.expr
			} else {
				errorft(rel.token(), "Bad type relbody %v", relbody)
			}
		} else {
			if rel.name != "_" {
				p.unresolvedRelations = append(p.unresolvedRelations, rel)
			}
		}
	} else {
		// foreign package
		relbody := p.scopes[pkg].get(rel.name)
		if relbody == nil {
			errorft(rel.token(), "name %s is not found in %s package", rel.name, pkg)
		}

		if relbody.gtype != nil {
			rel.gtype = relbody.gtype
		} else if relbody.expr != nil {
			rel.expr = relbody.expr
		} else {
			errorft(rel.token(), "Bad type relbody %v", relbody)
		}
	}
}

var typeId int

func (p *parser) parseTypeDecl() *DeclType {
	p.traceIn(__func__)
	defer p.traceOut(__func__)
	ptok := p.expectKeyword("type")

	newName := p.expectIdent()
	if p.peekToken().isKeyword("interface") {
		return p.parseInterfaceDef(newName)
	}

	gtype := p.parseType()
	r := &DeclType{
		tok:   ptok,
		name:  newName,
		gtype: gtype,
	}

	p.allNamedTypes = append(p.allNamedTypes, r)
	p.currentScope.setGtype(newName, gtype)
	return r
}

// https://golang.org/ref/spec#TopLevelDecl
// TopLevelDecl  = Declaration | FunctionDecl | MethodDecl .
func (p *parser) parseTopLevelDecl(nextToken *Token) *TopLevelDecl {
	p.traceIn(__func__)
	defer p.traceOut(__func__)

	if !nextToken.isTypeKeyword() {
		errorft(nextToken, "invalid token")
	}

	switch nextToken.sval {
	case "func":
		funcdecl := p.parseFuncDef()
		return &TopLevelDecl{funcdecl: funcdecl}
	case "var":
		vardecl := p.parseVarDecl()
		return &TopLevelDecl{vardecl: vardecl}
	case "const":
		constdecl := p.parseConstDecl()
		return &TopLevelDecl{constdecl: constdecl}
	case "type":
		typedecl := p.parseTypeDecl()
		return &TopLevelDecl{typedecl: typedecl}
	}

	TBI(nextToken, "")
	return nil
}

func (p *parser) parseTopLevelDecls() []*TopLevelDecl {
	p.traceIn(__func__)
	defer p.traceOut(__func__)

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

func (p *parser) isGlobal() bool {
	return p.currentScope == p.packageBlockScope
}

// https://golang.org/ref/spec#Source_file_organization
// Each source file consists of
// a package clause defining the package to which it belongs,
// followed by a possibly empty set of import declarations that declare packages whose contents it wishes to use,
// followed by a possibly empty set of declarations of functions, types, variables, and constants.
func (p *parser) parseSourceFile(bs *ByteStream, packageBlockScope *scope, importOnly bool) *SourceFile {

	p.clearLocalState()

	// initialize parser's status per file
	p.tokenStream = NewTokenStream(bs)
	p.packageBlockScope = packageBlockScope
	p.currentScope = packageBlockScope
	p.importedNames = map[identifier]bool{}

	packageClause := p.parsePackageClause()
	importDecls := p.parseImportDecls()

	// regsiter imported names
	for _, importdecl := range importDecls {
		for _, spec := range importdecl.specs {
			pkgName := getBaseNameFromImport(spec.path)
			p.importedNames[identifier(pkgName)] = true
		}
	}

	if importOnly {
		return &SourceFile{
			tok:           packageClause.tok,
			packageClause: packageClause,
			importDecls:   importDecls,
		}
	}

	// @TODO import external decls here

	topLevelDecls := p.parseTopLevelDecls()

	return &SourceFile{
		tok:           packageClause.tok,
		name:          bs.filename,
		packageClause: packageClause,
		importDecls:   importDecls,
		topLevelDecls: topLevelDecls,
	}
}

func (ast *StmtShortVarDecl) infer() {
	debugf("infering StmtShortVarDecl")
	var rightTypes []*Gtype
	for _, rightExpr := range ast.rights {
		switch rightExpr.(type) {
		case *ExprFuncallOrConversion:
			fcallOrConversion := rightExpr.(*ExprFuncallOrConversion)
			if fcallOrConversion.rel.gtype != nil {
				// Conversion
				rightTypes = append(rightTypes, fcallOrConversion.rel.gtype)
			} else {
				fcall := fcallOrConversion
				funcdef := fcall.getFuncDef()
				if funcdef == nil {
					errorft(fcall.token(), "funcdef of %s is not found", fcall.fname)
				}
				if funcdef == builtinLen {
					rightTypes = append(rightTypes, gInt)
				} else {
					for _, gtype := range fcall.getFuncDef().rettypes {
						rightTypes = append(rightTypes, gtype)
					}
				}
			}
		case *ExprMethodcall:
			fcall := rightExpr.(*ExprMethodcall)
			rettypes := fcall.getRettypes()
			for _, gtype := range rettypes {
				rightTypes = append(rightTypes, gtype)
			}
		case *ExprTypeAssertion:
			assertion := rightExpr.(*ExprTypeAssertion)
			rightTypes = append(rightTypes, assertion.gtype)
			rightTypes = append(rightTypes, gBool)
		case *ExprIndex:
			e := rightExpr.(*ExprIndex)
			gtype := e.getGtype()
			assertNotNil(gtype != nil, e.tok)
			rightTypes = append(rightTypes, gtype)
			//debugf("rightExpr.gtype=%s", gtype)
			secondGtype := rightExpr.(*ExprIndex).getSecondGtype()
			if secondGtype != nil {
				rightTypes = append(rightTypes, secondGtype)
			}
		default:
			if rightExpr == nil {
				errorft(ast.token(), "rightExpr is nil")
			}
			gtype := rightExpr.getGtype()
			if gtype == nil {
				errorft(ast.token(), "rightExpr %T gtype is nil", rightExpr)
			}
			//debugf("infered type %s", gtype)
			rightTypes = append(rightTypes, gtype)
		}
	}

	if len(ast.lefts) > len(rightTypes) {
		// @TODO this check is too loose.
		errorft(ast.tok, "number of lhs and rhs does not match")
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
		debugf("resolving %s ...", rel.name)
		p.tryResolve("", rel)
	}

	debugf("resolving methods ...")
	p.resolveMethods()
	debugf("inferring types ...")
	p.inferTypes()
}

// copy methods from p.nameTypes to gtype.methods of each type
func (p *parser) resolveMethods() {
	for typeName, methods := range p.methods {
		gtype := p.packageBlockScope.getGtype(typeName)
		if gtype == nil {
			debugf("%#v", p.packageBlockScope.idents)
			errorf("typaneme %s is not found in the package scope %s", typeName, p.currentPackageName)
		}
		gtype.methods = methods
	}
}

//  infer recursively all the types of global variables
func (variable *ExprVariable) infer() {
	debugf("infering ExprVariable")
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
		errorft(e.token(), "unexpected type %T", e)
	}
	vr, ok := rel.expr.(*ExprVariable)
	vr.infer() // recursive call
	variable.gtype = e.getGtype()
	//debugf("infered type=%s", variable.gtype)
}

func (p *parser) inferTypes() {
	debugf("infering globals")
	for _, variable := range p.globaluninferred {
		variable.infer()
	}
	debugf("infering locals")
	for _, ast := range p.localuninferred {
		ast.infer()
	}
}
