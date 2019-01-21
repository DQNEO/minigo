package main

import "fmt"

func eval() {
	fmt.Printf("%d\n", gint)
	fmt.Printf("%c\n", gchar)
	if gtrue {
		fmt.Printf("3\n")
	}
	if !gfalse {
		fmt.Printf("4\n")
	}

	fmt.Printf("%d\n", gstruct.parent_id) // 5
	fmt.Printf("%d\n", gstruct.id) // 6
	fmt.Printf("%d\n", gstruct.age + 7) // 7

}

func assign() {

}

func main() {
	eval()
	assign()
}


var gint int = 1
var gchar byte = '2'
var gtrue bool = true
var gfalse bool = false

var gstruct = Hobbit{
	parent_id:5,
	id:6,
}

type Hobbit struct {
	id int
	age int
	parent_id int
}
