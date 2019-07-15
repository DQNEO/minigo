package syscall

func BytePtrFromString(s string) *byte {
	var r *byte = s
	return r
}

func Open(path string, flag int, perm int) (int, error) {
	var fd int
	var _p0 *byte
	_p0 = BytePtrFromString(path)
	fd = open(_p0, flag)
	return fd, nil
}
