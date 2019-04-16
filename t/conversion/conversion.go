package main

import "fmt"

func f1() {
	var bytes []byte
	var s string
	s = string(bytes)
	fmt.Printf("%s0\n", s) // 0
	fmt.Printf("%d\n", len(bytes) + 1) // 1
	fmt.Printf("%d\n", len(s) + 2) // 2
}

func f2() {
	var s string
	fmt.Printf("%s3\n", s) // 3
	fmt.Printf("%d\n", len(s) + 4) // 4
}

func f3() {
	var s string = ""
	fmt.Printf("%s5\n", s) // 5
	fmt.Printf("%d\n", len(s) + 6) // 6
}

func f4() {
	var s string
	var bytes []byte
	bytes = []byte(s)
	fmt.Printf("%s7\n", string(bytes)) // 7
	fmt.Printf("%d\n", len(bytes) + 8) // 8
}

func f5() {
	var s string
	var bytes []byte
	bytes = []byte(s)
	if bytes == nil {
		fmt.Printf("ERROR")
	} else {
		fmt.Printf("9\n")
	}
}

func main() {
	f1()
	f2()
	f3()
	f4()
	f5()
}
