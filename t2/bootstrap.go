package runtime

var heap [10485760]byte
var heapTail *int

func malloc(size int) *int {
	if heapTail+ size > len(heap) + heap  {
		panic("malloc exceeds heap capacity")
		return 0
	}
	r := heapTail
	heapTail += size
	return r
}
