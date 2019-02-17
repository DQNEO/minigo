package main

import "fmt"

func f1() {
	var methods map[int][]string = map[int][]string{} // typeId : []methods
	methods[1] = []string{"1"}
	fmt.Printf("%s\n", methods[1][0])
}

func main() {
	f1()
}
