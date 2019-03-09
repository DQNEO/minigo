package main

import "fmt"

func f1() {
	i := recover()
	if i == nil {
		fmt.Printf("nil\n")
	}
}

func main() {
	f1()
}
