package runtime

import "unsafe"

func cstring2string(b *byte) string {
	var buf []byte
	for {
		if b == nil || *b == 0 {
			break
		}
		buf = append(buf, *b)
		p := uintptr(unsafe.Pointer(b)) + 1
		b = (*byte)(unsafe.Pointer(p))
	}
	return string(buf)
}

func eq(a string, b string) bool {
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

func cmpStrings(a string, b string, flag bool) bool {
	var ret bool
	ret = eq(a, b)
	if flag {
		return ret
	} else {
		return !ret
	}
}

func concat(as string, bs string) string {
	var r []byte
	for i := 0; i < len(as); i++ {
		r = append(r, as[i])
	}
	for i := 0; i < len(bs); i++ {
		r = append(r, bs[i])
	}
	return string(r)
}
