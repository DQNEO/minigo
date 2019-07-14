package main

import (
	"os"
	"strconv"
)

type bytes []byte

type identifier string

func fmtPrintf(format string, a... interface{}) {
	s := Sprintf(string(format), a...)
	os.Stdout.Write(bytes(s))
}

var _trash int
func Sprintf(format string, a... interface{}) string {
	var args []interface{}
	for _, x := range a {
		var y interface{}
		switch x.(type) {
		case bytes: // This case is not reached by 2nd gen compiler
			var tmpb bytes = x.(bytes)
			var tmpbytes []byte = []byte(tmpb)
			y = tmpbytes
		case identifier:   // This case is not reached by 2nd gen compiler
			var tmpgoident identifier = x.(identifier)
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
				panic("Unkown type to format:")
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

func fmtPrintln(s string) {
	writeln([]byte(s))
}

func writeln(s bytes) {
	var b []byte = []byte(s)
	b = append(b, '\n')
	os.Stdout.Write(b)
}

func concat(as string, bs string) string {
	a := bytes(as)
	b := bytes(bs)

	var r []byte
	for i:=0;i<len(a);i++ {
		r = append(r, a[i])
	}
	for i:=0;i<len(b);i++ {
		r = append(r, b[i])
	}
	return string(r)
}

func concat3(as string, bs string, cs string) string {
	a := bytes(as)
	b := bytes(bs)
	c := bytes(cs)
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
	return string(r)
}

func eq(as string, bs string) bool {
	if len(as) != len(bs) {
		return false
	}

	a := []byte(as)
	b := []byte(bs)
	for i:=0;i<len(a);i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
