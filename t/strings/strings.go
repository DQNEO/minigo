package main

import (
	"fmt"
	"strings"
)

func f1() {
	s := "main.go"
	suffix := ".go"
	if strings.HasSuffix(s, suffix) {
		fmt.Printf("1\n")
	}
}

func f2() {
	if strings.Contains("foo/bar", "/") {
		fmt.Printf("2\n")
	} else {
		fmt.Printf("ERROR")
	}
}

func f3() {
	s := strings.Split("foo/bar", "/")
	fmt.Printf("%d\n", len(s)+1) // 3
	fmt.Printf("%s\n", s[0])     // foo
	fmt.Printf("%s\n", s[1])     // bar
}

func main() {
	f1()
	f2()
	f3()
}
