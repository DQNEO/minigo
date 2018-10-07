package main

import "fmt"
import "io/ioutil"
import "os"
import "errors"

var source string
var sourceInex int

func getc() (byte,error) {
	if sourceInex >= len(source) {
		return 0, errors.New("EOF")
	}
	r := source[sourceInex]
	//fmt.Printf("%c",r)
	sourceInex++
	return r, nil
}

func ungetc() {
	sourceInex--
}

func is_number(c byte) bool {
	return '0' <= c && c  <= '9'
}

func is_punct(c byte) bool {
	switch c {
	case '+', '-', '(', ')', '=', '{','}','*','[',']',',',':','.','!', '<','>','&','|', '%', '/':
		return true
	default:
		return false
	}
}

func read_number(c0 byte) string {
	var chars = []byte{c0}
	for {
		c,err := getc()
		if err != nil {
			return string(chars)
		}
		if is_number(c) {
			chars = append(chars, c)
			continue
		} else {
			ungetc()
			return string(chars)
		}
	}
}

func is_name(b byte) bool {
	return b == '_' || is_alphabet(b)
}


func is_alphabet(b byte) bool {
	return ('a' <= b && b <= 'z') || ('A' <= b && b <= 'Z')
}

func read_name(c0 byte) string {
	var chars = []byte{c0}
	for {
		c,err := getc()
		if err != nil {
			return string(chars)
		}
		if is_name(c) {
			chars = append(chars, c)
			continue
		} else {
			ungetc()
			return string(chars)
		}
	}
}

func read_string() string {
	var chars = []byte{}
	for {
		c,err := getc()
		if err != nil {
			panic("invalid string literal")
		}
		if c == '\\' {
			c,err = getc()
			chars = append(chars, c)
			continue
		}
		if c != '"' {
			chars = append(chars, c)
			continue
		} else {
			return string(chars)
		}
	}
}

func expect(e byte) {
	c,err := getc()
	if err != nil {
		panic("unexpected EOF")
	}
	if c != e {
		fmt.Printf("char '%c' expected, but got '%c'\n", e, c)
		panic("unexpected char")
	}
}

func read_char() string {
	c,err := getc()
	if err != nil {
		panic("invalid char literal")
	}
	if c == '\\' {
		c,err = getc()
	}
	debugPrint("gotc:" +  string(c))
	expect('\'')
	return string([]byte{c})
}

func is_space(c byte) bool {
	return  c == ' ' || c == '\t'
}

func skip_space() {
	for {
		c,err:= getc()
		if err != nil {
			return
		}
		if is_space(c) {
			continue
		} else {
			ungetc()
			return
		}
	}
}

func readFile(filename string) string {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func tokenize(s string) []*Token {
	var r []*Token
	source = s
	for  {
		c, err := getc()
		if err != nil {
			return r
		}
		var tok *Token
		switch  {
		case c == 0:
			return r
		case c == '\n':
			tok = &Token{typ:"newline"}
		case is_number(c):
			sval := read_number(c)
			tok = &Token{typ: "number", sval: sval}
		case is_name(c):
			sval := read_name(c)
			tok = &Token{typ:"ident", sval:sval}
		case c == '\'':
			sval := read_char()
			tok = &Token{typ: "char", sval:sval}
		case c == '"':
			sval := read_string()
			tok = &Token{typ: "string", sval:sval}
		case c == ' ' || c == '\t' :
			skip_space()
			tok = &Token{typ: "space"}
		case is_punct(c):
			tok = &Token{typ: "punct", sval: fmt.Sprintf("%c", c)}
		default:
			fmt.Printf("c='%c'\n", c)
			panic("unknown char")
		}
		debugToken(tok)
		r = append(r, tok)
	}

	return r
}


func renderTokens(tokens []*Token) {
	debugPrint("==== Start Render Tokens ===")
	for _, tok := range tokens {
		if tok.typ == "newline" {
			fmt.Fprintf(os.Stderr, "\n")
		} else if tok.typ == "space" {
			fmt.Fprintf(os.Stderr, "  ")
		} else if tok.typ == "string" {
			fmt.Fprintf(os.Stderr, "\"%s\"", tok.sval)
		} else {
			fmt.Fprintf(os.Stderr, tok.sval)
		}
	}
	debugPrint("==== End Render Tokens ===")
}

func tokenizeFromFile(path string) {
	s := readFile(path)
	tokens = tokenize(s)
}

