package main


func f1() {
	i := recover()
	if i == nil {
		fmtPrintf(S("nil\n"))
	}
}

func main() {
	f1()
}
