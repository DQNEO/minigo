package main


func literal() {
	var u *User
	u = &User{
		id:  1,
		age: 2,
	}
	fmtPrintf("%d\n", u.id)
	fmtPrintf("%d\n", u.age)

	u = &User{
		id:  3,
		age: 4,
	}
	fmtPrintf("%d\n", u.id)
	fmtPrintf("%d\n", u.age)
}

func assign() {
	var u *User
	u = &User{
		id:  0,
		age: 4,
	}
	u.age = 5
	fmtPrintf("%d\n", u.age)
	u.age++
	fmtPrintf("%d\n", u.age)
	u.age = 8
	u.age--
	fmtPrintf("%d\n", u.age)
}

type S struct {
	dummy *int
	id    int
}

func f1() {
	var p *S
	p = &S{
		id: 123,
	}

	fmtPrintf("%d\n", p.id-115) // 8

	p.dummy = nil
	fmtPrintf("%d\n", p.id-114) // 9
}

func main() {
	literal()
	assign()
	return
	f1()
}

type User struct {
	id  int
	age int
}
