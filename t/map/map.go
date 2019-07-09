package main


var gmap map[string]bool

//var debug [6]int // r10, r11,...

func f1() {
	var lmap map[string]bool
	fmtPrintf(S("%d\n"), len(gmap)+1) // 1
	fmtPrintf(S("%d\n"), len(lmap)+2) // 2
}

func f2() {
	var lmap map[int]int = map[int]int{
		4: 7,
		5: 8,
		6: 9,
	}

	fmtPrintf(S("%d\n"), len(lmap)) // 3
	for i := range lmap {
		fmtPrintf(S("%d\n"), i) // 4,5,6
	}

	for _, v := range lmap {
		fmtPrintf(S("%d\n"), v) // 7,8,9
	}
}

func f3() {
	var lmap map[int]int = map[int]int{
		10: 11,
		12: 13,
	}

	lmap[14] = 15
	lmap[16] = 17
	for i, v := range lmap {
		fmtPrintf(S("%d\n"), i)
		fmtPrintf(S("%d\n"), v)
	}
}

func f4() {
	var lmap map[int]int = map[int]int{
		7:  17,
		9:  10,
		11: 12,
		0:  18,
	}

	fmtPrintf(S("%d\n"), lmap[0]) // 18

	fmtPrintf(S("%d\n"), lmap[999]+19) // 19
	lmap[9] = 21
	fmtPrintf(S("%d\n"), len(lmap)+16) // 20
	fmtPrintf(S("%d\n"), lmap[9])      // 21

	lmap[2] = 23
	fmtPrintf(S("%d\n"), len(lmap)+17) // 22
	fmtPrintf(S("%d\n"), lmap[2])      // 23

	var lmap2 map[int]int = map[int]int{
		0: 1,
		1: 1,
		2: 1,
		3: 1,
	}

	fmtPrintf(S("%d\n"), lmap[7]+7)   // 24
	fmtPrintf(S("%d\n"), lmap2[0]+24) // 25
}

func f5() {
	var lmap map[int]gostring = map[int]gostring{
		27: S("twenty seven"),
		26: S("twenty six"),
	}

	fmtPrintf(S("%s\n"), lmap[27])
	fmtPrintf(S("%s\n"), lmap[26])

	lmap[1] = S("one")
	fmtPrintf(S("%s\n"), lmap[1])

	for _, v := range lmap {
		fmtPrintf(S("%s\n"), v)
	}
}

// assign to an empty map
func f6() {
	var m map[int]int = map[int]int{}
	m[3] = 28
	fmtPrintf(S("%d\n"), m[3])
}

func main() {
	f1()
	f2()
	f3()
	f4()
	f5()
	f6()
}
