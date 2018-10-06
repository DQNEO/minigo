package main

import "fmt"
import "io/ioutil"

func readFile(filename string) string {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}
func main() {
	s := readFile("/dev/stdin")
	fmt.Println("\t.globl	main")
	fmt.Println("main:")
	fmt.Printf("\tmovl	$%s, %%eax\n", string(s))
	fmt.Println("\tret")
}
