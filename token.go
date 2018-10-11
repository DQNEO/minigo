package main

import "fmt"
import "os"

type Token struct {
	typ  string
	sval string
	filename string
	line int
	column int
}

var bs *ByteStream

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
	return bs.get()
}

func ungetc() {
	bs.unget()
}

func is_number(c byte) bool {
	return '0' <= c && c <= '9'
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
func tokenize() []*Token {
	var r []*Token
	for {
		c, err := getc()
		if err != nil {
			return r
		}
		var tok *Token
		switch c {
		case 0:
			return r
		case '\n':
			tok = makeToken("newline", "")
		case '0','1','2','3','4','5','6','7','8','9':
			sval := read_number(c)
			tok = makeToken( "number",  sval)
		case '_','A','B','C','D','E','F','G','H','I','J','K','L','M','N','O','P','Q','R','S','T','U','V','W','X','Y','Z',
			'a','b','c','d','e','f','g','h','i','j','k','l','m','n','o','p','q','r','s','t','u','v','w','x','y','z':
			sval := read_name(c)
			if in_array(sval, keywords) {
				tok = makeToken("keyword", sval)
			} else {
				tok = makeToken( "ident",  sval)
			}
		case '\'':
			sval := read_char()
			tok = makeToken( "char",  sval)
		case '"':
			sval := read_string()
			tok = makeToken( "string",  sval)
		case ' ','\t':
			skip_space()
			tok = makeToken( "space", string(c))
		case '(',')','[',']','{','}',',',';':
			tok = makeToken( "punct", string(c))
		case '!':
			c, _ := getc()
			if c == '=' {
				tok = makeToken( "punct", "!=")
			} else {
				ungetc()
				tok = makeToken( "punct", "!")
			}
		case '%':
			c, _ := getc()
			if c == '=' {
				tok = makeToken( "punct", "%=")
			} else {
				ungetc()
				tok = makeToken( "punct", "%")
			}
		case '*':
			c, _ := getc()
			if c == '=' {
				tok = makeToken( "punct", "*=")
			} else {
				ungetc()
				tok = makeToken( "punct", "*")
			}
		case ':':
			c, _ := getc()
			if c == '=' {
				tok = makeToken( "punct", ":=")
			} else {
				ungetc()
				tok = makeToken( "punct", ":")
			}
		case '=':
			c, _ := getc()
			if c == '=' {
				tok = makeToken( "punct", "==")
			} else {
				ungetc()
				tok = makeToken( "punct", "=")
			}
		case '^':
			c, _ := getc()
			if c == '=' {
				tok = makeToken( "punct", "^=")
			} else {
				ungetc()
				tok = makeToken( "punct", "^")
			}
		case '&':
			c, _ := getc()
			if c == '&' {
				tok = makeToken("punct", "&&")
			} else if c == '=' {
				tok = makeToken("punct", "&=")
			} else if c == '^' {
				c, _ := getc()
				if c == '=' {
					tok = makeToken("punct", "&^=")
				} else {
					ungetc()
					tok = makeToken( "punct", "&^")
				}
			} else {
				ungetc()
				tok = makeToken( "punct", "&")
			}
		case '+':
			c, _ = getc()
			if c == '+' {
				tok = makeToken("punct", "++")
			} else if c == '=' {
				tok = makeToken("punct", "+=")
			} else {
				ungetc()
				tok = makeToken( "punct", "+")
			}
		case '-':
			c, _ = getc()
			if c == '-' {
				tok = makeToken("punct", "--")
			} else if c == '=' {
				tok = makeToken("punct", "-=")
			} else {
				ungetc()
				tok = makeToken( "punct", "-")
			}
		case '|':
			c, _ = getc()
			if c == '=' {
				tok = makeToken("punct", "|=")
			} else if c == '|' {
				tok = makeToken("punct", "||")
			} else {
				ungetc()
				tok = makeToken( "punct", "|")
			}
		case '.':
			c, _ = getc()
			if c == '.' {
				c, _ = getc()
				if c == '.' {
					tok = makeToken("punct", "...")
				} else {
					panic("invalid token '..'")
				}
			} else {
				ungetc()
				tok = makeToken( "punct", ".")
			}
		case '>':
			c, _ = getc()
			if c == '=' {
				tok = makeToken("punct", ">=")
			} else if c == '>' {
				c, _ = getc()
				if c == '=' {
					tok = makeToken("punct", ">>=")
				} else {
					ungetc()
					tok = makeToken("punct", ">>")
				}
			} else {
				ungetc()
				tok = makeToken( "punct", ">")
			}
		case '<':
			c ,_ = getc()
			if c == '-' {
				tok = makeToken("punct", "<-")
			} else if c == '=' {
				tok = makeToken("punct", ">=")
			} else if c == '<' {
				c ,_ = getc()
				if c == '=' {
					tok = makeToken("punct", "<<=")
				} else {
					ungetc()
					tok = makeToken("punct", "<<")
				}
			} else {
				ungetc()
				tok = makeToken( "punct", "<")
			}
		case '/':
			c ,_ = getc()
			if c == '=' {
				tok = makeToken("punct", "/=")
			} else {
				ungetc()
				tok = makeToken("punct", "/")
			}
			// @TODO handle comments
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
			fmt.Fprintf(os.Stderr, "%s", tok.sval)
		}
	}
	debugPrint("==== End Render Tokens ===")
}

func tokenizeFromFile(path string) []*Token {
	bs = NewByteStream(path)
	return tokenize()
}
