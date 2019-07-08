package main


func f1() {
	// C style
	for i := 0; i < 10; i = i + 1 {
		fmtPrintf(S("%d\n"), i)
	}
}

func f2() {
	for i := 9; i < 20; i = i + 1 {
		if i == 9 {
			continue
		}
		if i == 16 {
			break
		}
		fmtPrintf(S("%d\n"), i)
	}
}

func f3() {
	var x int = 1
	for {
		if x == 0 {
			for {
				return
			}
			fmtPrintf(S("ERROR"))
			return
		}
		x = 16
		break
	}
	fmtPrintf(S("%d\n"), x)
}

func f4() {
	var i int = 17
	for ; i <= 19; i++ {
		fmtPrintf(S("%d\n"), i)
	}
}

func main() {
	f1()
	f2()
	f3()
	f4()
}
