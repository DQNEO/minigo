package main

// depends on libc
func strconv_Atoi(gs string) (int, error) {
	i, e := Atoi(gs)
	return i,e
}
