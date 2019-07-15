package main

import (
	"os"
	"fmt"
)

type identifier string

func Sprintf(format string, a... interface{}) string {
	var args []interface{}
	for _, x := range a {
		var y interface{}
		switch x.(type) {
		case identifier: // This case is not reached by 2nd gen compiler
			var tmpgoident identifier = x.(identifier)
			var tmpbytes2 []byte = []byte(tmpgoident)
			y = tmpbytes2
		default:
			y = x
		}
		args = append(args, y)
	}
	a = nil // unset
	return fmt.Sprintf(format, args...)
}


func write(s []byte) {
	var b []byte = []byte(s)
	os.Stdout.Write(b)
}

func writeln(s []byte) {
	var b []byte = []byte(s)
	b = append(b, '\n')
	os.Stdout.Write(b)
}
