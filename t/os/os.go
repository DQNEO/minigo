package main

import "os"

func f1() {
	var a = "hello stdout\n"
	var a2 []byte = []byte(a)
	os.Stdout.Write(a2)

	var b = "hello stderr\n"
	var b2 []byte = []byte(b)
	os.Stderr.Write(b2)
}

func main() {
	f1()
}
