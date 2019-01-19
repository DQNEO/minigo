package main

import "fmt"

const MYBUFSIZ = 1024
const O_RDONLY = 0

func main() {
	var fd int
	var buf [1024]byte

	fd = open("/etc/lsb-release", O_RDONLY)
	read(fd, buf, MYBUFSIZ)
	fmt.Printf("%s", buf)
}
