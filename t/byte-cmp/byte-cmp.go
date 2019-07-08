package main


func f1() {
	var x byte = 'a'
	var e byte = 'e'

	if e <= 'z' {
		fmtPrintf(S("1\n"))
	} else {
		fmtPrintf(S("%s\n"), x)
	}
}

func f2() {
	var c1 byte = 'p'
	var c2 byte = 'a'

	if 'a' <= c1 && c1 <= 'z' {
		fmtPrintf(S("2\n"))
	}

	if 'a' <= c2 && c2 <= 'z' {
		fmtPrintf(S("3\n"))
	}
}

func main() {
	f1()
	f2()
}
