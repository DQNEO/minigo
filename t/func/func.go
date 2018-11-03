package main

import "fmt"

func plus(a int, b int) int {
	return a + b
}

func fvoid() {

}

func fvoidsemicolon() {
	;
}

func fvoidreturn() {
	return
}

func f1() {
	fmt.Printf("%d\n", plus(0, 1))
	fmt.Printf("%d\n", plus(1, 1))
	fmt.Printf("%d\n", plus(2, 1))
	fmt.Printf("%d\n", plus(2, 2))
}

func main() {
	f1()
	fvoid()
	fvoidsemicolon()
	fvoidreturn()
}
