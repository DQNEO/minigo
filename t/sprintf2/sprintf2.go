package main

import (
	"os"
)

func receiveSliceInVariadic(format []byte, a... interface{})  {
	var i0 interface{} = a[0]
	var b []byte
	var ok bool

	b, ok = i0.([]byte)

	fmtPrintf(S("ok=%d\n"), ok)
	fmtPrintf(S("b=%s,len=%d,cap=%d\n"), b,len(b), cap(b))
}

func f0() {
	receiveSliceInVariadic([]byte("%s\n"), []byte("abc"))
}

func f1() {
	var b bytes
	b = Sprintf(bytes("hello\n"))
	os.Stdout.Write(b)

	b = Sprintf(bytes("%s\n"), []byte("world"))
	os.Stdout.Write(b)

	b = Sprintf(bytes("%s\n"), bytes("world"))
	os.Stdout.Write(b)

	b = Sprintf(bytes("left %s right\n"), bytes("center"))
	os.Stdout.Write(b)

	b = Sprintf(bytes("%s center right\n"), bytes("left"))
	os.Stdout.Write(b)

	b = Sprintf(bytes("left center %s"), bytes("right\n"))
	os.Stdout.Write(b)

	b = Sprintf(bytes("%s center %s\n"), bytes("left"), bytes("right"))
	os.Stdout.Write(b)

	var i int

	i = 123
	b = Sprintf(bytes("123=%d\n"), i)
	os.Stdout.Write(b)

	i = 4567
	b = Sprintf(bytes("%s=%d\n"), bytes("4567"), i)
	os.Stdout.Write(b)
}

func f2() {
	var b bytes
	b = Sprintf(bytes("push %%rax\n"))
	os.Stdout.Write(b)
}


func main() {
	fmtPrintf(S("--- f0 ---\n"))
	f0()
	fmtPrintf(S("--- f1 ---\n"))
	f1()
	fmtPrintf(S("--- f2 ---\n"))
	f2()
}
