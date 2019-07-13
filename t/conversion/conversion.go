package main


func f1() {
	var vbytes []byte
	var s bytes
	s = bytes(vbytes)
	fmtPrintf(S("%s0\n"), s)           // 0
	fmtPrintf(S("%d\n"), len(vbytes)+1) // 1
	fmtPrintf(S("%d\n"), len(s)+2)     // 2
}

func f2() {
	var s bytes
	fmtPrintf(S("%s3\n"), bytes(s))       // 3
	fmtPrintf(S("%d\n"), len(s)+4) // 4
}

func f3() {
	var s bytes = S("")
	fmtPrintf(S("%s5\n"), s)       // 5
	fmtPrintf(S("%d\n"), len(s)+6) // 6
}

func f4() {
	var s bytes
	var vbytes []byte
	vbytes = []byte(s)
	fmtPrintf(S("%s7\n"), bytes(vbytes)) // 7
	fmtPrintf(S("%d\n"), len(vbytes)+8)   // 8
}

func f5() {
	var s bytes
	var bytes []byte
	bytes = []byte(s)
	if bytes == nil {
		fmtPrintf(S("9\n"))
	} else {
		fmtPrintf(S("ERROR"))
	}
}

func main() {
	f1()
	f2()
	f3()
	f4()
	f5()
}
