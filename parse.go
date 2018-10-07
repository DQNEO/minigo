package main

import "strconv"

var tokens []*Token
var tokenIndex int


func readToken() *Token {

	if tokenIndex <= len(tokens)-1 {
		r := tokens[tokenIndex]
		tokenIndex++
		return r
	}
	return nil
}

func parseUnaryExpr() *Ast {
	tok := readToken()
	if tok.typ == "space" {
		tok = readToken()
	}
	ival, _ := strconv.Atoi(tok.sval)
	return &Ast{
		typ: "uop",
		operand: &Ast{
			typ:  "int",
			ival: ival,
		},
	}
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
		} else {
			debugToken(tok)
			errorf("unknown token=%v\n", tok.sval)
		}
	}

	return ast
}


