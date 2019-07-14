package main


func receiveVarg(s string, a ...interface{}) {
	fmtPrintf("-\n")
	fmtPrintf("%d\n", len(a))
	fmtPrintf("%s\n", a[0])
	fmtPrintf("%d\n", a[1])
}

func receiveIfcSlice(s string, a []interface{}) {
	fmtPrintf("-\n")
	fmtPrintf("%d\n", len(a))
	fmtPrintf("%s\n", a[0])
	fmtPrintf("%d\n", a[1])
}

var format string = string("format-%s-%d\n")

func f1() {
	receiveVarg(format, "hello", 123)
}

func f2() {
	var a []interface{}
	a = append(a, "hello")
	a = append(a, 123)
	receiveIfcSlice(format, a)
}

func f3() {
	var a []interface{}
	a = append(a, "hello")
	a = append(a, 123)
	receiveVarg(format, a...)
}

func main() {
	f1()
	f2()
	f3()
}
