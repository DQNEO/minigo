package main

import "fmt"

type Hobbit struct {
	id int
	age int
}

func return_ptr() *Hobbit {
	var ptr = &Hobbit{
		id:1,
		age:123,
	}
	return ptr
}

func dummy() {
	var ptr = &Hobbit{
		id:0,
		age:123,
	}
	fmt.Printf("%d\n", ptr.id)
}

func f1() {
	p := return_ptr()
	dummy()
	fmt.Printf("%d\n", p.id)
}

func main() {
	f1()
}
