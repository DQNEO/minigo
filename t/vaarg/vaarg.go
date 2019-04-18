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
	a = append(a, s)
	a = append(a, i)
	receiveVarg(format, a)
}

func f2() {
	var a []interface{}
	var s string = "hello"
	var i int = 123
	a = append(a, s)
	a = append(a, i)
	receiveIfcSlice(format, a)
}

func main() {
	f1()
	f2()
}
