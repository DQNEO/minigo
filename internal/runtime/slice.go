package runtime

type slice struct {
	ptr uintptr
	ln int
	cap int
}

func makeSlice(elmSize int, len int, cap int) slice

func copySlice1(src []byte, dst []byte) {
	for i := 0; i < len(src); i++ {
		dst[i] = src[i]
	}
}

func copySlice8(src []int, dst []int) {
	for i := 0; i < len(src); i++ {
		dst[i] = src[i]
	}
}

func copySlice24(src []interface{}, dst []interface{}) {
	for i := 0; i < len(src); i++ {
		dst[i] = src[i]
	}
}

func copySlice(size int, src []interface{}, dst []interface{}) {
	switch size {
	case 1:
		copySlice1(src, dst)
	case 8:
		copySlice8(src, dst)
	case 24:
		copySlice24(src, dst)
	}
}

func append1(x []byte, elm byte) slice {
	var z []byte
	xlen := len(x)
	zlen := xlen + 1

	if cap(x) >= zlen {
		z = x[:zlen]
	} else {
		var newcap int
		if xlen == 0 {
			newcap = 1
		} else {
			newcap = xlen * 2
		}
		z = makeSlice(1, zlen, newcap)
		copySlice(1, x,z)
	}

	z[xlen] = elm
	return z
}

func append8(x []int, elm int) slice {
	var z []int
	xlen := len(x)
	zlen := xlen + 1

	if cap(x) >= zlen {
		z = x[:zlen]
	} else {
		var newcap int
		if xlen == 0 {
			newcap = 1
		} else {
			newcap = xlen * 2
		}
		z = makeSlice(8, zlen, newcap)
		copySlice(8, x,z)
	}

	z[xlen] = elm
	return z
}

func append24(x []interface{}, elm interface{}) slice {
	//dumpInterface(elm)
	var z []interface{}
	xlen := len(x)
	zlen := xlen + 1

	if cap(x) >= zlen {
		z = x[:zlen]
	} else {
		var newcap int
		if xlen == 0 {
			newcap = 1
		} else {
			newcap = xlen * 2
		}
		z = makeSlice(24, zlen, newcap)
		copySlice(24, x,z)
	}

	z[xlen] = elm
	return slice(z)
}
