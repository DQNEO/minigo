package main

import "fmt"

func f1() {
	var s []byte = []byte("abcde")
	var sub []byte
	sub = s[1:3]
	fmt.Printf("%d\n", len(sub)-1) // 1
	if string(sub) == "bc" {
		fmt.Printf("2\n")
	}
}

func f2() {
	var s []byte = []byte("main.go")
	var suffix string = ".go"
	if len(s) == 7 {
		fmt.Printf("3\n")
	}
	if len(suffix) == 3 {
		fmt.Printf("4\n")
	}
	var suf2 []byte
	suf2 = s[4:]
	if string(suf2) == ".go" {
		fmt.Printf("5\n")
	}

	if len(s) >= len(suffix) {
		fmt.Printf("6\n")
	}

	low := len(s) - len(suffix)
	fmt.Printf("%d\n", low+3) //7

	// strings.HasSuffix
	suff3 := s[len(s)-len(suffix):]
	if string(suff3) == suffix {
		fmt.Printf("8\n") // 8
	}
}

func f3() {
	s := "hello"
	subs := s[0:3]
	if subs == "hel" {
		fmt.Println("9")
	}

	if s[3:] == "lo" {
		fmt.Println("10")
	}
}

func main() {
	f1()
	f2()
	f3()
}
