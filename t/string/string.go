package main


var message gostring = gostring("hello")

func f1() {
	fmtPrintf(S("%s\n"), message)
}

func f2() {

	var mybytes []byte
	mybytes = []byte(message)

	fmtPrintf(S("%c"), mybytes[0])
	fmtPrintf(S("%c"), mybytes[1])
	fmtPrintf(S("%c"), mybytes[2])
	fmtPrintf(S("%c"), mybytes[3])
	fmtPrintf(S("%c"), mybytes[4])
	fmtPrintf(S("\n"))
}

var gfoo gostring = gostring("foo")

func f3() {
	foo := S("foo")
	if eq(S("foo"), S("foo")) {
		fmtPrintf(S("1\n"))
	}
	if eq(foo ,foo) {
		fmtPrintf(S("2\n"))
	}
	if eq(S("foo") , foo) {
		fmtPrintf(S("3\n"))
	}
	if eq(foo, S("foo")) {
		fmtPrintf(S("4\n"))
	}
	if eq(foo, gfoo) {
		fmtPrintf(S("5\n"))
	}
}

func f4() {
	s1 := S("aaa")
	if !eq(s1, S("bbb")) {
		fmtPrintf(S("6\n"))
	}

	if !eq(s1, S("")) {
		fmtPrintf(S("7\n"))
	}
}

type mystring gostring

func f5() {
	s := S("8")
	ms := mystring(s)
	fmtPrintf(S("%s\n"), ms) // 8
}

func f6() {
	fmtPrintf(S("%d\n"), len("123456789")) // 9
	s := S("0123456789")
	fmtPrintf(S("%d\n"), len(s)) // 10
}

func f7() {
	s := S(`11
12`)
	fmtPrintf(S("%s\n"), s) // 11 12
}

func f8() {
	s := S("dummy")
	if eq(s, S("")) {
		fmtPrintf(S("ERROR\n"))
	}
}

func f9() {
	var s2 gostring
	fmtPrintf(S("%s"), s2)
	if eq(s2, S("")) {
		fmtPrintf(S("13\n"))
	}
}

func main() {
	f1()
	f2()
	f3()
	f4()
	f5()
	f6()
	f7()
	f8()
	f9()
}
