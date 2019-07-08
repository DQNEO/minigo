package main


func f1() {
	var s string = "abcde"
	var sub string
	sub = s[1:3]
	fmtPrintf(S("%d\n"), len(sub)-1) // 1
	if sub == "bc" {
		fmtPrintf(S("2\n"))
	}
}

func f2() {
	var s = "main.go"
	var suffix = ".go"
	if len(s) == 7 {
		fmtPrintf(S("3\n"))
	}
	if len(suffix) == 3 {
		fmtPrintf(S("4\n"))
	}
	var suf2 string
	suf2 = s[4:]
	if suf2 == ".go" {
		fmtPrintf(S("5\n"))
	}

	if len(s) >= len(suffix) {
		fmtPrintf(S("6\n"))
	}

	low := len(s) - len(suffix)
	fmtPrintf(S("%d\n"), low+3) //7

	// strings.HasSuffix
	var suff3 string
	suff3 = s[len(s)-len(suffix):]
	if suff3 == suffix {
		fmtPrintf(S("8\n")) // 8
	}
}

func main() {
	f1()
	f2()
}
