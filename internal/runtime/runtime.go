package iruntime

var runtimeArgc int
var runtimeArgv *int
var Args []string //  Do not remove this. Actually this is real os.Args.

var heap [640485760]byte
var heapTail *int

const intSize = 8

func init() {
	// set head address of heap
	heapTail = heap + 0
}

func malloc(size int) *int {
	if heapTail+ size > len(heap) + heap  {
		panic("malloc exceeds heap capacity")
		return 0
	}
	r := heapTail
	heapTail += size
	return r
}

func panic(msg string) {
	printf("panic: %s\n", msg)
	exit(1)
}

func reportMemoryUsage() {
	printf("# memory-usage %d\n", getMemoryUsage())
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

func strcopy(src string, dest string, slen int) string {
	for i:=0; i < slen ; i++ {
		dest[i] = src[i]
	}
	dest[slen] = 0
	return dest
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

func strcat(a []byte, b []byte) string {
	var c []byte
	for i:=0;i<len(a);i++ {
		c = append(c, a[i])
	}
	for i:=0;i<len(b);i++ {
		c = append(c, b[i])
	}
	return string(c)
}

const MiniGo int = 1
