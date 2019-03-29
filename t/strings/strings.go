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

func main() {
	f1()
}
