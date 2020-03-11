package main

import (
	"fmt"
	"syscall"
)

const O_RDONLY = 0

func f1() int {
	var fd int
	fd, _ = syscall.Open("t/min/min.go", O_RDONLY, 0)
	return fd
}

func f2() int {
	var fd int
	fd, _ = syscall.Open("/var/noexists.txt", O_RDONLY, 0)
	return fd
}

func main() {
	var fd int
	fd = f1()
	fmt.Printf("%d\n", fd) // 3

	fd = f2()
	if fd < 0 {
		fmt.Printf("error ok\n")
	}
}
