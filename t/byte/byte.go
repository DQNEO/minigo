package main


var hello = [5]byte{'h', 'e', 'l', 'l', 'o'}

func ghello() {
	fmtPrintf("%c", hello[0])
	fmtPrintf("%c", hello[1])
	fmtPrintf("%c", hello[2])
	fmtPrintf("%c", hello[3])
	fmtPrintf("%c", hello[4])
	fmtPrintf("%s", S("\n"))

	s := bytes(hello[:])
	fmtPrintf("%s\n", s)

}

func lworld() {
	var world = [5]byte{'w', 'o', 'r', 'l', 'd'}
	fmtPrintf("%c", world[0])
	fmtPrintf("%c", world[1])
	fmtPrintf("%c", world[2])
	fmtPrintf("%c", world[3])
	fmtPrintf("%c", world[4])
	fmtPrintf("%s", S("\n"))

	b := world[:]
	fmtPrintf("%s\n", bytes(b))
}

func fappend() {
	var chars []byte
	chars = append(chars, '7')
	chars = append(chars, '8')
	fmtPrintf("%d\n", len(chars)+4) // 6
	fmtPrintf("%c\n", chars[0])     // 7
	fmtPrintf("%c\n", chars[1])     // 8
	fmtPrintf("9\n")                // 9

	chars[0] = '1'
	chars[1] = '0'
	fmtPrintf("%s\n", chars) // 10
}

func main() {
	ghello()
	lworld()
	fappend()
}
