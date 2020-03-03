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
	var plainstring []byte
	var numPercent int
	var inPercent bool
	var argIndex int

	//var sign byte
	bformat := []byte(format)
	for _, c := range bformat {
		if !inPercent  {
			if  c == '%' {
				inPercent = true
				blocks = append(blocks, string(plainstring))
				plainstring = nil
				numPercent++
			} else {
				plainstring = append(plainstring, c)
			}
		} else if inPercent {
			if c == '%' {
				plainstring = append(plainstring, c)
				inPercent = false
				continue
			}
			inPercent = false
			plainstring = nil
			var s string
			arg := args[argIndex]
			argIndex++
			switch arg.(type) {
			case string:
				s = arg.(string)
			case []byte:
				bf := arg.([]byte)
				s = string(bf)
			case byte:
				b := arg.(byte)
				bf := []byte{b}
				s = string(bf)
			case int:
				var _int int
				_int = arg.(int)
				s = strconv.Itoa(_int)
			case bool: // "%v"
				if arg.(bool) {
					s = "true"
				} else {
					s = "false"
				}
			default:
				panic("Unkown type to format:")
			}
			blocks = append(blocks, s)
		}
	}
	blocks = append(blocks, string(plainstring))

	var r []byte
	for _, block := range blocks {
		bf := []byte(block)
		for _, c := range bf {
			r = append(r, c)
		}
	}

	return string(r)
}
