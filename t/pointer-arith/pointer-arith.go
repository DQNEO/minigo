package main

import (
	"fmt"
	"unsafe"
)

func f0() {
	var c byte = '3'
	var d byte = '4'
	cp := &c
	*cp = '1'
	fmt.Printf("%c\n", c)
	fmt.Printf("%c\n", d)
}

func f1() {
	buf := []byte("123456789")
	addr := &buf[0]
	fmt.Printf("%s\n", buf)
	p := uintptr(unsafe.Pointer(addr))
	p = p + 2
	up := unsafe.Pointer(p)
	var pb *byte
	pb = (*byte)(up)
	*pb = '1'
	fmt.Printf("%s\n", buf)
}

func main() {
	f0()
	f1()
}
