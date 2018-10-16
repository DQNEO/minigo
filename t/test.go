package main

var ga int

func f0() {

}

func fa() {
	printf("%d\n", ga)// => 0
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
	printf("%d\n", 2 - 1)// this is a comment
	printf("%d\n", 1 + 1)
	printf("%d\n", 1 + 1 + 1) // this is another comment //
	printf("%d\n", 2 * 2)
	printf("%d\n", 2 * 3 - 1)
	printf("%d\n", 1 + 1 * 5)
}

func fc() {
	var i int
	i = 3
	printf("%d\n", i + 4)
}

func fd() {
	var j int = 2
	printf("%d\n", j * 4)
}

func fe() {
	var a int = 5
	var b int = 4
	printf("%d\n", a + b)
}

var gb int = 10

func ff() {
	printf("%d\n", gb)
}

func fg(a int, b int) {
	printf("%d\n",  a + b)
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
	println("hello world")
}

