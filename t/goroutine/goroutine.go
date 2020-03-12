package main

import (
	"fmt"
)

func say(s string) {
	fmt.Println(s)
}

func main() {
	go say("I am goroutine")
	say("hello")
}
