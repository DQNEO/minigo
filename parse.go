package main

import "strconv"
import "fmt"

type Ast struct {
	typ     string
	// int
	ival    int
	// unary
	operand *Ast
	// binop
	op string
	left    *Ast
	right   *Ast
	// string
	sval   string
	slabel string
	// compound
	stmts []*Ast
	// funcall
	fname string
	args []*Ast
	// funcdef
	body *Ast
	// package
	pkgname string
	// imports
	packages []string
}

type TokenStream struct {
	tokens []*Token
	index int
}

var ts *TokenStream

func (ts *TokenStream) readToken() *Token{
	if ts.index <= len(ts.tokens)-1 {
		r := ts.tokens[ts.index]
		ts.index++
		return r
	}
	return nil
}

func (ts *TokenStream) unreadToken() {
	ts.index--
}

func (stream *TokenStream) getToken(i int) interface{} {
	return ts.tokens[i]
}

func readToken() *Token {
	return ts.readToken()
}

func unreadToken() {
	ts.unreadToken()
}

func (tok *Token) isPunct(punct string) bool {
	return tok != nil && tok.typ == "punct"  && tok.sval == punct
}

func skipSpaceToken() {
	for {
		tok := readToken()
		if tok == nil {
			return
		}
		if tok.typ == "space" {
			continue
		} else {
			unreadToken()
			return
		}

	}
}

func expectPunct(punct string) {
	tok := readToken()
	if tok.typ != "punct" {
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
		}
	}
}

func parseIdentOrFuncall(name string) *Ast {
	tok := readToken()
	if tok.isPunct("(") {
		args := readFuncallArgs()
		return &Ast{
			typ: "funcall",
			fname : name,
			args:args,
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
		slabel: fmt.Sprintf("L%d",stringIndex),
	}
	stringIndex++
	strings = append(strings, ast)
	return ast
}

func parseUnaryExpr() *Ast {
	tok := readToken()
	if tok == nil {
		return nil
	}
	if tok.typ == "space" {
		tok = readToken()
	}

	switch tok.typ {
	case "string":
		return newAstString(tok.sval)
	case "ident":
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

func parseExpr() *Ast {
	ast := parseUnaryExpr()
	for {
		tok := readToken()
		if tok == nil || tok.typ == "newline" {
			return ast
		}
		if tok.typ == "space" {
			continue
		}
		if tok.typ != "punct" {
			return ast
		}
		if tok.sval == "+" || tok.sval == "*" || tok.sval == "-" {
			right := parseUnaryExpr()
			return &Ast{
				typ:   "binop",
				op: tok.sval,
				left:  ast,
				right: right,
			}
		} else if tok.sval == "," || tok.sval == ")"{ // end of funcall argument
			unreadToken()
			return ast
		} else {
			debugToken(tok)
			errorf("unable to handle token=\"%s\"\n", tok.sval)
		}
	}

	return ast
}

func parseStmt() *Ast {
	return parseExpr()
}

func parseCompoundStmt() []*Ast {
	var r []*Ast
	for {
		tok := readToken()
		if tok.isPunct("}") {
			return r
		}
		if tok.typ == "newline" {
			continue
		}
		unreadToken()
		stmt := parseStmt()
		if stmt == nil {
			errorf("internal error")
		}
		r = append(r, stmt)
	}
	return nil
}

func parseFuncDef() *Ast {
	skipSpaceToken()
	fname := readToken()
	if fname.typ != "ident" {
		errorf("identifer expected, but got %v", fname)
	}
	expectPunct("(")
	expectPunct(")")
	skipSpaceToken()
	// expect Type
	expectPunct("{")
	stmts := parseCompoundStmt()

	return &Ast{
		typ: "funcdef",
		fname: fname.sval,
		body : &Ast{
			typ:"compound",
			stmts:stmts,
		},
	}
}

func parseTopLevels() []*Ast {
	var r []*Ast
	for {
		tok := readToken()
		if tok == nil {
			return r
		}
		if tok.typ == "newline" {
			continue
		}
		if tok.typ == "ident" && tok.sval == "package" {
			skipSpaceToken()
			tok = readToken()
			assert(tok.typ == "ident", "expect ident")
			ast := &Ast{
				typ: "package",
				pkgname: tok.sval,
			}
			readToken() // expect newline
			r = append(r, ast)
			continue
		} else if tok.typ == "ident"  && tok.sval == "func" {
			ast := parseFuncDef()
			r = append(r, ast)
			continue
		} else {
			errorf("unknown token %v", tok)
		}
	}
	return r
}

func parse(t *TokenStream) []*Ast {
	ts = t
	return parseTopLevels()
}
