package runtime

var argv []*byte

func runtime_args() []string {
	var r []string
	for _, a := range argv {
		// Convert *byte to string
		var s string = cstring2string(a)
		r = append(r, s)
	}
	return r
}

