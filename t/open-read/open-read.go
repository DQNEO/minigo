package main

import (
	"os"
	"syscall"
)

import "fmt"

const MYBUFSIZ = 1024
const O_RDONLY = 0

var buf [1024]byte

func main() {
	var fd int
	var b []byte = buf[0:len(buf)]
	fname := os.Args[1]

	fd, _ = syscall.Open(fname, O_RDONLY, 0)
	ln, _ := syscall.Read(fd, b)
	s := buf[0:ln]
	fmt.Printf("%s", s)
}
