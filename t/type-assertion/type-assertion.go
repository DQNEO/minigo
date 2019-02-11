package main

import "fmt"

func do(flg bool) int {
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

	dp, ok := i.(*Point)
	if ok {
		return dp.calc()
	}

	dp2, ok := i.(*Point2)
	if ok {
		return dp2.calc()
	}
	return 0
}

type myInterface interface {
	calc() int
}

type Point struct {
	a int
	b int
}

type Point2 struct {
	c int
	d int
}

func (p *Point) calc() int {
	return p.a + p.b
}

func (p *Point2) calc() int {
	return p.c * p.d
}

func main() {
	var x int
	x = do(true)
	fmt.Printf("%d\n", x - 2) // 1

	x = do(false)
	fmt.Printf("%d\n", x - 4) // 2
}
