package syscall

func BytePtrFromString(s string) *byte {
	var r *byte = s
	return r
}

func Open(name string, flag int, perm int) (int, error) {
	var fd int
	var pchar *byte = BytePtrFromString(name)
	fd = open(pchar, flag)
	return fd, nil
}

