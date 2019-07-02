package main

import (
	"io/ioutil"
	"strings"
)

func ioutil_ReadFile(filename gostring) ([]byte, error) {
	return ioutil.ReadFile(string(filename))
}

func strings_Split(s gostring, sep gostring) []gostring {
	css := strings.Split(string(s), string(sep))
	return convertCstringsToGostrings(css)
}

func strings_HasSuffix(s gostring, suffix gostring) bool {
	return strings.HasSuffix(string(s), string(suffix))
}
