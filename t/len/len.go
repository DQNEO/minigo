package main


func f0() {
	var x []int
	l := len(x)
	fmtPrintf(S("%d\n"), l)
}

func f1() {
	var a [1]int
	fmtPrintf(S("%d\n"), len(a))

	var b [2]int
	fmtPrintf(S("%d\n"), len(b))

	var c []int = b[:]
	fmtPrintf(S("%d\n"), len(c)+1) // 3

	c = b[0:1]
	fmtPrintf(S("%d\n"), len(c)+3) // 4

	c = b[1:2]
	fmtPrintf(S("%d\n"), len(c)+4) // 5

	var d []int = []int{1, 2, 3, 4, 5, 6}
	fmtPrintf(S("%d\n"), len(d)) // 6
}

func f2() {
	type Hobbit struct {
		id    int
		items []int
	}
	var h Hobbit
	h.items = []int{1}
	fmtPrintf(S("%d\n"), len(h.items)+6)          // 7
	var x int = len([]byte{'a', 'b'})
	fmtPrintf(S("%d\n"), x+6) //8
	var array [10]int
	fmtPrintf(S("%d\n"), len(array[2:7])+4) // 9
}

func f3() {
	var array = [...]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}
	fmtPrintf(S("%d\n"), len(array))
}

func receive_strings(a bytes, b bytes) {
	fmtPrintf(S("%d\n"), len(a))
	fmtPrintf(S("%d\n"), len(b))
}

func f4() {
	var hello bytes = S("01234567890")
	fmtPrintf(S("%d\n"), len(hello))
	s1 := S("012345678901")
	s2 := S("0123456789012")
	receive_strings(s1, s2)
}

func main() {
	f0()
	f1()
	f2()
	f3()
	f4()
}
