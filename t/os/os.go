package main

import "os"

func f1() {
	var a = S("hello stdout\n")
	var a2 []byte = []byte(a)
	os.Stdout.Write(a2)

	var b = S("hello stderr\n")
	var b2 []byte = []byte(b)
	os.Stderr.Write(b2)
}

func main() {
	f1()
}
