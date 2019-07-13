package main

import (
	"os"
)

const MYBUFSIZ = 1024
const O_RDONLY = 0

var buf [1024]byte

func main() {
	var fd int

	fname := os.Args[1]
	var filename *byte = fname
	fd = open(filename, O_RDONLY)
	ln := read(fd, buf, MYBUFSIZ)
	s := buf[0:ln]
	fmtPrintf(S("%s"), s)
}
