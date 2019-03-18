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
}

func main() {
	f1()
}
