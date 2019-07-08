package main


/*
type myinterface interface  {
	x()
	y()
}

type mystruct struct {
	x int
	y int
}
*/
var ga int

type myint int
type mymyint myint
type bool int
type mybool bool

func f0() {
}

func fa() {
	fmtPrintf(S("%d\n"), ga) // => 0
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
	fmtPrintf(S("%d\n"), 1) // this is a comment
	fmtPrintf(S("%d\n"), 4-2)
	fmtPrintf(S("%d\n"), 1+1+1) // this is another comment //
	fmtPrintf(S("%d\n"), 1*2+2)
	fmtPrintf(S("%d\n"), 2*3-1)
	fmtPrintf(S("%d\n"), 9-1-2)
}

func fc() {
	var i int
	i = 7
	fmtPrintf(S("%d\n"), i)
}

func fd() {
	var j int = 2
	fmtPrintf(S("%d\n"), j*4)
}

func fe() {
	var a int = 5
	var b int = 4
	fmtPrintf(S("%d\n"), a+b)
}

func ff() {
	//fmtPrintf(S("%v%v\n"), true, false)
}

func fg(a int, b int) {
	fmtPrintf(S("%d\n"), a+b)
}

var gc int
var gd int = 10
var ge = 2

func fh() {
	fmtPrintf(S("%d\n"), gc+gd+ge)
}

const c0 int = 1
const c1 = 2

func fi() {
	const c1 int = 3
	const c2 = 9
	fmtPrintf(S("%d\n"), c0+c1+c2)
}

var gb int = 14

func f14() {
	fmtPrintf(S("%d\n"), gb)
}

var garbage int

func f15() {
	var sum int
	var i int
	var a int
	var b int
	/*
		for a, b := range []int{1,2,3} {
			sum = a + b
		}
	*/
	i = sum + a + b
	sum = i
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
	f15()
	fmtPrintf(S("hello world\n"))
}
