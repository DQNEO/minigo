package main

import (
	"os"
)

func debugf(format bytes, v ...interface{}) {
	if !debugMode {
		return
	}

	var indents []byte

	for i := 0; i < debugNest; i++ {
		indents = append(indents, ' ')
		indents = append(indents, ' ')
	}
	os.Stderr.Write(indents)
	s2 := Sprintf(string(format), v...)
	var b []byte = []byte(s2)
	b = append(b, '\n')
	os.Stderr.Write(b)
}

var debugNest int

// States "To Be Implemented"
func TBI(tok *Token, format string, v ...interface{}) {
	errorft(tok, concat(S("(To Be Implemented) "), bytes(format)), v...)
}

// errorf with a position token
func errorft(tok *Token, format bytes, v ...interface{}) {

	var tokString bytes
	if tok != nil {
		tokString = tok.String()
	}
	gs := concat3(format,S("\n "), tokString)
	errorf(string(gs), v...)
}

func errorf(format string, v ...interface{}) {
	s := Sprintf(format, v...)
	os.Stderr.Write(concat(s, S("\n")))
	panic(S(""))
}

func assert(cond bool, tok *Token, format bytes, v ...interface{}) {
	if !cond {
		s := Sprintf(string(format), v...)
		var toks bytes = S("")
		if tok != nil {
			toks = tok.String()
		}
		msg := concat3(S("assertion failed: "), s,  toks)
		os.Stderr.Write(msg)
		os.Stderr.Write(S("\n"))
		panic(S(""))
	}
}

func assertNotReached(tok *Token) {
	msg := concat(S("assertNotReached "), tok.String())
	os.Stderr.Write(msg)
	panic(S(""))
}

func assertNotNil(cond bool, tok *Token) {
	assert(cond, tok, S("should not be nil"))
}
