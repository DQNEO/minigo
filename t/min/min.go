package main

import (
	"os"
)

func f1() {
	s := "0\n"
	os.Stdout.Write([]byte(s))
}

func main() {
	f1()
}
