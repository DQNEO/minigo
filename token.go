package main

import (
	"fmt"
	"os"
	"strconv"
)

// https://golang.org/ref/spec#Keywords
var keywords = []string{
	"break", "default", "func", "interface", "select", "case", "defer", "go", "map",
	"struct", "chan", "else", "goto", "package", "switch", "const", "fallthrough",
	"if", "range", "type", "continue", "for", "import", "return", "var",
}

type identifier string

type TokenType string

const (
	T_EOF      TokenType = "EOF"
	T_INT      TokenType = "int"
	T_STRING   TokenType = "string"
	T_CHAR     TokenType = "char"
	T_IDENT    TokenType = "ident"
	T_PUNCT    TokenType = "punct"
	T_KEYWORWD TokenType = "keyword"
)

type Token struct {
	typ      TokenType
	sval     gostring
	filename string
	line     int
	column   int
}

type TokenStream struct {
	tokens []*Token
	index  int
}

func NewTokenStream(bs *ByteStream) *TokenStream {
	tokens := Tokenize(bs)
	assert(len(tokens) > 0, nil, "tokens should have length")
	return &TokenStream{
		tokens: tokens,
		index:  0,
	}

}

func (ts *TokenStream) isEnd() bool {
	return ts.index > len(ts.tokens)-1
}

func (tok *Token) String() string {
	var sval cstring = ""
	if len(tok.sval) > 0 {
		sval = cstring(tok.sval)
	}
	return fmt.Sprintf("(\"%s\" at %s:%d:%d)",
		sval, tok.filename, tok.line, tok.column)
}

func (tok *Token) isEOF() bool {
	return tok != nil && tok.typ == T_EOF
}

func (tok *Token) isPunct(s string) bool {
	gs := gostring(s)
	return tok != nil && tok.typ == T_PUNCT && eqGostring(tok.sval, gs)
}

func (tok *Token) isKeyword(s string) bool {
	gs := gostring(s)
	return tok != nil && tok.typ == T_KEYWORWD && eqGostring(tok.sval,gs)
}

func (tok *Token) isIdent(s string) bool {
	gs := gostring(s)
	return tok != nil && tok.typ == T_IDENT && eqGostring(tok.sval,gs)
}

func (tok *Token) getIdent() identifier {
	if !tok.isTypeIdent() {
		errorft(tok, "ident expeced, but got %v", tok)
	}
	return identifier(tok.sval)
}

func (tok *Token) getIntval() int {
	val, _ := strconv.Atoi(string(tok.sval))
	return val
}

func (tok *Token) isTypePunct() bool {
	return tok != nil && tok.typ == T_PUNCT
}

func (tok *Token) isTypeKeyword() bool {
	return tok != nil && tok.typ == T_KEYWORWD
}

func (tok *Token) isTypeInt() bool {
	return tok != nil && tok.typ == T_INT
}

func (tok *Token) isTypeChar() bool {
	return tok != nil && tok.typ == T_CHAR
}

func (tok *Token) isTypeString() bool {
	return tok != nil && tok.typ == T_STRING
}

func (tok *Token) isTypeIdent() bool {
	return tok != nil && tok.typ == T_IDENT
}

func (tok *Token) isSemicolon() bool {
	return tok.isPunct(";")
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

func (tok *Token) dump() {
	var s string = fmt.Sprintf("tok: line=%d, type=%s, sval=\"%s\"\n", tok.line, tok.typ, tok.sval)
	var b []byte = []byte(s)
	os.Stderr.Write(b)
}
