package main

import (
	"fmt"
	"os"
)

func debugln(s string) {
	if !debugMode {
		return
	}
	fmt.Printf("# %s\n", s)
}

func debugf(format string, v ...interface{}) {
	if !debugMode {
		return
	}

	var indents []byte

	for i := 0; i < debugNest; i++ {
		indents = append(indents, ' ')
		indents = append(indents, ' ')
	}

	var format2 string = string(indents)+format+"\n"
	s2 := fmt.Sprintf(format2, v...)
	var b []byte = []byte(s2)
	os.Stderr.Write(b)
}

func dumpToken(tok *Token) {
	var s string = fmt.Sprintf("tok: type=%-8s, sval=\"%s\"\n", tok.typ, tok.sval)
	os.Stderr.Write([]byte(s))
}

var debugNest int

// States "To Be Implemented"
func TBI(tok *Token, format string, v ...interface{}) {
	errorft(tok, "(To Be Implemented) "+format, v...)
}

// errorf with a position token
func errorft(tok *Token, format string, v ...interface{}) {
	var tokString string
	if tok != nil {
		tokString = tok.String()
	}
	errorf(format+"\n "+tokString, v...)
}

func errorf(format string, v ...interface{}) {
	/*
		currentTokenIndex := ts.index - 1
		fmt.Printf("%v %v %v\n",
			ts.getToken(currentTokenIndex-2), ts.getToken(currentTokenIndex-1), ts.getToken(currentTokenIndex))
	*/
	//s += bs.location() + ": "
	s := fmt.Sprintf(format, v...)
	panic(s)
}

func assert(cond bool, tok *Token, msg string) {
	if !cond {
		if GENERATION == 2 {
			print("assertion failed:")
			panic(msg)
		} else {
			panic(fmt.Sprintf("assertion failed: %s %s", msg, tok))
		}
	}
}

func assertNotNil(cond bool, tok *Token) {
	assert(cond, tok, "should not be nil")
}
