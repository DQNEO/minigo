package main

import (
	"os"
	"fmt"
)

func main() {
	fmt.Printf("%d\n", len(os.Args))
	fmt.Printf("%s\n", []byte(os.Args[1]))
}
