package main


func typeswtch(flg bool) int {
	var r int
	var p *Point = &Point{
		a: 1,
		b: 2,
	}

	var p2 *Point2 = &Point2{
		c: 2,
		d: 3,
	}

	var i myInterface
	if flg {
		i = p
	} else {
		i = p2
	}
	switch i.(type) {
	case *Point:
		r = 1
	case *Point2:
		r = 2
	default:
		r = 3
	}

	return r
}

func f3() {
	var x int
	x = typeswtch(true)
	fmtPrintf(S("%d\n"), x) // 1

	x = typeswtch(false)
	fmtPrintf(S("%d\n"), x) // 2
}

func f4() {
	var i myInterface

	switch i.(type) {
	case *Point:
		fmtPrintf(S("ERROR\n"))
	default:
		fmtPrintf(S("3\n"))
	}

	switch i.(type) {
	case nil:
		fmtPrintf(S("4\n"))
	default:
		fmtPrintf(S("ERROR\n"))
	}

	switch i.(type) {
	case nil:
		fmtPrintf(S("5\n"))
	default:
		fmtPrintf(S("ERROR"))
	}

}

type myInterface interface {
	sum() int
}

type Point struct {
	a int
	b int
}

type Point2 struct {
	c int
	d int
}

func (p *Point) sum() int {
	return p.a + p.b
}

func (p *Point2) sum() int {
	return p.c + p.d
}

func main() {
	f3()
	f4()
}
