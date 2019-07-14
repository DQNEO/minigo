package main


var message bytes = bytes("hello")

func f1() {
	fmtPrintf("%s\n", message)
}

func f2() {

	var mybytes []byte
	mybytes = []byte(message)

	fmtPrintf("%c", mybytes[0])
	fmtPrintf("%c", mybytes[1])
	fmtPrintf("%c", mybytes[2])
	fmtPrintf("%c", mybytes[3])
	fmtPrintf("%c", mybytes[4])
	fmtPrintf("\n")
}

var gfoo string = string("foo")

func f3() {
	foo := "foo"
	if eq("foo", "foo") {
		fmtPrintf("1\n")
	}
	if eq(foo ,foo) {
		fmtPrintf("2\n")
	}
	if eq("foo" , foo) {
		fmtPrintf("3\n")
	}
	if eq(foo, "foo") {
		fmtPrintf("4\n")
	}
	if eq(foo, gfoo) {
		fmtPrintf("5\n")
	}
}

func f4() {
	s1 := "aaa"
	if !eq(s1, "bbb") {
		fmtPrintf("6\n")
	}

	if !eq(s1, "") {
		fmtPrintf("7\n")
	}
}

type mystring bytes

func f5() {
	s := "8"
	ms := mystring(s)
	fmtPrintf("%s\n", ms) // 8
}

func f6() {
	fmtPrintf("%d\n", len("123456789")) // 9
	s := "0123456789"
	fmtPrintf("%d\n", len(s)) // 10
}

func f7() {
	s := `11
12`
	fmtPrintf("%s\n", s) // 11 12
}

func f8() {
	s := "dummy"
	if eq(s, "") {
		fmtPrintf("ERROR\n")
	}
}

func f9() {
	var s2 string
	fmtPrintf("%s", s2)
	if eq(s2, "") {
		fmtPrintf("13\n")
	}
}

func f10() {
	var str0 string
	fmtPrintf("1%s4\n", str0)

	var str1 string = ""
	fmtPrintf("1%s5\n", str1)

	str2 := ""
	fmtPrintf("1%s6\n", str2)

	fmtPrintf("%d\n", len(str0) + len(str1) + 17) // 17

	str3 := "abc\n"

	fmtPrintf("%d\n", len(str3) + 14) // 18
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
	f10()
}
