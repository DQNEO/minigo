package main

import (
	"fmt"
	"os"
)

const MYBUFSIZ = 1024
const O_RDONLY = 0

func main() {
	var fd int
	var buf [1024]byte

	fname := os.Args[1]
	fd = open(fname, O_RDONLY)
	read(fd, buf, MYBUFSIZ)
	fmt.Printf("%s", buf)
}
