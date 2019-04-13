package main

func doPrintf(format string, a []interface{}) {
	switch len(a) {
	case 0:
			printf(format)
	case 1:
			var a0 interface{} = a[0]
			printf(format, *a0)
	case 2:
			var a0 interface{} = a[0]
			var a1 interface{} = a[1]
			printf(format, *a0, *a1)
	case 3:
			var a0 interface{} = a[0]
			var a1 interface{} = a[1]
			var a2 interface{} = a[2]
			printf(format, *a0, *a1, *a2)
	case 4:
		var a0 interface{} = a[0]
		var a1 interface{} = a[1]
		var a2 interface{} = a[2]
		var a3 interface{} = a[3]
		printf(format, *a0, *a1, *a2, *a3)
	default:
		printf("ERROR: doPrintf cannot handle more than 4 params")
	}
}

func f0() {
	var a []interface{}
	doPrintf("hello\n", a)
}


func f1() {
	var a []interface{}
	var i int  = 123
	var ifc interface{}
	ifc = i
	a = append(a, ifc)
	doPrintf("%d\n", a)
}

func f2() {
	var a []interface{}
	var i int = 123
	var i2 int  = 456
	var ifc interface{}
	var ifc2 interface{}
	ifc = i
	ifc2 = i2
	a = nil
	a = append(a, ifc)
	a = append(a, ifc2)
	doPrintf("%d %d\n", a)
}

func f3() {
	var a []interface{}
	var s string  = "hello"
	var s2 string  = "world"
	var ifc interface{}
	var ifc2 interface{}
	ifc = s
	ifc2 = s2
	a = append(a, ifc)
	a = append(a, ifc2)
	doPrintf("%s %s\n", a)
}

func f4() {
	var a []interface{}
	var s string  = "hello"
	var i int = 123
	var ifc interface{}
	var ifc2 interface{}
	ifc = s
	ifc2 = i
	a = append(a, ifc)
	a = append(a, ifc2)
	doPrintf("%s %d\n", a)
}

func f5() {
	var a []interface{}
	var s string  = "hello"
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
	doPrintf("%s %d %d\n", a)
}

func main() {
	f0()
	f1()
	f2()
	f3()
	f4()
	f5()
}
