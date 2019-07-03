package main

import (
	"os"
)

func debugf(format gostring, v ...interface{}) {
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
	errorft(tok, "(To Be Implemented) "+format, v...)
}

// errorf with a position token
func errorft(tok *Token, format string, v ...interface{}) {
	var format2 gostring = gostring(format)
	var tokString gostring
	if tok != nil {
		tokString = tok.String()
	}
	gs := concat3(format2,S("\n "), tokString)
	errorf(gs, v...)
}

func errorf(format gostring, v ...interface{}) {
	s := Sprintf(format, v...)
	os.Stderr.Write(concat(s, S("\n")))
	panic("")
}

func assert(cond bool, tok *Token, format string, v ...interface{}) {
	if !cond {
		s := Sprintf(gostring(format), v...)
		msg := concat3(S("assertion failed: "), s,  tok.String())
		os.Stderr.Write(msg)
		panic("")
	}
}

func assertNotReached(tok *Token) {
	msg := concat(S("assertNotReached "), tok.String())
	os.Stderr.Write(msg)
	panic("")
}

func assertNotNil(cond bool, tok *Token) {
	assert(cond, tok, "should not be nil")
}
