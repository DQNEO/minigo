package main


func f1() {
	var x interface{}
	var i int = 1
	x = i

	var z int
	switch x.(type) {
	case nil:
		z = -1
	case int:
		z = 1
	case gostring:
		z = 2
	default:
		z = 5
	}

	fmtPrintf(S("%d\n"), z)

	var s gostring = S("hello")
	x = s
	switch x.(type) {
	case nil:
		z = -1
	case int:
		z = 1
	case gostring:
		z = 2
	default:
		z = 5
	}
	fmtPrintf(S("%d\n"), z)
}

func main() {
	f1()
}
