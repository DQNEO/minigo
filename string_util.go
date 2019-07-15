package main

import (
	"os"
	"fmt"
)

type identifier string

func fmtPrintf(format string, a... interface{}) {
	s := Sprintf(string(format), a...)
	os.Stdout.Write([]byte(s))
}

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

func concat(as string, bs string) string {
	a := []byte(as)
	b := []byte(bs)

	var r []byte
	for i:=0;i<len(a);i++ {
		r = append(r, a[i])
	}
	for i:=0;i<len(b);i++ {
		r = append(r, b[i])
	}
	return string(r)
}

func concat3(as string, bs string, cs string) string {
	a := []byte(as)
	b := []byte(bs)
	c := []byte(cs)
	var r []byte
	for i:=0;i<len(a);i++ {
		r = append(r, a[i])
	}
	for i:=0;i<len(b);i++ {
		r = append(r, b[i])
	}
	for i:=0;i<len(c);i++ {
		r = append(r, c[i])
	}
	return string(r)
}

func eq(as string, bs string) bool {
	if len(as) != len(bs) {
		return false
	}

	a := []byte(as)
	b := []byte(bs)
	for i:=0;i<len(a);i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
