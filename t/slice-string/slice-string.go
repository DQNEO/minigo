package main

import "fmt"

func f1() {
	var s string = "abcde"
	var sub string
	sub = s[1:3]
	fmt.Printf("%d\n", len(sub) - 1) // 1
	if sub == "bc" {
		fmt.Printf("2\n")
	}
}


func main() {
	f1()
}
