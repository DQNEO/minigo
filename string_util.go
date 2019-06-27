package main

type gostring []byte
type cstring string

func catGostrings(a []byte, b []byte) []byte {
	var c []byte
	for i:=0;i<len(a);i++ {
		c = append(c, a[i])
	}
	for i:=0;i<len(b);i++ {
		c = append(c, b[i])
	}
	return c
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


