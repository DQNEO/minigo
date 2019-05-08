package strngs

func HasSuffix(s string, suffix string) bool {
	if len(s) >= len(suffix) {
		var suf string
		suf = s[len(s)-len(suffix):]
		return suf == suffix
	}
	return false
}

func Contains(s string, substr string) bool {
	bytes := []byte(s)
	bsub := []byte(substr)
	var in bool
	var index int
	for _, b := range bytes {
		if !in && b == bsub[0] {
			in = true
			index = 0
		}

		if in {
			if b == bsub[index] {
				if index == len(bsub) - 1 {
					return true
				}
			} else {
				in = false
				index = 0
			}
		}
	}

	return false
}

// "foo/bar", "/" => []string{"foo", "bar"}
func Split(s string, sep string) []string {
	if len(sep) > 1  {
		panic("no supported")
	}
	seps := []byte(sep)
	sepchar := seps[0]
	bytes := []byte(s)
	var buf []byte
	var r []string
	for _, b := range bytes {
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
