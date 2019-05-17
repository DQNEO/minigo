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

func (tok *Token) String() string {
	return fmt.Sprintf("(\"%s\" at %s:%d:%d)",
		tok.sval, tok.filename, tok.line, tok.column)
}

func (tok *Token) isEOF() bool {
	return tok != nil && tok.typ == T_EOF
}

func (tok *Token) isPunct(s string) bool {
	return tok != nil && tok.typ == T_PUNCT && tok.sval == s
}

func (tok *Token) isKeyword(s string) bool {
	return tok != nil && tok.typ == T_KEYWORWD && tok.sval == s
}

func (tok *Token) isIdent(s string) bool {
	return tok != nil && tok.typ == T_IDENT && tok.sval == s
}

func (tok *Token) getIdent() identifier {
	if !tok.isTypeIdent() {
		errorft(tok, "ident expeced, but got %v", tok)
	}
	return identifier(tok.sval)
}

func (tok *Token) getIntval() int {
	val, _ := strconv.Atoi(tok.sval)
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

type Tokenizer struct {
	bs *ByteStream
}

func (tn *Tokenizer) read_number(c0 byte) string {
	var chars = []byte{c0}
	for {
		c, err := tn.bs.get()
		if err != nil {
			return string(chars)
		}
		if tn.isUnicodeDigit(c) {
			chars = append(chars, c)
			continue
		} else {
			tn.bs.unget()
			return string(chars)
		}
	}
}

// https://golang.org/ref/spec#unicode_letter
func (tn *Tokenizer) isUnicodeLetter(b byte) bool {
	// tentative implementation
	return ('a' <= b && b <= 'z') || ('A' <= b && b <= 'Z')
}

// https://golang.org/ref/spec#unicode_digit
func (tn *Tokenizer) isUnicodeDigit(c byte) bool {
	// tentative implementation
	return '0' <= c && c <= '9'
}

// https://golang.org/ref/spec#Letters_and_digits
func (tn *Tokenizer) isLetter(b byte) bool {
	return tn.isUnicodeLetter(b) || b == '_'
}

// https://golang.org/ref/spec#Identifiers
func (tn *Tokenizer) readIdentifier(c0 byte) string {
	var chars = []byte{c0}
	for {
		c, err := tn.bs.get()
		if err != nil {
			return string(chars)
		}
		if tn.isLetter(c) || tn.isUnicodeDigit(c) {
			chars = append(chars, c)
			continue
		} else {
			tn.bs.unget()
			return string(chars)
		}
	}
}

func (tn *Tokenizer) read_string() string {
	var chars []byte
	for {
		c, err := tn.bs.get()
		if err != nil {
			panic("invalid string literal")
		}
		if c == '\\' {
			chars = append(chars, c)
			c, err = tn.bs.get()
			chars = append(chars, c)
			continue
		}
		if c == '\n' {
			chars = append(chars, '\\')
			chars = append(chars, 'n')
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

func (tn *Tokenizer) read_raw_string() string {
	var chars []byte
	for {
		c, err := tn.bs.get()
		if err != nil {
			panic("invalid string literal")
		}
		if c == '\\' {
			// @FIXME
			chars = append(chars, c)
			c, err = tn.bs.get()
			chars = append(chars, c)
			continue
		}
		if c == '"' {
			chars = append(chars, '\\')
			chars = append(chars, c)
			continue
		}
		if c != '`' {
			if c == '\n' {
				chars = append(chars, '\\')
				chars = append(chars, 'n')
			} else {
				chars = append(chars, c)
			}
			continue
		} else {
			return string(chars)
		}
	}
}

func (tn *Tokenizer) read_char() string {
	c, err := tn.bs.get()
	if err != nil {
		panic("invalid char literal")
	}
	if c == '\\' {
		var sval string
		c, err = tn.bs.get()
		switch c {
		case 'n':
			sval = "\n"
		case 't':
			sval = "\t"
		case 'r':
			sval = "\r"
		case '\\':
			sval = "\\"
		case '\'':
			sval = "'"
		default:
			errorf("unexpected char 1:%c", c)
		}

		end, _ := tn.bs.get()
		if end != '\'' {
			errorf("unexpected char 2:%c", end)
		}
		return sval
	}
	end, _ := tn.bs.get()
	if end != '\'' {
		errorf("unexpected char:%c", end)
	}
	return string([]byte{c})
}

func (tn *Tokenizer) isSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\r'
}

func (tn *Tokenizer) skipSpace() {
	for {
		c, err := tn.bs.get()
		if err != nil {
			return
		}
		if tn.isSpace(c) {
			continue
		} else {
			tn.bs.unget()
			return
		}
	}
}

func (tn *Tokenizer) makeToken(typ TokenType, sval string) *Token {
	return &Token{
		typ:      typ,
		sval:     sval,
		filename: tn.bs.filename,
		line:     tn.bs.line,
		column:   tn.bs.column,
	}
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

var semicolonToken = Token{
	typ:  "punct", // @FIXME: should be T_PUNCT
	sval: ";",
}

// https://golang.org/ref/spec#Semicolons
func (tn *Tokenizer) autoSemicolonInsert(last *Token) bool {
	if last.isTypeIdent() {
		return true
	}
	if last.typ == T_INT || last.typ == T_STRING || last.typ == T_CHAR {
		return true
	}
	if last.isKeyword("break") || last.isKeyword("continue") || last.isKeyword("fallthrough") || last.isKeyword("return") {
		return true
	}

	if last.isPunct("++") || last.isPunct("--") || last.isPunct(")") || last.isPunct("]") || last.isPunct("}") {
		return true
	}

	return false
}

func (tn *Tokenizer) skipLine() {
	for {
		c, err := tn.bs.get()
		if err != nil || c == '\n' {
			tn.bs.unget()
			return
		}
	}
}

func (tn *Tokenizer) skipBlockComment() {
	var hasReadAsterisk bool

	for {
		c, err := tn.bs.get()
		if err != nil {
			errorf("premature end of block comment")
		}
		if c == '*' {
			hasReadAsterisk = true
		} else if hasReadAsterisk && c == '/' {
			return
		} else {
			hasReadAsterisk = false
		}
	}
}

func isIn(c byte, set []byte) bool {
	for _, c2 := range set {
		if c == c2 {
			return true
		}
	}
	return false
}

func Tokenize(bs *ByteStream) []*Token {

	var tn = &Tokenizer{
		bs:bs,
	}

	var r []*Token
	for {
		c, err := tn.bs.get()
		if err != nil {
			return r
		}
		var tok *Token
		switch c {
		case 0: // no need?
			return r
		case '\n':
			// Insert semicolon
			if len(r) > 0 {
				last := r[len(r)-1]
				if tn.autoSemicolonInsert(last) {
					r = append(r, &semicolonToken)
				}
			}
			continue
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			sval := tn.read_number(c)
			tok = tn.makeToken(T_INT, sval)
		case '_', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
			'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z':
			sval := tn.readIdentifier(c)
			if in_array(sval, keywords) {
				tok = tn.makeToken(T_KEYWORWD, sval)
			} else {
				tok = tn.makeToken(T_IDENT, sval)
			}
		case '\'':
			sval := tn.read_char()
			tok = tn.makeToken(T_CHAR, sval)
		case '"':
			sval := tn.read_string()
			tok = tn.makeToken(T_STRING, sval)
		case '`':
			sval := tn.read_raw_string()
			tok = tn.makeToken(T_STRING, sval)
		case ' ', '\t':
			tn.skipSpace()
			continue
		case '/':
			c, _ = tn.bs.get()
			if c == '/' {
				tn.skipLine()
				continue
			} else if c == '*' {
				tn.skipBlockComment()
				continue
			} else if c == '=' {
				tok = tn.makeToken(T_PUNCT, "/=")
			} else {
				tn.bs.unget()
				tok = tn.makeToken(T_PUNCT, "/")
			}
		case '(', ')', '[', ']', '{', '}', ',', ';':
			tok = tn.makeToken(T_PUNCT, string([]byte{c}))
		case '!':
			c, _ := tn.bs.get()
			if c == '=' {
				tok = tn.makeToken(T_PUNCT, "!=")
			} else {
				tn.bs.unget()
				tok = tn.makeToken(T_PUNCT, "!")
			}
		case '%':
			c, _ := tn.bs.get()
			if c == '=' {
				tok = tn.makeToken(T_PUNCT, "%=")
			} else {
				tn.bs.unget()
				tok = tn.makeToken(T_PUNCT, "%")
			}
		case '*':
			c, _ := tn.bs.get()
			if c == '=' {
				tok = tn.makeToken(T_PUNCT, "*=")
			} else {
				tn.bs.unget()
				tok = tn.makeToken(T_PUNCT, "*")
			}
		case ':':
			c, _ := tn.bs.get()
			if c == '=' {
				tok = tn.makeToken(T_PUNCT, ":=")
			} else {
				tn.bs.unget()
				tok = tn.makeToken(T_PUNCT, ":")
			}
		case '=':
			c, _ := tn.bs.get()
			if c == '=' {
				tok = tn.makeToken(T_PUNCT, "==")
			} else {
				tn.bs.unget()
				tok = tn.makeToken(T_PUNCT, "=")
			}
		case '^':
			c, _ := tn.bs.get()
			if c == '=' {
				tok = tn.makeToken(T_PUNCT, "^=")
			} else {
				tn.bs.unget()
				tok = tn.makeToken(T_PUNCT, "^")
			}
		case '&':
			c, _ := tn.bs.get()
			if c == '&' {
				tok = tn.makeToken(T_PUNCT, "&&")
			} else if c == '=' {
				tok = tn.makeToken(T_PUNCT, "&=")
			} else if c == '^' {
				c, _ := tn.bs.get()
				if c == '=' {
					tok = tn.makeToken(T_PUNCT, "&^=")
				} else {
					tn.bs.unget()
					tok = tn.makeToken(T_PUNCT, "&^")
				}
			} else {
				tn.bs.unget()
				tok = tn.makeToken(T_PUNCT, "&")
			}
		case '+':
			c, _ = tn.bs.get()
			if c == '+' {
				tok = tn.makeToken(T_PUNCT, "++")
			} else if c == '=' {
				tok = tn.makeToken(T_PUNCT, "+=")
			} else {
				tn.bs.unget()
				tok = tn.makeToken(T_PUNCT, "+")
			}
		case '-':
			c, _ = tn.bs.get()
			if c == '-' {
				tok = tn.makeToken(T_PUNCT, "--")
			} else if c == '=' {
				tok = tn.makeToken(T_PUNCT, "-=")
			} else {
				tn.bs.unget()
				tok = tn.makeToken(T_PUNCT, "-")
			}
		case '|':
			c, _ = tn.bs.get()
			if c == '=' {
				tok = tn.makeToken(T_PUNCT, "|=")
			} else if c == '|' {
				tok = tn.makeToken(T_PUNCT, "||")
			} else {
				tn.bs.unget()
				tok = tn.makeToken(T_PUNCT, "|")
			}
		case '.':
			c, _ = tn.bs.get()
			if c == '.' {
				c, _ = tn.bs.get()
				if c == '.' {
					tok = tn.makeToken(T_PUNCT, "...")
				} else {
					panic("invalid token '..'")
				}
			} else {
				tn.bs.unget()
				tok = tn.makeToken(T_PUNCT, ".")
			}
		case '>':
			c, _ = tn.bs.get()
			if c == '=' {
				tok = tn.makeToken(T_PUNCT, ">=")
			} else if c == '>' {
				c, _ = tn.bs.get()
				if c == '=' {
					tok = tn.makeToken(T_PUNCT, ">>=")
				} else {
					tn.bs.unget()
					tok = tn.makeToken(T_PUNCT, ">>")
				}
			} else {
				tn.bs.unget()
				tok = tn.makeToken(T_PUNCT, ">")
			}
		case '<':
			c, _ = tn.bs.get()
			if c == '-' {
				tok = tn.makeToken(T_PUNCT, "<-")
			} else if c == '=' {
				tok = tn.makeToken(T_PUNCT, "<=")
			} else if c == '<' {
				c, _ = tn.bs.get()
				if c == '=' {
					tok = tn.makeToken(T_PUNCT, "<<=")
				} else {
					tn.bs.unget()
					tok = tn.makeToken(T_PUNCT, "<<")
				}
			} else {
				tn.bs.unget()
				tok = tn.makeToken(T_PUNCT, "<")
			}
		default:
			fmt.Printf("c=%d\n", c)
			panic("unknown char")
		}
		if debugToken {
			tok.dump()
		}
		r = append(r, tok)
	}

	return r
}

func (tok *Token) dump() {
	var s string = fmt.Sprintf("tok: line=%d, type=%s, sval=\"%s\"\n", tok.line, tok.typ, tok.sval)
	var b []byte = []byte(s)
	os.Stderr.Write(b)
}
