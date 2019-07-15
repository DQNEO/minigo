package syscall

func Open(name string, flag int, perm int) (int, error) {
	var fd int
	var pchar *byte = name
	fd = open(pchar, flag)
	return fd, nil
}

