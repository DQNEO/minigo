package main

import (
	"os"
)

func main() {
	fmtPrintf("%d\n", len(os.Args))
	fmtPrintf("%s\n", []byte(os.Args[1]))
}
