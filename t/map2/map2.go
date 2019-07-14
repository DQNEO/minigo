package main


func f1() {
	var lmap map[int]int = map[int]int{
		1: 2,
		3: 4,
	}

	for i, v := range lmap {
		fmtPrintf("%d\n", i)
		fmtPrintf("%d\n", v)
	}

	fmtPrintf("%d\n", lmap[1]+3) // 5
	fmtPrintf("%d\n", lmap[3]+2) // 6

	lmap[7] = 8
	fmtPrintf("%d\n", lmap[4]+7) // 7
	fmtPrintf("%d\n", lmap[7])   // 8
}

func f2() {

	var lmap map[string]string = map[string]string{
		"15": "10",
		"17": "11",
	}

	fmtPrintf("9%s\n", lmap["noexists"])
	fmtPrintf("%s\n", lmap["15"]) // 10
	fmtPrintf("%s\n", lmap["17"]) // 11

	fmtPrintf("%d\n", len(lmap)+10) // 12

	var lenmap int

	lmap["19"] = "13"
	v19, ok := lmap["19"]
	if ok {
		fmtPrintf("%s\n", v19) // 13
	} else {
		fmtPrintf("ERROR\n")
	}

	lenmap = len(lmap) // 3
	fmtPrintf("%d\n", lenmap+11) // 14

	lmap["15"] = "16"
	lmap["17"] = "18"
	lmap["19"] = "20"
	for k, v := range lmap {
		fmtPrintf("%s\n%s\n", k, v) // 15,16,17,18,19,20
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
	fmtPrintf("%d\n", val) // 21
	if ok {
		fmtPrintf("%d\n", 22)
	}

	val, ok = lmap[2]
	if !ok {
		fmtPrintf("%d\n", 23)
	}
	fmtPrintf("%d\n", val+24) //24
}

var gkeyFoo2 bytes = bytes("keyfoo")

func f4() {
	gkeyFoo := "keyfoo"

	var lmap map[string]string = map[string]string{
		gkeyFoo:   "26",
		"keybar": "valuebar",
	}

	var ok bool
	var v bytes
	v, ok = lmap[string(gkeyFoo2)]
	if ok {
		fmtPrintf("%d\n", 25)
		fmtPrintf("%s\n", v) // 26
	} else {
		fmtPrintf("ERROR\n")
	}

	v, ok = lmap["noexits"]
	if !ok {
		fmtPrintf("%d\n", 27)
	}
	fmtPrintf("28%s\n", v) // 28
}

func main() {
	f1()
	f2()
	f3()
	f4()
}
