package main


func f1() {
	var h Hobbit = Hobbit{
		age:    1,
		height: 2,
	}

	fmtPrintf(S("%d\n"), h.age)
	fmtPrintf(S("%d\n"), h.height)

	var h2 Hobbit = h
	fmtPrintf(S("%d\n"), h2.age+2) // 3

	h.height = 100
	fmtPrintf(S("%d\n"), h2.height+2) // 4
}

func f2() {
	var h Hobbit = Hobbit{
		age:    1,
		height: 2,
	}

	var p *Hobbit = &h

	var h3 Hobbit = *p
	fmtPrintf(S("%d\n"), h3.age+4)    // 5
	fmtPrintf(S("%d\n"), h3.height+4) // 6
}

type Hobbit struct {
	age    int
	height int
}

func main() {
	f1()
	f2()
}
