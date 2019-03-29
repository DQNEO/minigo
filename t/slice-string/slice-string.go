package main

import "fmt"

func f1() {
	var s string = "abcde"
	var sub string
	sub = s[1:3]
	fmt.Printf("%d\n", len(sub)-1) // 1
	if sub == "bc" {
		fmt.Printf("2\n")
	}
}

func f2() {
	var s = "main.go"
	var suffix = ".go"
	if len(s) == 7 {
		fmt.Printf("3\n")
	}
	if len(suffix) == 3 {
		fmt.Printf("4\n")
	}
	var suf2 string
	suf2 = s[4:]
	if suf2 == ".go" {
		fmt.Printf("5\n")
	}

	if len(s) >= len(suffix) {
		fmt.Printf("6\n")
	}

	low := len(s) - len(suffix)
	fmt.Printf("%d\n", low+3) //7

	// strings.HasSuffix
	var suff3 string
	suff3 = s[len(s)-len(suffix):]
	if suff3 == suffix {
		fmt.Printf("8\n") // 8
	}
}

func main() {
	f1()
	f2()
}
