package main

import "fmt"

func f1() int {
	a := &A{}
	var lefts []Expr = []Expr{a}
	return len(lefts)
}

type Expr interface {
	emit()
}

type A struct {
	i int
}

func (a *A) emit() {
}

func main() {
	l := f1()
	fmt.Printf("%d\n", l)
}
