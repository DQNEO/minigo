package fmt

import (
	"os"
	"strconv"
	"io"
)

func Print(s string) {
	os.Stdout.Write([]byte(s))
}

func Println(s string) {
	b := []byte(s)
	b = append(b, '\n')
	os.Stdout.Write(b)
}

func Fprintf(w io.Writer, format string, a ...interface{}) {
	p := newPrinter()
	p.doPrintf(format, a)
	w.Write(p.buf)
	p.free()
}

func Sprintf(format string, a ...interface{}) string {
	p := newPrinter()
	p.doPrintf(format, a)
	s := string(p.buf)
	p.free()
	return s
}

func Printf(format string, a ...interface{}) {
	Fprintf(os.Stdout, format, a...)
}

type printer struct {
	buf []byte
}

func newPrinter() *printer {
	p := &printer{}
	return p
}

func (p *printer) free() {
	p.buf = p.buf[0:0]
}

func (p *printer) doPrintf(format string, a []interface{}) {

	var blocks []string
	var plainstring []byte
	var numPercent int
	var inPercent bool
	var argIndex int

	//var sign byte
	for _, c := range []byte(format) {
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
				inPercent = false
				plainstring = append(plainstring, c)
				continue
			}
			inPercent = false
			plainstring = nil
			arg := a[argIndex]
			argIndex++
			s := p.printArg(arg, c)
			blocks = append(blocks, s)
		}
	}
	blocks = append(blocks, string(plainstring))

	var r []byte
	for _, block := range blocks {
		for _, c := range []byte(block) {
			r = append(r, c)
		}
	}
	p.buf = r
}

func (p *printer) printArg(arg interface{}, c byte) string {
	var s string
	switch arg.(type) {
	case string: // for %s
		s = arg.(string)
	case []byte: // for %s
		s = string(arg.([]byte))
	case byte: // for %c
		switch c {
		case 'c':
			b := arg.(byte)
			s = string([]byte{b})
		case 'd':
			b := arg.(byte)
			i := int(b)
			s = strconv.Itoa(i)
		default:
			panic("unknown format flag")
		}
	case int: // for %d
		s = strconv.Itoa(arg.(int))
	case uint16: // for %d
		s = strconv.Itoa(int(arg.(uint16)))
	case bool: // for %v
		if arg.(bool) {
			s = "true"
		} else {
			s = "false"
		}
	default:
		//panic(fmt.Sprintf("%T\n", arg))
		panic("Unkown type to format:")
	}
	return s
}
