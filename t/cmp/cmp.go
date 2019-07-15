package main

import "fmt"

func f1() {
	var l = 1
	var g = 2

	if 1 != 1 {
		fmt.Printf("Error\n")
	}
	if l < g {
		fmt.Printf("%d\n", 1)
	}
	if l > g {
		fmt.Printf("Error\n")
	}
	fmt.Printf("%d\n", 2)
	if 1 == l {
		fmt.Printf("%d\n", 3)
	}

	if g == 2 {
		fmt.Printf("%d\n", 4)
	}

	if 1 <= l {
		fmt.Printf("%d\n", 5)
	}

	if g >= 2 {
		fmt.Printf("%d\n", 6)
	}
}

func f2() {
	if 1 == 0 || 1 == 1 {
		fmt.Printf("7\n")
	} else {
		fmt.Printf("ERROR\n")
	}

	if 1 == 1 && 1 == 0 {
		fmt.Printf("ERROR\n")
	} else {
		fmt.Printf("8\n")
	}
}

func f3() {
	var flg bool
	flg = true
	if flg {
		fmt.Printf("9\n")
	}
	if !flg {
		fmt.Printf("ERROR\n")
	}
	flg = false
	if !flg {
		fmt.Printf("10\n")
	}
}

func f4() {
	if 0 > 20-1 {
		fmt.Printf("ERROR\n")
	} else {
		fmt.Printf("11\n")
	}
}
func main() {
	f1()
	f2()
	f3()
	f4()
}
