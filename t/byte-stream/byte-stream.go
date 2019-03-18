package main

import (
	"fmt"
	"io/ioutil"
)

type ByteStream struct {
	filename  string
	source    []byte
}

func readFile(filename string) []byte {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return bytes
}

func f1() {
	path := "t/min/min.go"
	s := readFile(path)
	bs := ByteStream{
		filename:  path,
		source:    s,
	}
	bsp := &bs

	len1 := len(bs.source)
	len2 := len(bsp.source)

	fmt.Printf("%d\n", len1 - 64) // 1
	fmt.Printf("%d\n", len2 - 63) // 2
}

func main() {
	f1()
}
