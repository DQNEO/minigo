package main


func f1() {
	var a = "abc"
	var b = "defg"
	var x string
	x = a + b
	fmtPrintf(S("%s\n"), x)
}

func f2() {
	spaces := "> "
	for i := 0; i < 3; i++ {
		spaces += "xx"
	}
	fmtPrintf(S("%s\n"), spaces)
}

var seq int = 0

func foo() string {
	seq++
	return "foo"
}

func f3() {
	label := foo() + foo() + "bar"
	fmtPrintf(S("%s\n"), label) // "foofoobar"
	fmtPrintf(S("%d\n"), seq)   // 2
}

func main() {
	f1()
	f2()
	f3()
}
