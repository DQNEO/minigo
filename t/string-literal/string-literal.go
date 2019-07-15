package main

import "fmt"

func f1() {
	var  format string = "lea \\varname+\\offset(%%rip), %%rax"
	s := fmt.Sprintf(format)
	fmt.Println(s)
}

func main() {
	f1()
}
