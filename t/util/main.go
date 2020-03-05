package main

import "fmt"
import "path"

func f1() {
	baseName := path.Base("foo/bar")
	fmt.Printf("%s\n", baseName) // bar
}

func main() {
	f1()
}
