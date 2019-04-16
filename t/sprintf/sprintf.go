package main

import "fmt"


func myPrintf(format string, a []interface{}) {
	var s string = fmt.Sprintf(format, a)
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

/*
func dumpToken(tok *Token) {
	s := fmt.Sprintf("tok: type=%-8s, sval=\"%s\"", tok.typ, tok.sval)
	debugf(s)
}
*/

func test_dumpToken() {
	format := "tok: type=%-4s,sval=\"%s\"\n"
	var s1 string = "int"
	var s2 string = "123"
	var ifcs []interface{}
	var ifc1 interface{} = s1
	var ifc2 interface{} = s2
	ifcs = append(ifcs, ifc1)
	ifcs = append(ifcs, ifc2)
	var s string =  fmt.Sprintf(format, ifcs)
	fmt.Printf(s)
}

func main() {
	f0()
	f1()
	f2()
	f3()
	f4()
	f5()
	test_dumpToken()
}
