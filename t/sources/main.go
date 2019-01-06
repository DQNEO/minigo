package main

import "fmt"

var var2 int = 2

func main() {
	fmt.Printf("%s\n", const1)
	fmt.Printf("%d\n", var2)
	fmt.Printf("%d\n", const3)
	func4()
	func5() // mutuall dependent
}

func func5sub() {
	println("5")
}
