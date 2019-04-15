package main

import "fmt"

func stdout_write() {
	s := "hello world\n"
	n := write(1, s, len(s))
	fmt.Printf("%d\n", n)
}

func stderr_write() {
	s := "hello stderr\n"
	n := write(2, s, len(s))
	fmt.Printf("%d\n", n)
}

func f1() {
	stdout_write()
	stderr_write()
}

func main() {
	f1()
}
