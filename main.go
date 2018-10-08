package main

import "fmt"
import (
	"os"
)

var debugMode = false

func debugPrint(s string) {
	fmt.Fprintf(os.Stderr, "# %s\n", s)
}

func debugPrintVar(name string, v interface{}) {
	debugPrint(fmt.Sprintf("%s=%v", name, v))
}

func debugToken(tok *Token) {
	debugPrint(fmt.Sprintf("tok: type=%-8s, sval=\"%s\"", tok.typ, tok.sval))
}

func debugAst(name string, ast *Ast) {
	switch ast.typ {
	case "funcall":
		debugPrintVar("funcall", ast)
		for _, arg := range ast.args {
			debugPrintVar("arg", arg)
		}
	case  "int" :
		debugPrintVar(name, fmt.Sprintf("%d", ast.ival))
	case "uop" :
		debugAst(name, ast.operand)
	case "binop" :
		debugPrintVar("ast.binop", ast.typ)
		debugAst("left", ast.left)
		debugAst("right", ast.right)
	default:
		debugPrintVar(name, ast)
	}
}

func errorf(format string, v ...interface{}) {
	/*
	currentTokenIndex := ts.index - 1
	fmt.Printf("%v %v %v\n",
		ts.getToken(currentTokenIndex-2), ts.getToken(currentTokenIndex-1), ts.getToken(currentTokenIndex))
	*/
	s := fmt.Sprintf(format, v...)
	panic(s)
}

func assert(cond bool, msg string) {
	if !cond {
		panic(fmt.Sprintf("assertion failed: %s", msg))
	}
}


func main() {
	debugMode = false

	var sourceFile string
	if len(os.Args) > 1 {
		sourceFile = os.Args[1] + ".go"
	} else {
		sourceFile = "/dev/stdin"
	}

	// tokenize
	tokens := tokenizeFromFile(sourceFile)
	assert(len(tokens) > 0, "tokens should have length")

	if debugMode {
		renderTokens(tokens)
	}

	t := &TokenStream{
		tokens: tokens,
		index: 0,
	}
	// parse
	asts := parse(t)

	if debugMode {
		debugPrint("==== Dump Ast ===")
		debugAst("root", asts[1])
	}

	// generate
	generate(asts)
}
