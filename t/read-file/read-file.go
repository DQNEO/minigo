package main

import (
	"io/ioutil"
)

func f1() {
	filename := "t/data/sample.txt"
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(S("err"))
	}
	if len(bytes) != 276 { // This is the size of the target file
		panic(S("Error: size does not match"))
	}
}

func f2() {
	filename := "t/data/gen.go.txt"
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(S("err"))
	}
	fmtPrintf(S("%s"), bytes)
	if len(bytes) != 83801 { // This is the size of the target file
		panic(S("Error: size does not match"))
	}
}

func main() {
	f1()
	f2()
}
