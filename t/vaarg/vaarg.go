package main

import "fmt"

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
	var a []interface{}
	var s string = "hello"
	var i int = 123
	var ifc interface{} = s
	var ifc2 interface{} = i
	a = append(a, ifc)
	a = append(a, ifc2)
	receiveVarg(format, a)
}

func f2() {
	var a []interface{}
	var s string = "hello"
	var i int = 123
	var ifc interface{} = s
	var ifc2 interface{} = i
	a = append(a, ifc)
	a = append(a, ifc2)
	receiveIfcSlice(format, a)
}

func main() {
	f1()
	f2()
}
