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

func test_array_literal() {
	var l int
	l = f1()
	fmt.Printf("%d\n", l) // 1

	l = f2()
	fmt.Printf("%d\n", l+1) // 2
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

func f4() {
	var array [3]Stmt
	var ifc Stmt
	var ifc2 Stmt

	var n *Number

	n = &Number{
		id:4,
	}
	ifc = n
	fmt.Printf("%d\n", ifc.getId()) // 4

	array[1] = ifc
	fmt.Printf("%d\n", array[1].getId() + 1) // 5

	ifc2 = array[1]
	fmt.Printf("%d\n", ifc2.getId() + 2) // 6
}


/*
func f5() {
	var stmts []Stmt
	var s Stmt
	var n *Number
	n = &Number{
		id:123,
	}
	s = n
	fmt.Printf("%d\n", s.getId())

	stmts = append(stmts, s)
}
*/

func main() {
	test_array_literal()
	f3()
	f4()
}
