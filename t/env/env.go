package main

import (
	"fmt"
	"os"
)

func main() {
	envFOO := os.Getenv("FOO")
	fmt.Printf("%s\n", envFOO)
}
