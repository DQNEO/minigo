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
	keyFoo := "15"
	var lmap map[string]string = map[string]string{
		keyFoo: "10",
		"17":   "11",
	}

	fmtPrintf(S("9%s\n"), lmap["noexists"])
	fmtPrintf(S("%s\n"), lmap["15"]) // 10
	fmtPrintf(S("%s\n"), lmap["17"]) // 11

	fmtPrintf(S("%d\n"), len(lmap)+10) // 12

	var lenmap int

	lmap["19"] = "13"
	lenmap = len(lmap) // 3

	fmtPrintf(S("%s\n"), lmap["19"]) // 13

	fmtPrintf(S("%d\n"), lenmap+11) // 14

	lmap["15"] = "16"
	lmap["17"] = "18"
	lmap["19"] = "20"
	for k, v := range lmap {
		fmtPrintf(S("%s\n%s\n"), k, v) // 15,16,17,18,19,20
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

var keyFoo2 string = "keyfoo"

func f4() {
	keyFoo := "keyfoo"
	var lmap map[string]string = map[string]string{
		keyFoo:   "26",
		"keybar": "valuebar",
	}

	var ok bool
	var v string
	v, ok = lmap[keyFoo2]
	if ok {
		fmtPrintf(S("%d\n"), 25)
	}
	fmtPrintf(S("%s\n"), v) // 26

	v, ok = lmap["noexits"]
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
