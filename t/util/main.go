package main

import "fmt"

func f1() {
	baseName := getBaseNameFromImport([]byte("foo/bar"))
	fmt.Printf("%s\n", baseName) // bar
}

func main() {
	f1()
}
