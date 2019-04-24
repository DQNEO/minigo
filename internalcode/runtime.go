package runtime

var runtimeArgc int
var runtimeArgv *int
var Args []string //  Do not remove this. Actually this is real os.Args.

var heap [10485760]byte
var heapTail *int

var heapSlice1 [10485760]byte
var heapSlice1Tail *byte

var heapSlice8 [10485760]int
var heapSlice8Tail *int

var heapSlice24 [10485760]interface{}
var heapSlice24Tail *interface{}
const intSize = 8

func init() {
	// set head address of heap
	heapTail = heap + 0
	heapSlice1Tail = heapSlice1 + 0
	heapSlice8Tail = heapSlice8 + 0
	heapSlice24Tail = heapSlice24 + 0
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

func makeSlice1(newLen int, newCap int) []byte {
	var r []byte
	low := (heapSlice1Tail - heapSlice1) + 1
	r = heapSlice1[low:low+newLen:low+newCap]
	heapSlice1Tail += (newCap + 1)
	return r
}

func makeSlice8(newLen int, newCap int) []int {
	var r []int
	low := (heapSlice8Tail- heapSlice8) / intSize + 1
	r = heapSlice8[low:low+newLen:low+newCap]
	heapSlice8Tail += (newCap + 1)* intSize
	return r
}

func makeSlice24(newLen int, newCap int) []interface{} {
	var r []interface{}
	low := (heapSlice24Tail- heapSlice24) / 24 + 1
	r = heapSlice24[low:low+newLen:low+newCap]
	heapSlice24Tail += (newCap + 1)* 24
	return r
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
			newcap = 8
		} else {
			newcap = xlen * 2
		}
		z = makeSlice1(zlen, newcap)
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
			newcap = 8
		} else {
			newcap = xlen * 2
		}
 		z = makeSlice8(zlen, newcap)
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
			newcap = 8
		} else {
			newcap = xlen * 2
		}
		z = makeSlice24(zlen, newcap)
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

const MiniGo int = 1

// builin functions
// https://golang.org/ref/spec#Predeclared_identifiers

// Functions:
//	append cap close complex copy delete imag len
//	make new panic print println real recover

func make(x interface{}) interface{} {
}

func panic(x interface{}) {
	switch x.(type) {
	case string:
		s := x.(string)
		printf("panic:%s\n", s)
	default:
		printf("panic:\n")
	}
	exit(1)
}

func println(s interface{}) {
	if s < 4096 {
		// regard it as int
		printf("%d\n", s)
	} else {
		printf("%s\n", s)
	}
}

func print(x interface{}) {
	printf(x)
}

func recover() interface{} {
	return nil
}

type error interface {
	Error() string
}
