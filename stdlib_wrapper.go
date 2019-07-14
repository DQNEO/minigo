package main

// depends on libc
func strconv_Atoi(gs bytes) (int, error) {
	i, e := Atoi(gs)
	return i,e
}

func strings_Congtains(s string, substr string) bool {
	return Contains(string(s), string(substr))
}
