package main

import "fmt"

var garray [3]int = [3]int{1,2,0}

func f1() {
	garray[2] = 3
	fmt.Printf("%d\n", garray[0])
	fmt.Printf("%d\n", garray[1])
	fmt.Printf("%d\n", garray[2])
}

func f2() {
	var larray [4]int = [4]int{4,5,0,7}
	larray[2] = 6
	fmt.Printf("%d\n", larray[0])
	fmt.Printf("%d\n", larray[1])
	fmt.Printf("%d\n", larray[2])
	fmt.Printf("%d\n", larray[3])
}

func f3() {
	var lbytes [6]byte = [6]byte{'?','e','l','l','o', 10}
	lbytes[0] = 'H'
	fmt.Printf("%c", lbytes[0])
	fmt.Printf("%c", lbytes[1])
	fmt.Printf("%c", lbytes[2])
	fmt.Printf("%c", lbytes[3])
	fmt.Printf("%c", lbytes[4])
	fmt.Printf("%c", lbytes[5])

	fmt.Printf("%s", lbytes)
}

func main() {
	f1()
	f2()
	f3()
}
