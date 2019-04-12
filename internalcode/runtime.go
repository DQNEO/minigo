package runtime

var heap [10485760]int
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
	if heapIndex == 0 {
		heapIndex = (heap + 0)
	}
	low := (heapIndex - heap) / intSize
	r = heap[low:low+newLen:low+newCap]
	heapIndex += newCap * intSize
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
		z = makeSlice(zlen, newcap)
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
 		z = makeSlice(zlen, newcap)
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
		z = makeSlice(zlen, newcap)
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
