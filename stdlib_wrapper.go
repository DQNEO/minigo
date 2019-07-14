package main

// depends on libc
func strconv_Atoi(gs bytes) (int, error) {
	i, e := Atoi(gs)
	return i,e
}

func strings_HasSuffix(s string, suffix string) bool {
	return HasSuffix(s, (suffix))
}

func strings_Congtains(s bytes, substr bytes) bool {
	return Contains(s, substr)
}
