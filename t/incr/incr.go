package main


func local() {
	var i int = 1
	fmtPrintf("%d\n", i)
	i++
	fmtPrintf("%d\n", i)
	i = 4
	i--
	fmtPrintf("%d\n", i)
}

var j int = 4

func global() {
	fmtPrintf("%d\n", j)
	j++
	fmtPrintf("%d\n", j)
	j = 7
	j--
	fmtPrintf("%d\n", j)
}

func pointerderef() {
	var a int = 6
	var b *int = &a
	*b++
	fmtPrintf("%d\n", a)
	*b = 9
	*b--
	fmtPrintf("%d\n", *b)
}

func main() {
	local()
	global()
	pointerderef()
}
