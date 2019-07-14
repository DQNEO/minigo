package main


func f1() {
	var a = "abc"
	var b = "defg"
	var x string
	x = concat(a , b)
	fmtPrintf("%s\n", x)
}

func f2() {
	spaces := "> "
	for i := 0; i < 3; i++ {
		spaces = concat(spaces, "xx")
	}
	fmtPrintf("%s\n", spaces)
}

var seq int = 0

func foo() string {
	seq++
	return "foo"
}

func f3() {
	label := concat3(foo() , foo() , "bar")
	fmtPrintf("%s\n", label) // "foofoobar"
	fmtPrintf("%d\n", seq)   // 2
}

func main() {
	f1()
	f2()
	f3()
}
