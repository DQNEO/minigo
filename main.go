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
	astFile := parse(t)

	if debugAst {
		debugPrint("==== Dump Ast Start ===")
		for _, toplevel := range astFile.asts {
			dumpAst(toplevel)
		}
		debugPrint("==== Dump Ast End ===")
	}

	// generate
	generate(astFile)
}
