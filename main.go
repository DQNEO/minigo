package main

import "fmt"

func main() {
	fmt.Println("\t.globl	main")
	fmt.Println("main:")
	fmt.Println("\tmovl	$0, %eax")
	fmt.Println("\tret")
}
