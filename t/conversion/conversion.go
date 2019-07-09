package main


func f1() {
	var bytes []byte
	var s gostring
	s = gostring(bytes)
	fmtPrintf(S("%s0\n"), s)           // 0
	fmtPrintf(S("%d\n"), len(bytes)+1) // 1
	fmtPrintf(S("%d\n"), len(s)+2)     // 2
}

func f2() {
	var s string
	fmtPrintf(S("%s3\n"), gostring(s))       // 3
	fmtPrintf(S("%d\n"), len(s)+4) // 4
}

func f3() {
	var s string = ""
	fmtPrintf(S("%s5\n"), gostring(s))       // 5
	fmtPrintf(S("%d\n"), len(s)+6) // 6
}

func f4() {
	var s string
	var bytes []byte
	bytes = []byte(s)
	fmtPrintf(S("%s7\n"), gostring(bytes)) // 7
	fmtPrintf(S("%d\n"), len(bytes)+8)   // 8
}

func f5() {
	var s string
	var bytes []byte
	bytes = []byte(s)
	if bytes == nil {
		fmtPrintf(S("ERROR"))
	} else {
		fmtPrintf(S("9\n"))
	}
}

func main() {
	f1()
	f2()
	f3()
	f4()
	f5()
}
