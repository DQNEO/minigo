package main


var hello = [5]byte{'h', 'e', 'l', 'l', 'o'}

func ghello() {
	fmtPrintf(S("%c"), hello[0])
	fmtPrintf(S("%c"), hello[1])
	fmtPrintf(S("%c"), hello[2])
	fmtPrintf(S("%c"), hello[3])
	fmtPrintf(S("%c"), hello[4])
	fmtPrintf(S("%s"), S("\n"))

	s := bytes(hello[:])
	fmtPrintf(S("%s\n"), s)

}

func lworld() {
	var world = [5]byte{'w', 'o', 'r', 'l', 'd'}
	fmtPrintf(S("%c"), world[0])
	fmtPrintf(S("%c"), world[1])
	fmtPrintf(S("%c"), world[2])
	fmtPrintf(S("%c"), world[3])
	fmtPrintf(S("%c"), world[4])
	fmtPrintf(S("%s"), S("\n"))

	b := world[:]
	fmtPrintf(S("%s\n"), bytes(b))
}

func fappend() {
	var chars []byte
	chars = append(chars, '7')
	chars = append(chars, '8')
	fmtPrintf(S("%d\n"), len(chars)+4) // 6
	fmtPrintf(S("%c\n"), chars[0])     // 7
	fmtPrintf(S("%c\n"), chars[1])     // 8
	fmtPrintf(S("9\n"))                // 9

	chars[0] = '1'
	chars[1] = '0'
	fmtPrintf(S("%s\n"), chars) // 10
}

func main() {
	ghello()
	lworld()
	fappend()
}
