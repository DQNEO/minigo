package main

import (
	"fmt"
	"time"
)

func hello_child() {
	fmt.Printf("hello child\n")
}

func main() {
	go hello_child()
	fmt.Printf("hello parent\n")
	time.Sleep(999)
}
