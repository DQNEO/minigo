package main

func receiveVarg(s string, a ...interface{}) {
	puts("-")
	println(len(a))
	println(*a[0])
	println(*a[1])
}

func receiveIfcSlice(s string, a []interface{}) {
	puts("-")
	println(len(a))
	println(*a[0])
	println(*a[1])
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

func main() {
	f1()
	f2()
}
