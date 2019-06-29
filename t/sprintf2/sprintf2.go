package main

import (
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


func main() {
	f0()
}
