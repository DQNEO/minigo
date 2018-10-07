package main

import "fmt"
import (
	"os"
)

var debugMode = false

type Token struct {
	typ   string
	sval  string
}

type Ast struct {
	typ     string
	ival    int
	operand *Ast
	// binop
	op string
	left    *Ast
	right   *Ast
	// string
	label string
	// funcall
	fname string
	args []*Ast
}

func debugPrint(s string) {
	fmt.Fprintf(os.Stderr, "# %s\n", s)
}

func debugPrintVar(name string, v interface{}) {
	debugPrint(fmt.Sprintf("%s=%v", name, v))
}

func debugToken(tok *Token) {
	debugPrint(fmt.Sprintf("tok: type= %7s, sval=\"%s\"", tok.typ, tok.sval))
}

func debugAst(name string, ast *Ast) {
	if ast.typ == "int" {
		debugPrintVar(name, fmt.Sprintf("%d", ast.ival))
	} else if ast.typ == "uop" {
		debugAst(name, ast.operand)
	} else if ast.typ == "binop" {
		debugPrintVar("ast.binop", ast.typ)
		debugAst("left", ast.left)
		debugAst("right", ast.right)
	}
}

func errorf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v)
	panic(s)
}

func assert(cond bool, msg string) {
	if !cond {
		panic(fmt.Sprintf("assertion failed: %s", msg))
	}
}


func main() {
	debugMode = true

	var sourceFile string
	if len(os.Args) > 1 {
		sourceFile = os.Args[1] + ".go"
	} else {
		sourceFile = "/dev/stdin"
	}

	// tokenize
	tokenizeFromFile(sourceFile)
	assert(len(tokens) > 0, "tokens should have length")

	if debugMode {
		renderTokens(tokens)
	}

	// parse
	tokenIndex = 0
	expr := parseExpr()

	if debugMode {
		debugPrint("==== Dump Ast ===")
		debugAst("root", expr)
	}

	// generate
	generate(expr)
}
