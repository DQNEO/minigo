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
	path := "/etc/hosts"
	s := readFile(path)
	bs := ByteStream{
		filename:  path,
		source:    s,
	}
	bsp := &bs

	len1 := len(bs.source)
	len2 := len(bsp.source)

	fmt.Printf("len1=%d\n", len1)
	fmt.Printf("len2=%d\n", len2)
}

func main() {
	f1()
}
