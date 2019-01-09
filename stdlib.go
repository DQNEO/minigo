package main

var fmtCode  = `
package fmt

func Printf(format string, param interface{}) {
	printf(format, param)
}

func Sprintf(format string, param interface{}) string {
}

func Fprintf(file interface{}, format string, param ...interface{}) {
}

func Println(s string) {
}

`
var errorsCode = `
package errors

func New() {
}

`
var ioutilCode = `
package ioutil

func ReadFile(filename string) ([]byte, error) {
	return 0, 0
}
`

var builtinCode  = `
package builtin

const MiniGo int = 1

// println should be a "Predeclared identifiers"
// https://golang.org/ref/spec#Predeclared_identifiers
func println(s string) {
	puts(s)
}
`

var osCode = `
package os

const Stderr = 2

var Args []string

func Exit(i int) {
}

`

var stringsCode = `
package strngs

func HasSuffix(s string) bool {
}

func Contains(s string) bool {
}

func Split(s string, x string) []string {
}
`
var runtimeCode = `
package runtime

func Caller(n int) (interface{}, interface{},interface{},interface{}) {
}

func FuncForPC(x int) int {
}

`
var strconvCode = `
package strconv

func Atoi() {}

`

type pkgsource struct {
	name identifier
	code string
}
var pkgsources []pkgsource = []pkgsource{
	pkgsource{
		name: "fmt",
		code: fmtCode,
	},
	pkgsource{
		name: "ioutil",
		code: ioutilCode,
	},
	pkgsource{
		name: "errors",
		code: errorsCode,
	},
	pkgsource{
		name: "os",
		code: osCode,
	},
	pkgsource{
		name: "strings",
		code: stringsCode,
	},
	pkgsource{
		name: "runtime",
		code: runtimeCode,
	},
	pkgsource{
		name: "strconv",
		code: strconvCode,
	},
}
