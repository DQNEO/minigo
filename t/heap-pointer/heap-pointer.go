package main

import "fmt"

type Hobbit struct {
	id int
	age int
}

func return_ptr() *Hobbit {
	var ptr = &Hobbit{
		id:2,
		age:3,
	}
	return ptr
}

func dummy1() {
	var ptr = &Hobbit{
		id:0,
		age:4,
	}
	fmt.Printf("%d\n", ptr.id) // 0
}

func dummy2() {
	var array [4]int = [...]int{0,0,0,10}
	fmt.Printf("%d\n", array[3] - 9) // 1
}
func f1() {
	p := return_ptr()
	dummy1()
	dummy2()
	fmt.Printf("%d\n", p.id)
	fmt.Printf("%d\n", p.age)
}

func main() {
	f1()
}
