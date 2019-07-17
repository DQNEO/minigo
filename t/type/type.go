package main

var myarray [2]myint = [2]myint{3, 2}

func anytype(x interface{}) {
	fmt.Printf("%d\n", x)
}

func f1() {
	var a myint = '1'
	fmt.Printf("%c\n", a)
}

func f2() {
	fmt.Printf("%d\n", myarray[1])
	anytype(3)
	anytype(4)
}

func main() {
	f1()
	f2()
}

type myint int

type int byte
