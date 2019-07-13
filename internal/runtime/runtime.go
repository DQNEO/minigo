package iruntime

var libcArgc int
var libcArgv *byte

var heap [640485760]byte
var heapTail *int

const intSize = 8

func init() {
	// set head address of heap
	heapTail = heap
}

func malloc(size int) *int {
	if heapTail+ size > len(heap) + heap  {
		panic([]byte("malloc exceeds heap capacity"))
		return 0
	}
	r := heapTail
	heapTail += size
	return r
}

// This is a copy from stconv
func itoa(i int) []byte {
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

func panic(msg []byte) {
	printstring([]byte("panic: "))
	printstring(msg)
	printstring([]byte("\n"))
	exit(1)
}

func reportMemoryUsage() {
	printstring([]byte("# memory-usage: "))
	i := getMemoryUsage()
	s := itoa(i)
	printstring([]byte(s))
	printstring([]byte("\n"))
}

func getMemoryUsage() int {
	return heapTail - heap
}

func append1(x []byte, elm byte) []byte {
	var z []byte
	xlen := len(x)
	zlen := xlen + 1

	if cap(x) >= zlen {
		z = x[:zlen]
	} else {
		var newcap int
		if xlen == 0 {
			newcap = 1
		} else {
			newcap = xlen * 2
		}
		z = makeSlice(zlen, newcap, 1)
		for i:=0;i<xlen;i++ {
			z[i] = x[i]
		}
	}

	z[xlen] = elm
	return z
}

func append8(x []int, elm int) []int {
	var z []int
	xlen := len(x)
	zlen := xlen + 1

	if cap(x) >= zlen {
		z = x[:zlen]
	} else {
		var newcap int
		if xlen == 0 {
			newcap = 1
		} else {
			newcap = xlen * 2
		}
 		z = makeSlice(zlen, newcap, 8)
		for i:=0;i<xlen;i++ {
			z[i] = x[i]
		}
	}

	z[xlen] = elm
	return z
}


func append24(x []interface{}, elm interface{}) []interface{} {
	//dumpInterface(elm)
	var z []interface{}
	xlen := len(x)
	zlen := xlen + 1

	if cap(x) >= zlen {
		z = x[:zlen]
	} else {
		var newcap int
		if xlen == 0 {
			newcap = 1
		} else {
			newcap = xlen * 2
		}
		z = makeSlice(zlen, newcap, 24)
		for i:=0;i<xlen;i++ {
			z[i] = x[i]
		}
	}

	z[xlen] = elm
	return z
}

func eqGostringInternal(a []byte, b []byte) bool {
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

func eqGostrings(a []byte, b []byte, eq bool) bool {
	var ret bool
	ret = eqGostringInternal(a,b)
	if eq {
		return ret
	} else {
		return !ret
	}
}

const MiniGo int = 1
