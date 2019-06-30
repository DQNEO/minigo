package main

import "fmt"

func receiveVarg(s string, a ...interface{}) {
	fmt.Printf("-\n")
	fmt.Printf("%d\n", len(a))
	fmt.Printf("%s\n", a[0])
	fmt.Printf("%d\n", a[1])
}

func receiveIfcSlice(s string, a []interface{}) {
	fmt.Printf("-\n")
	fmt.Printf("%d\n", len(a))
	fmt.Printf("%s\n", a[0])
	fmt.Printf("%d\n", a[1])
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
