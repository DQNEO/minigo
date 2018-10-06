package main

import "fmt"
import "io/ioutil"

func main() {
	contents, err := ioutil.ReadFile("/dev/stdin")
	if err != nil {
		panic(err)
	}
	fmt.Println("\t.globl	main")
	fmt.Println("main:")
	fmt.Printf("\tmovl	$%s, %%eax\n", string(contents))
	fmt.Println("\tret")
}
