package fmt

import (
	"strconv"
	"os"
)

func Println(s string) {
	var b []byte = []byte(s)
	b = append(b, '\n')
	os.Stdout.Write(b)
}

var _fmt_trash int
func Sprintf(format string, args... interface{}) string {
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
	_fmt_trash = i
	_fmt_trash = j
	return string(r)
}
