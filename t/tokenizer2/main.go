package main


var GENERATION int = 2

var debugMode = true
var debugToken = false

func f1() {
	filename := "t/data/gen.go.txt"
	bs := NewByteStreamFromFile(filename)
	tokens := Tokenize(bs)
	expectedLen := 17977
	if len(tokens) == expectedLen {
		fmtPrintf(S("1\n"))
	} else {
		panic(S("ERROR"))
	}
	/*
		for _, tok := range tokens {
			fmtPrintf(S("%s:%s\n", string(tok.typ)), tok.sval)
		}
	*/
}

func main() {
	f1()
}
