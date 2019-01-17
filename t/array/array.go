package main

import "fmt"

var garray = [3]int{1, 2, 0}

func f1() {
	garray[2] = 3
	fmt.Printf("%d\n", garray[0])
	fmt.Printf("%d\n", garray[1])
	fmt.Printf("%d\n", garray[2])
}

func f2() {
	var i4 = 4
	var i5 = 5
	var i0 = 0
	var i7 = 7

	var larray = [4]int{i4, i5, i0, i7}
	larray[2] = 6
	fmt.Printf("%d\n", larray[0])
	fmt.Printf("%d\n", larray[1])
	fmt.Printf("%d\n", larray[2])
	fmt.Printf("%d\n", larray[3])
	larray[3]++
	fmt.Printf("%d\n", larray[3])
	larray[3] = 10
	larray[3]--
	fmt.Printf("%d\n", larray[3])
}

func f3() {
	var lbytes = [6]byte{'?', 'e', 'l', 'l', 'o', 10}
	lbytes[0] = 'H'
	fmt.Printf("%c", lbytes[0])
	fmt.Printf("%c", lbytes[1])
	fmt.Printf("%c", lbytes[2])
	fmt.Printf("%c", lbytes[3])
	fmt.Printf("%c", lbytes[4])
	fmt.Printf("%c", lbytes[5])

	fmt.Printf("%s", lbytes)
}

func assign() {
	var a [3]int = [3]int{10, 11, 12}
	var b [3]int = a
	a[1] = 99
	fmt.Printf("%d\n", b[0])
	fmt.Printf("%d\n", b[1])
	fmt.Printf("%d\n", b[2])
	var c [3]int
	c = b
	fmt.Printf("%d\n", c[0]+3)
	fmt.Printf("%d\n", c[1]+3)
	fmt.Printf("%d\n", c[2]+3)
}

func assignStrctField() {
	bilbo := Hobbit{
		dishes: [3]int{1, 2, 3},
	}
	var dishes [3]int
	dishes = bilbo.dishes
	//bilbo.dishes[2] = 0
	fmt.Printf("%d\n", dishes[2]+13) // 16
}

func main() {
	f1()
	f2()
	assign()
	assignStrctField()
	f3()
}

type Hobbit struct {
	id     int
	dishes [3]int
}
