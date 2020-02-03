package syscall

func BytePtrFromString(s string) *byte {
	bs := []byte(s)
	var r *byte = &bs[0]
	return r
}

// 64-bit system call numbers
const __x64_sys_read = 0
const __x64_sys_write = 1
const __x64_sys_open = 2

func Open(path string, flag int, perm int) (int, error) {
	var fd int
	var _p0 *byte
	_p0 = BytePtrFromString(path)
	fd = syscall(__x64_sys_open, _p0, flag)
	return fd, nil
}

func Write(fd int, b []byte) (int, error) {
	var addr *byte = &b[0]
	var n int
	n = syscall(__x64_sys_write, fd, addr, len(b))
	return n, nil
}

func Read(fd int, b []byte) (int, error) {
	var ptr *byte
	ptr = &b[0]
	var nread int
	nread = syscall(__x64_sys_read, fd, ptr, cap(b))
	return nread, nil
}
