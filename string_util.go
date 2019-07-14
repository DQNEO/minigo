package main

import (
	"os"
	"strconv"
)

type bytes []byte

type identifier string

type goidentifier []byte

func S(s string) bytes {
	return bytes(s)
}

func fmtPrintf(format bytes, a... interface{}) {
	r := Sprintf(string(format), a...)
	write(bytes(r))
}
var _trash int
func Sprintf(format string, a... interface{}) string {
	var args []interface{}
	for _, x := range a {
		var y interface{}
		switch x.(type) {
		case bytes: // This case is not reached by 2nd gen compiler
			var tmpgostring bytes = x.(bytes)
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
	var blocks []bytes
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
				var _args string
				_args = arg.(string)
				blocks = append(blocks, bytes(_args))
			case []byte:
				var _arg []byte
				_arg = arg.([]byte)
				blocks = append(blocks, _arg)
			case byte:
				var _argByte byte
				_argByte = arg.(byte)
				bts := []byte{_argByte}
				g := bytes(bts)
				blocks = append(blocks, g)
			case int:
				var _argInt int
				_argInt = arg.(int)
				b := bytes(strconv.Itoa(_argInt))
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
				panic(S("Unkown type to format"))
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
	return string(r)
}

func write(s bytes) {
	var b []byte = []byte(s)
	os.Stdout.Write(b)
}

func fmtPrintln(s bytes) {
	writeln(s)
}

func writeln(s bytes) {
	var b []byte = []byte(s)
	b = append(b, '\n')
	os.Stdout.Write(b)
}

func concat(a bytes, b bytes) bytes {
	var r []byte
	for i:=0;i<len(a);i++ {
		r = append(r, a[i])
	}
	for i:=0;i<len(b);i++ {
		r = append(r, b[i])
	}
	return r
}

func concat3(a bytes, b bytes, c bytes) bytes {
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

func eq(a bytes, b bytes) bool {
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
