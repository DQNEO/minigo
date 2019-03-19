package main

import (
	"fmt"
)

var debugMode = true
var debugToken = true

func tokenize2(_bs *ByteStream) []*Token {
	bs = _bs
	var r []*Token
	for {
		c, err := getc()
		if err != nil {
			return r
		}
		var tok *Token
		switch c {
		case 0: // no need?
			return r
		case '\n':
			fmt.Printf("newline:'%c'\n", c)
			continue
			// Insert semicolon
			if len(r) > 0 {
				last := r[len(r)-1]
				if autoSemicolonInsert(last) {
					r = append(r, semicolon)
				}
			}
			continue
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			sval := read_number(c)
			tok = makeToken(T_INT, sval)
		case '_', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
			'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z':
			fmt.Printf("alphabet:'%c'\n", c)
			sval := readIdentifier(c)
			fmt.Printf("alphabet sval:%s\n", sval)
			continue
			if in_array(sval, keywords) {
				tok = makeToken(T_KEYWORWD, sval)
			} else {
				tok = makeToken(T_IDENT, sval)
			}
		case '\'':
			sval := read_char()
			tok = makeToken(T_CHAR, sval)
		case '"':
			sval := read_string()
			tok = makeToken(T_STRING, sval)
		case '`':
			sval := read_raw_string()
			tok = makeToken(T_STRING, sval)
		case ' ', '\t':
			fmt.Printf("space:'%c'\n", c)
			continue
			skipSpace()
			continue
		case '/':
			c, _ = getc()
			if c == '/' {
				skipLine()
				continue
			} else if c == '*' {
				skipBlockComment()
				continue
			} else if c == '=' {
				tok = makeToken(T_PUNCT, "/=")
			} else {
				ungetc()
				tok = makeToken(T_PUNCT, "/")
			}
		case '(', ')', '[', ']', '{', '}', ',', ';':
			tok = makeToken(T_PUNCT, string(c))
		case '!':
			c, _ := getc()
			if c == '=' {
				tok = makeToken(T_PUNCT, "!=")
			} else {
				ungetc()
				tok = makeToken(T_PUNCT, "!")
			}
		case '%':
			c, _ := getc()
			if c == '=' {
				tok = makeToken(T_PUNCT, "%=")
			} else {
				ungetc()
				tok = makeToken(T_PUNCT, "%")
			}
		case '*':
			c, _ := getc()
			if c == '=' {
				tok = makeToken(T_PUNCT, "*=")
			} else {
				ungetc()
				tok = makeToken(T_PUNCT, "*")
			}
		case ':':
			c, _ := getc()
			if c == '=' {
				tok = makeToken(T_PUNCT, ":=")
			} else {
				ungetc()
				tok = makeToken(T_PUNCT, ":")
			}
		case '=':
			c, _ := getc()
			if c == '=' {
				tok = makeToken(T_PUNCT, "==")
			} else {
				ungetc()
				tok = makeToken(T_PUNCT, "=")
			}
		case '^':
			c, _ := getc()
			if c == '=' {
				tok = makeToken(T_PUNCT, "^=")
			} else {
				ungetc()
				tok = makeToken(T_PUNCT, "^")
			}
		case '&':
			c, _ := getc()
			if c == '&' {
				tok = makeToken(T_PUNCT, "&&")
			} else if c == '=' {
				tok = makeToken(T_PUNCT, "&=")
			} else if c == '^' {
				c, _ := getc()
				if c == '=' {
					tok = makeToken(T_PUNCT, "&^=")
				} else {
					ungetc()
					tok = makeToken(T_PUNCT, "&^")
				}
			} else {
				ungetc()
				tok = makeToken(T_PUNCT, "&")
			}
		case '+':
			c, _ = getc()
			if c == '+' {
				tok = makeToken(T_PUNCT, "++")
			} else if c == '=' {
				tok = makeToken(T_PUNCT, "+=")
			} else {
				ungetc()
				tok = makeToken(T_PUNCT, "+")
			}
		case '-':
			c, _ = getc()
			if c == '-' {
				tok = makeToken(T_PUNCT, "--")
			} else if c == '=' {
				tok = makeToken(T_PUNCT, "-=")
			} else {
				ungetc()
				tok = makeToken(T_PUNCT, "-")
			}
		case '|':
			c, _ = getc()
			if c == '=' {
				tok = makeToken(T_PUNCT, "|=")
			} else if c == '|' {
				tok = makeToken(T_PUNCT, "||")
			} else {
				ungetc()
				tok = makeToken(T_PUNCT, "|")
			}
		case '.':
			c, _ = getc()
			if c == '.' {
				c, _ = getc()
				if c == '.' {
					tok = makeToken(T_PUNCT, "...")
				} else {
					panic("invalid token '..'")
				}
			} else {
				ungetc()
				tok = makeToken(T_PUNCT, ".")
			}
		case '>':
			c, _ = getc()
			if c == '=' {
				tok = makeToken(T_PUNCT, ">=")
			} else if c == '>' {
				c, _ = getc()
				if c == '=' {
					tok = makeToken(T_PUNCT, ">>=")
				} else {
					ungetc()
					tok = makeToken(T_PUNCT, ">>")
				}
			} else {
				ungetc()
				tok = makeToken(T_PUNCT, ">")
			}
		case '<':
			c, _ = getc()
			if c == '-' {
				tok = makeToken(T_PUNCT, "<-")
			} else if c == '=' {
				tok = makeToken(T_PUNCT, "<=")
			} else if c == '<' {
				c, _ = getc()
				if c == '=' {
					tok = makeToken(T_PUNCT, "<<=")
				} else {
					ungetc()
					tok = makeToken(T_PUNCT, "<<")
				}
			} else {
				ungetc()
				tok = makeToken(T_PUNCT, "<")
			}
		default:
			fmt.Printf("c='%c'\n", c)
			panic("unknown char")
		}
		if debugToken {
			dumpToken(tok)
		}
		//r = append(r, tok)
	}

	return r
}

func f1() {
	path := "t/min/min.go"
	s := readFile(path)
	_bs := ByteStream{
		filename:  path,
		source:    s,
		nextIndex: 0,
		line:      1,
		column:    0,
	}
	bs := &_bs
	len1 := len(bs.source)

	fmt.Printf("%d\n", len1-64) // 1
	var c byte
	c, _ = bs.get()
	fmt.Printf("%d\n", c-'p'+2)        // 2
	fmt.Printf("%d\n", bs.nextIndex+2) // 3
	c, _ = bs.get()
	fmt.Printf("%d\n", c-'a'+4)        // 4
	fmt.Printf("%d\n", bs.nextIndex+3) // 5
}

func f2() {
	var chars []byte
	chars = append(chars, '7')
	chars = append(chars, '8')
	fmt.Printf("%d\n", len(chars) + 4) // 6
	fmt.Printf("%c\n", chars[0]) // 7
	fmt.Printf("%c\n", chars[1]) // 8
	fmt.Printf("9\n") // 9

	chars[0] = '1'
	chars[1] = '0'
	fmt.Printf("%s\n", chars) // 10
}

func f3() {
	path := "t/min/min.go"
	s := readFile(path)
	_bs := ByteStream{
		filename:  path,
		source:    s,
		nextIndex: 0,
		line:      1,
		column:    0,
	}
	bs = &_bs

	var c byte
	c, _ = bs.get()
	ident := readIdentifier(c)
	if ident == "package" {
		fmt.Printf("11\n")
	} else {
		println("error")
	}
}



func main() {
	f1()
	f2()
	f3()
}
