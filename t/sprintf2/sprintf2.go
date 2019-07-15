package main

import (
	"os"
	"fmt"
)

func receiveSliceInVariadic(format []byte, a... interface{})  {
	var i0 interface{} = a[0]
	var b []byte
	var ok bool

	b, ok = i0.([]byte)

	fmt.Printf("ok=%d\n", ok)
	fmt.Printf("b=%s,len=%d,cap=%d\n", b,len(b), cap(b))
}

func f0() {
	receiveSliceInVariadic([]byte("%s\n"), []byte("abc"))
}

func f1() {
	var b []byte
	b = Sprintf([]byte("hello\n"))
	os.Stdout.Write(b)

	b = Sprintf([]byte("%s\n"), []byte("world"))
	os.Stdout.Write(b)

	b = Sprintf([]byte("%s\n"), []byte("world"))
	os.Stdout.Write(b)

	b = Sprintf([]byte("left %s right\n"), []byte("center"))
	os.Stdout.Write(b)

	b = Sprintf([]byte("%s center right\n"), []byte("left"))
	os.Stdout.Write(b)

	b = Sprintf([]byte("left center %s"), []byte("right\n"))
	os.Stdout.Write(b)

	b = Sprintf([]byte("%s center %s\n"), []byte("left"), []byte("right"))
	os.Stdout.Write(b)

	var i int

	i = 123
	b = Sprintf([]byte("123=%d\n"), i)
	os.Stdout.Write(b)

	i = 4567
	b = Sprintf([]byte("%s=%d\n"), []byte("4567"), i)
	os.Stdout.Write(b)
}

func f2() {
	var b []byte
	b = Sprintf([]byte("push %%rax\n"))
	os.Stdout.Write(b)
}


func main() {
	fmt.Printf("--- f0 ---\n")
	f0()
	fmt.Printf("--- f1 ---\n")
	f1()
	fmt.Printf("--- f2 ---\n")
	f2()
}
