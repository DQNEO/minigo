package main

import "strconv"
import "fmt"

var tokens []*Token

type TokenStream struct {
	tokens []*Token
	index int
}

var ts *TokenStream

func (ts *TokenStream) readToken() *Token{
	if ts.index <= len(tokens)-1 {
		r := tokens[ts.index]
		ts.index++
		return r
	}
	return nil
}

func (ts *TokenStream) unreadToken() {
	ts.index--
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

func expectPunct(punct string) {
	tok := readToken()
	if tok.isPunct(punct) {
		errorf("punct %s expected but got %v", punct, tok)
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
			debugAst("right", right)
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


