package main

import "fmt"

func f1() int {
	a := &A{}
	var lefts []Expr = []Expr{a}
	return len(lefts)
}


func f2() int {
	var lefts []Expr = []Expr{&A{}}
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
	var l int
	l = f1()
	fmt.Printf("%d\n", l) // 1

	l = f2()
	fmt.Printf("%d\n", l + 1) // 2
}
