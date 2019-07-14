package main

import "strings"

// depends on libc
func strconv_Atoi(gs bytes) (int, error) {
	i, e := Atoi(gs)
	return i,e
}

func strings_Congtains(s string, substr string) bool {
	return strings.Contains(string(s), string(substr))
}
