package main

import (
	"os"
)

type bytes []byte
type switchexpr string
type osarg string

type identifier string

type goidentifier bytes

func S(s string) bytes {
	return bytes(s)
}

func fmtPrintf(gos bytes, a... interface{}) {
	r := Sprintf(gos, a...)
	write(r)
}
var _trash int
func Sprintf(format []byte, a... interface{}) []byte {
	var args []interface{}
	for _, x := range a {
		var y interface{}
		switch x.(type) {
		case bytes: // This case is not reached by 2nd gen compiler
			var tmpgostring bytes = x.(bytes)
			var tmpbytes []byte = []byte(tmpgostring)
			y = tmpbytes
		case goidentifier:   // This case is not reached by 2nd gen compiler
			var tmpgoident goidentifier = x.(goidentifier)
			var tmpbytes2 []byte = []byte(tmpgoident)
			y = tmpbytes2
		default:
			y = x
		}
		args = append(args, y)
	}
	a = nil // unset

	var r []byte
	var blocks []bytes
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
			arg := args[argIndex]
			switch arg.(type) {
			case string:
				var _args string
				_args = arg.(string)
				blocks = append(blocks, bytes(_args))
			case []byte:
				var _arg []byte
				_arg = arg.([]byte)
				blocks = append(blocks, _arg)
			case byte:
				var _argByte byte
				_argByte = arg.(byte)
				bts := []byte{_argByte}
				g := bytes(bts)
				blocks = append(blocks, g)
			case int:
				var _argInt int
				_argInt = arg.(int)
				b := bytes(Itoa(_argInt))
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
				panic(S("Unkown type to format"))
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
	_trash = i
	_trash = j
	return r
}

func write(s bytes) {
	var b []byte = []byte(s)
	os.Stdout.Write(b)
}

func fmtPrintln(s bytes) {
	writeln(s)
}

func writeln(s bytes) {
	var b []byte = []byte(s)
	b = append(b, '\n')
	os.Stdout.Write(b)
}

func concat(a bytes, b bytes) bytes {
	var r []byte
	for i:=0;i<len(a);i++ {
		r = append(r, a[i])
	}
	for i:=0;i<len(b);i++ {
		r = append(r, b[i])
	}
	return r
}

func concat3(a bytes, b bytes, c bytes) bytes {
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

func eq(a bytes, b bytes) bool {
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

// "foo/bar", "/" => []bytes{"foo", "bar"}
func Split(s bytes, sep bytes) []bytes {
	if len(sep) > 1  {
		panic(S("no supported"))
	}
	seps := []byte(sep)
	sepchar := seps[0]
	vbytes := []byte(s)
	var buf []byte
	var r []bytes
	for _, b := range vbytes {
		if b == sepchar {
			r = append(r, bytes(buf))
			buf = nil
		} else {
			buf = append(buf, b)
		}
	}
	r = append(r, bytes(buf))

	return r
}

// Contains reports whether substr is within s.
func Contains(s bytes, substr bytes) bool {
	return Index(s, substr) >= 0
}

func Index(s bytes, substr bytes) int {
	bytes := []byte(s)
	bsub := []byte(substr)
	var in bool
	var subIndex int
	var r int = -1 // not found
	for i, b := range bytes {
		if !in && b == bsub[0] {
			in = true
			r = i
			subIndex = 0
		}

		if in {
			if b == bsub[subIndex] {
				if subIndex == len(bsub) - 1 {
					return r
				}
			} else {
				in = false
				r = -1
				subIndex = 0
			}
		}
	}

	return -1
}

func HasSuffix(s bytes, suffix bytes) bool {
	if len(s) >= len(suffix) {
		var low int =  len(s)-len(suffix)
		var lensb int = len(s)
		var suf []byte
		sb := []byte(s)
		suf = sb[low:lensb]  // lensb is required
		return eq(bytes(suf) ,suffix)
	}
	return false
}

func Itoa(i int) []byte {
	var r []byte
	var tmp []byte
	var isMinus bool

	// open(2) returs  0xffffffff 4294967295 on error.
	// I don't understand this yet.
	if i > 2147483648 {
		i = i - 2147483648*2
	}


	if i < 0 {
		i = i * -1
		isMinus = true
	}
	for i>0 {
		mod := i % 10
		tmp = append(tmp, byte('0') + byte(mod))
		i = i /10
	}

	if isMinus {
		r = append(r, '-')
	}

	for j:=len(tmp)-1;j>=0;j--{
		r = append(r, tmp[j])
	}

	if len(r) == 0 {
		return []byte{'0'}
	}
	return r
}

func Atoi(gs bytes) (int, error) {
	if len(gs) == 0 {
		return 0,nil
	}
	var b byte
	var n int
	var i int
	var isMinus bool
	for i, b = range gs {
		if b == '.' {
			return 0,nil // @FIXME all no number should return error
		}
		if b == '-' {
			isMinus = true
			continue
		}
		var x byte = b - byte('0')
		n  = n * 10
		n = n + int(x)
	}
	if isMinus {
		n = -n
	}
	_trash = i
	return n, nil
}
