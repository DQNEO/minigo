package main

import "fmt"

func multi(a int, b int, c int, d int) (int, int, int, int) {
	return a, b, c, d
}

func multireverse(a int, b int, c int, d int) (int, int, int, int) {
	return d, c, b, a
}

func getMulti() {
	var i int = 0
	var j int = 0
	var k int = 0
	var l int = 0
	i, j, k, l = multi(1, 2, 3, 4)
	fmt.Printf("%d\n", i)
	fmt.Printf("%d\n", j)
	fmt.Printf("%d\n", k)
	fmt.Printf("%d\n", l)

	a, b, c, d := multireverse(8, 7, 6, 5)
	fmt.Printf("%d\n", a)
	fmt.Printf("%d\n", b)
	fmt.Printf("%d\n", c)
	fmt.Printf("%d\n", d)
}

var bbytes = [3]byte{'a', 'b', 'c'}

func ReadFile() ([]byte, int) {
	return bbytes[:], 9
}

func getReadFile() {
	b, i := ReadFile()
	fmt.Printf("%d\n", i)
	fmt.Printf("%c\n", b[0])
	fmt.Printf("%c\n", b[1])
	fmt.Printf("%c\n", b[2])
}

func main() {
	getMulti()
	getReadFile()
}
