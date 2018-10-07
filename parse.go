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

func unreadToken() {
	tokenIndex--
}

func parseIdentOrFuncall(name string) *Ast {
	tok := readToken()
	if tok != nil && tok.typ == "punct" && tok.sval == "(" {
		arg1 := parseExpr()
		readToken() // expect ","
		arg2 := parseExpr()
		readToken() // expect "):
		return &Ast{
			typ: "funcall",
			fname : name,
			args:[]*Ast{arg1, arg2},
		}
	}

	errorf("TBD")
	return nil
}

func parseUnaryExpr() *Ast {
	tok := readToken()
	if tok.typ == "space" {
		tok = readToken()
	}

	switch tok.typ {
	case "string":
		return &Ast{
			typ: "string",
			sval: tok.sval,
			label:"L0",
		}
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


