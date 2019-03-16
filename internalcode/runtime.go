package runtime

var heap [1048576]int
var heapIndex int
const intSize = 8

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
	heapIndex += newCap * intSize
	return r
}

func append(x []int, elm int) []int {
	var z []int
	zlen := len(x) + 1
	if cap(x) >= zlen {
		z = x[:zlen]
	} else {
		newcap := len(x) * 2
		z = makeSlice(zlen, newcap)
		for i:=0;i<len(x);i++ {
			z[i] = x[i]
		}
	}

	z[len(x)] = elm
	return z
}

func strcopy(src string, dest string, slen int) string {
	for i:=0; i < slen ; i++ {
		dest[i] = src[i]
	}
	dest[slen] = 0
	return dest
}
