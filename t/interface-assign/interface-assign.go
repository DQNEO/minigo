package main

import "fmt"

func f0() *int {
	p := &Point{
		x:1,
		y:2,
	}
	x := &p.x
	return x
}

func f1() string {
	f := &StmtFor{
		rng: &ForRangeClause{
			invisibleMapCounter:&ExprVariable{
				id:1,
			},
		},
	}
	mapCounter := &Relation{
		name: "2",
		expr: f.rng.invisibleMapCounter,
	}
	return mapCounter.name
}

func main() {
	p := f0()
	fmt.Printf("%d\n", *p)
	s := f1()
	fmt.Printf("%s\n", s)
}

type Point struct {
	x int
	y int
}

type StmtFor struct {
	rng   *ForRangeClause
}

type ForRangeClause struct {
	invisibleMapCounter *ExprVariable
}

type ExprVariable struct {
	id int
}

func (e *ExprVariable) f() {
}

type Relation struct {
	name string
	expr Expr
}

type Expr interface {
	f()
}
