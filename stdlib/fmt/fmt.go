package fmt

import (
	"os"
	"strconv"
)

func Println(s string) {
	var b []byte = []byte(s)
	b = append(b, '\n')
	os.Stdout.Write(b)
}

func Printf(format string, a ...interface{}) {
	s := Sprintf(string(format), a...)
	os.Stdout.Write([]byte(s))
}

var _fmt_trash int

func Sprintf(format string, args ...interface{}) string {
	var blocks []string
	var bs []byte
	var f []byte = []byte(format)
	var numPercent int
	var inPercent bool
	var argIndex int

	//var sign byte
	for _, c := range f {
		if !inPercent && c == '%' {
			inPercent = true
			blocks = append(blocks, string(bs))
			bs = nil
			numPercent++
			continue
		}
		if inPercent {
			if c == '%' {
				bs = append(bs, c)
				inPercent = false
				continue
			}
			arg := args[argIndex]
			switch arg.(type) {
			case string:
				s := arg.(string)
				blocks = append(blocks, s)
			case []byte:
				bf := arg.([]byte)
				blocks = append(blocks, string(bf))
			case byte:
				b := arg.(byte)
				bf := []byte{b}
				blocks = append(blocks,  string(bf))
			case int:
				var _int int
				_int = arg.(int)
				b := string(strconv.Itoa(_int))
				blocks = append(blocks, b)
			case bool: // "%v"
				var s string
				if arg.(bool) {
					s = "true"
				} else {
					s = "false"
				}
				blocks = append(blocks, s)
			default:
				panic("Unkown type to format:")
			}
			argIndex++
			inPercent = false
			bs = nil
			continue
		}
		bs = append(bs, c)
	}
	blocks = append(blocks, string(bs))

	var r []byte
	for _, block := range blocks {
		bf := []byte(block)
		for _, c := range bf {
			r = append(r, c)
		}
	}

	return string(r)
}
