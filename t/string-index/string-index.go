package main

import "fmt"

func f1() {
	var s []byte = []byte("543210")
	var c byte = s[5]
	fmt.Printf("%c\n", c)
	fmt.Printf("%c\n", s[4])
	fmt.Printf("%c\n", s[3])
	fmt.Printf("%c\n", s[2])
	fmt.Printf("%c\n", s[1])
	fmt.Printf("%c\n", s[0])
}

func main() {
	f1()
}
