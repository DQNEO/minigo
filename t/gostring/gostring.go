package main

import "fmt"

type gostring []byte

func f1() {
	var gs gostring = gostring("hi")
	dumpSlice(gs)
}

func main() {
	f1()
}

