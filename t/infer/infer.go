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
	fmtPrintf(S("%d\n"), a)
	fmtPrintf(S("%c\n"), b)
	fmtPrintf(S("%s\n"), c)
	fmtPrintf(S("%d\n"), d.id)
	fmtPrintf(S("%d\n"), e.id)
	fmtPrintf(S("%d\n"), f)
	fmtPrintf(S("%d\n"), g)
}

type User struct {
	id  int
	age int
}

func getIntInt() (int, int) {
	return 6, 7
}
