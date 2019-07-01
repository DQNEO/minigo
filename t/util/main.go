package main

import (
	"fmt"
)

type gostring []byte

func f1() {
	baseName := getBaseNameFromImport(gostring("foo/bar"))
	fmt.Printf("%s\n", string(baseName)) // bar
}

func main() {
	f1()
}
