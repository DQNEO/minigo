package main


var GlobalInt int = 1
var GlobalPtr *int = &GlobalInt

var pp1 *Point = &Point{x: 1, y: 2}
var pp2 *Point = &Point{x: 3, y: 4}

type Point struct {
	x int
	y int
}

func f1() {
	fmtPrintf("%d\n", *GlobalPtr) // 1
}

func f2() {
	fmtPrintf("%d\n", pp1.y) // 2
	fmtPrintf("%d\n", pp2.x) // 3
}

type Gtype struct {
	typ  int
	size int
}

var gIntE = Gtype{typ: 7, size: 8}
var gInt = &gIntE

type DeclFunc struct {
	tok      string
	rettypes []*Gtype
}

var builtinLenGlobal = &DeclFunc{
	tok:      "tok",
	rettypes: []*Gtype{&gIntE, &gIntE},
}

func f3() {
	retTypes := builtinLenGlobal.rettypes
	fmtPrintf("%d\n", len(retTypes)+2) // 4
	var gi *Gtype = retTypes[0]
	fmtPrintf("%d\n", gi.typ-2) // 5
}

/*
func f4() {
	var builtinLenLocal = &DeclFunc{
		tok:"tok",
		rettypes: []*Gtype{&gIntE,&gIntE},
	}

	retTypes := builtinLenLocal.rettypes
	fmtPrintf("%d\n", len(retTypes) + 4) // 6
	var gi *Gtype = retTypes[0]
	fmtPrintf("%d\n", gi.size - 1) // 7
}
*/

func main() {
	f1()
	f2()
	f3()
	//f4()
}
