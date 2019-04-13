package main

import "fmt"

func test_write() {
	s := "hello world\n"
	n := write(1, s, len(s))
	fmt.Printf("%d\n", n)
}

func stderr_write() {
	s := "hello stderr\n"
	write(2, s, len(s))
}

func f1() {
	test_write()
	stderr_write()
}

func main() {
	f1()
}
