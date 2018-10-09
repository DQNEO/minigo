package main

import "fmt"
import "io/ioutil"
import "os"
import "errors"

type byteStream struct {
	filename  string
	source    string
	nextIndex int
	line      int
	column    int
}

type Token struct {
	typ  string
	sval string
	filename string
	line int
	column int
}

var bs *byteStream

func (bs *byteStream) getc() (byte, error) {
	if bs.nextIndex >= len(bs.source) {
		return 0, errors.New("EOF")
	}
	r := bs.source[bs.nextIndex]
	if r == '\n' {
		bs.line++
		bs.column = 1
	}
	bs.nextIndex++
	bs.column++
	return r, nil
}

func (bs *byteStream) ungetc() {
	bs.nextIndex--
	r := bs.source[bs.nextIndex]
	if r == '\n' {
		bs.line--
	}
}

func (tok *Token) String() string {
	return fmt.Sprintf("(%s \"%s\" %s:%d:%d)",
		tok.typ, tok.sval, tok.filename, tok.line, tok.column)
}

func (tok *Token) isPunct(s string) bool {
	return tok != nil && tok.typ == "punct" && tok.sval == s
}

func (tok *Token) isKeyword(s string) bool {
	return tok != nil && tok.typ == "keyword" && tok.sval == s
}

func (tok *Token) isIdent(s string) bool {
	return tok != nil && tok.typ == "ident" && tok.sval == s
}

func (tok *Token) isTypePunct() bool {
	return tok != nil && tok.typ == "punct"
}

func (tok *Token) isTypeKeyword() bool {
	return tok != nil && tok.typ == "keyword"
}

func (tok *Token) isTypeString() bool {
	return tok != nil && tok.typ == "string"
}

func (tok *Token) isTypeIdent() bool {
	return tok != nil && tok.typ == "ident"
}

func (tok *Token) isTypeSpace() bool {
	return tok != nil && tok.typ == "space"
}

func (tok *Token) isTypeNewline() bool {
	return tok != nil && tok.typ == "newline"
}

func getc() (byte, error) {
	return bs.getc()
}

func ungetc() {
	bs.ungetc()
}

func is_number(c byte) bool {
	return '0' <= c && c <= '9'
}

/**

 Operators and punctuation
 https://golang.org/ref/spec#Operators_and_punctuation

+    &     +=    &=     &&    ==    !=    (    )
-    |     -=    |=     ||    <     <=    [    ]
*    ^     *=    ^=     <-    >     >=    {    }
/    <<    /=    <<=    ++    =     :=    ,    ;
%    >>    %=    >>=    --    !     ...   .    :
     &^          &^=

 */
func is_punct(c byte) bool {
	switch c {
	case '+', '-', '(', ')', '=', '{', '}', '*', '[', ']', ',', ':', '.', '!', '<', '>', '&', '|', '%', '/':
		return true
	default:
		return false
	}
}

func read_number(c0 byte) string {
	var chars = []byte{c0}
	for {
		c, err := getc()
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
		c, err := getc()
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
		c, err := getc()
		if err != nil {
			panic("invalid string literal")
		}
		if c == '\\' {
			chars = append(chars, c)
			c, err = getc()
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
	c, err := getc()
	if err != nil {
		panic("unexpected EOF")
	}
	if c != e {
		fmt.Printf("char '%c' expected, but got '%c'\n", e, c)
		panic("unexpected char")
	}
}

func read_char() string {
	c, err := getc()
	if err != nil {
		panic("invalid char literal")
	}
	if c == '\\' {
		c, err = getc()
	}
	debugPrint("gotc:" + string(c))
	expect('\'')
	return string([]byte{c})
}

func is_space(c byte) bool {
	return c == ' ' || c == '\t'
}

func skip_space() {
	for {
		c, err := getc()
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

func makeToken(typ string, sval string) *Token {
	return &Token{
		typ: typ,
		sval: sval,
		filename: bs.filename,
		line: bs.line,
		column:bs.column,
	}
}

// https://golang.org/ref/spec#Keywords
var keywords = []string{
	"break", "default", "func", "interface", "select", "case", "defer", "go", "map", "struct", "chan", "else", "goto", "package", "switch", "const", "fallthrough", "if", "range", "type", "continue", "for", "import", "return", "var",
}

// util
func in_array(item string, list []string) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}

func tokenize() []*Token {
	var r []*Token
	for {
		c, err := getc()
		if err != nil {
			return r
		}
		var tok *Token
		switch {
		case c == 0:
			return r
		case c == '\n':
			tok = makeToken("newline", "")
		case is_number(c):
			sval := read_number(c)
			tok = makeToken( "number",  sval)
		case is_name(c):
			sval := read_name(c)
			if in_array(sval, keywords) {
				tok = makeToken("keyword", sval)
			} else {
				tok = makeToken( "ident",  sval)
			}
		case c == '\'':
			sval := read_char()
			tok = makeToken( "char",  sval)
		case c == '"':
			sval := read_string()
			tok = makeToken( "string",  sval)
		case c == ' ' || c == '\t':
			skip_space()
			tok = makeToken( "space", "")
		case is_punct(c):
			tok = makeToken( "punct", string([]byte{c}))
		default:
			fmt.Printf("c='%c'\n", c)
			panic("unknown char")
		}
		if debugMode {
			debugToken(tok)
		}
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

func tokenizeFromFile(path string) []*Token {
	s := readFile(path)
	bs = &byteStream{
		filename:path,
		source:    s,
		nextIndex: 0,
		line:      1,
		column:    0,
	}
	return tokenize()
}
