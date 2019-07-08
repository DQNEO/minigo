package main

import (
	"os"
)

func main() {
	fmtPrintf(S("%d\n"), len(os.Args))
	fmtPrintf(S("%s\n"), os.Args[1])
}
