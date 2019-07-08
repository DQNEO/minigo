package main

func f1() {
	s := "main.go"
	suffix := ".go"
	if HasSuffix(S(s), S(suffix)) {
		fmtPrintf(S("1\n"))
	} else {
		fmtPrintf(S("ERROR\n"))
	}
}

func f2() {
	if Contains(S("foo/bar"), S("/")) {
		fmtPrintf(S("2\n"))
	} else {
		fmtPrintf(S("ERROR"))
	}
}

func f3() {
	s := Split(S("foo/bar"), S("/"))
	fmtPrintf(S("%d\n"), len(s)+1) // 3
	fmtPrintf(S("%s\n"), s[0])     // foo
	fmtPrintf(S("%s\n"), s[1])     // bar
}

func main() {
	f1()
	f2()
	f3()
}
