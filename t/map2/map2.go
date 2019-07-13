package main


func f1() {
	var lmap map[int]int = map[int]int{
		1: 2,
		3: 4,
	}

	for i, v := range lmap {
		fmtPrintf(S("%d\n"), i)
		fmtPrintf(S("%d\n"), v)
	}

	fmtPrintf(S("%d\n"), lmap[1]+3) // 5
	fmtPrintf(S("%d\n"), lmap[3]+2) // 6

	lmap[7] = 8
	fmtPrintf(S("%d\n"), lmap[4]+7) // 7
	fmtPrintf(S("%d\n"), lmap[7])   // 8
}

func f2() {

	var lmap map[string]bytes = map[string]bytes{
		string("15"): S("10"),
		string("17"):   S("11"),
	}

	fmtPrintf(S("9%s\n"), lmap[string("noexists")])
	fmtPrintf(S("%s\n"), lmap[string("15")]) // 10
	fmtPrintf(S("%s\n"), lmap[string("17")]) // 11

	fmtPrintf(S("%d\n"), len(lmap)+10) // 12

	var lenmap int

	lmap[string("19")] = S("13")
	v19, ok := lmap[string("19")]
	if ok {
		fmtPrintf(S("%s\n"), v19) // 13
	} else {
		fmtPrintf(S("ERROR\n"))
	}

	lenmap = len(lmap) // 3
	fmtPrintf(S("%d\n"), lenmap+11) // 14

	lmap[string("15")] = S("16")
	lmap[string("17")] = S("18")
	lmap[string("19")] = S("20")
	for k, v := range lmap {
		fmtPrintf(S("%s\n%s\n"), S(k), v) // 15,16,17,18,19,20
	}
}

func f3() {
	var lmap map[int]int = map[int]int{
		1: 2,
		3: 21,
	}
	var ok bool
	var val int
	val, ok = lmap[3]
	fmtPrintf(S("%d\n"), val) // 21
	if ok {
		fmtPrintf(S("%d\n"), 22)
	}

	val, ok = lmap[2]
	if !ok {
		fmtPrintf(S("%d\n"), 23)
	}
	fmtPrintf(S("%d\n"), val+24) //24
}

var gkeyFoo2 bytes = bytes("keyfoo")

func f4() {
	gkeyFoo := bytes("keyfoo")

	var lmap map[string]bytes = map[string]bytes{
		string(gkeyFoo):   S("26"),
		string("keybar"): S("valuebar"),
	}

	var ok bool
	var v bytes
	v, ok = lmap[string(gkeyFoo2)]
	if ok {
		fmtPrintf(S("%d\n"), 25)
		fmtPrintf(S("%s\n"), v) // 26
	} else {
		fmtPrintf(S("ERROR\n"))
	}

	v, ok = lmap[string("noexits")]
	if !ok {
		fmtPrintf(S("%d\n"), 27)
	}
	fmtPrintf(S("28%s\n"), v) // 28
}

func main() {
	f1()
	f2()
	f3()
	f4()
}
