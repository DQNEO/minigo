package main

import "os"

func main() {

	var sourceFile string
	sourceFile = "/dev/stdin"

	if len(os.Args) > 1 {
		opt := os.Args[1]
		if opt == "-v" {
			debugMode = true
		}
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
