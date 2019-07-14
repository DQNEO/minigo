package main

// depends on libc
func strconv_Atoi(gs bytes) (int, error) {
	i, e := Atoi(gs)
	return i,e
}

func strings_Split(s bytes, sep bytes) []string {
	return Split(string(s), string(sep))
}

func strings_HasSuffix(s bytes, suffix bytes) bool {
	return HasSuffix(s, suffix)
}

func strings_Congtains(s bytes, substr bytes) bool {
	return Contains(s, substr)
}
