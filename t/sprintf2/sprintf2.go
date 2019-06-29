package main

import (
	"fmt"
	"os"
	"strconv"
)

func receiveSliceInVariadic(format []byte, a... interface{})  {
	var i0 interface{} = a[0]
	var b []byte
	var ok bool

	b, ok = i0.([]byte)

	fmt.Printf("ok=%d\n", ok)
	fmt.Printf("b=%s,len=%d,cap=%d\n", b,len(b), cap(b))
}

func f0() {
	receiveSliceInVariadic([]byte("%s\n"), []byte("abc"))
}

var trash int
func GoSprintf2(format gostring, a... interface{}) gostring {
	var r []byte
	var blocks [][]byte
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
		if c == '%' {
			inPercent = true
			blocks = append(blocks, str)
			str = nil
			numPercent++
			continue
		}
		if inPercent {
			//sign = c
			arg := a[argIndex]
			//dumpInterface(arg)
			switch arg.(type) {
			case []byte:
				var _arg []byte
				_arg = arg.([]byte)
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

func f1() {
	var b gostring
	b = GoSprintf2(gostring("hello\n"))
	os.Stdout.Write(b)

	b = GoSprintf2(gostring("%s\n"), gostring("world"))
	os.Stdout.Write(b)

	b = GoSprintf2(gostring("left %s right\n"), gostring("center"))
	os.Stdout.Write(b)

	b = GoSprintf2(gostring("%s center right\n"), gostring("left"))
	os.Stdout.Write(b)

	b = GoSprintf2(gostring("left center %s"), gostring("right\n"))
	os.Stdout.Write(b)

	b = GoSprintf2(gostring("%s center %s\n"), gostring("left"), gostring("right"))
	os.Stdout.Write(b)

	var i int

	i = 123
	b = GoSprintf2(gostring("123=%d\n"), i)
	os.Stdout.Write(b)

	i = 4567
	b = GoSprintf2(gostring("%s=%d\n"), gostring("4567"), i)
	os.Stdout.Write(b)
}




func main() {
	fmt.Printf("--- f0 ---\n")
	f0()
	fmt.Printf("--- f1 ---\n")
	f1()
}
