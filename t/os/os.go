package main

import "os"

func f1() {
	var a = "hello stdout\n"
	os.Stdout.Write(a)

	var b = "hello stderr\n"
	os.Stderr.Write(b)
}

func main() {
	f1()
}
