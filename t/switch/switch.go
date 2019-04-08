package main

import "fmt"

func swtch(x int) int {
	var y int
	switch x {
	case 1:
		y = 1
	case 2, 3:
		y = 2
	case 5:
		y = 5
		y = 2
	default:
		y = 7
	}

	return y
}

func f1() {
	var i int
	i = swtch(1)
	fmt.Printf("%d\n", i) // 1
	i = swtch(2)
	fmt.Printf("%d\n", i) // 2
	i = swtch(3)
	fmt.Printf("%d\n", i+1) // 3
	i = swtch(999)
	fmt.Printf("%d\n", i-3) // 4
}

func swtch2(x int) int {
	var y int
	switch x {
	case 1:
		y = 1
	}

	return y
}

func f2() {
	i := swtch2(3)
	fmt.Printf("%d\n", i+5) // 5
}

func f3() {
	switch {
	case 1+1 == 3:
		println("Error")
	case 1+1 == 2:
		println(6)
	default:
		println("Error")
	}
}

func main() {
	f1()
	f2()
	f3()
}
