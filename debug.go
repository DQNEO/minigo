package main

import "fmt"
import "os"

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

	if GENERATION == 2 {
		//fmt.Printf("# ")
		//fmt.Printf(string(indents))
		//fmt.Printf("%s\n", format)
		//return
	}

	var s string = string(indents)+format+"\n"
	os.Stderr.Write([]byte(s))
	//fmt.Fprintf(os.Stderr, s, v...)
}

func dumpToken(tok *Token) {
	var frmt string = "tok: type=%-8s, sval=\"%s\"\n"
	var a []interface{}
	var ifc interface{} = tok.typ
	var ifc2 interface{} = tok.sval
	a = nil
	a = append(a, ifc)
	a = append(a, ifc2)

	var s string = fmt.Sprintf(frmt, a)
	os.Stderr.Write([]byte(s))
	//debugf(s)
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
	if GENERATION == 2 {
		panic(format)
	} else {
		errorf(format+"\n "+tokString, v...)
	}
}

func errorf(format string, v ...interface{}) {
	/*
		currentTokenIndex := ts.index - 1
		fmt.Printf("%v %v %v\n",
			ts.getToken(currentTokenIndex-2), ts.getToken(currentTokenIndex-1), ts.getToken(currentTokenIndex))
	*/
	var s string
	//s += bs.location() + ": "
	s += fmt.Sprintf(format, v...)
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
