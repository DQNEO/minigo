package main


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
	fmtPrintf("%d\n", l) // 1

	l = f2()
	fmtPrintf("%d\n", l+1) // 2
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
	fmtPrintf("%d\n", j) // 3
}

func f4() {
	var array [3]Stmt
	var ifc Stmt
	var ifc2 Stmt

	var n *Number

	n = &Number{
		id: 4,
	}
	ifc = n
	fmtPrintf("%d\n", ifc.getId()) // 4

	array[1] = ifc
	fmtPrintf("%d\n", array[1].getId()+1) // 5

	ifc2 = array[1]
	fmtPrintf("%d\n", ifc2.getId()+2) // 6
}

func f5() {
	var stmts []Stmt
	var s Stmt
	var s2 Stmt
	var n *Number
	var id int
	n = &Number{
		id: 7,
	}
	s = n
	fmtPrintf("%d\n", s.getId()) // 7
	stmts = append(stmts, s)
	fmtPrintf("%d\n", len(stmts)+7) // 8
	s2 = stmts[0]
	id = s2.getId()
	fmtPrintf("%d\n", id+2)               // 9
	fmtPrintf("%d\n", stmts[0].getId()+3) // 10
}

func main() {
	test_array_literal()
	f3()
	f4()
	f5()
}
