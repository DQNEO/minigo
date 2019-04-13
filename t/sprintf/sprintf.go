package main

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

func myPrintf(format string, a []interface{}) {
	var s string = doPrintf(format, a)
	printf(s)
}

func f0() {
	var a []interface{}
	myPrintf("hello\n", a)
}

func f1() {
	var a []interface{}
	var i int = 123
	var ifc interface{}
	ifc = i
	a = append(a, ifc)
	myPrintf("%d\n", a)
}

func f2() {
	var a []interface{}
	var i int = 123
	var i2 int = 456
	var ifc interface{}
	var ifc2 interface{}
	ifc = i
	ifc2 = i2
	a = nil
	a = append(a, ifc)
	a = append(a, ifc2)
	myPrintf("%d %d\n", a)
}

func f3() {
	var a []interface{}
	var s string = "hello"
	var s2 string = "world"
	var ifc interface{}
	var ifc2 interface{}
	ifc = s
	ifc2 = s2
	a = append(a, ifc)
	a = append(a, ifc2)
	myPrintf("%s %s\n", a)
}

func f4() {
	var a []interface{}
	var s string = "hello"
	var i int = 123
	var ifc interface{}
	var ifc2 interface{}
	ifc = s
	ifc2 = i
	a = append(a, ifc)
	a = append(a, ifc2)
	myPrintf("%s %d\n", a)
}

func f5() {
	var a []interface{}
	var s string = "hello"
	var i int = 123
	var i2 int = 456
	var ifc interface{}
	var ifc2 interface{}
	var ifc3 interface{}
	ifc = s
	ifc2 = i
	ifc3 = i2
	a = append(a, ifc)
	a = append(a, ifc2)
	a = append(a, ifc3)
	myPrintf("%s %d %d\n", a)
}

func main() {
	f0()
	f1()
	f2()
	f3()
	f4()
	f5()
}
