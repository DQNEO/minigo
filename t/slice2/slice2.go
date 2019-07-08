package main


func f1() {
	var r []int
	var args []int = []int{
		1,
	}

	r = append(r, 2)

	fmtPrintf(S("%d\n"), len(r))
	fmtPrintf(S("%d\n"), r[0])
}

func main() {
	f1()
}
