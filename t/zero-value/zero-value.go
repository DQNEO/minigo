package main


func strct() {
	var h = Hobbit{}
	fmtPrintf("%d\n", h.id+1)  // 0
	fmtPrintf("%d\n", h.age+2) // 0

	var h2 = Hobbit{
		id: 3,
	}
	fmtPrintf("%d\n", h2.id)    // 3
	fmtPrintf("%d\n", h2.age+4) // 0

	var h3 = Hobbit{
		age: 6,
	}
	fmtPrintf("%d\n", h3.id+5) // 0
	fmtPrintf("%d\n", h3.age)  // 6

	var h4 Hobbit
	fmtPrintf("%d\n", h4.id+7)  // 0
	fmtPrintf("%d\n", h4.age+8) // 0
}

func array() {
	var array1 [2]int
	fmtPrintf("%d\n", array1[0]+9)  // 0
	fmtPrintf("%d\n", array1[1]+10) // 0

	var array2 [2]int = [2]int{}
	fmtPrintf("%d\n", array2[0]+11) // 0
	fmtPrintf("%d\n", array2[1]+12) // 0

	var array3 [2]int = [2]int{3}
	fmtPrintf("%d\n", array3[0]+10) // 3
	fmtPrintf("%d\n", array3[1]+14) // 0
}

func primitives() {
	var i int
	fmtPrintf("%d\n", i+15) // 0
	var b byte

	var c int = b+16
	fmtPrintf("%d\n", c) // 0

	var bol bool
	if !bol {
		fmtPrintf("%d\n", 17)
	}
	var slice []int
	if slice == nil {
		fmtPrintf("%d\n", 18)
	}
}

func main() {
	strct()
	array()
	primitives()
}

type Hobbit struct {
	id  int
	age int
}
