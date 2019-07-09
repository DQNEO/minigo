package main


func f1() {
	var a = S("abc")
	var b = S("defg")
	var x gostring
	x = concat(a , b)
	fmtPrintf(S("%s\n"), x)
}

func f2() {
	spaces := S("> ")
	for i := 0; i < 3; i++ {
		spaces = concat(spaces, S("xx"))
	}
	fmtPrintf(S("%s\n"), spaces)
}

var seq int = 0

func foo() gostring {
	seq++
	return S("foo")
}

func f3() {
	label := concat3(foo() , foo() , S("bar"))
	fmtPrintf(S("%s\n"), label) // "foofoobar"
	fmtPrintf(S("%d\n"), seq)   // 2
}

func main() {
	f1()
	f2()
	f3()
}
