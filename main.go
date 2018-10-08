package main

import (
	"os"
)

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
		index:  0,
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
