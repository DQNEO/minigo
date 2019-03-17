package main

// https://golang.org/ref/spec#Predeclared_identifiers

// Functions:
//	append cap close complex copy delete imag len
//	make new panic print println real recover

var builtinCode = `
package builtin

const MiniGo int = 1

func make(x interface{}) interface{} {
}

func panic(x interface{}) {
	puts(x)
	exit(1)
}

func println(s string) {
	puts(s)
}

func print(s string) {
	printf(s)
}

func recover() interface{} {
}

type error interface {
	Error() string
}
`
