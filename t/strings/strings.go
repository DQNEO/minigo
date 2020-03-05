package main

import "strings"
import "fmt"

func f1() {
	s := "main.go"
	suffix := ".go"
	if strings.HasSuffix(s, suffix) {
		fmt.Printf("1\n")
	} else {
		fmt.Printf("ERROR\n")
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

func f4() {
	target := "foo bar buz"
	if !strings.HasPrefix(target, "foo") {
		panic("error")
	}

	if strings.HasPrefix(target, " ") {
		panic("error")
	}

	if strings.HasPrefix(target, "buz") {
		panic("error")
	}
	fmt.Printf("4\n")
}
func main() {
	f1()
	f2()
	f3()
	f4()
}
