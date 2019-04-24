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
