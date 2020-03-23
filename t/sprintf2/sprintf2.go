package main

import (
	"fmt"
	"os"
)

func receiveSliceInVariadic(a ...interface{}) {
	var i0 interface{} = a[0]
	var b []byte
	var ok bool

	b, ok = i0.([]byte)

	fmt.Printf("ok=%v\n", ok)
	fmt.Printf("b=%s,len=%d,cap=%d\n", b, len(b), cap(b))
}

func f0() {
	receiveSliceInVariadic([]byte("abc"))
}

func Write(s string) {
	os.Stdout.Write([]byte(s))
}
func f1() {
	var b string
	b = fmt.Sprintf("hello\n")
	Write(b)

	b = fmt.Sprintf("%s\n", "world")
	Write(b)

	b = fmt.Sprintf("%s\n", "world")
	Write(b)

	b = fmt.Sprintf("left %s right\n", "center")
	Write(b)

	b = fmt.Sprintf("%s center right\n", "left")
	Write(b)

	b = fmt.Sprintf("left center %s", "right\n")
	Write(b)

	b = fmt.Sprintf("%s center %s\n", "left", "right")
	Write(b)

	var i int

	i = 123
	b = fmt.Sprintf("123=%d\n", i)
	Write(b)

	i = 4567
	b = fmt.Sprintf("%s=%d\n", "4567", i)
	Write(b)
}

func f2() {
	var b string
	b = fmt.Sprintf("pushq %%rax\n")
	Write(b)
}

func main() {
	fmt.Printf("--- f0 ---\n")
	f0()
	fmt.Printf("--- f1 ---\n")
	f1()
	fmt.Printf("--- f2 ---\n")
	f2()
}
