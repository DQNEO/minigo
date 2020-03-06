package main

import (
	"./util"
	"fmt"
)

func f1() {
	strings := []string{"foo", "bar", "buz"}
	i := util.Index("bar", strings)
	fmt.Printf("%s\n", strings[i])

	if util.InArray("foo", strings) {
		fmt.Printf("foo\n")
	}
}

func main() {
	f1()
}
