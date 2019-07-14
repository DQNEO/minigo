package main


const a int = b
const b int = 1

const x int = iota
const iota int = 2

const sum int = 1 + 2

const (
	c = 4
	d = 5
)

func f1() {
	fmtPrintf("%d\n", 0)
}

func f2() {
	fmtPrintf("%d\n", a)
}

func main() {
	f1()
	f2()
	fmtPrintf("%d\n", x)
	fmtPrintf("%d\n", sum)
}
