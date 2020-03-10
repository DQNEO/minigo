package runtime

// Stupid memory allocator.

const heapSize uintptr = 640485760

var heapHead uintptr
var heapPtr uintptr
var heapTail uintptr

func heapInit() {
	heapHead = brk(0)
	heapTail = brk(heapHead + heapSize)
	heapPtr = heapHead
}

func malloc(size uintptr) uintptr {
	if heapPtr+size > heapTail {
		panic([]byte("malloc exceeds heap capacity"))
		return 0
	}
	r := heapPtr
	heapPtr += size
	return r
}
