package main

import "fmt"

func local() {
	var i int = 1
	fmt.Printf("%d\n", i)
	i++
	fmt.Printf("%d\n", i)
	i = 4
	i--
	fmt.Printf("%d\n", i)
}

var j int = 4

func global() {
	fmt.Printf("%d\n", j)
	j++
	fmt.Printf("%d\n", j)
	j = 7
	j--
	fmt.Printf("%d\n", j)
}
func main() {
	local()
	global()
}
