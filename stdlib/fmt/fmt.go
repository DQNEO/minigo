package fmt

import (
	"os"
	"strconv"
)

type gostring []byte

// used only in tests
func Printf(format string, a ...interface{}) {
	b := Sprintf(gostring(format), a...)
	os.Stdout.Write(b)
}

var trash int
func Sprintf(format gostring, a... interface{}) []byte {
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
				var s string
				var bytes []byte
				s = arg.(string)
				bytes = []byte(s)
				blocks = append(blocks, bytes)
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
			case byte:
				var _argByte byte
				_argByte = arg.(byte)
				bts := []byte{_argByte}
				g := gostring(bts)
				blocks = append(blocks, bts)
			case int:
				var _argInt int
				_argInt = arg.(int)
				b := gostring(strconv.Itoa(_argInt))
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
