package main

// depends on libc
func strconv_Atoi(gs bytes) (int, error) {
	i, e := Atoi(gs)
	return i,e
}

func strings_HasSuffix(s bytes, suffix bytes) bool {
	return HasSuffix(string(s), string(suffix))
}

func strings_Congtains(s bytes, substr bytes) bool {
	return Contains(s, substr)
}
