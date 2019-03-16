package strngs

func HasSuffix(s string, suffix string) bool {
	if len(s) >= len(suffix) {
		var suf string
		suf = s[len(s)-len(suffix):]
		return suf == suffix
	}
	return false
}

func Contains(s string) bool {
}

func Split(s string, x string) []string {
}
