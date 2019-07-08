package main


func f1() {
	var x map[int]map[int]int = map[int]map[int]int{}
	fmtPrintf(S("%d\n"), x[0][0])
}

func f2() {
	var mi MapIntInt = map[int]int{
		5: 1,
	}

	fmtPrintf(S("%d\n"), mi[5])

	var x map[int]map[int]int = map[int]map[int]int{
		111: map[int]int{
			11: 2,
		},
		112: map[int]int{
			12: 3,
		},
	}

	y := x[111]
	z := y[11]
	fmtPrintf(S("%d\n"), z)

	y = x[112]
	z = y[12]
	fmtPrintf(S("%d\n"), z)
}

type MapIntInt map[int]int

func main() {
	f1()
	f2()
}
