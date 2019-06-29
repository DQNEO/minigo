package main

import (
	"fmt"
	"os"
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
	var b gostring
	b = GoSprintf2(gostring("hello\n"))
	os.Stdout.Write(b)

	b = GoSprintf2(gostring("%s\n"), []byte("world"))
	os.Stdout.Write(b)

	b = GoSprintf2(gostring("%s\n"), gostring("world"))
	os.Stdout.Write(b)

	b = GoSprintf2(gostring("left %s right\n"), gostring("center"))
	os.Stdout.Write(b)

	b = GoSprintf2(gostring("%s center right\n"), gostring("left"))
	os.Stdout.Write(b)

	b = GoSprintf2(gostring("left center %s"), gostring("right\n"))
	os.Stdout.Write(b)

	b = GoSprintf2(gostring("%s center %s\n"), gostring("left"), gostring("right"))
	os.Stdout.Write(b)

	var i int

	i = 123
	b = GoSprintf2(gostring("123=%d\n"), i)
	os.Stdout.Write(b)

	i = 4567
	b = GoSprintf2(gostring("%s=%d\n"), gostring("4567"), i)
	os.Stdout.Write(b)
}

func f2() {
	var b gostring
	b = GoSprintf2(gostring("push %%rax\n"))
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
