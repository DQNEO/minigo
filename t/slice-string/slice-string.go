package main


func f1() {
	var s gostring = S("abcde")
	var sub gostring
	sub = s[1:3]
	fmtPrintf(S("%d\n"), len(sub)-1) // 1
	if eq(sub, S("bc")) {
		fmtPrintf(S("2\n"))
	}
}

func f2() {
	var s gostring = S("main.go")
	var suffix gostring = S(".go")
	if len(s) == 7 {
		fmtPrintf(S("3\n"))
	}
	if len(suffix) == 3 {
		fmtPrintf(S("4\n"))
	}
	var suf2 gostring
	suf2 = s[4:]
	if eq(suf2, S(".go")) {
		fmtPrintf(S("5\n"))
	}

	if len(s) >= len(suffix) {
		fmtPrintf(S("6\n"))
	}

	low := len(s) - len(suffix)
	fmtPrintf(S("%d\n"), low+3) //7

	// strings.HasSuffix
	var suff3 gostring
	suff3 = s[len(s)-len(suffix):]
	if eq(gostring(suff3), gostring(suffix)) {
		fmtPrintf(S("8\n")) // 8
	}
}

func main() {
	f1()
	f2()
}
