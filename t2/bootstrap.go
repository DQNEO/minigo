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

