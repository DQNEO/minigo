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
