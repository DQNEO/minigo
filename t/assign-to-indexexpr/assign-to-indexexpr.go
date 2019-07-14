package main


func f1() {
	var methods map[int][][]byte = map[int][][]byte{} // typeId : []methods
	methods[1] = [][]byte{[]byte("1")}
	fmtPrintf("%s\n", methods[1][0])
}

func main() {
	f1()
}
