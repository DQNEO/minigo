package main


func main() {
	f, g := getIntInt()
	a := 1
	b := 2 // should be '2'
	c := "3"
	d := User{
		id:  4,
		age: 0,
	}
	e := &User{
		id:  5,
		age: 0,
	}
	fmtPrintf("%d\n", a)
	fmtPrintf("%c\n", b)
	fmtPrintf("%s\n", c)
	fmtPrintf("%d\n", d.id)
	fmtPrintf("%d\n", e.id)
	fmtPrintf("%d\n", f)
	fmtPrintf("%d\n", g)
}

type User struct {
	id  int
	age int
}

func getIntInt() (int, int) {
	return 6, 7
}
