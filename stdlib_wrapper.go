package main

import (
	"io/ioutil"
	"strings"
	"strconv"
)

// depends on libc
func ioutil_ReadFile(filename gostring) ([]byte, error) {
	return ioutil.ReadFile(string(filename))
}

// depends on libc
func strconv_Atoi(gs gostring) (int, error) {
	i, e := strconv.Atoi(string(gs))
	return i,e
}

func strings_Split(s gostring, sep gostring) []gostring {
	return Split(s, sep)
}

func strings_HasSuffix(s gostring, suffix gostring) bool {
	return strings.HasSuffix(string(s), string(suffix))
}

func strings_Congtains(s gostring, substr gostring) bool {
	return Contains(s, substr)
}
