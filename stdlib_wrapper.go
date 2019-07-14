package main

// depends on libc
func strconv_Atoi(gs bytes) (int, error) {
	i, e := Atoi(string(gs))
	return i,e
}
