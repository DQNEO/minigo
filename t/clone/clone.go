package main

import (
	"fmt"
	"time"
)

func hello() {
	fmt.Printf("hello child\n")
}

func main() {
	go hello()
	fmt.Printf("hello parent\n")
	time.Sleep(999)
}
