package main

import "fmt"

type Expr interface {
	emit() int
}

type A struct {
	i int
}

func (a *A) emit() int {
	return a.i
}

func return_ifc() Expr {
	a := &A{i:7}
	return a
}

func f1() {
	x := return_ifc()
	fmt.Printf("%d\n", x.emit() - 6)
}


/*
func f2() int {
	var lefts []Expr = []Expr{&A{}}
	return len(lefts)
}

type Number struct {
	id int
}

func (n *Number) getId() int {
	return n.id
}

type Stmt interface {
	getId() int
}

func f3() {
	var array [3]interface{}
	var i int = 3
	var ifc interface{}
	ifc = i
	array[1] = ifc

	var j int
	j = ifc.(int)
	fmt.Printf("%d\n", j) // 3
}
*/

func main() {
	f1()
}
