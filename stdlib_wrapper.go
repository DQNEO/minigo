package main

import (
	"io/ioutil"
)

// depends on libc
func ioutil_ReadFile(filename bytes) ([]byte, error) {
	return ioutil.ReadFile(string(filename))
}

// depends on libc
func strconv_Atoi(gs bytes) (int, error) {
	i, e := Atoi((gs))
	return i,e
}

func strings_Split(s bytes, sep bytes) []bytes {
	return Split(s, sep)
}

func strings_HasSuffix(s bytes, suffix bytes) bool {
	return HasSuffix(s, suffix)
}

func strings_Congtains(s bytes, substr bytes) bool {
	return Contains(s, substr)
}
