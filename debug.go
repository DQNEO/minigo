package main

import (
	"os"
)

func debugf(format string, v ...interface{}) {
	if !debugMode {
		return
	}

	var indents []byte

	for i := 0; i < debugNest; i++ {
		indents = append(indents, ' ')
		indents = append(indents, ' ')
	}
	os.Stderr.Write(indents)
	s2 := Sprintf(format, v...)
	var b []byte = []byte(s2)
	b = append(b, '\n')
	os.Stderr.Write(b)
}

var debugNest int

// States "To Be Implemented"
func TBI(tok *Token, format string, v ...interface{}) {
	format2 := concat(S("(To Be Implemented) "), bytes(format))
	errorft(tok, format2, v...)
}

// errorf with a position token
func errorft(tok *Token, format string, v ...interface{}) {

	var tokString bytes
	if tok != nil {
		tokString = tok.String()
	}
	gs := concat3(bytes(format),S("\n "), tokString)
	errorf(gs, v...)
}

func errorf(format string, v ...interface{}) {
	s := Sprintf(format, v...)
	b := bytes(s)
	b = append(b, '\n')
	os.Stderr.Write(b)
	panic(S(""))
}

func assert(cond bool, tok *Token, format string, v ...interface{}) {
	if !cond {
		s := Sprintf(string(format), v...)
		var toks bytes = S("")
		if tok != nil {
			toks = tok.String()
		}
		msg := concat3(S("assertion failed: "), bytes(s),  toks)
		os.Stderr.Write([]byte(msg))
		os.Stderr.Write(S("\n"))
		panic(S(""))
	}
}

func assertNotReached(tok *Token) {
	msg := concat(S("assertNotReached "), tok.String())
	os.Stderr.Write([]byte(msg))
	panic(S(""))
}

func assertNotNil(cond bool, tok *Token) {
	assert(cond, tok, "should not be nil")
}
