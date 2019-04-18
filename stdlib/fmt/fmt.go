package fmt

func Printf(format string, param interface{}) {
	printf(format, param)
}

func Sprintf(format string, param []interface{}) string {
	return doPrintf(format, param)
}

func Println(s string) {
}

var pbuf [1024]byte

func doPrintf(format string, a []interface{}) string {
	var a0 interface{}
	var a1 interface{}
	var a2 interface{}
	var a3 interface{}
	var numred int
	switch len(a) {
	case 0:
		numred = sprintf(pbuf, format)
	case 1:
		a0 = a[0]
		numred = sprintf(pbuf, format, *a0)

	case 2:
		a0 = a[0]
		a1 = a[1]
		numred = sprintf(pbuf, format, *a0, *a1)
	case 3:
		a0 = a[0]
		a1 = a[1]
		a2 = a[2]
		numred = sprintf(pbuf, format, *a0, *a1, *a2)
	case 4:
		a0 = a[0]
		a1 = a[1]
		a2 = a[2]
		a3 = a[3]
		numred = sprintf(pbuf, format, *a0, *a1, *a2, *a3)
	default:
		printf("ERROR: doPrintf cannot handle more than 4 params")
	}

	// copy string to heap area
	var b []slice
	b = makeSlice(numred+1, numred+1)
	strcopy(pbuf, b, numred)

	// return heap
	return b
}
