package main

import (
	"os"
)

// https://golang.org/ref/spec#Keywords
var keywords = []gostring{
	gostring("break"),
	gostring("default"),
	gostring("func"),
	gostring("interface"),
	gostring("select"),
	gostring("case"),
	gostring("defer"),
	gostring("go"),
	gostring("map"),
	gostring("struct"),
	gostring("chan"),
	gostring("else"),
	gostring("goto"),
	gostring("package"),
	gostring("switch"),
	gostring("const"),
	gostring("fallthrough"),
	gostring("if"),
	gostring("range"),
	gostring("type"),
	gostring("continue"),
	gostring("for"),
	gostring("import"),
	gostring("return"),
	gostring("var"),
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

func typeToGostring (typ TokenType) gostring {
	switch typ {
	case T_EOF:
			return S("EOF")
	case T_INT:
		return S("int")
	case T_STRING:
		return S("string")
	case T_CHAR:
		return S("char")
	case T_IDENT:
		return S("ident")
	case T_PUNCT:
		return S("punct")
	case T_KEYWORWD:
		return S("keyword")
	}

	return S("")
}

type Token struct {
	typ      TokenType
	sval     gostring
	filename gostring
	line     int
	column   int
}

type TokenStream struct {
	tokens []*Token
	index  int
}

func NewTokenStream(bs *ByteStream) *TokenStream {
	tokens := Tokenize(bs)
	assert(len(tokens) > 0, nil, S("tokens should have length"))
	return &TokenStream{
		tokens: tokens,
		index:  0,
	}

}

func (ts *TokenStream) isEnd() bool {
	return ts.index > len(ts.tokens)-1
}

func (tok *Token) getSval() gostring {
	if len(tok.sval) > 0 {
		return tok.sval
	}
	return S("")
}

func (tok *Token) String() gostring {
	sval := tok.getSval()
	gs := Sprintf(S("(\"%s\" at %s:%d:%d)"),
		sval, gostring(tok.filename), tok.line, tok.column)
	return gs
}

func (tok *Token) isEOF() bool {
	return tok != nil && tok.typ == T_EOF
}

func (tok *Token) isPunct(s gostring) bool {
	gs := gostring(s)
	return tok != nil && tok.typ == T_PUNCT && eq(tok.sval, gs)
}

func (tok *Token) isKeyword(s gostring) bool {
	gs := gostring(s)
	return tok != nil && tok.typ == T_KEYWORWD && eq(tok.sval,gs)
}

func (tok *Token) isIdent(s gostring) bool {
	gs := gostring(s)
	return tok != nil && tok.typ == T_IDENT && eq(tok.sval,gs)
}

func (tok *Token) getIdent() goidentifier {
	if !tok.isTypeIdent() {
		errorft(tok, S("ident expeced, but got %v"), tok)
	}
	return goidentifier(tok.sval)
}

func (tok *Token) getIntval() int {
	val, _ := strconv_Atoi(tok.sval)
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
	s := Sprintf(S("tok: line=%d, type=%s, sval=\"%s\"\n"),
		tok.line, typeToGostring(tok.typ), sval)
	var b []byte = []byte(s)
	os.Stderr.Write(b)
}
