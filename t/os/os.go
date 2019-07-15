package main

import (
	"fmt"
	"os"
)

func f0() {
	var f *os.File
	f = os.Stderr
	fmt.Printf("%d\n", f.innerFile.fd.Sysfd)
}

func f1() {
	var a = "hello stdout\n"
	var a2 []byte = []byte(a)
	os.Stdout.Write(a2)

	var b = "hello stderr\n"
	var b2 []byte = []byte(b)
	os.Stderr.Write(b2)
}

func main() {
	f0()
	f1()
}
