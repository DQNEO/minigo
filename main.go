package main

import (
	"os"
)

var debugAst = false
var debugToken = false

func parseOpts(args []string) {
	for _,opt := range args {
		if opt == "-t" {
			debugToken = true
		}
		if opt == "-a" {
			debugAst = true
		}
	}
}

func main() {

	var sourceFile string
	sourceFile = "/dev/stdin"
	if len(os.Args) > 1 {
		parseOpts(os.Args[1:len(os.Args)])
	}

	// tokenize
	tokens := tokenizeFromFile(sourceFile)
	assert(len(tokens) > 0, "tokens should have length")

	if debugToken {
		renderTokens(tokens)
	}

	t := &TokenStream{
		tokens: tokens,
		index:  0,
	}
	// parse
	asts := parse(t)

	if debugAst {
		debugPrint("==== Dump Ast ===")
		dumpAst("root", asts[1])
	}

	// generate
	generate(asts)
}
