package main

import (
	"fmt"
	"os"
)

type gostring []byte
type cstring string

type identifier string
type packageName identifier

type goidentifier gostring

func S(s string) gostring {
	return gostring(s)
}

func GoSprintf(format gostring, param ...interface{}) gostring {
	s := fmt.Sprintf(string(format), param...)
	return gostring(s)
}

func write(s gostring) {
	var b []byte = []byte(s)
	os.Stdout.Write(b)
}

func writeln(s gostring) {
	var b []byte = []byte(s)
	b = append(b, '\n')
	os.Stdout.Write(b)
}

func convertCstringsToGostrings(cstrings []string) []gostring {
	var r []gostring
	for _, cs := range cstrings {
		r = append(r, gostring(cs))
	}

	return r
}

func concat(a gostring, b gostring) gostring {
	var r []byte
	for i:=0;i<len(a);i++ {
		r = append(r, a[i])
	}
	for i:=0;i<len(b);i++ {
		r = append(r, b[i])
	}
	return r
}

func concat3(a gostring, b gostring, c gostring) gostring {
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
	return r
}

func eq(a gostring, b cstring) bool {
	return eqGostrings(a, gostring(b))
}

func eqGostrings(a gostring, b gostring) bool {
	if len(a) != len(b) {
		return false
	}

	for i:=0;i<len(a);i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func eqCstring(a cstring, b cstring) bool {
	return a == b
}


