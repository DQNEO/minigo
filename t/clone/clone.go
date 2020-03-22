package main

import (
	"fmt"
	"time"
)

var i int = 0

func hello() {
	i++
	fmt.Printf("hello %d\n", i)
}

func main() {
	go hello()
	time.Sleep(999)
	hello()
}
