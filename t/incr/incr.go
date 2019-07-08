package main


func local() {
	var i int = 1
	fmtPrintf(S("%d\n"), i)
	i++
	fmtPrintf(S("%d\n"), i)
	i = 4
	i--
	fmtPrintf(S("%d\n"), i)
}

var j int = 4

func global() {
	fmtPrintf(S("%d\n"), j)
	j++
	fmtPrintf(S("%d\n"), j)
	j = 7
	j--
	fmtPrintf(S("%d\n"), j)
}

func pointerderef() {
	var a int = 6
	var b *int = &a
	*b++
	fmtPrintf(S("%d\n"), a)
	*b = 9
	*b--
	fmtPrintf(S("%d\n"), *b)
}

func main() {
	local()
	global()
	pointerderef()
}
