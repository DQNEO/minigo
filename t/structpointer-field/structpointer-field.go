package main

import "fmt"

type IrRoot struct {
	vars           []int
	funcs          []int
	stringLiterals []string
}

func f1() {

	root := &IrRoot{}

	var vars []int = []int{3}

	root.vars = vars

	fmt.Printf("%d\n", len(vars))        // 1
	fmt.Printf("%d\n", len(root.vars)+1) // 2
	fmt.Printf("%d\n", root.vars[0])     // 3
}

type IrRoot2 struct {
	x  interface{}
	id int
}

func f2() {
	root := &IrRoot2{}

	var i int = 4
	var x interface{} = i
	root.x = x
	x = root.x
	var i2 int
	var ok bool
	i2, ok = x.(int)
	if !ok {
		fmt.Printf("ERROR\n")
	}
	fmt.Printf("%d\n", i2) // 4
}

func main() {
	f1()
	f2()
}
