package main


func f0() {
	var x []int
	l := cap(x)
	fmtPrintf("%d\n", l)
}

func f1() {
	var a [1]int
	fmtPrintf("%d\n", cap(a))

	var b [2]int
	fmtPrintf("%d\n", cap(b))

	var c []int = b[:]
	fmtPrintf("%d\n", cap(c)+1) // 3

	c = b[0:1]
	fmtPrintf("%d\n", cap(c)+2) // 4

	c = b[1:2]
	fmtPrintf("%d\n", cap(c)+4) // 5

	var d []int = []int{1, 2, 3, 4, 5, 6}
	fmtPrintf("%d\n", cap(d)) // 6
}

func f2() {
	type Hobbit struct {
		id    int
		items []int
	}
	var h Hobbit
	h.items = []int{1}
	fmtPrintf("%d\n", cap(h.items)+6)          // 7
	fmtPrintf("%d\n", cap([]byte{'a', 'b'})+6) //8
	var array [10]int
	fmtPrintf("%d\n", cap(array[2:7])+1) // 9
}

func f3() {
	var array = [...]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}
	fmtPrintf("%d\n", cap(array)) // 10
}

/*

https://golang.org/ref/spec#Slice_expressions

After slicing the array a

a := [5]int{1, 2, 3, 4, 5}
s := a[1:4]
the slice s has type []int, length 3, capacity 4, and elements

s[0] == 2
s[1] == 3
s[2] == 4

*/
func f4() {
	var a = [5]int{1, 2, 3, 4, 5}
	var s []int = a[1:4]
	fmtPrintf("%d\n", len(s)+8) // 11
	fmtPrintf("%d\n", cap(s)+8) // 12
	fmtPrintf("%d\n", s[0]+11)  // 13
}

//  Full slice expressions
//
//  For an array, pointer to array, or slice a (but not a string), the primary expression
//
//  a[low : high : max]
//  constructs a slice of the same type, and with the same length and elements as the simple slice expression a[low : high]. Additionally, it controls the resulting slice's capacity by setting it to max - low. Only the first index may be omitted; it defaults to 0. After slicing the array a
//
//  a := [5]int{1, 2, 3, 4, 5}
//  t := a[1:3:5]
//  the slice t has type []int, length 2, capacity 4, and elements
//
//  t[0] == 2
//  t[1] == 3

func f5() {
	var a = [5]int{1, 2, 3, 4, 5}
	var s []int = a[1:3:4]
	fmtPrintf("%d\n", len(s)+12) // 14
	fmtPrintf("%d\n", cap(s)+12) // 15
	fmtPrintf("%d\n", s[0]+14)   // 16
}

func main() {
	f0()
	f1()
	f2()
	f3()
	f4()
	f5()
}
