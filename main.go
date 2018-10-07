package main

import "fmt"
import (
	"io/ioutil"
	"strings"
	"regexp"
	"strconv"
)

type Token struct {
	typ  string
	sval string
}

type Ast struct {
	typ     string
	ival    int
	operand *Ast
}

var tokens []*Token
var tokenIndex int

func readFile(filename string) string {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func tokinize(s string) []*Token {
	var r []*Token
	trimed := strings.Trim(s, "\n")
	chars := strings.Split(trimed, " ")
	var regexNumber = regexp.MustCompile(`^[0-9]+$`)
	for _, char := range chars {
		debugPrint("char", char)
		var tok *Token
		if regexNumber.MatchString(char) {
			tok = &Token{typ: "number", sval: strings.Trim(char, " \n")}
		}

		r = append(r, tok)
	}

	return r
}

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
	return ast
}

func generate(ast *Ast) {
	fmt.Println("\t.globl	main")
	fmt.Println("main:")
	emitAst(ast)
	fmt.Println("\tret")
}

func emitAst(ast *Ast) {
	if ast.typ == "uop" {
		emitUop(ast)
	}
}

func emitUop(ast *Ast) {
	fmt.Printf("\tmovl	$%d, %%eax\n", ast.operand.ival)
}

func debugPrint(name string, v interface{}) {
	fmt.Printf("# %s=%v\n", name, v)
}

func debugTokens(tokens []*Token) {
	for _, tok := range tokens {
		debugPrint("tok", tok)
	}
}

func debugAst(ast *Ast) {
	if ast.typ == "uop" {
		debugPrint("ast.uop", ast.operand)
	}
}

func main() {
	s := readFile("/dev/stdin")
	tokens = tokinize(s)
	tokenIndex = 0
	debugTokens(tokens)
	ast := parseExpr()
	debugAst(ast)
	generate(ast)
}
