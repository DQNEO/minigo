package fmt

import (
	"os"
	"github.com/DQNEO/minigo/stdlib/reflect"
	"github.com/DQNEO/minigo/stdlib/strconv"
	"github.com/DQNEO/minigo/stdlib/io"
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

	var buf string
	var tmpbuf []byte
	var inPercent bool
	var argIndex int

	for _, c := range []byte(format) {
		if inPercent {
			if c == '%' {
				tmpbuf = append(tmpbuf, c)
			} else {
				tmpbuf = nil
				arg := a[argIndex]
				argIndex++
				s := p.printArg(arg, c)
				buf = buf + s
			}
			inPercent = false
		} else {
			if  c == '%' {
				inPercent = true
				buf = buf + string(tmpbuf)
				tmpbuf = nil
			} else {
				tmpbuf = append(tmpbuf, c)
			}
		}
	}
	buf = buf + string(tmpbuf)

	p.buf = []byte(buf)
}

func (p *printer) printArg(arg interface{}, verb byte) string {

	switch verb {
	case 'T':
		return reflect.TypeOf(arg).String()
	case 'p':
		return "[==ADDRESS==]"
	}

	var s string
	switch arg.(type) {
	case string: // for %s
		s = arg.(string)
	case []byte: // for %s
		s = string(arg.([]byte))
	case byte: // for %c or %d
		switch verb {
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
