package runtime

var heap [1048576]int
var heapIndex int

func malloc(size int) int {
	if heapIndex == 0 {
		heapIndex = (heap + 0)
	}
	if heapIndex + size - heap > len(heap) {
		return 0
	}
	r := heapIndex
	heapIndex += size
	return r
}

func makeSlice(newLen int, newCap int) []int {
	var r []int
	r = heap[heapIndex:newLen:heapIndex+newCap]
	heapIndex += newCap * 8
	return r
}
