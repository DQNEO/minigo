package main

import (
	"os"
)

func f1() {
	bytes := bytes("0\n")
	os.Stdout.Write(bytes)
}

func main() {
	f1()
}
