package main


const (
	a0 int = iota
	a1
	a2
)

const (
	b0 int = iota
	b1     = iota
	b2     = iota
)

const (
	c0 int = iota
	c1
	c2
)

const (
	d0 int = 7
	d1
	d2 = iota
	d3
)

func f0() {
	fmtPrintf("%d\n", a0)
}

func f1() {
	fmtPrintf("%d\n", a1)
	fmtPrintf("%d\n", a2)
}

func f3() {
	fmtPrintf("%d\n", b0)
	fmtPrintf("%d\n", b1)
	fmtPrintf("%d\n", b2)
}

func main() {
	f0()
	f1()
	f3()
	fmtPrintf("%d\n", c0)
	fmtPrintf("%d\n", c1)
	fmtPrintf("%d\n", c2)
	fmtPrintf("%d\n", d0)
	fmtPrintf("%d\n", d1)
	fmtPrintf("%d\n", d2)
	fmtPrintf("%d\n", d3)
}
