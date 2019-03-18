package main

import (
	"fmt"
	"io/ioutil"
)


func f1() {
	filename := "/etc/resolv.conf"
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		panic("err")
	}
	fmt.Printf("%s", bytes)
	if len(bytes) != 68 { // 68 is the size of the target file
		panic("Error: size does not match")
	}
}

func main() {
	f1()
}
