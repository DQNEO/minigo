package main

import (
	"os"
	"strconv"
)

// https://golang.org/ref/spec#Keywords
var keywords = []string{
	"break",
	"default",
	"func",
	"interface",
	"select",
	"case",
	"defer",
	"go",
	"map",
	"struct",
	"chan",
	"else",
	"goto",
	"package",
	"switch",
	"const",
	"fallthrough",
	"if",
	"range",
	"type",
	"continue",
	"for",
	"import",
	"return",
	"var",
}

type TokenType int

const (
	T_EOF      TokenType = iota
	T_INT
	T_STRING
	T_CHAR
	T_IDENT
	T_PUNCT
	T_KEYWORWD
	)

func typeToGostring (typ TokenType) string {
	switch typ {
	case T_EOF:
			return "EOF"
	case T_INT:
		return "int"
	case T_STRING:
		return "string"
	case T_CHAR:
		return "char"
	case T_IDENT:
		return "ident"
	case T_PUNCT:
		return "punct"
	case T_KEYWORWD:
		return "keyword"
	}

	return ""
}

type Token struct {
	typ      TokenType
	sval     string
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

func (tok *Token) getSval() string {
	if len(tok.sval) > 0 {
		return tok.sval
	}
	return ""
}

func (tok *Token) String() string {
	sval := tok.getSval()
	gs := Sprintf("(\"%s\" at %s:%d:%d)",
		sval, bytes(tok.filename), tok.line, tok.column)
	return gs
}

func (tok *Token) isEOF() bool {
	return tok != nil && tok.typ == T_EOF
}

func (tok *Token) isPunct(s bytes) bool {
	gs := bytes(s)
	return tok != nil && tok.typ == T_PUNCT && eq(bytes(tok.sval), gs)
}

func (tok *Token) isKeyword(s bytes) bool {
	gs := bytes(s)
	return tok != nil && tok.typ == T_KEYWORWD && eq(bytes(tok.sval),gs)
}

func (tok *Token) isIdent(s bytes) bool {
	gs := bytes(s)
	return tok != nil && tok.typ == T_IDENT && eq(bytes(tok.sval),gs)
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
	return tok.isPunct(S(";"))
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
	sval := tok.getSval()
	s := Sprintf("tok: line=%d, type=%s, sval=\"%s\"\n",
		tok.line, typeToGostring(tok.typ), sval)
	var b []byte = []byte(s)
	os.Stderr.Write(b)
}
