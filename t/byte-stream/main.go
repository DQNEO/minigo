package main


func f1() {
	path := S("t/min/min.go")
	s := readFile(path)
	_bs := ByteStream{
		filename:  bytes(path),
		source:    s,
		nextIndex: 0,
		line:      1,
		column:    0,
	}
	bs := &_bs
	len1 := len(bs.source)

	fmtPrintf("%d\n", len1-115) // 1

	var c byte
	c, _ = bs.get()
	fmtPrintf("%d\n", c-'p'+2)        // 2
	fmtPrintf("%d\n", bs.nextIndex+2) // 3
	c, _ = bs.get()
	fmtPrintf("%d\n", c-'a'+4)        // 4
	fmtPrintf("%d\n", bs.nextIndex+3) // 5
}

func main() {
	f1()
}
