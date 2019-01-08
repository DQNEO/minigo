package main

var fmtCode  = `
package fmt

func Printf(format string, param interface{}) {
	printf(format, param)
}

func Sprintf(format string, param interface{}) string {
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
}
