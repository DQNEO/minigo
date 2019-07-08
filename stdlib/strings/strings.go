package strings

func HasSuffix(s string, suffix string) bool {
	if len(s) >= len(suffix) {
		var suf string
		suf = s[len(s)-len(suffix):]
		return suf == suffix
	}
	return false
}
