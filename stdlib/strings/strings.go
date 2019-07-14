package strings

// "foo/bar", "/" => []bytes{"foo", "bar"}
func Split(ss string, ssep string) []string {
	s := []byte(ss)
	sep := []byte(ssep)
	if len(sep) > 1  {
		panic("no supported")
	}
	seps := []byte(sep)
	sepchar := seps[0]
	vbytes := []byte(s)
	var buf []byte
	var r []string
	for _, b := range vbytes {
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

func HasSuffix(ss string, ssuffix string) bool {
	s := []byte(ss)
	suffix := []byte(ssuffix)
	if len(s) >= len(suffix) {
		var low int =  len(s)-len(suffix)
		var lensb int = len(s)
		var suf []byte
		sb := []byte(s)
		suf = sb[low:lensb]  // lensb is required
		return eq2([]byte(suf) , suffix)
	}
	return false
}

func eq2(a []byte, b []byte) bool {
	if len(a) != len(b) {
		return false
	}

	for i:=0;i<len(a);i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}