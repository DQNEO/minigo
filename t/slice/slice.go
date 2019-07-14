package main


func f1() {
	var array [3]int = [3]int{1, 2, 3}
	fmtPrintf("%d\n", array[0]-1) // 1

	var slice []int

	slice = array[:] // {1,2,3}
	fmtPrintf("%d\n", slice[0])
	fmtPrintf("%d\n", slice[1])
	fmtPrintf("%d\n", slice[2])

	slice = array[:3] // {1,2,3}
	fmtPrintf("%d\n", slice[0]+3)
	fmtPrintf("%d\n", slice[1]+3)
	fmtPrintf("%d\n", slice[2]+3)

	slice = array[1:3] // {2,3}
	fmtPrintf("%d\n", slice[0]+5)
	fmtPrintf("%d\n", slice[1]+5)

	slice = array[2:3] // {3}
	fmtPrintf("%d\n", slice[0]+6)

	slice = array[2:] // {3}
	fmtPrintf("%d\n", slice[0]+7)

	var slice2 []int = array[1:3] // {2,3}
	fmtPrintf("%d\n", slice2[1]+8)
}

func f2() {

	var slice []int = []int{1, 2, 3}
	fmtPrintf("%d\n", slice[2]+9) // 12

	var slice2 []int
	slice2 = []int{4, 5, 6}
	fmtPrintf("%d\n", slice2[2]+7) // 13

	var slice3 []int = slice2
	fmtPrintf("%d\n", slice3[2]+8) // 14

	var slice4 []int
	slice4 = slice3
	fmtPrintf("%d\n", slice4[2]+9) //15

	bilbo := Hobbit{
		id:    0,
		items: nil,
	}
	if bilbo.items == nil {
		fmtPrintf("%d\n", 16)
	}

	bilbo = Hobbit{
		id:    0,
		items: []int{1, 2, 3},
	}

	fmtPrintf("%d\n", bilbo.items[2]+14) // 17
	bilbo.items = []int{15, 16, 17, 18}
	fmtPrintf("%d\n", bilbo.items[3]) // 18

}

func f3() {
	var array [3]int = [3]int{1, 2, 3}
	var slice = array[1:3]
	slice[1] = 19
	fmtPrintf("%d\n", slice[1])   // 19
	fmtPrintf("%d\n", array[2]+1) // 20
}

var gslice = []int{1, 3, 5}

func f4() {
	gslice[1] = 21
	fmtPrintf("%d\n", gslice[1]) // 21
}

func main() {
	f1()
	f2()
	f3()
	f4()
}

type Hobbit struct {
	id    int
	items []int
}
