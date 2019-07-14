package main


// for range test
func f1() {
	var array1 [3]int = [3]int{9, 9, 9}
	var array2 [3]int = [3]int{4, 6, 8}

	var v int
	var i int
	for i = range array1 {
		fmtPrintf("%d\n", i)
	}

	for i, v = range array2 {
		fmtPrintf("%d\n", i*2+3)
		fmtPrintf("%d\n", v)
	}
}

func f2() {
	bilbo := Hobbit{
		id:    1,
		age:   111,
		items: [3]int{9, 10, 11},
	}
	for _, v := range bilbo.items {
		fmtPrintf("%d\n", v)
	}
}

func f3() {
	var slice = []int{112, 113, 114}
	for _, v := range slice {
		fmtPrintf("%d\n", v-100)
	}
}

func f4() {
	var slice []int
	for i := range slice {
		fmtPrintf("error %d\n", i)
	}
}

func f5() {
	var slice []int = nil
	for i, v := range slice {
		fmtPrintf(S("error %d,%d\n", i), v)
	}

}

func f6() {
	var array [5]int = [5]int{1, 2, 3, 4, 5}
	for _, v := range array {
		if v == 1 {
			continue
		}
		if v == 4 {
			break
		}
		fmtPrintf("%d\n", v+13) // 15,16
	}

}

func f7() {
	var slice = []int{1, 1, 1}
	var i int
	for i, _ = range slice {
		fmtPrintf("%d\n", i+17) // 17,18,19
	}
	fmtPrintf("%d\n", i+18) // 20
}

func main() {
	f1()
	f2()
	f3()
	f4()
	f5()
	f6()
	f7()
}

type Hobbit struct {
	id    int
	age   int
	items [3]int
}
