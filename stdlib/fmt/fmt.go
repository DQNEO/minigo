package fmt

import (
	"os"
	"strconv"
)

type gostring []byte

// used only in tests
func Println(format string) {
	b := gostring(format)
	b = append(b, '\n')
	os.Stdout.Write(b)
}

// used only in tests
func Printf(format string, a ...interface{}) {
	b := Sprintf(gostring(format), a...)
	os.Stdout.Write(b)
}

var _fmt_trash int
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
			arg := a[argIndex]
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
			case gostring:  // This case is not reached by 2nd gen compiler
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
				dumpInterface(arg)
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
	_fmt_trash = i
	_fmt_trash = j
	return r
}
