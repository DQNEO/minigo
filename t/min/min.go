package main

import (
	"os"
)

func f1() {
	bytes := gostring("0\n")
	os.Stdout.Write(bytes)
}

func main() {
	f1()
}
