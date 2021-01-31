package strings

// "foo/bar", "/" => []string{"foo", "bar"}
func Split(s string, ssep string) []string {
	if len(ssep) > 1 {
		panic("no supported")
	}
	sepchar := ssep[0]
	var buf []byte
	var r []string
	for _, b := range []byte(s) {
		if b == sepchar {
			r = append(r, string(buf))
			buf = nil
		} else {
			buf = append(buf, b)
		}
	}
	r = append(r, string(buf))

	return r
}

func HasPrefix(s string, prefix string) bool {
	for i, bp := range []byte(prefix) {
		if bp != s[i] {
			return false
		}
	}
	return true
}

func HasSuffix(s string, suffix string) bool {
	if len(s) >= len(suffix) {
		var low int = len(s) - len(suffix)
		var lensb int = len(s)
		var suf []byte
		sb := []byte(s)
		suf = sb[low:lensb] // lensb is required
		return eq2(suf, []byte(suffix))
	}
	return false
}

func eq2(a []byte, b []byte) bool {
	if len(a) != len(b) {
		return false
	}

	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// Contains reports whether substr is within s.
func Contains(s string, substr string) bool {
	return Index(s, substr) >= 0
}

func Index(s string, substr string) int {
	var in bool
	var subIndex int
	var r int = -1 // not found
	for i, b := range []byte(s) {
		if !in && b == substr[0] {
			in = true
			r = i
			subIndex = 0
		}

		if in {
			if b == substr[subIndex] {
				if subIndex == len(substr)-1 {
					return r
				}
			} else {
				in = false
				r = -1
				subIndex = 0
			}
		}
	}

	return -1
}

// search index of the specified char from backward
func LastIndexByte(s string, c byte) int {
	for i:=len(s)-1;i>=0;i-- {
		if s[i] == c {
			return i
		}
	}
	// not found
	return -1
}
