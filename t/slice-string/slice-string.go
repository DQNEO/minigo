package main


func f1() {
	var s bytes = S("abcde")
	var sub bytes
	sub = s[1:3]
	fmtPrintf("%d\n", len(sub)-1) // 1
	if eq(sub, S("bc")) {
		fmtPrintf("2\n")
	}
}

func f2() {
	var s bytes = S("main.go")
	var suffix bytes = S(".go")
	if len(s) == 7 {
		fmtPrintf("3\n")
	}
	if len(suffix) == 3 {
		fmtPrintf("4\n")
	}
	var suf2 bytes
	suf2 = s[4:]
	if eq(suf2, S(".go")) {
		fmtPrintf("5\n")
	}

	if len(s) >= len(suffix) {
		fmtPrintf("6\n")
	}

	low := len(s) - len(suffix)
	fmtPrintf("%d\n", low+3) //7

	// strings.HasSuffix
	var suff3 bytes
	suff3 = s[len(s)-len(suffix):]
	if eq(bytes(suff3), bytes(suffix)) {
		fmtPrintf("8\n") // 8
	}
}

func main() {
	f1()
	f2()
}
