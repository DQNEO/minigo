package strconv

var __strconv_trash int

func Atoi(gs string) (int, error) {
	bts := []byte(gs)
	if len(bts) == 0 {
		return 0,nil
	}
	var b byte
	var n int
	var i int
	var isMinus bool
	for i, b = range bts {
		if b == '.' {
			return 0,nil // @FIXME all no number should return error
		}
		if b == '-' {
			isMinus = true
			continue
		}
		var x byte = b - byte('0')
		n  = n * 10
		n = n + int(x)
	}
	if isMinus {
		n = -n
	}
	__strconv_trash = i
	return n, nil
}
