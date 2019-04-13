package main

import "fmt"

func test_write() {
	s := "hello world\n"
	n := write(1, s, len(s))
	fmt.Printf("%d\n", n)
}

func f1() {
	test_write()
}

func main() {
	f1()
}
