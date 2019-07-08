package main


func structfield() {
	bilbo := Hobbit{
		id:    1,
		age:   2,
		items: [3]int{3, 4, 5},
	}

	fmtPrintf(S("%d\n"), bilbo.id)
	fmtPrintf(S("%d\n"), bilbo.age)
	fmtPrintf(S("%d\n"), bilbo.items[0])
	fmtPrintf(S("%d\n"), bilbo.items[1])
	fmtPrintf(S("%d\n"), bilbo.items[2])
}

type Hobbit struct {
	id    int
	age   int
	items [3]int
}

func main() {
	structfield()
}
