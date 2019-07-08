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
	var b gostring
	b = Sprintf(gostring("hello\n"))
	os.Stdout.Write(b)

	b = Sprintf(gostring("%s\n"), []byte("world"))
	os.Stdout.Write(b)

	b = Sprintf(gostring("%s\n"), gostring("world"))
	os.Stdout.Write(b)

	b = Sprintf(gostring("left %s right\n"), gostring("center"))
	os.Stdout.Write(b)

	b = Sprintf(gostring("%s center right\n"), gostring("left"))
	os.Stdout.Write(b)

	b = Sprintf(gostring("left center %s"), gostring("right\n"))
	os.Stdout.Write(b)

	b = Sprintf(gostring("%s center %s\n"), gostring("left"), gostring("right"))
	os.Stdout.Write(b)

	var i int

	i = 123
	b = Sprintf(gostring("123=%d\n"), i)
	os.Stdout.Write(b)

	i = 4567
	b = Sprintf(gostring("%s=%d\n"), gostring("4567"), i)
	os.Stdout.Write(b)
}

func f2() {
	var b gostring
	b = Sprintf(gostring("push %%rax\n"))
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
