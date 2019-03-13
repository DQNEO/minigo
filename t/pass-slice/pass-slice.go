package main

import "fmt"

func handle_slice(s []int) {
	fmt.Printf("%d\n", len(s)-2) // 1
	fmt.Printf("%d\n", cap(s)-3) // 2
	fmt.Printf("%d\n", s[0]-12)  // 3
	fmt.Printf("%d\n", s[1]-10)  // 4
	s[2] = 3
	fmt.Printf("%d\n", s[2]+2) // 5
}

var array [5]int = [...]int{15, 14, 13, 12, 11}

func f1() {
	var s []int = array[0:3]
	handle_slice(s)
	fmt.Printf("%d\n", len(s)+3) // 6
	fmt.Printf("%d\n", cap(s)+2) // 7
	fmt.Printf("%d\n", s[0]-7)   // 8
	fmt.Printf("%d\n", s[1]-5)   // 9
	fmt.Printf("%d\n", s[2]+7)   // 10
}

func copy_slice(x []int) []int {
	var s []int
	x[0] = 11
	s = x
	return s
}

var array2 [5]int = [...]int{15, 14, 13, 12, 11}

func f2() {
	var slice []int = array2[0:2]
	var slice2 []int
	slice2 = copy_slice(slice)
	fmt.Printf("%d\n", slice2[0]) // 11
	fmt.Printf("%d\n", len(slice2) + 10) // 12
	fmt.Printf("%d\n", cap(slice2) + 8) // 13
}

var array3 [5]int = [...]int{1,2,3,4,16}

func f3() {
	var slice []int = array3[0:2]
	var slice2 []int = slice[0:5]
	fmt.Printf("%d\n", len(slice2) + 9) // 14
	fmt.Printf("%d\n", cap(slice2) + 10) // 15
	fmt.Printf("%d\n", slice2[4]) // 16
}

func append(x []int, elm int) []int {
	var s []int
	zlen := len(x) + 1
	if cap(x) >= zlen {
		s = x[:zlen]
		s[len(x)] = elm
	}

	return s
}

var array4 [5]int = [...]int{1,2,3,4,5}

func f4() {
	var slice []int = array4[0:1]
	var slice2 []int
	slice2 = append(slice, 18)
	fmt.Printf("%d\n", len(slice2) + 15) //17
	fmt.Printf("%d\n", slice2[1]) // 18
}

func main() {
	f1()
	f2()
	f3()
	f4()
}
