package main

import (
	"os"
	"strconv"
)

type gostring []byte
type cstring string

type identifier string

type goidentifier gostring

func S(s string) gostring {
	return gostring(s)
}

func fmtPrintf(gos gostring, a... interface{}) {
	r := Sprintf(gos, a...)
	write(r)
}
var _trash int
func Sprintf(format []byte, a... interface{}) []byte {
	var args []interface{}
	for _, x := range a {
		var y interface{}
		switch x.(type) {
		case gostring:      // This case is not reached by 2nd gen compiler
			var tmpgostring gostring = x.(gostring)
			var tmpbytes []byte = []byte(tmpgostring)
			y = tmpbytes
		case goidentifier:   // This case is not reached by 2nd gen compiler
			var tmpgoident goidentifier = x.(goidentifier)
			var tmpbytes2 []byte = []byte(tmpgoident)
			y = tmpbytes2
		default:
			y = x
		}
		args = append(args, y)
	}
	a = nil // unset

	var r []byte
	var blocks []gostring
	var str []byte
	var f []byte = []byte(format)
	var c byte
	var i int
	var j int
	var numPercent int
	var inPercent bool
	var argIndex int
	//var sign byte
	for i, c = range f {
		if ! inPercent && c == '%' {
			inPercent = true
			blocks = append(blocks, str)
			str = nil
			numPercent++
			continue
		}
		if inPercent {
			if c == '%' {
				str = append(str,c)
				inPercent = false
				continue
			}
			arg := args[argIndex]
			switch arg.(type) {
			case string:
				var s string
				var bytes []byte
				s = arg.(string)
				bytes = []byte(s)
				blocks = append(blocks, bytes)
			case []byte:
				var _arg []byte
				_arg = arg.([]byte)
				blocks = append(blocks, _arg)
			case byte:
				var _argByte byte
				_argByte = arg.(byte)
				bts := []byte{_argByte}
				g := gostring(bts)
				blocks = append(blocks, g)
			case int:
				var _argInt int
				_argInt = arg.(int)
				b := gostring(strconv.Itoa(_argInt))
				blocks = append(blocks, b)
			case bool: // "%v"
				var _argBool bool
				_argBool = arg.(bool)
				var b []byte
				if _argBool {
					b = []byte("true")
				} else{
					b = []byte("false")
				}
				blocks = append(blocks, b)
			default:
				panic("Unkown type to format")
			}
			argIndex++
			inPercent = false
			str = nil
			continue
		}
		str = append(str,c)
	}
	blocks = append(blocks, str)
	for i, str = range blocks {
		for j, c = range str {
			r = append(r, c)
		}
	}
	_trash = i
	_trash = j
	return r
}

func write(s gostring) {
	var b []byte = []byte(s)
	os.Stdout.Write(b)
}

func fmtPrintln(s gostring) {
	writeln(s)
}

func writeln(s gostring) {
	var b []byte = []byte(s)
	b = append(b, '\n')
	os.Stdout.Write(b)
}

func convertCstringsToGostrings(cstrings []string) []gostring {
	var r []gostring
	for _, cs := range cstrings {
		r = append(r, gostring(cs))
	}

	return r
}

func concat(a gostring, b gostring) gostring {
	var r []byte
	for i:=0;i<len(a);i++ {
		r = append(r, a[i])
	}
	for i:=0;i<len(b);i++ {
		r = append(r, b[i])
	}
	return r
}

func concat3(a gostring, b gostring, c gostring) gostring {
	var r []byte
	for i:=0;i<len(a);i++ {
		r = append(r, a[i])
	}
	for i:=0;i<len(b);i++ {
		r = append(r, b[i])
	}
	for i:=0;i<len(c);i++ {
		r = append(r, c[i])
	}
	return r
}

func eq(a gostring, b gostring) bool {
	if len(a) != len(b) {
		return false
	}

	for i:=0;i<len(a);i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
