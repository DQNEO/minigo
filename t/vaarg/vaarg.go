package main

import "fmt"
import "os"

func f0() {
	debugf("hello debug with vaargs")
}

func receiveVarg(s string, a ...interface{}) {
	println("-")
	println(len(a))
	println(*a[0])
	println(*a[1])
}

func receiveIfcSlice(s string, a []interface{}) {
	println("-")
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

func f3() {
	var a []interface{}
	a = append(a, "hello")
	a = append(a, 123)
	receiveVarg(format, a...)
}

func debugf(format string, v ...interface{}) {
	var indents []byte

	var format2 string = string(indents) + format + "\n"
	s2 := fmt.Sprintf(format2, v)
	var b []byte = []byte(s2)
	os.Stdout.Write(b)
}

func main() {
	f0()
	f1()
	f2()
	f3()
}
