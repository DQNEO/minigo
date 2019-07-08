package main


func receiveVarg(s string, a ...interface{}) {
	fmtPrintf(S("-\n"))
	fmtPrintf(S("%d\n"), len(a))
	fmtPrintf(S("%s\n"), a[0])
	fmtPrintf(S("%d\n"), a[1])
}

func receiveIfcSlice(s string, a []interface{}) {
	fmtPrintf(S("-\n"))
	fmtPrintf(S("%d\n"), len(a))
	fmtPrintf(S("%s\n"), a[0])
	fmtPrintf(S("%d\n"), a[1])
}

var format string = "format-%s-%d\n"

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
