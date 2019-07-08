package main


func f1() {
	var p *Point
	p = &Point{
		x: 1,
		y: 2,
	}
	sum := p.sum()
	fmtPrintf(S("%d\n"), sum-2) // 1
}

func f2() {
	var myInterface MyInterface
	ptr := &Point{
		x: 2,
		y: 3,
	}
	myInterface = ptr
	sum := myInterface.sum()
	fmtPrintf(S("%d\n"), sum-3) // 2

	var myInterface2 MyInterface
	myInterface2 = myInterface
	diff := myInterface2.diff()
	fmtPrintf(S("%d\n"), diff+2) // 3
}

func f3() {
	var myInterface MyInterface
	ptr := &Asset{
		money: 2,
		stock: 3,
	}
	myInterface = ptr
	sum := myInterface.sum()
	fmtPrintf(S("%d\n"), sum-1) // 4

	diff := myInterface.diff()
	fmtPrintf(S("%d\n"), diff+4) // 5
}

func f4(bol bool) {
	var myInterface MyInterface
	point := &Point{
		x: 2,
		y: 4,
	}

	asset := &Asset{
		money: 2,
		stock: 6,
	}

	if bol {
		myInterface = point
	} else {
		myInterface = asset
	}

	sum := myInterface.sum()

	fmtPrintf(S("%d\n"), sum) // 6, 8

	diff := myInterface.diff()
	fmtPrintf(S("%d\n"), diff+5) // 7, 9
}

var gpoint = Point{
	x: 6,
	y: 4,
}

var gptr *Point

func return_interface() MyInterface {
	var myInterface MyInterface
	gptr = &gpoint
	myInterface = gptr
	sum := myInterface.sum()

	fmtPrintf(S("%d\n"), sum) // 10
	return myInterface
}

func f5() {
	var myif MyInterface = return_interface()
	fmtPrintf(S("%d\n"), myif.sum()+1) // 11
}

func main() {
	f1()
	f2()
	f3()
	f4(true)
	f4(false)
	f5()
}

type MyInterface interface {
	sum() int
	diff() int
}

type Point struct {
	x int
	y int
}

func (p *Point) sum() int {
	return p.x + p.y
}

func (p *Point) diff() int {
	return p.y - p.x
}

type Asset struct {
	money int
	stock int
}

func (a *Asset) sum() int {
	return a.money + a.stock
}

func (a *Asset) diff() int {
	return a.stock - a.money
}
