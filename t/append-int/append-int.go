package main


func handle_slice(s []int) {
	fmtPrintf(S("%d\n"), len(s)-2) // 1
	fmtPrintf(S("%d\n"), cap(s)-3) // 2
	fmtPrintf(S("%d\n"), s[0]-12)  // 3
	fmtPrintf(S("%d\n"), s[1]-10)  // 4
	s[2] = 3
	fmtPrintf(S("%d\n"), s[2]+2) // 5
}

var array [5]int = [...]int{15, 14, 13, 12, 11}

func f1() {
	var s []int = array[0:3]
	handle_slice(s)
	fmtPrintf(S("%d\n"), len(s)+3) // 6
	fmtPrintf(S("%d\n"), cap(s)+2) // 7
	fmtPrintf(S("%d\n"), s[0]-7)   // 8
	fmtPrintf(S("%d\n"), s[1]-5)   // 9
	fmtPrintf(S("%d\n"), s[2]+7)   // 10
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
	fmtPrintf(S("%d\n"), slice2[0])      // 11
	fmtPrintf(S("%d\n"), len(slice2)+10) // 12
	fmtPrintf(S("%d\n"), cap(slice2)+8)  // 13
}

var array3 [5]int = [...]int{1, 2, 3, 4, 16}

func f3() {
	var slice []int = array3[0:2]
	var slice2 []int = slice[0:5]
	fmtPrintf(S("%d\n"), len(slice2)+9)  // 14
	fmtPrintf(S("%d\n"), cap(slice2)+10) // 15
	fmtPrintf(S("%d\n"), slice2[4])      // 16
}

var array4 [5]int = [...]int{1, 2, 3, 4, 5}

func f4() {
	var slice []int = array4[0:1]
	var slice2 []int
	slice2 = append(slice, 18)
	fmtPrintf(S("%d\n"), len(slice2)+15) //17
	fmtPrintf(S("%d\n"), slice2[1])      // 18
}

func f5() {
	var array [2]int = [...]int{22, 23}
	var tmp []int = array[0:2]
	var s []int
	s = append(tmp, 19)
	fmtPrintf(S("%d\n"), s[2])               // 19
	fmtPrintf(S("%d\n"), cap(s) /* 4 */ +16) // 20
	s = append(s, 1)
	s = append(s, 1)
	fmtPrintf(S("%d\n"), cap(s) /* 8 */ +13) // 21
	fmtPrintf(S("%d\n"), s[0])               // 22
	fmtPrintf(S("%d\n"), s[1])               // 23
}

func main() {
	f1()
	f2()
	f3()
	f4()
	f5()
}
