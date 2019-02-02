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
}

func println(s string) {
	puts(s)
}

func recover() interface{} {
}

type error interface {
	Error() string
}
`

var internalRuntimeCode = `
package runtime
// Runtime
var heap [1048576]int
var heapIndex int

func malloc(size int) int {
	if heapIndex == 0 {
		heapIndex = (heap + 0)
	}
	if heapIndex + size - heap > len(heap) {
		return 0
	}
	r := heapIndex
	heapIndex += size
	return r
}
`
