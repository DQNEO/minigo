package main

import "fmt"

type gostring []byte
type cstring string

func GoSprintf(format gostring, param ...interface{}) gostring {
	s := fmt.Sprintf(string(format), param...)
	return gostring(s)
}

func convertCstringsToGostrings(cstrings []string) []gostring {
	var r []gostring
	for _, cs := range cstrings {
		r = append(r, gostring(cs))
	}

	return r
}

func catGostrings(a gostring, b gostring) gostring {
	var c []byte
	for i:=0;i<len(a);i++ {
		c = append(c, a[i])
	}
	for i:=0;i<len(b);i++ {
		c = append(c, b[i])
	}
	return c
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


