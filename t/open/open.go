package main

import "fmt"

const O_RDONLY = 0

func f1() int {
	var fd int
	fd = open("/etc/hosts", O_RDONLY)
	return fd
}

func f2() int {
	var fd int
	fd = open("/etc/foobar", O_RDONLY)
	return fd
}

func main() {
	var fd int
	fd = f1()
	fmt.Printf("%d\n", fd) // 3

	fd = f2()
	fmt.Printf("%d\n", fd) // -1
}
