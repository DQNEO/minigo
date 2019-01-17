package main

import "fmt"

var message = "hello"

func f1() {
	fmt.Printf("%s\n", message)
}

func f2() {

	var mybytes []byte
	mybytes = []byte(message)

	fmt.Printf("%c", mybytes[0])
	fmt.Printf("%c", mybytes[1])
	fmt.Printf("%c", mybytes[2])
	fmt.Printf("%c", mybytes[3])
	fmt.Printf("%c", mybytes[4])
	println("")
}

func main() {
	f1()
	f2()
}
