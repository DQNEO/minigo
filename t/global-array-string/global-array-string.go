package main

import "fmt"

var messages = [2]string{"hello", "world"}

func main() {
	fmt.Println(string(messages[1]))
}
