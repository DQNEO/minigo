package main

import (
	"io/ioutil"
	"fmt"
)

func f1() {
	filename := "t/data/sample.txt"
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		panic("err")
	}
	if len(bytes) != 276 { // This is the size of the target file
		panic("Error: size does not match")
	}
}

func f2() {
	filename := "t/data/gen.go.txt"
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		panic("err")
	}
	fmt.Printf("%s", bytes)
	if len(bytes) != 83801 { // This is the size of the target file
		panic("Error: size does not match")
	}
}

func main() {
	f1()
	f2()
}
