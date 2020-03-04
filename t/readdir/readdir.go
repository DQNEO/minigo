package main

import (
	"fmt"
	"os"
)

func f1() {
	dir := "t/data"
	f, err := os.Open(dir)
	if err != nil {
		panic(err)
	}
	names, err := f.Readdirnames(-1)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%d\n", len(names))
}

func main() {
	f1()
}
