package main

import "fmt"

type myinterface interface  {
	x()
	y()
}

type mystruct struct {
	x int
	y int
}

var ga int

type myint int
type mymyint myint
type bool int
type mybool bool

func f0() {
}

func fa() {
	fmt.Printf("%d\n", ga)// => 0
}


/* this is
a
block
  comment
* /
*
/

*/

func fb() {
	fmt.Printf("%d\n", 2 - 1)// this is a comment
	fmt.Printf("%d\n", 1 + 1)
	fmt.Printf("%d\n", 1 + 1 + 1) // this is another comment //
	fmt.Printf("%d\n", 2 * 2)
	fmt.Printf("%d\n", 2 * 3 - 1)
	fmt.Printf("%d\n", 1 + 1 * 5)
}

func fc() {
	var i int
	i = 3
	fmt.Printf("%d\n", i + 4)
}

func fd() {
	var j int = 2
	fmt.Printf("%d\n", j * 4)
}

func fe() {
	var a int = 5
	var b int = 4
	fmt.Printf("%d\n", a + b)
}

func ff() {
	fmt.Printf("%d%d\n", true, false)
}

func fg(a int, b int) {
	fmt.Printf("%d\n",  a + b)
}

var gc int
var gd int = 10
var ge = 2

func fh() {
	fmt.Printf("%d\n", gc + gd + ge)
}

const c0 int = 1
const c1 = 2

func fi() {
	const c1 int = 3
	const c2  = 9
	fmt.Printf("%d\n", c0 + c1 + c2)
}

var gb int = 14

func f14() {
	fmt.Printf("%d\n", gb)
}

func main() {
	f0()
	fa()
	fb()
	fc()
	fd()
	fe()
	ff()
	fg(5, 6)
	fh()
	fi()
	f14()
	println("hello world")
}

