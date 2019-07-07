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
		fmt.Println("7")
	} else {
		fmt.Println("ERROR")
	}

	if 1 == 1 && 1 == 0 {
		fmt.Println("ERROR")
	} else {
		fmt.Println("8")
	}
}

func f3() {
	var flg bool
	flg = true
	if flg {
		fmt.Println("9")
	}
	if !flg {
		fmt.Println("ERROR")
	}
	flg = false
	if !flg {
		fmt.Println("10")
	}
}

func f4() {
	if 0 > 20-1 {
		fmt.Println("ERROR")
	} else {
		fmt.Println("11")
	}
}
func main() {
	f1()
	f2()
	f3()
	f4()
}
