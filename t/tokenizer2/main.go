package main

var GENERATION int = 2

var debugMode = true
var debugToken = true

func f1() {
	filename := "t/data/gen.go.txt"
	s := readFile(filename)
	_bs := ByteStream{
		filename:  filename,
		source:    s,
		nextIndex: 0,
		line:      1,
		column:    0,
	}
	bs = &_bs
	tokens := tokenize(bs)
	expectedLen := 17977
	if len(tokens) == expectedLen {
		println("1")
	} else {
		panic("ERROR")
	}
	/*
		for _, tok := range tokens {
			fmt.Printf("%s:%s\n", string(tok.typ), tok.sval)
		}
	*/
}

func main() {
	f1()
}
