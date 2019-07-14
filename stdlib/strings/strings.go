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

