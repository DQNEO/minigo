package main


func f1() {
	path := S("t/min/min.go")
	s := readFile(path)
	_bs := ByteStream{
		filename:  gostring(path),
		source:    s,
		nextIndex: 0,
		line:      1,
		column:    0,
	}
	bs := &_bs
	len1 := len(bs.source)

	fmtPrintf(S("%d\n"), len1-111) // 1

	var c byte
	c, _ = bs.get()
	fmtPrintf(S("%d\n"), c-'p'+2)        // 2
	fmtPrintf(S("%d\n"), bs.nextIndex+2) // 3
	c, _ = bs.get()
	fmtPrintf(S("%d\n"), c-'a'+4)        // 4
	fmtPrintf(S("%d\n"), bs.nextIndex+3) // 5
}

func main() {
	f1()
}
