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

/*
var methods [16]int = [16]int{}

func findmethod() int {
	return methods[1]
}
*/

