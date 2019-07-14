package main

import (
	"os"
	"strconv"
)

type identifier string

func fmtPrintf(format string, a... interface{}) {
	s := Sprintf(string(format), a...)
	os.Stdout.Write([]byte(s))
}

var _trash int
func Sprintf(format string, a... interface{}) string {
	var args []interface{}
	for _, x := range a {
		var y interface{}
		switch x.(type) {
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
	var blocks []string
	var bs []byte
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
			blocks = append(blocks, string(bs))
			bs = nil
			numPercent++
			continue
		}
		if inPercent {
			if c == '%' {
				bs = append(bs,c)
				inPercent = false
				continue
			}
			arg := args[argIndex]
			switch arg.(type) {
			case string:
				var _args string
				_args = arg.(string)
				blocks = append(blocks, _args)
			case []byte:
				var _arg []byte
				_arg = arg.([]byte)
				blocks = append(blocks, string(_arg))
			case byte:
				var _argByte byte
				_argByte = arg.(byte)
				bts := []byte{_argByte}
				g := string(bts)
				blocks = append(blocks, g)
			case int:
				var _argInt int
				_argInt = arg.(int)
				b := string(strconv.Itoa(_argInt))
				blocks = append(blocks, b)
			case bool: // "%v"
				var _argBool bool
				_argBool = arg.(bool)
				var b string
				if _argBool {
					b = "true"
				} else{
					b = "false"
				}
				blocks = append(blocks, b)
			default:
				panic("Unkown type to format:")
			}
			argIndex++
			inPercent = false
			bs = nil
			continue
		}
		bs = append(bs,c)
	}
	blocks = append(blocks, string(bs))
	var ss string
	for i, ss = range blocks {
		var bb []byte = []byte(ss)
		for j, c = range bb {
			r = append(r, c)
		}
	}
	_trash = i
	_trash = j
	return string(r)
}

func write(s []byte) {
	var b []byte = []byte(s)
	os.Stdout.Write(b)
}

func fmtPrintln(s string) {
	writeln([]byte(s))
}

func writeln(s []byte) {
	var b []byte = []byte(s)
	b = append(b, '\n')
	os.Stdout.Write(b)
}

func concat(as string, bs string) string {
	a := []byte(as)
	b := []byte(bs)

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
	a := []byte(as)
	b := []byte(bs)
	c := []byte(cs)
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
