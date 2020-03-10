package runtime

// Stupid memory allocator.

const heapSize uintptr = 640485760

var heapHead uintptr
var heapCurrent uintptr
var heapTail uintptr

func heapInit() {
	heapHead = brk(0)
	heapTail = brk(heapHead + heapSize)
	heapCurrent = heapHead
}

func malloc(size uintptr) uintptr {
	if heapCurrent+size > heapTail {
		panic([]byte("malloc exceeds heap capacity"))
		return 0
	}
	r := heapCurrent
	heapCurrent += size
	return r
}
