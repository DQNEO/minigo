package strconv

var __strconv_trash int

func Atoi(gs string) (int, error) {
	if len(gs) == 0 {
		return 0, nil
	}
	var b byte
	var n int
	var i int
	var isMinus bool
	for i, b = range []byte(gs) {
		if b == '.' {
			return 0, nil // @FIXME all no number should return error
		}
		if b == '-' {
			isMinus = true
			continue
		}
		var x byte = b - byte('0')
		n = n * 10
		n = n + int(x)
	}
	if isMinus {
		n = -n
	}
	__strconv_trash = i
	return n, nil
}

func Itoa(i int) string {
	var r []byte
	var tmp []byte

	if i < 0 {
		i = i * -1
		r = append(r, '-')
	}
	for i > 0 {
		mod := i % 10
		tmp = append(tmp, byte('0')+byte(mod))
		i = i / 10
	}

	for j := len(tmp) - 1; j >= 0; j-- {
		r = append(r, tmp[j])
	}

	if len(r) == 0 {
		return "0"
	}
	return string(r)
}
