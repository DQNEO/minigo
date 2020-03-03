package main

import (
	"fmt"
	"os"
)

func f1() {
	dir := "."
	f, err := os.Open(dir)
	if err != nil {
		panic(err)
	}
	names, err := f.Readdirnames(-1)
	if err != nil {
		panic(err)
	}
	for _, n := range names {
		fmt.Printf("%s\n", n)
	}
}

func main() {
	f1()
}
