package main

import (
	"os"
	"strconv"
)

type gostring []byte
type cstring string

type identifier string
type packageName identifier

type goidentifier gostring

func S(s string) gostring {
	return gostring(s)
}

var trash int
func Sprintf(format gostring, a... interface{}) gostring {
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
	for i,c = range f {
		//fmt.Printf("# c=%c\n", c)
		if !inPercent && c == '%' {
			inPercent = true
			//fmt.Printf("# in Percent \n")
			blocks = append(blocks, str)
			str = nil
			numPercent++
			continue
		}
		if inPercent {
			if c == '%' {
				str = append(str,c)
				inPercent = false
				//fmt.Printf("# outof Percent \n")
				continue
			}
			//fmt.Printf("# check arg for c=%c\n", c)
			arg := a[argIndex]
			//dumpInterface(arg)
			switch arg.(type) {
			case string:
				panic("string should not be passed to Sprintf")
			case []byte:
				//fmt.Printf("# byte\n")
				var _arg []byte
				_arg = arg.([]byte)
				blocks = append(blocks, _arg)
			case gostring:  // This case is not reached by 2nd gen compiler
				// fmt.Printf("# gostring\n")
				var _arg []byte
				_arg = arg.(gostring)
				blocks = append(blocks, _arg)
			case int:
				var _argInt int
				_argInt = arg.(int)
				var s string
				s = strconv.Itoa(_argInt)
				b := []byte(s)
				blocks = append(blocks, b)
			}
			argIndex++
			inPercent = false
			str = nil
			continue
		}
		str = append(str,c)
	}
	//fmt.Printf("# blocks:%v", blocks)
	blocks = append(blocks, str)
	for i, str = range blocks {
		for j, c = range str {
			r = append(r, c)
		}
	}
	trash = i
	trash = j
	return r
}

func write(s gostring) {
	var b []byte = []byte(s)
	os.Stdout.Write(b)
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

func eq(a gostring, b cstring) bool {
	return eqGostrings(a, gostring(b))
}

func eqGostrings(a gostring, b gostring) bool {
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

func eqCstring(a cstring, b cstring) bool {
	return a == b
}


