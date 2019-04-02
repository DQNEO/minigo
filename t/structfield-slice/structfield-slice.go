package main

import "fmt"

type Gtype struct {
	fields      []*Gtype
}

func f1() {
	var gtype *Gtype = &Gtype{}
	var x []*Gtype
	x = gtype.fields
	fmt.Printf("%d\n", len(x))
}

func main() {
	f1()
}
