package main


func f1() {
	var methods map[int][]gostring = map[int][]gostring{} // typeId : []methods
	methods[1] = []gostring{S("1")}
	fmtPrintf(S("%s\n"), methods[1][0])
}

func main() {
	f1()
}
