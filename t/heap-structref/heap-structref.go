package main


type Hobbit struct {
	id  int
	age int
}

func return_ptr() *Hobbit {
	var ptr = &Hobbit{
		id:  2,
		age: 3,
	}
	return ptr
}

func dummy1() {
	var ptr = &Hobbit{
		id:  0,
		age: 4,
	}
	fmtPrintf("%d\n", ptr.id) // 0
}

func dummy2() {
	var array [4]int = [...]int{0, 1, 0, 0}
	fmtPrintf("%d\n", array[1]) // 1
}
func f1() {
	p := return_ptr()
	dummy1()
	dummy2()
	fmtPrintf("%d\n", p.id)  // 2
	fmtPrintf("%d\n", p.age) // 3
}

func main() {
	f1()
}
