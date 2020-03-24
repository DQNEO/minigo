package runtime

type slice struct {
	ptr uintptr
	ln int
	cap int
}

func makeSlice(elmSize int, len int, cap int) slice

func append1(x []byte, elm byte) []byte {
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
		for i := 0; i < xlen; i++ {
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
			newcap = 1
		} else {
			newcap = xlen * 2
		}
		z = makeSlice(8, zlen, newcap)
		for i := 0; i < xlen; i++ {
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
			newcap = 1
		} else {
			newcap = xlen * 2
		}
		z = makeSlice(24, zlen, newcap)
		for i := 0; i < xlen; i++ {
			z[i] = x[i]
		}
	}

	z[xlen] = elm
	return z
}
