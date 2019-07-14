package main


func f1() {
	var methods map[int][]bytes = map[int][]bytes{} // typeId : []methods
	methods[1] = []bytes{S("1")}
	fmtPrintf("%s\n", methods[1][0])
}

func main() {
	f1()
}
