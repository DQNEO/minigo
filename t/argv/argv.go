package main

import (
	"os"
)

func main() {
	fmtPrintf("%d\n", len(os.Args))
	fmtPrintf("%s\n", bytes(os.Args[1]))
}
